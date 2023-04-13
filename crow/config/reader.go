package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Balancer struct {
		Enabled       bool
		BalancerValue int
	}
	Bruteforce struct {
		Timeout             int
		GenerateHost        bool
		DictionaryPath      string
		GenererateHostCount int
		MaxAttempts         int
		ServersMode         int
		Delay               int
		ResultFormat        string
	}
	RandomPassword struct {
		Enabled              bool
		RandomPasswordLen    int
		RandomCustomPassword string
	}

	Payload struct {
		Enabled bool
		Payload string
	}
}

func ReadConfig() (*Config, error) {

	var conf Config

	_, err := toml.DecodeFile("cfg.toml", &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil

}
