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

//https://stackoverflow.com/questions/23367857/accurate-calculation-of-cpu-usage-given-in-percentage-in-linux
//
//PrevNonIdle = prevuser + prevnice + prevsystem + previrq + prevsoftirq + prevsteal
//NonIdle = user + nice + system + irq + softirq + steal
//
//PrevTotal = PrevIdle + PrevNonIdle
//Total = Idle + NonIdle
//
//# differentiate: actual value minus the previous one
//totald = Total - PrevTotal
//idled = Idle - PrevIdle
//
//CPU_Percentage = (totald - idled)/totald

func GetCpuLoad(interval time.Duration) error {
	dat, err := openFile("/proc/stat")
	if err != nil {
		return err
	}

	totalPre, err := calc(strings.Split(dat, "\n")[0])
	if err != nil {
		return err
	}
	fmt.Println(totalPre)

	time.Sleep(interval)

	dat, err = openFile("/proc/stat")
	if err != nil {
		return err
	}
	totalPost, err := calc(strings.Split(dat, "\n")[0])
	if err != nil {
		return err
	}
	fmt.Println(totalPost)

	//PrevIdle = previdle + previowait
	preIdle := totalPre.Idle + totalPre.IOWait

	//Idle = idle + iowait
	idle := totalPost.Idle + totalPost.IOWait

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
