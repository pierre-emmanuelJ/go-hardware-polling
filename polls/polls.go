package polls

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Metric struct {
	Name   string `json:"name"`
	Metric string `json:"metric"`
}

type Metrics struct {
	Timestamp string    `json:"timestamp"`
	Metrics   []*Metric `json:"metrics"`
}

type CpuInfos struct {
	previousIdleTime  int
	previousTotalTime int
}

func Pollsinfos(partition, iNetwork string, cpuInfos *CpuInfos) error {

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

func getCPUPercentage(cpuInfos *CpuInfos) (*Metric, error) {

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

	cpuTimes, err := getCPUTimes(scanner.Text())
	if err != nil {
		return nil, err
	}

	idleTime, totalTime := 0, 0

	if err := getCPUIdleTimes(&idleTime, &totalTime, cpuTimes); err != nil {
		return nil, err
	}

	println(idleTime, totalTime)
	result := ""
	idleTimeDelta := idleTime - cpuInfos.previousIdleTime
	totalTimeDelta := totalTime - cpuInfos.previousTotalTime
	utilization := 100.0 * (1.0 - idleTimeDelta/totalTimeDelta)

	result += strconv.FormatFloat(float64(utilization), 'E', -1, 64)
	result += " %"

	println(idleTime, totalTime)
	cpuInfos.previousIdleTime = idleTime
	cpuInfos.previousTotalTime = totalTime

	return &Metric{Name: "cpu user", Metric: result}, nil
}

func getCPUIdleTimes(idleTime, totalTime *int, cpuTimes []int) error {
	if len(cpuTimes) < 4 {
		return nil
	}
	*idleTime = cpuTimes[3]

	for _, i := range cpuTimes {
		*totalTime += i
	}
	return nil
}

func getCPUTimes(s string) ([]int, error) {

	times := strings.Split(s, " ")
	var res []int
	for _, time := range times {

		if time == "" || time == "cpu" {
			continue
		}
		timeNum, err := strconv.Atoi(time)
		if err != nil {
			return nil, err
		}
		res = append(res, timeNum)
	}
	return res, nil
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
