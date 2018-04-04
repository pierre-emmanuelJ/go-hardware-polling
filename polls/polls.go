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
		return nil, nil
	}

	cpuTimes, err := cpu.GetCPUTimes(scanner.Text())
	if err != nil {
		return nil, err
	}

	idleTime, totalTime := 0, 0

	if err := cpu.GetCPUIdleTimes(&idleTime, &totalTime, cpuTimes); err != nil {
		return nil, err
	}

	println(idleTime, totalTime)
	result := ""
	idleTimeDelta := idleTime - cpuInfos.PreviousIdleTime
	totalTimeDelta := totalTime - cpuInfos.PreviousTotalTime
	utilization := 100.0 * (1.0 - idleTimeDelta/totalTimeDelta)

	result += strconv.FormatFloat(float64(utilization), 'E', -1, 64)
	result += " %"

	println(idleTime, totalTime)
	cpuInfos.PreviousIdleTime = idleTime
	cpuInfos.PreviousTotalTime = totalTime

	return &Metric{Name: "cpu user", Metric: result}, nil
}

func getNetStat() (*Metric, error) {
	return &Metric{}, nil

}

func getDiskUsage() (*Metric, error) {
	return &Metric{}, nil

}

func getMemUsage() (*Metric, error) {
	return &Metric{}, nil

}
