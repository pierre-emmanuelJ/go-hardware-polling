package polls

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pierre-emmanuelJ/go-exercises/cpu"
	"github.com/pierre-emmanuelJ/go-exercises/memory"
)

type Metric struct {
	Name   string `json:"name"`
	Metric string `json:"metric"`
}

type Metrics struct {
	Timestamp string    `json:"timestamp"`
	Metrics   []*Metric `json:"metrics"`
}

func Pollsinfos(partition, iNetwork string, cpuInfos *cpu.CPUInfos) error {

	av, err := getLoadAv()
	if err != nil {
		return err
	}

	cpu, err := getCPUPercentage(cpuInfos)
	if err != nil {
		return err
	}

	netStat, err := getNetStat()
	if err != nil {
		return err
	}

	disk, err := getDiskUsage()
	if err != nil {
		return err
	}

	memUse, err := getMemUsage()
	if err != nil {
		return err
	}

	metrics := &Metrics{Timestamp: time.Now().String(), Metrics: []*Metric{av, cpu, netStat, disk, memUse}}

	b, err := json.Marshal(metrics)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Printf("%s\n", b)
	return nil
}

func getLoadAv() (*Metric, error) {

	file, err := os.Open("/proc/loadavg")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	scanner.Scan()
	return &Metric{Name: "Load average", Metric: scanner.Text()}, nil
}

func getCPUPercentage(cpuInfos *cpu.CPUInfos) (*Metric, error) {

	file, err := os.Open("/proc/stat")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "cpu") {
			break
		}
		return nil, fmt.Errorf("No cpu infos in this file or invalide file")
	}

	cpuTimes, err := cpu.GetCPUTimes(scanner.Text())
	if err != nil {
		return nil, err
	}

	idleTime, totalTime := 0, 0

	if err := cpu.GetCPUIdleTimes(&idleTime, &totalTime, cpuTimes); err != nil {
		return nil, err
	}

	idleTimeDelta := 1.0 * (idleTime - cpuInfos.PreviousIdleTime)
	totalTimeDelta := 1.0 * (totalTime - cpuInfos.PreviousTotalTime)
	utilization := (1000.0*(totalTimeDelta-idleTimeDelta)/totalTimeDelta + 5) / 10

	result := fmt.Sprintf("%.2f", float64(utilization))

	cpuInfos.PreviousIdleTime = idleTime
	cpuInfos.PreviousTotalTime = totalTime

	return &Metric{Name: "Cpu user", Metric: result}, nil
}

func getNetStat() (*Metric, error) {
	return &Metric{}, nil

}

func getDiskUsage() (*Metric, error) {
	return &Metric{}, nil

}

func getMemUsage() (*Metric, error) {

	memTotal, err := memory.GetMemoryInfoByKey("MemTotal")
	if err != nil {
		return nil, err
	}
	memAvailable, err := memory.GetMemoryInfoByKey("MemAvailable")
	if err != nil {
		return nil, err
	}

	memUtilization := (memTotal / (memTotal - memAvailable) * 100)

	return &Metric{Name: "Memory utilization", Metric: strconv.FormatInt(int64(memUtilization), 10)}, nil

}
