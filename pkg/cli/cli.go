package cli

import (
	"fmt"
	"gwop/pkg/clitool"
	"gwop/pkg/util"
	"io"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

var (
	cliMenuState   string // main, agent
	cliCreateState int
	prompt         *readline.Instance
	payloadOptions *clitool.PayloadOptions
)

// Shell is the exported function to start the command line interface
func Shell() {
	cliMenuState = "main"
	cliCreateState = 0
	payloadOptions = &clitool.PayloadOptions{
		TargetOS:        0,
		TargetFramework: 0,
		Payload:         "",
	}

	p, err := readline.NewEx(&readline.Config{
		Prompt:              "\033[31mGWOP»\033[0m ",
		InterruptPrompt:     "^C",
		EOFPrompt:           "exit",
		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})

	if err != nil {
		color.Red("[!]There was an error with the provided input")
		color.Red(err.Error())
	}
	prompt = p

	defer func() {
		err := prompt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.SetOutput(prompt.Stderr())

	for {
		line, err := prompt.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line) // get the prompt line
		cmd := strings.Fields(line)    // get the command from the line

		// let's figure out the requested command
		if len(cmd) > 0 {
			// First things first, let's break out if any of the states has a quit/exit request in it
			if cmd[0] == "exit" || cmd[0] == "quit" {
				fmt.Println("Cleaning up and shutting down")
				syscall.Kill(syscall.Getpid(), syscall.SIGINT) // possibly derp but might work
			}

			switch cliMenuState {
			case "main":
				switch cmd[0] {
				case "help":
					printHelpMainMenu()
				case "create":
					setStateCreate(0)
				case "targets":
					printListTargets()
				case "frameworks":
					printListFrameworks()
				case "payloads":
					printListPayloads()
				}
			case "create":
				switch cliCreateState {
				case 0:
					if !util.IsAnInteger(cmd[0]) {
						fmt.Println("Sorry, input was not an integer value")
						setStateCreate(cliCreateState)
					} else {
						switch cmd[0] {
						case "1":

						}
						setStateCreate(1)
					}
				case 1:
					if !util.IsAnInteger(cmd[0]) {
						fmt.Println("Sorry, input was not an integer value")
						setStateCreate(cliCreateState)
					} else {
						switch cmd[0] {
						case "1":
							payloadOptions.TargetFramework = 0
						case "2":
							payloadOptions.TargetFramework = 1
						}
						setStateCreate(2)
					}
				case 2:
					if !util.IsAnInteger(cmd[0]) {
						fmt.Println("Sorry, input was not an integer value")
						setStateCreate(cliCreateState)
					} else {
						payloadArg, err := clitool.ConvertUserInputToPayload(0)
						if err != nil {
							fmt.Println("Argument not recognised, please try again")
							setStateCreate(2)
						}
						payloadOptions.Payload = payloadArg
						setStateCreate(3)
					}
				}
			}
		}
	}
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func printHelpMainMenu() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetHeader([]string{"Command", "Description"})

	data := [][]string{
		{"create", "Begins the interactive dropper creation process"},
		{"targets", "Lists the current target operating systems"},
		{"frameworks", "Lists the current frameworks/C2 platforms available"},
		{"payloads", "Lists all possible payloads"},
		{"help", "Prints ths menu"},
		{"exit", "Exit and close the GWOP cli tool"},
		{"quit", "Exit and close the GWOP cli tool"},
	}

	table.AppendBulk(data)
	fmt.Println()
	table.Render()
	fmt.Println()
}

func printListTargets() {
	// TODO
	fmt.Println("TODO")
}

func printListFrameworks() {
	// TODO
	fmt.Println("TODO")
}

func printListPayloads() {
	// TODO
	fmt.Println("TODO")
}

// SetStateMainMenu sets the Cli state back to "main"
func setStateMainMenu() {
	cliMenuState = "main"
	prompt.SetPrompt("\033[31mGWOP»\033[0m ")
}

func setStateCreate(stage int) {
	cliMenuState = "create"
	cliCreateState = stage

	switch cliCreateState {
	case 0:
		fmt.Println("Please choose your target operating system:")

		// TODO: list, loopable so the logic switch above can check the answer easily

	case 1:
		fmt.Println("Please choose your target framework:")

		// TODO: list, loopable so the logic switch above can check the answer easily
	case 2:
		fmt.Println("Please choose your payload:")

		// TODO: list, loopable so the logic switch above can check the answer easily
	case 3:
		fmt.Println("Would you like to start the listener?")
	}
}

func setStateGeneratePayload() {
	cliMenuState = "generate"

	targetOS, _ := clitool.ConvertUserInputToOperatingSystem(payloadOptions.TargetOS)
	targetFramework, _ := clitool.ConvertUserInputToFramework(payloadOptions.TargetFramework)
	targetPayload, _ := clitool.ConvertUserInputToPayload(payloadOptions.Payload)

	fmt.Println("Payload ready to generate with following args:")
	fmt.Println("\tTarget OS: %s", targetOS)
	fmt.Println("\tTarget Framework: %s", targetFramework)
	fmt.Println("\tPayload: %s", targetPayload)
	fmt.Println("\nAre these correct?")
}
