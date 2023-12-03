package consul

import (
	"fmt"
	"net"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type LookupParams struct {
	cluster    string
	limit      int
	tag        string
	nocache    bool
	consulHash string
}

func (o LookupParams) ConvertToUrl() *url.Values {
	uv := &url.Values{}
	if o.limit > 0 {
		uv.Set("limit", strconv.Itoa(o.limit))
	}
	if o.cluster != "" {
		uv.Set("cluster", o.cluster)
	}
	if o.tag != "" {
		uv.Set("tag", o.tag)
	}
	return uv
}

type LookupOption func(params *LookupParams)

func WithCluster(cluster string) LookupOption {
	return func(params *LookupParams) {
		params.cluster = cluster
	}
}
func WithLimit(n int) LookupOption {
	return func(params *LookupParams) {
		params.limit = n
	}
}

func WithNoCache(nocache bool) LookupOption {
	return func(params *LookupParams) {
		params.nocache = nocache
	}
}

func WithTag(tag string) LookupOption {
	return func(params *LookupParams) {
		params.tag = tag
	}
}

type Tags map[string]string

func (t Tags) ToString() string {
	if len(t) == 0 {
		return "[]"
	}
	keys := make([]string, 0, len(t))
	for k := range t {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		ki := keys[i]
		kj := keys[j]
		vi, viExist := t[ki]
		vj, vjExist := t[kj]

		if viExist == vjExist {
			return vi >= vj
		} else {
			if viExist {
				return true
			}
			return false
		}
	})

	tags := "["
	for _, k := range keys {
		if strings.HasPrefix(k, "__up_") {
			continue
		}
		if t[k] != "" {
			tags = fmt.Sprintf("%s\"%s\":\"%s\",", tags, k, t[k])
		} else {
			tags = fmt.Sprintf("%s\"%s\",", tags, k)
		}
	}
	return tags[0:len(tags)-1] + "]"
}

type Endpoint struct {
	Host    string
	Port    int
	Addr    string
	Cluster string
	Env     string
	Weight  int
	Tags    Tags
}

type Endpoints []Endpoint

func (endpoints Endpoints) FilterCluster(name string) Endpoints {
	ret := make([]Endpoint, 0, len(endpoints))
	for _, e := range endpoints {
		if e.Cluster == name {
			ret = append(ret, e)
		}
	}
	return ret
}

func (endpoints Endpoints) Parse() {}

type ConsulEndpoint struct {
	Host string
	Port int
	Tags Tags
}

func (e *ConsulEndpoint) Parse() Endpoint {
	ret := Endpoint{
		Host:    e.Host,
		Port:    e.Port,
		Addr:    net.JoinHostPort(e.Host, strconv.Itoa(e.Port)),
		Tags:    e.Tags,
		Cluster: GetFromMapOrDefault(e.Tags, "cluster", ""),
		Env:     GetFromMapOrDefault(e.Tags, "env", ""),
	}

	if w, err := strconv.Atoi(e.Tags["weight"]); err == nil {
		ret.Weight = w
	} else {
		ret.Weight = DefaultWeight
	}
	return ret
}

type CachedEndpoints struct {
	UpdatedAt  time.Time
	Endpoints  Endpoints
	ConsulHash string
}
