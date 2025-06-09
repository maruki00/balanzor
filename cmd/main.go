package main

import (
	"github.com/maruki00/balanzor"
	"github.com/maruki00/balanzor/types"
)

func main() {
	cfg, err := types.NewConfig("config.yaml")
	if err != nil {
		panic(err)
	}
	lb := balanzor.NewBalanzor(cfg.Servers, cfg.Algo, "0.0.0.0:9999", "/lb")
	lb.Run()
}
