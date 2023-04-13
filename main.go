package main

import (
	"crowbrute/crow"
	"crowbrute/crow/config"
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	config, err := config.ReadConfig()
	if err != nil {
		fmt.Println("The cfg.toml is not filled in correctly")
		return
	}

	servers, err := crow.ReadServers()

	if err != nil {
		fmt.Println("Can't find the servers.txt file")
		return
	}

	wg.Add(1)
	go crow.StartCore(config, servers)

	wg.Wait()
}
