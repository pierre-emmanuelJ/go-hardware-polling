package cpu

import (
	"fmt"
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

func GetCPUIdleTimes(idleTime, totalTime *int, cpuTimes []int) error {
	if len(cpuTimes) < minValidCPUStats {
		return fmt.Errorf("Invalide cpu line missing some stats")
	}
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
