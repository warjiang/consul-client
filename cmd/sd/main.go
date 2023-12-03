package main

import (
	"os"

	"github.com/warjiang/consul-client/cmd/sd/cmd"
)

func main() {
	c := cmd.NewServiceDiscoveryCommand()
	if err := c.Execute(); err != nil {
		//fmt.Println(err)
		os.Exit(1)
	}
}
