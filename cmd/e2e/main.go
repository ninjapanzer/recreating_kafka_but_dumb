package e2e

import (
	"flag"
	"fmt"
)

func main() {
	multiplexProducer := flag.Bool("start-multiplex-producer", false, "Starts a multiplexed producer")
	multiplexConsumer := flag.Bool("start-multiplex-consumer", false, "Starts a multiplexed consumer")

	// Parse the flags from the command line
	flag.Parse()

	if true == *multiplexProducer {
		StartMultiplexProducer()
	} else if true == *multiplexConsumer {
		StartMultiplexConsumer()
	} else {
		fmt.Print(flag.ErrHelp)
	}
}
