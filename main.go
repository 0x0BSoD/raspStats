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

func startScheduler() {
	gocron.Every(1).Minute().Do(scheduledTask)
	_, _time := gocron.NextRun()
	fmt.Println("Next Run: ", _time)
	<-gocron.Start()
}

func scheduledTask() {
	uptime, err := stats.GetUptime()
	poe(err)

	loadAvg, err := stats.GetLoadAvg()
	poe(err)

	db, err := DBConn()
	poe(err)

	err = StoreItem(db, DBItem{
		Uptime:  uptime,
		LoadAvg: loadAvg,
	})
	poe(err)

	_, err = GetAllItems(db)
	poe(err)

	err = db.Close()
	poe(err)
}

func main() {
	//startScheduler()
	for {
		err := stats.GetCpuLoad(1 * time.Second)
		poe(err)
	}

}
