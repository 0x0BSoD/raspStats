// https://www.kernel.org/doc/Documentation/filesystems/proc.txt
//https://stackoverflow.com/questions/23367857/accurate-calculation-of-cpu-usage-given-in-percentage-in-linux

package stats

import (
	"fmt"
	"github.com/0x0bsod/strNorm"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Time units are in USER_HZ (typically hundredths of a second)
type cpuItem struct {
	User      int // normal processes executing in user mode
	Nice      int // niced processes executing in user mode
	System    int // processes executing in kernel mode
	Idle      int // twiddling thumbs
	IOWait    int // waiting for I/O to complete https://www.kernel.org/doc/Documentation/filesystems/proc.txt
	Irq       int // servicing interrupts
	SoftIrq   int // servicing softirqs
	Steal     int // involuntary wait
	Guest     int // running a normal guest
	GuestNice int // running a niced guest
}

func convertIt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return i
}

func createStruct(data string) (cpuItem, error) {
	normString := strNorm.Normalize(data)
	items := strings.Split(normString, " ")[1:]

	return cpuItem{
		User:      convertIt(items[0]),
		Nice:      convertIt(items[1]),
		System:    convertIt(items[2]),
		Idle:      convertIt(items[3]),
		IOWait:    convertIt(items[4]),
		Irq:       convertIt(items[5]),
		SoftIrq:   convertIt(items[6]),
		Steal:     convertIt(items[7]),
		Guest:     convertIt(items[8]),
		GuestNice: convertIt(items[9]),
	}, nil
}

func calc(totalPre, totalPost cpuItem) float64 {
	preIdle := totalPre.Idle + totalPre.IOWait
	idle := totalPost.Idle + totalPost.IOWait

	preNonIdle := totalPre.User + totalPre.Nice + totalPre.System + totalPre.Irq + totalPre.SoftIrq + totalPre.Steal
	nonIdle := totalPost.User + totalPost.Nice + totalPost.System + totalPost.Irq + totalPost.SoftIrq + totalPost.Steal

	prevTotal := preIdle + preNonIdle
	total := idle + nonIdle

	totalD := total - prevTotal
	idled := idle - preIdle

	percents := ((float64(totalD) - float64(idled)) / float64(totalD)) * 100.0

	return percents
}

// GetCpuLoad return calculated cpu load from /proc/stat
func GetCpuLoad(interval time.Duration) error {
	coresPre := make([]cpuItem, runtime.NumCPU())
	coresPost := make([]cpuItem, runtime.NumCPU())

	// first probe
	dat, err := openFile("/proc/stat")
	if err != nil {
		return err
	}
	cpuStatsFirst := strings.Split(dat, "\n")
	// total
	totalPre, err := createStruct(cpuStatsFirst[0])
	if err != nil {
		return err
	}
	// other CPU's
	for idx, i := range cpuStatsFirst[1:] {
		if strings.HasPrefix(i, "cpu") {
			pre, err := createStruct(i)
			if err != nil {
				return err
			}
			coresPre[idx] = pre
		}
	}
	// ===========================================================

	// wait
	time.Sleep(interval)

	// second probe
	dat, err = openFile("/proc/stat")
	if err != nil {
		return err
	}
	cpuStatsSecond := strings.Split(dat, "\n")
	// total
	totalPost, err := createStruct(cpuStatsSecond[0])
	if err != nil {
		return err
	}
	// other CPU's
	for idx, i := range cpuStatsSecond[1:] {
		if strings.HasPrefix(i, "cpu") {
			post, err := createStruct(i)
			if err != nil {
				return err
			}
			coresPost[idx] = post
		}
	}
	// ===========================================================

	fmt.Printf("Overall CPU usage: %.2f%%\n", calc(totalPre, totalPost))

	for i := 0; i < runtime.NumCPU(); i++ {
		fmt.Printf("CPU %d usage: %.2f%%\n", i, calc(coresPre[i], coresPost[i]))
	}
	fmt.Println("---")

	return nil
}
