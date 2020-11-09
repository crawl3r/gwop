package clitool

import "fmt"

// PayloadOptions is populated and passed into the clitool startImplantCreationProcess to help guide the flow
type PayloadOptions struct {
	TargetOS        int // 0 == windows, 1 == linux, 2 == macosx
	TargetFramework int // 0 == msf
	Payload         int // will map to the list presented in the CLI for easy scaling
	Lhost           string
	Lport           string
}

// TODO: replace these with loaded JSON on start - easy scalable etc

// OperatingSystemChoices ...
var OperatingSystemChoices = []string{
	"Windows",
	"Linux",
	"MacOSX",
}

// FrameworkChoices ...
var FrameworkChoices = []string{
	"Metasploit/MSFvenom",
}

// PayloadChoices ...
var PayloadChoices = [][]string{
	{
		"linux/x86/meterpreter/reverse_tcp",
		"windows/meterpreter/reverse_tcp",
	},
}

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

// ConvertUserInputToOperatingSystem is used to convert the user input (int) to the string value for the gen tools
func ConvertUserInputToOperatingSystem(val int) (string, error) {
	if val > len(OperatingSystemChoices) {
		return "", fmt.Errorf("No operating system choice found with input")
	}

	return OperatingSystemChoices[val], nil
}

// ConvertUserInputToFramework is used to convert the user input (int) to the string value for the gen tools
func ConvertUserInputToFramework(val int) (string, error) {
	if val > len(FrameworkChoices) {
		return "", fmt.Errorf("No framework choice found with input")
	}

	return FrameworkChoices[val], nil
}

// ConvertUserInputToPayload is used to convert the user input (int) to the string value for the gen tools
func ConvertUserInputToPayload(frameworkVal int, val int) (string, error) {
	if frameworkVal > len(PayloadChoices) {
		return "", fmt.Errorf("No payloads found for framework choice")
	}

	if val > len(PayloadChoices[frameworkVal]) {
		return "", fmt.Errorf("No payload choice found with input")
	}

	return PayloadChoices[frameworkVal][val], nil
}
