package stats

import (
	"strconv"
	"strings"
)

// Uptime contain uptime and idle time, idel for all cores
type Uptime struct {
	Uptime int `json:"uptime"`
	Idle   int `json:"idle"`
}

// GetUptime return formatted data from /proc/uptime
func GetUptime() (Uptime, error) {
	dat, err := openFile("/proc/uptime")
	if err != nil {
		return Uptime{}, err
	}

	strData := strings.Split(dat, " ")

	up, err := strconv.Atoi(strings.Split(strData[0], ".")[0])
	if err != nil {
		return Uptime{}, err
	}

	idle, err := strconv.Atoi(strings.Trim(strings.Split(strData[1], ".")[0], "\n"))
	if err != nil {
		return Uptime{}, err
	}

	return Uptime{
		Uptime: up,
		Idle:   idle,
	}, nil
}
