package consul

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ConsulClient struct {
	client *http.Client
	lock   sync.RWMutex
	cache  map[string]*CachedEndpoints
}

func NewConsulClient(addr string) *ConsulClient {
	dialer := net.Dialer{Timeout: 3 * time.Second}
	httpClient := &http.Client{Timeout: 500 * time.Millisecond}
	httpClient.Transport = &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return dialer.DialContext(ctx, "tcp", addr)
		},
	}
	return NewConsulClientWithHttpClient(httpClient)
}

func NewConsulClientWithHttpClient(client *http.Client) *ConsulClient {
	return &ConsulClient{
		client: client,
		cache:  make(map[string]*CachedEndpoints),
	}
}

func (c *ConsulClient) Lookup(name string, opts ...LookupOption) (Endpoints, error) {
	t := GetNow()

	params := LookupParams{}
	for _, op := range opts {
		op(&params)
	}

	key := strings.Join([]string{name, params.cluster}, "|")

	c.lock.RLock()
	item := c.cache[key]
	c.lock.RUnlock()
	if item != nil {
		if t.Sub(item.UpdatedAt) < CacheTime && !params.nocache {
			return item.Endpoints, nil
		}
		if len(item.ConsulHash) > 0 {
			params.consulHash = item.ConsulHash
		}
	}

	endpoints, hash, err := c.lookup(name, params)
	if err != nil && err != ErrNoEndpoint {
		if item != nil {
			newItem := *item
			newItem.UpdatedAt = t
			c.lock.Lock()
			c.cache[key] = &newItem
			c.lock.Unlock()
		}
		if err == ErrNotModified {
			return item.Endpoints, nil
		}
		if params.nocache {
			return nil, err
		}
		if item != nil {
			return item.Endpoints, nil
		}
		return nil, err
	}
	ret := make(Endpoints, len(endpoints))
	for i, e := range endpoints {
		ret[i] = e.Parse()
	}
	if params.cluster != "" {
		ret = ret.FilterCluster(params.cluster)
	}
	c.lock.Lock()
	c.cache[key] = &CachedEndpoints{
		Endpoints:  ret,
		ConsulHash: hash,
		UpdatedAt:  t,
	}
	c.lock.Unlock()
	return ret, nil
}

func (c *ConsulClient) lookup(name string, cfg LookupParams) ([]ConsulEndpoint, string, error) {
	uv := cfg.ConvertToUrl()
	uv.Set("name", name)
	u := "http://127.0.0.1:8500/v1/lookup/name?" + uv.Encode()
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, "", err
	}
	if len(cfg.consulHash) > 0 {
		req.Header.Set(HeaderConsulHash, cfg.consulHash)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		if s := err.Error(); strings.Contains(s, "connect: connection refused") {
			return nil, "", ErrConnectionRefused
		}
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotModified {
			return nil, "", ErrNotModified
		}
		b, _ := io.ReadAll(resp.Body)
		return nil, "", errors.New(string(b))
	}
	hash := resp.Header.Get(HeaderConsulHash)

	endpoints := make([]ConsulEndpoint, 0, 2)
	if err := json.NewDecoder(resp.Body).Decode(&endpoints); err != nil {
		return nil, "", err
	}
	if len(endpoints) == 0 {
		return nil, "", ErrNoEndpoint
	}
	return endpoints, hash, nil
}
