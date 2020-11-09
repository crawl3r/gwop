package main

import (
	"fmt"
	"gwop/pkg/cli"
	"gwop/pkg/clitool"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	quitEarly := false

	if quitEarly {
		os.Exit(0)
	}
}

func main() {
	clitool.LoadJSONData()
	go cli.Shell() //  we pretty much get "stuck" in here, until a quit/exit/signal is issued

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Press CTRL-C to exit or type 'exit' || 'quit'.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
