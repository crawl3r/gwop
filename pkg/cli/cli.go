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

	// make sure the Cli state is set up AFTER the prompt is created otherwise we seg. Coding is hard.
	setStateMainMenu()

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
			if cmd[0] == "exit" || cmd[0] == "quit" || cmd[0] == "q" {
				fmt.Println("[*] Cleaning up and shutting down")
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
					printListTargetOperatingSystems()
				case "frameworks":
					printListFrameworks()
				case "payloads":
					if len(cmd) < 2 {
						fmt.Println("[!] Please supply a target Framework ID (i.e 'payloads 1')")
						printListFrameworks()
					} else {
						if !util.IsAnInteger(cmd[1]) {
							fmt.Println("[!] Target framework ID was not recognised (i.e 'payloads 1')")
						}
						frameworkID, _ := strconv.Atoi(cmd[1])
						if frameworkID-1 >= 0 && frameworkID-1 < len(clitool.GetFrameworks(-1)) {
							printListPayloads(frameworkID - 1)
						} else {
							fmt.Println("[!] Target framework ID was not recognised (i.e 'payloads 1')")
						}
					}

				default:
					fmt.Println("[!] Sorry, input was not recognised")
				}
			case "create":
				switch cliCreateState {
				case 0:
					if !util.IsAnInteger(cmd[0]) {
						fmt.Println("[!] Sorry, input was not an integer value")
						setStateCreate(cliCreateState)
					} else {
						val, _ := strconv.Atoi(cmd[0])
						_, err := clitool.ConvertUserInputToOperatingSystem(val - 1)
						if err != nil {
							fmt.Println("[!] Input was not recognised, please try again")
							setStateCreate(cliCreateState)
						} else {
							payloadOptions.TargetOS = val - 1
							setStateCreate(1)
						}
					}
				case 1:
					if !util.IsAnInteger(cmd[0]) {
						fmt.Println("[!] Sorry, input was not an integer value")
						setStateCreate(cliCreateState)
					} else {
						val, _ := strconv.Atoi(cmd[0])
						_, err := clitool.ConvertUserInputToFramework(val - 1)
						if err != nil {
							fmt.Println("[!] Input was not recognised, please try again")
							setStateCreate(cliCreateState)
						} else {
							payloadOptions.TargetFramework = val - 1
							setStateCreate(2)
						}
					}
				case 2:
					if !util.IsAnInteger(cmd[0]) {
						fmt.Println("[!] Sorry, input was not an integer value")
						setStateCreate(cliCreateState)
					} else {
						val, _ := strconv.Atoi(cmd[0])
						_, err := clitool.ConvertUserInputToPayload(payloadOptions.TargetFramework, payloadOptions.TargetOS, val-1)
						if err != nil {
							fmt.Println("[!] Input was not recognised, please try again")
							setStateCreate(cliCreateState)
						} else {
							payloadOptions.Payload = val - 1
							setStateCreate(3)
						}
					}
				case 3:
					if !util.IsLegalIPAddress(cmd[0]) {
						fmt.Println("[!] Sorry, input was not a legal IPV4 address")
						setStateCreate(cliCreateState)
					} else {
						payloadOptions.Lhost = cmd[0]
						setStateCreate(4)
					}
				case 4:
					if !util.IsLegalPortNumber(cmd[0]) {
						fmt.Println("[!] Sorry, input was not a legal port value")
						setStateCreate(cliCreateState)
					} else {
						payloadOptions.Lport = cmd[0]
						setStateGeneratePayload()
					}
				}
			case "generate":
				switch cmd[0] {
				case "y":
					fmt.Println("[*] Generating implant... please wait")
					clitool.StartImplantCreationProcess(payloadOptions)
					setStateStartListener()
				case "n":
					fmt.Println("[*] Quitting back to main menu")
					setStateMainMenu()
				default:
					fmt.Println("[!] Sorry, input was unknown")
				}
			case "listener":
				switch cmd[0] {
				case "y":
					fmt.Println("[*] Starting listener... please wait")
					clitool.StartListenerProcess(payloadOptions)

					// TODO: do we quit out and just focus another terminal window that is now our listener? Or overwrite the current output with listener?
				case "n":
					fmt.Println("[*] Quitting back to main menu")
					setStateMainMenu()
				default:
					fmt.Println("[!] Sorry, input was unknown")
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
		{"payloads n", "Lists all possible payloads (n = target OS)"},
		{"help", "Prints ths menu"},
		{"exit", "Exit and close the GWOP cli tool"},
		{"quit", "Exit and close the GWOP cli tool"},
	}

	table.AppendBulk(data)
	fmt.Println()
	table.Render()
	fmt.Println()
}

