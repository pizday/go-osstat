package main

import (
	"fmt"
	"os"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/disk"
	"github.com/mackerelio/go-osstat/loadavg"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/mackerelio/go-osstat/network"
	"github.com/mackerelio/go-osstat/uptime"
)

func main() {
	memory, err := memory.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	cpu, err := cpu.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	disk, err := disk.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	uptime, err := uptime.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	network, err := network.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	loadavg, err := loadavg.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	fmt.Printf("memory total: %d bytes\n", memory.Total)
	fmt.Printf("memory used: %d bytes\n", memory.Used)
	fmt.Printf("memory cached: %d bytes\n", memory.Cached)
	fmt.Printf("memory free: %d bytes\n", memory.Free)

	fmt.Printf("cpu total: %d bytes\n", cpu.Total)
	fmt.Printf("cpu user: %d bytes\n", cpu.User)

	fmt.Printf("cpu system: %d bytes\n", cpu.System)
	fmt.Printf("cpu Idle: %d bytes\n", cpu.Idle)

	fmt.Printf("disk total: %d bytes\n", disk.Total)
	// fmt.Printf("disk free: %d bytes\n", disk.free)
	// fmt.Printf("disk used: %d bytes\n", disk.used)

	fmt.Printf("uptime : %.2f minutes\n", uptime.Minutes())
	for _, v := range network {
		fmt.Println("network", v.Name, v.RxBytes, v.TxBytes)

	}
	fmt.Println(*loadavg)

}
