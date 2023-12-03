package consul

import "math/rand"

var (
	consulClient *ConsulClient
	random       *rand.Rand
)

func init() {
	random = GetRandom()
	consulClient = NewConsulClient(GetAddr())
}

func Lookup(name string, opts ...LookupOption) (Endpoints, error) {
	return consulClient.Lookup(name, opts...)
}
