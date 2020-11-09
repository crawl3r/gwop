package clitool

import "fmt"

// PayloadOptions is populated and passed into the clitool startImplantCreationProcess to help guide the flow
type PayloadOptions struct {
	TargetOS        int // 0 == windows, 1 == linux, 2 == macosx
	TargetFramework int // 0 == msf
	Payload         int // will map to the list presented in the CLI for easy scaling
}

func startImplantCreationProcess(opts *PayloadOptions) {
	fmt.Printf("Generating payload with defined args:\nOS: %d\nFramework: %d\nPayload: %d\n", opts.TargetOS, opts.TargetFramework, opts.Payload)
}

// TODO: this will take the options and utilise the target tool to create the payload (msfvenom for example)
func generatePayload() {
	fmt.Println("Payload generated")
}

// This will take a script template for an implant and inject the data required (shellcode, key etc)
func generateImplantScript() {
	fmt.Println("Implant script created and ready for compilation")
}

func buildAndStoreImplant( /* This will need to know the target platform for compilation reasons */ ) {
	fmt.Println("Implant compiled and ready")
}

// ConvertUserInputToOperatingSystem is used to convert the user input (int) to the string value for the gen tools
func ConvertUserInputToOperatingSystem(val int) (string, error) {
	return "", nil
}

// ConvertUserInputToFramework is used to convert the user input (int) to the string value for the gen tools
func ConvertUserInputToFramework(val int) (string, error) {
	return "", nil
}

// ConvertUserInputToPayload is used to convert the user input (int) to the string value for the gen tools
func ConvertUserInputToPayload(val int) (string, error) {
	return "", nil
}
