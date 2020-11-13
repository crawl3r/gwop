package clitool

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

// PayloadOptions is populated and passed into the clitool startImplantCreationProcess to help guide the flow
type PayloadOptions struct {
	TargetOS        int // 0 == windows, 1 == linux, 2 == macosx
	TargetFramework int // 0 == msf
	Payload         int // will map to the list presented in the CLI for easy scaling
	Lhost           string
	Lport           string
}

var jsonFilePath = "data/data.json"

// AllLoadedData ...
type AllLoadedData struct {
	OperatingSystems []OperatingSystem `json:"operatingsystems"`
	Frameworks       []Framework       `json:"frameworks"`
}

// OperatingSystem ...
type OperatingSystem struct {
	Name           string `json:"name"`
	GoArchitecture string `json:"goarchi"`
}

// Framework ...
type Framework struct {
	Name         string    `json:"name"`
	Generator    string    `json:"generator"`
	GeneratorCmd string    `json:"generatorcommand"`
	Payloads     []Payload `json:"payloads"`
}

// Payload ...
type Payload struct {
	OperatingSystem string   `json:"operatingsystem"`
	Options         []string `json:"options"`
}

var loadedData AllLoadedData

// StartImplantCreationProcess is called by the cli upon completion of the 'create' state
func StartImplantCreationProcess(opts *PayloadOptions) {
	shellcode := generatePayload(opts)
	generateImplantScript(shellcode)
	if !compileAndStoreImplant(opts) {
		fmt.Println("Implant failed to compile")
	}
}

// TODO: this will take the options and utilise the target tool to create the payload (msfvenom for example)
func generatePayload(opts *PayloadOptions) string {
	framework, _ := ConvertUserInputToFramework(opts.TargetFramework)
	fmt.Printf("[*] Generating payload with target framework: %s\n", framework)

	fmt.Println("[*] Payload generated")
	return "deadbeef"
}

// This will take a script template for an implant and inject the data required (shellcode, key etc)
func generateImplantScript(shellcode string) {
	// logic
	// load the implant.template text file into memory
	// replace the key values with real values
	// blit the text file to a Go script that is ready to be compiled (./cmd/implant_gen/main.go)

	// 1) Load the file into memory line by line
	templateFile, err := os.Open("data/implant.template")
	lines := []string{}

	if err != nil {
		log.Fatal(err)
	}

	defer templateFile.Close()

	scanner := bufio.NewScanner(templateFile)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("---- loaded script ----")
	for _, l := range lines {
		fmt.Println(l)
	}
	fmt.Println("-----------------------")

	// 2) edit the target script variables here, might change to a switch case based on the OS
	fmt.Println("[*] Editing implant script data")
	for i := 0; i < len(lines); i++ {
		l := lines[i]
		if strings.Contains(l, "<--HEXSC-->") {
			fmt.Println("[*] Found hex shellcode key. Adding shellcode to template")
			lines[i] = strings.Replace(l, "<--HEXSC-->", shellcode, -1)
		}
	}

	// 3) blit text to a go script ready for compilation
	fmt.Println("[*] Blitting implant file to Go script")
	implantFile, err := os.Create("cmd/implant_gen/main.go")
	if err != nil {
		log.Fatal(err)
	}

	for _, l := range lines {
		_, err = implantFile.WriteString(l + "\n")
		if err != nil {
			implantFile.Close()
			log.Fatal(err)
		}
	}
	err = implantFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("[*] Implant script created and ready for compilation (cmd/implant_gen/main.go)")
}

