package crow

import (
	"crowbrute/crow/config"
	"sync"
)

type Core struct {
	wg           sync.WaitGroup
	Config       *config.Config
	BalancerChan chan bool
	Servers      []string
}

var Divc = make(chan int)

func StartCore(config *config.Config, servers []string) {

	var c = &Core{
		Config:       config,
		BalancerChan: make(chan bool),
		Servers:      servers,
		wg:           sync.WaitGroup{},
	}

	go c.Balancer()

	c.Start()

}
