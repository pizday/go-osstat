package main

import (
	"github.com/mackerelio/go-osstat/network"
)

type networkGenerator struct {
	networks []network.NetworkStats
	err      error
}

func (gen *networkGenerator) Get() {
	networks, err := network.Get()
	if err != nil {
		gen.err = err
		return
	}
	gen.networks = networks
}

func (gen *networkGenerator) Error() error {
	return gen.err
}

func (gen *networkGenerator) Print(out chan<- value) {
	for _, network := range gen.networks {
		out <- value{"network." + network.Name + ".rx_bytes", network.RxBytes, "bytes"}
		out <- value{"network." + network.Name + ".tx_bytes", network.TxBytes, "bytes"}
	}
}
