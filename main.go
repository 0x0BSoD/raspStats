package main

import "fmt"

import "github.com/0x0bsod/raspStats/stats"

// poe - Panic on error
func poe(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	uptime, err := stats.GetUptime()
	poe(err)

	loadAvg, err := stats.GetLoadAvg()
	poe(err)

	fmt.Printf("%+v\n", uptime)
	fmt.Printf("%+v\n", loadAvg)
}
