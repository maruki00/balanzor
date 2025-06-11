package main

import (
	"github.com/maruki00/balanzor"
)

func main() {
	// backends := []types.Server{}

	servers := []string{"http://localhost:9090", "http://localhost:9091", "http://localhost:9092"}
	algo := "round-roubin"
	lb := balanzor.NewBalanzor(servers, algo, "0.0.0.0:9999", "/lb")
	lb.Run()
}
