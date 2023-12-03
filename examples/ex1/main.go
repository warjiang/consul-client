package main

import (
	"fmt"

	"github.com/warjiang/consul-client/pkg/consul"
)

func main() {
	fmt.Println("hello world")
	ets, err := consul.Lookup("vpc.demo.nginx")
	if err != nil {
		panic(err)
	}
	for _, et := range ets {
		fmt.Printf("endpoints %+v\n", et)
	}
}
