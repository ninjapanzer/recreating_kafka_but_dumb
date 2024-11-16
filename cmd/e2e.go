package main

import (
	"flag"
	"fmt"
	"go_stream_events/cmd/e2e"
)

func main() {
	multiplexProducer := flag.Bool("start-multiplex-producer", false, "Starts a multiplexed producer")
	multiplexConsumer := flag.Bool("start-multiplex-consumer", false, "Starts a multiplexed consumer")

	// Parse the flags from the command line
	flag.Parse()

	if true == *multiplexProducer {
		e2e.StartMultiplexProducer()
	} else if true == *multiplexConsumer {
		e2e.StartMultiplexConsumer()
	} else {
		fmt.Print(flag.ErrHelp)
	}
}
