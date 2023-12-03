package consul

import (
	"errors"
	"time"
)

var (
	ErrNoEndpoint        = errors.New("no endpoint")
	ErrNotModified       = errors.New("consul not modified")
	ErrConnectionRefused = errors.New("connection refused (try set CONSUL_HTTP_HOST in env)")
)

const (
	HeaderConsulHash = "X-Consul-Result-Hash"
	CacheTime        = 15 * time.Second
	DefaultWeight    = 50
)
