package cpu

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const idle = 3
const iowait = 4
const minValidCPUStats = 5
const endCPUList = 8

type CPUInfos struct {
	PreviousIdleTime  int
	PreviousTotalTime int
}

func GetCPULine() (string, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if line, isOk := getCPU(scanner.Text(), "cpu"); isOk {
			return line, nil
		}
	}
	return "", fmt.Errorf("No cpu infos in this file or invalide file")
}

func getCPU(s, key string) (string, bool) {

	cpuInfos := strings.Split(s, " ")

	if len(cpuInfos) < minValidCPUStats {
		return "", false
	}

	if cpuInfos[0] == key {
		return s, true
	}
	return "", false
}

func GetCPUIdleTimes(idleTime, totalTime *int, cpuTimes []int) error {
	*idleTime = cpuTimes[idle]

	for index, i := range cpuTimes {

		if index >= endCPUList {
			break
		}
		*totalTime += i
	}
	return nil
}

func GetCPUTimes(s string) ([]int, error) {

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
