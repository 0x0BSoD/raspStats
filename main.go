package main

import (
	"fmt"
	"time"
)

import (
	"github.com/0x0bsod/raspStats/stats"
	"github.com/jasonlvhit/gocron"
)

// poe - Panic on error
func poe(err error) {
	if err != nil {
		panic(err)
	}
}

func startScheduler(sessWS *socket) {
	gocron.Every(1).Minute().Do(scheduledTask, sessWS)
	_, _time := gocron.NextRun()
	fmt.Println("Next Run: ", _time)
	<-gocron.Start()
}

func scheduledTask(sessWS *socket) {
	uptime, err := stats.GetUptime()
	poe(err)

	loadAvg, err := stats.GetLoadAvg()
	poe(err)

	cpuLoad, err := stats.GetCpuLoad(1 * time.Second)
	poe(err)

	db, err := DBConn()
	poe(err)

	err = StoreItem(db, DBItem{
		Uptime:  uptime,
		LoadAvg: loadAvg,
		CpuLoad: cpuLoad,
	})
	poe(err)

	err = db.Close()
	poe(err)
}

func main() {
	Server()
}
