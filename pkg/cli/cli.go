package cli

import (
	"fmt"
	"gwop/pkg/clitool"
	"gwop/pkg/util"
	"io"
	"log"
	"os"
	"strconv"
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
		Payload:         0,
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
						val, _ := strconv.Atoi(cmd[0])
						_, err := clitool.ConvertUserInputToOperatingSystem(val - 1)
						if err != nil {
							fmt.Println("Input was not recognised, please try again")
							setStateCreate(cliCreateState)
						} else {
							payloadOptions.TargetOS = val - 1
							setStateCreate(1)
						}
					}
				case 1:
					if !util.IsAnInteger(cmd[0]) {
						fmt.Println("Sorry, input was not an integer value")
						setStateCreate(cliCreateState)
					} else {
						val, _ := strconv.Atoi(cmd[0])
						_, err := clitool.ConvertUserInputToFramework(val - 1)
						if err != nil {
							fmt.Println("Input was not recognised, please try again")
							setStateCreate(cliCreateState)
						} else {
							payloadOptions.TargetFramework = val - 1
							setStateCreate(2)
						}
					}
				case 2:
					if !util.IsAnInteger(cmd[0]) {
						fmt.Println("Sorry, input was not an integer value")
						setStateCreate(cliCreateState)
					} else {
						val, _ := strconv.Atoi(cmd[0])
						_, err := clitool.ConvertUserInputToPayload(payloadOptions.TargetFramework, val-1)
						if err != nil {
							fmt.Println("Input was not recognised, please try again")
							setStateCreate(cliCreateState)
						} else {
							payloadOptions.Payload = val - 1
							setStateGeneratePayload()
						}
					}
				case 3:
				case 4:
				}
			}
		}
	}
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// block Ctrl + Z feature please. Ctrl+c is used to back out (as specified by Cli itself)
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
		prompt.SetPrompt("\033[31mGWOP|OS»\033[0m")

		for i, v := range clitool.OperatingSystemChoices {
			fmt.Printf("%d - %s\n", i+1, v)
		}
	case 1:
		fmt.Println("Please choose your target framework:")
		prompt.SetPrompt("\033[31mGWOP|Framework»\033[0m")

		for i, v := range clitool.FrameworkChoices {
			fmt.Printf("%d - %s\n", i+1, v)
		}
	case 2:
		fmt.Println("Please choose your payload:")
		prompt.SetPrompt("\033[31mGWOP|Payload»\033[0m")

		for i, v := range clitool.PayloadChoices {
			fmt.Printf("%d - %s\n", i+1, v)
		}
	case 3:
		// TODO: payload options LHOST
		fmt.Println("Please specify the listener host IP:")
		prompt.SetPrompt("\033[31mGWOP|lhost\033[0m")
	case 4:
		// TODO: payload options LPORT
		fmt.Println("Please specify the listener host port:")
		prompt.SetPrompt("\033[31mGWOP|lport\033[0m")
	}
}

func setStateGeneratePayload() {
	cliMenuState = "generate"

	targetOS, _ := clitool.ConvertUserInputToOperatingSystem(payloadOptions.TargetOS)
	targetFramework, _ := clitool.ConvertUserInputToFramework(payloadOptions.TargetFramework)
	targetPayload, _ := clitool.ConvertUserInputToPayload(payloadOptions.TargetFramework, payloadOptions.Payload)

	fmt.Println("Payload ready to generate with following args:")
	fmt.Printf("\tTarget OS: %s\n", targetOS)
	fmt.Printf("\tTarget Framework: %s\n", targetFramework)
	fmt.Printf("\tPayload: %s\n", targetPayload)
	fmt.Println("\nShall we generate the implant with these options?")
}
