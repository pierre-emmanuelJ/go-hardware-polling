package polls

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
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

func Pollsinfos(partition, iNetwork string) error {

	av, err := getLoadAv()
	if err != nil {
		return err
	}

	cpu, err := getCpuPercentage()
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

func getCpuPercentage() (*Metric, error) {

	return &Metric{}, nil

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