func printListTargetOperatingSystems() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetHeader([]string{"#", "Operating System"})

	data := [][]string{}

	for i, opSys := range clitool.GetOperatingSystems() {
		lineData := []string{}
		lineData = append(lineData, strconv.Itoa(i))
		lineData = append(lineData, opSys.Name)
		data = append(data, lineData)
	}

	table.AppendBulk(data)
	fmt.Println()
	table.Render()
	fmt.Println()
}

func printListFrameworks() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetHeader([]string{"#", "Frameworks", "Operating Systems"})

	data := [][]string{}
	for i, f := range clitool.GetFrameworks(-1) {
		lineData := []string{}
		lineData = append(lineData, strconv.Itoa(i+1))
		lineData = append(lineData, f.Name)

		// OSes available
		oses := clitool.GetFrameworkOperatingSystemOptions(i)
		osToPrint := ""
		for i, s := range oses {
			osToPrint += s
			if i < len(oses)-1 {
				osToPrint += ", "
			}
		}
		lineData = append(lineData, osToPrint)
		data = append(data, lineData)
	}

	table.AppendBulk(data)
	fmt.Println()
	table.Render()
	fmt.Println()
}

// Unlike the other prints, this one could be pretty hectic output due to the different frameworks
// It is worth having the user do 'payloads 1' for listing Windows
// 'payloads 2' for listing Linux etc
func printListPayloads(frameworkID int) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetHeader([]string{"#", "Payloads"})

	data := [][]string{}
	for i, p := range clitool.GetPayloads(frameworkID, -1) {
		lineData := []string{}
		lineData = append(lineData, strconv.Itoa(i+1))
		lineData = append(lineData, p)

		data = append(data, lineData)
	}

	table.AppendBulk(data)
	fmt.Println()
	table.Render()
	fmt.Println()
}

// SetStateMainMenu sets the Cli state back to "main"
func setStateMainMenu() {
	cliMenuState = "main"
	cliCreateState = 0
	payloadOptions = &clitool.PayloadOptions{
		TargetOS:        -1,
		TargetFramework: -1,
		Payload:         -1,
	}
	fmt.Println("Type 'create' to get started, or 'help' for other options")
	prompt.SetPrompt("\033[31mGWOP»\033[0m ")
}

func setStateCreate(stage int) {
	cliMenuState = "create"
	cliCreateState = stage

	switch cliCreateState {
	case 0:
		fmt.Println("Please choose your target operating system:")
		prompt.SetPrompt("\033[31mGWOP|OS»\033[0m ")

		for i, v := range clitool.GetOperatingSystems() {
			fmt.Printf("%d - %s\n", i+1, v.Name)
		}
	case 1:
		fmt.Println("Please choose your target framework:")
		prompt.SetPrompt("\033[31mGWOP|Framework»\033[0m ")

		for i, v := range clitool.GetFrameworks(payloadOptions.TargetOS) {
			fmt.Printf("%d - %s\n", i+1, v.Name)
		}
	case 2:
		fmt.Println("Please choose your payload:")
		prompt.SetPrompt("\033[31mGWOP|Payload»\033[0m ")

		for i, v := range clitool.GetPayloads(payloadOptions.TargetFramework, payloadOptions.TargetOS) {
			fmt.Printf("%d - %s\n", i+1, v)
		}
	case 3:
		fmt.Println("Please specify the listener host IP:")
		prompt.SetPrompt("\033[31mGWOP|lhost\033[0m ")
	case 4:
		fmt.Println("Please specify the listener host port:")
		prompt.SetPrompt("\033[31mGWOP|lport\033[0m ")
	}
}

func setStateGeneratePayload() {
	cliMenuState = "generate"

	targetOS, _ := clitool.ConvertUserInputToOperatingSystem(payloadOptions.TargetOS)
	targetFramework, _ := clitool.ConvertUserInputToFramework(payloadOptions.TargetFramework)
	targetPayload, _ := clitool.ConvertUserInputToPayload(payloadOptions.TargetFramework, payloadOptions.TargetOS, payloadOptions.Payload)

	fmt.Println("Payload ready to generate with following args:")
	fmt.Printf("\tTarget OS: %s\n", targetOS)
	fmt.Printf("\tTarget Framework: %s\n", targetFramework)
	fmt.Printf("\tPayload: %s\n", targetPayload)
	fmt.Printf("\tListener IP: %s\n", payloadOptions.Lhost)
	fmt.Printf("\tListener Port: %s\n", payloadOptions.Lport)
	fmt.Printf("\nShall we generate the implant with these options?\n")

	fmt.Println("y - yes")
	fmt.Println("n - no (quits back to main menu)")

	prompt.SetPrompt("\033[31mGWOP|generate\033[0m ")
}

func setStateStartListener() {
	cliMenuState = "listener"

	fmt.Printf("Payload generated, start listener?\n")
	fmt.Println("y - yes")
	fmt.Println("n - no (quits back to main menu)")

	prompt.SetPrompt("\033[31mGWOP|listener\033[0m ")
}
