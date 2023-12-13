package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Conf struct {
	Peers            []string `yaml:"Peers"`
	TotalNodes       int      `yaml:"TotalNodes"`
	LeaderTimeout    int      `yaml:"LeaderTimeout"`
	CandidateTimeout int      `yaml:"CandidateTimeout"`
	// Port             string   `yaml:"Port"`
}

var conf *Conf

func initConf(path string) *Conf {
	conf = new(Conf)
	configFile, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(configFile, conf)
	if err != nil {
		fmt.Printf("Error parsing config file: %s\n", err)
		os.Exit(1)
	}
	return conf
}

// func initConf() *Conf {
// 	conf = new(Conf)
// 	conf.Peers[0] = "10.1.8.1"
// 	conf.Peers[1] = "10.1.8.2"
// 	conf.Peers[2] = "10.1.8.3"
// 	conf.TotalNodes = 3
// 	conf.LeaderTimeout = 1000
// 	conf.CandidateTimeout = 1000
// 	return conf
// }
