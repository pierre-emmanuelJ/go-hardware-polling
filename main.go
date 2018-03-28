package main

import (
	"flag"
	"time"

	"github.com/pierre-emmanuelJ/go-exercises/polls"
)

func main() {

	var interval = flag.Int("i", 1, "interval in seconds at which to poll")
	var partition = flag.String("p", "", "partition to poll")
	var iNetwork = flag.String("n", "", "network interface to poll")

	flag.Parse()

	for true {
		time.Sleep(time.Duration(*interval) * time.Second)
		polls.Pollsinfos(*partition, *iNetwork)
	}

}
