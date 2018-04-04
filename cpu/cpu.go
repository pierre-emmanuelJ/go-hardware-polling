package cpu

import (
	"strconv"
	"strings"
)

type CPUInfos struct {
	PreviousIdleTime  int
	PreviousTotalTime int
}

func GetCPUIdleTimes(idleTime, totalTime *int, cpuTimes []int) error {
	if len(cpuTimes) < 4 {
		return nil
	}
	*idleTime = cpuTimes[3]

	for _, i := range cpuTimes {
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
