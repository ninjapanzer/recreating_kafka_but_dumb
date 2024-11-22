package main

import (
	"flag"
	"fmt"
	"go_stream_events/cmd/e2e/tools"
)

func main() {
	multiplexProducer := flag.Bool("start-multiplex-producer", false, "Starts a multiplexed producer")
	multiplexConsumer := flag.Bool("start-multiplex-consumer", false, "Starts a multiplexed consumer")

	// Parse the flags from the command line
	flag.Parse()

	if true == *multiplexProducer {
		tools.StartMultiplexProducer()
	} else if true == *multiplexConsumer {
		tools.StartMultiplexConsumer()
	} else {
		fmt.Print(flag.ErrHelp)
	}
}
