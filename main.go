package main

import (
	"flag"
	"log"
	"time"

	"github.com/pierre-emmanuelJ/go-exercises/polls"
)

func main() {

	var interval = flag.Int("i", 1, "interval in seconds at which to poll")
	var partition = flag.String("p", "", "partition to poll")
	var iNetwork = flag.String("n", "", "network interface to poll")

	flag.Parse()

	// if *partition == "" && *iNetwork == "" {
	// 	println(" -i: interval in seconds at which to poll\n",
	// 		"-p: partition to poll\n",
	// 		"-n: network interface to poll")
	// 	return
	// }

	cpuInfos := &polls.CpuInfos{}
	for true {
		time.Sleep(time.Duration(*interval) * time.Second)
		err := polls.Pollsinfos(*partition, *iNetwork, cpuInfos)
		if err != nil {
			log.Fatal(err)
		}
	}

}
