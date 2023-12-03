package consul

import (
	"encoding/json"
	"math/rand"
	"net"
	"os"
	"time"
)

func GetNow() time.Time {
	return time.Now()
}

func GetRandom() *rand.Rand {
	// avoid use rand.Seed(time.Now().UnixNano())
	src := rand.NewSource(time.Now().UnixNano())
	return rand.New(src)
}

func GetAddr() string {
	host := GetEnvOrDefault("CONSUL_HTTP_HOST", "127.0.0.1")
	port := GetEnvOrDefault("CONSUL_HTTP_PORT", "8500")
	return net.JoinHostPort(host, port)
}

func GetEnvOrDefault(name string, defaultValue string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return defaultValue
}

func GetFromMapOrDefault[T any](m map[string]T, key string, defaultValue T) T {
	v, exist := m[key]
	if exist {
		return v
	}
	return defaultValue
}

func MarshalWithDefault(o any, defaultValue string) string {
	buff, err := json.Marshal(o)
	if err != nil {
		return defaultValue
	}
	return string(buff)
}
