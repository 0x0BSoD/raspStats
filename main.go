package main

import (
	"fmt"
	"github.com/gorilla/websocket"
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

func startScheduler(ws *websocket.Conn) {
	gocron.Every(1).Minute().Do(scheduledTask, ws)
	_, _time := gocron.NextRun()
	fmt.Println("Next Run: ", _time)
	<-gocron.Start()
}

func scheduledTask(ws *websocket.Conn) {
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
