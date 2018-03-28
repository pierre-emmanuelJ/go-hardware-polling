package polls

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Metric struct {
	Name   string  `json:"name"`
	Metric float64 `json:"metric"`
}

type Metrics struct {
	Timestamp string   `json:"timestamp"`
	Metrics   []Metric `json:"metrics"`
}

func Pollsinfos(partition, iNetwork string) {

	metrics := &Metrics{Timestamp: time.Now().String()}

	// getLoadAv()
	// getCpuPercentage()
	// getNetStat()
	// getDiskUsage()
	// getMemUsage()

	b, err := json.Marshal(metrics)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
}
