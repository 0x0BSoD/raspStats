package stats

import (
	"strconv"
	"strings"
)

// LoadAvg result
type LoadAvg struct {
	RunningProc int      `json:"running_proc"`
	TotalProc   int      `json:"total_proc"`
	Load        loadItem `json:"load"`
}

type loadItem struct {
	OneM     float64 `json:"one_m"`
	FiveM    float64 `json:"five_m"`
	FifteenM float64 `json:"fifteen_m"`
}

// GetLoadAvg return formatted data from /proc/loadavg
func GetLoadAvg() (LoadAvg, error) {
	dat, err := openFile("/proc/loadavg")
	if err != nil {
		return LoadAvg{}, err
	}

	strData := strings.Split(strings.Replace(dat, "\n", "", 1), " ")

	running, _ := strconv.Atoi(strings.Split(strData[3], "/")[0])
	total, _ := strconv.Atoi(strings.Split(strData[3], "/")[1])

	o, _ := strconv.ParseFloat(strData[0], 64)
	v, _ := strconv.ParseFloat(strData[1], 64)
	t, _ := strconv.ParseFloat(strData[2], 64)

	return LoadAvg{
		RunningProc: running,
		TotalProc:   total,
		Load: loadItem{
			OneM:     o,
			FiveM:    v,
			FifteenM: t,
		},
	}, nil
}
