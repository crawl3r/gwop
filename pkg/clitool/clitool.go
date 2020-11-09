package clitool

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
	Name string `json:"name"`
}

// Framework ...
type Framework struct {
	Name         string   `json:"name"`
	Generator    string   `json:"generator"`
	GeneratorCmd string   `json:"generatorcommand"`
	Payloads     []string `json:"payloads"`
}

var loadedData AllLoadedData

// StartImplantCreationProcess is called by the cli upon completion of the 'create' state
func StartImplantCreationProcess(opts *PayloadOptions) {
	fmt.Printf("Generating payload with defined args:\nOS: %d\nFramework: %d\nPayload: %d\n", opts.TargetOS, opts.TargetFramework, opts.Payload)

	generatePayload()
	generateImplantScript()
	compileAndStoreImplant()
}

// TODO: this will take the options and utilise the target tool to create the payload (msfvenom for example)
func generatePayload() {
	fmt.Println("Payload generated")
}

// This will take a script template for an implant and inject the data required (shellcode, key etc)
func generateImplantScript() {
	fmt.Println("Implant script created and ready for compilation")
}

func compileAndStoreImplant( /* This will need to know the target platform for compilation reasons */ ) {
	fmt.Println("Implant compiled and ready")
}

// StartListenerProcess is called by the cli on user demand and will start a listener related to the payload
func StartListenerProcess(opts *PayloadOptions) {
	fmt.Println("Starting listener for target payload")
}

// GetOperatingSystems is a getter for the slice of OS data loaded from JSON
func GetOperatingSystems() []OperatingSystem {
	return loadedData.OperatingSystems
}

// GetFrameworks is a getter for the slice of Framework data loaded from JSON
func GetFrameworks() []Framework {
	return loadedData.Frameworks
}

// GetPayloads is a getter for the specific payloads from the already selected Framework, loaded from JSON
func GetPayloads(framework int) []string {
	return loadedData.Frameworks[framework].Payloads
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
func ConvertUserInputToPayload(frameworkVal int, val int) (string, error) {
	if frameworkVal > len(loadedData.Frameworks) {
		return "", fmt.Errorf("No payloads found for framework choice")
	}

	if val > len(loadedData.Frameworks[frameworkVal].Payloads) {
		return "", fmt.Errorf("No payload choice found with input")
	}

	return loadedData.Frameworks[frameworkVal].Payloads[val], nil
}

// LoadJSONData is called at the start of the program (main.go) to populate our data here
func LoadJSONData() {
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Successfully opened %s\n", jsonFilePath)

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &loadedData)
}