func compileAndStoreImplant(opts *PayloadOptions) bool {
	// create a system call argument to one liner compile the script depending on the target architecture
	targetOs := getGoArchitectureForOS(opts.TargetOS)

	// Environment variables for current build
	// first we cache the current ones
	fmt.Println("[*] Cacheing current build environment variables")
	currentGoOsEnv := os.Getenv(targetOs)
	currentGoArchEnv := os.Getenv("GOARCH")

	// now we set the environment variables for the new build
	osEnvErr := os.Setenv("GOOS", targetOs)
	if osEnvErr != nil {
		log.Fatal("Error setting GOOS env var:", osEnvErr)
	}
	osEnvErr = os.Setenv("GOARCH", "amd64") // TODO: extend this to allow other architectures
	if osEnvErr != nil {
		log.Fatal("Error setting GOARCH env var:", osEnvErr)
	}

	// set up compile script
	compileScript, err := os.OpenFile("helpers/compile.sh", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer compileScript.Close()

	fileExt := ""
	if opts.TargetOS == 0 {
		fileExt = ".exe"
	}
	line := "go build -ldflags \"-s -w\" -o out/implant-" + targetOs + fileExt + " cmd/implant_dev/main.go"
	_, err = compileScript.WriteString(line)
	if err != nil {
		compileScript.Close()
		log.Fatal(err)
	}

	compileScript.Close()

	// TODO: don't be a moron and use file I/O just to build and execute a one liner...
	_, err = exec.Command("/bin/sh", "helpers/compile.sh").Output()
	if err != nil {
		fmt.Printf("error %s", err)
	}

	fmt.Println("[*] Restoring environment variables")
	// set the environment variables back to their previous values
	osEnvErr = os.Setenv("GOOS", currentGoOsEnv)
	if osEnvErr != nil {
		log.Fatal("Error setting GOOS env var back:", osEnvErr)
	}
	osEnvErr = os.Setenv("GOARCH", currentGoArchEnv)
	if osEnvErr != nil {
		log.Fatal("Error setting GOARCH env var back:", osEnvErr)
	}

	fmt.Printf("[*] Implant compiled and ready. Stored in 'out/implant-%s%s'\n", targetOs, fileExt)
	return true
}

// StartListenerProcess is called by the cli on user demand and will start a listener related to the payload
func StartListenerProcess(opts *PayloadOptions) {
}

// GetOperatingSystems is a getter for the slice of OS data loaded from JSON
func GetOperatingSystems() []OperatingSystem {
	return loadedData.OperatingSystems
}

// GetFrameworks is a getter for the slice of Framework data loaded from JSON
func GetFrameworks(targetOs int) []Framework {
	if targetOs == -1 {
		return loadedData.Frameworks
	}

	possibleFrameworks := []Framework{}
	osName, _ := ConvertUserInputToOperatingSystem(targetOs)

	for i := range loadedData.Frameworks {
		isLegalFramework := false
		currentFramework := loadedData.Frameworks[i]
		for j := range currentFramework.Payloads {
			currentPayloadSelection := currentFramework.Payloads[j]
			if currentPayloadSelection.OperatingSystem == osName {
				isLegalFramework = true
			}
		}

		if isLegalFramework {
			possibleFrameworks = append(possibleFrameworks, currentFramework)
		}
	}

	return possibleFrameworks
}

// GetPayloads is a getter for the specific payloads from the already selected Framework, loaded from JSON
func GetPayloads(framework int, opSys int) []string {
	possibleFrameworks := GetFrameworks(opSys)
	return possibleFrameworks[framework].Payloads[opSys].Options
}

// ConvertUserInputToOperatingSystem is used to convert the user input (int) to the string value for the gen tools
func ConvertUserInputToOperatingSystem(val int) (string, error) {
	if val > len(loadedData.OperatingSystems) {
		return "", fmt.Errorf("No operating system choice found with input")
	}

	return loadedData.OperatingSystems[val].Name, nil
}

// ConvertUserInputToFramework is used to convert the user input (int) to the string value for the gen tools
func ConvertUserInputToFramework(val int) (string, error) {
	if val > len(loadedData.Frameworks) {
		return "", fmt.Errorf("No framework choice found with input")
	}

	return loadedData.Frameworks[val].Name, nil
}

// ConvertUserInputToPayload is used to convert the user input (int) to the string value for the gen tools
func ConvertUserInputToPayload(frameworkVal int, opSysVal int, val int) (string, error) {
	if frameworkVal > len(loadedData.Frameworks) {
		return "", fmt.Errorf("No payloads found for framework choice")
	}

	if opSysVal > len(loadedData.Frameworks[frameworkVal].Payloads) {
		return "", fmt.Errorf("No payloads found for operating system choice")
	}

	if val > len(loadedData.Frameworks[frameworkVal].Payloads[opSysVal].Options) {
		return "", fmt.Errorf("No payload choice found with input")
	}

	return loadedData.Frameworks[frameworkVal].Payloads[opSysVal].Options[val], nil
}

func getGoArchitectureForOS(osChoice int) string {
	return loadedData.OperatingSystems[osChoice].GoArchitecture
}

// LoadJSONData is called at the start of the program (main.go) to populate our data here
func LoadJSONData() {
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("[*] Successfully opened %s\n", jsonFilePath)

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &loadedData)
}
