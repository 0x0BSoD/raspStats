// https://www.kernel.org/doc/Documentation/filesystems/proc.txt
//https://stackoverflow.com/questions/23367857/accurate-calculation-of-cpu-usage-given-in-percentage-in-linux

package stats

import (
	"fmt"
	"github.com/0x0bsod/strNorm"
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

func calc(data string) (cpuItem, error) {
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

// GetCpuLoad return calculated cpu load from /proc/stat
func GetCpuLoad(interval time.Duration) error {
	dat, err := openFile("/proc/stat")
	if err != nil {
		return err
	}

	// first probe
	totalPre, err := calc(strings.Split(dat, "\n")[0])
	if err != nil {
		return err
	}

	// wait
	time.Sleep(interval)

	// second probe
	dat, err = openFile("/proc/stat")
	if err != nil {
		return err
	}
	totalPost, err := calc(strings.Split(dat, "\n")[0])
	if err != nil {
		return err
	}

	preIdle := totalPre.Idle + totalPre.IOWait
	idle := totalPost.Idle + totalPost.IOWait

	preNonIdle := totalPre.User + totalPre.Nice + totalPre.System + totalPre.Irq + totalPre.SoftIrq + totalPre.Steal
	nonIdle := totalPost.User + totalPost.Nice + totalPost.System + totalPost.Irq + totalPost.SoftIrq + totalPost.Steal

	prevTotal := preIdle + preNonIdle
	total := idle + nonIdle

	totalD := total - prevTotal
	idled := idle - preIdle

	percents := ((float64(totalD) - float64(idled)) / float64(totalD)) * 100.0

	fmt.Printf("Percents: %.2f%%\n", percents)

	//for _, i := range strData[1:] {
	//	if strings.HasPrefix(i, "cpu") {
	//		seconds, err := calc(i)
	//		if err != nil {
	//			return err
	//		}
	//
	//		fmt.Println(i)
	//		fmt.Println(seconds)
	//	}
	//}

	return nil
}
