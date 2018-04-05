package polls

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/pierre-emmanuelJ/go-exercises/cpu"
	"github.com/pierre-emmanuelJ/go-exercises/memory"
	"github.com/pierre-emmanuelJ/go-exercises/net"
	part "github.com/pierre-emmanuelJ/go-exercises/partition"
)

const NetInterfaceDown = 1
const NetInterfaceUp = 9
const ValidNetInterface = 10

type Metric struct {
	Name   string `json:"name"`
	Metric string `json:"metric"`
}

type Metrics struct {
	Timestamp string    `json:"timestamp"`
	Metrics   []*Metric `json:"metrics"`
}

func Pollsinfos(partition, iNetwork string, cpuInfos *cpu.CPUInfos, netInfos *net.NetInfos) error {

	metrics := []*Metric{}

	av, err := getLoadAv()
	if err != nil {
		return err
	}

	metrics = append(metrics, av)

	cpu, err := getCPUPercentage(cpuInfos)
	if err != nil {
		return err
	}

	metrics = append(metrics, cpu)

	if iNetwork != "" {
		netStat, err := getNetStat(iNetwork, netInfos)
		if err != nil {
			return err
		}
		metrics = append(metrics, netStat)
	}

	if partition != "" {
		disk, err := getDiskUsage(partition)
		if err != nil {
			return err
		}
		metrics = append(metrics, disk)
	}

	memUse, err := getMemUsage()
	if err != nil {
		return err
	}

	metrics = append(metrics, memUse)

	metricsList := &Metrics{Timestamp: time.Now().String(), Metrics: metrics}

	b, err := json.Marshal(metricsList)
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

	cpuLine, err := cpu.GetCPULine()
	if err != nil {
		return nil, err
	}

	cpuTimes, err := cpu.GetCPUTimes(cpuLine)
	if err != nil {
		return nil, err
	}

	idleTime, totalTime := 0, 0

	if err := cpu.GetCPUIdleTimes(&idleTime, &totalTime, cpuTimes); err != nil {
		return nil, err
	}

	var utilization float64

	idleTimeDelta := float64(idleTime - cpuInfos.PreviousIdleTime)
	totalTimeDelta := float64(totalTime - cpuInfos.PreviousTotalTime)
	utilization = (1000*(totalTimeDelta-idleTimeDelta)/totalTimeDelta + 5) / 10

	result := fmt.Sprintf("%.2f", utilization)

	cpuInfos.PreviousIdleTime = idleTime
	cpuInfos.PreviousTotalTime = totalTime

	return &Metric{Name: "Cpu user", Metric: result}, nil
}

func getNetStat(iNet string, netInfos *net.NetInfos) (*Metric, error) {

	lineInfos, err := net.GetInterfaceLine(iNet)
	if err != nil {
		return nil, err
	}

	if len(lineInfos) < ValidNetInterface {
		//TODO implem Error
		return nil, fmt.Errorf("Invalide interface line infos")
	}

	down, err := strconv.ParseInt(lineInfos[NetInterfaceDown], 10, 64)
	if err != nil {
		return nil, err
	}

	up, err := strconv.ParseInt(lineInfos[NetInterfaceUp], 10, 64)
	if err != nil {
		return nil, err
	}

	resDown := down - netInfos.PrevDown
	resUp := up - netInfos.PrevUp

	netInfos.PrevDown = down
	netInfos.PrevUp = up

	return &Metric{Name: fmt.Sprintf("Interface: %s up/down bytes", iNet), Metric: fmt.Sprintf("%v/%v", resUp, resDown)}, nil
}

func getDiskUsage(partition string) (*Metric, error) {

	partitionPath, err := part.GetPartitionMountPointPath(partition)
	if err != nil {
		return nil, err
	}

	statfs := &syscall.Statfs_t{}

	if err := syscall.Statfs(partitionPath, statfs); err != nil {
		return nil, err
	}

	ret := float64(statfs.Bavail) / float64(statfs.Blocks) * 100

	ret = 100 - ret

	return &Metric{Name: fmt.Sprintf("partition %s", partition), Metric: fmt.Sprintf("%.2f", ret)}, nil

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

	memUtilization := memAvailable / memTotal * 100

	memUtilization = 100 - memUtilization

	return &Metric{Name: "Memory utilization", Metric: fmt.Sprintf("%.2f", memUtilization)}, nil
}
