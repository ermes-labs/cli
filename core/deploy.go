package core

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"yourapp/infrastructure" // Replace with your actual package path
)

// Deploy function in Go
func Deploy(functionName, infrastructureFileName, inEvery string, inAreas, exceptIn []string, faasCliArguments []string) error {
	var infra infrastructure.Infrastructure

	// Read and parse the infrastructure file.
	infraBytes, err := os.ReadFile(infrastructureFileName)
	// Check for errors.
	if err != nil {
		// Compose the error message.
		return fmt.Errorf("Error reading infrastructure file: %w", err)
	}

	// Parse the infrastructure JSON.
	err = json.Unmarshal(infraBytes, &infra)
	// Check for errors.
	if err != nil {
		// Compose the error message.
		return fmt.Errorf("Error parsing infrastructure JSON: %w", err)
	}

	// TODO: Implement autoFill() method in Go.
	// Assuming autoFill() method logic is implemented in Go.
	infra.AutoFill()

	// Re-encode to JSON string.
	infraBytes, err = json.Marshal(infra)
	// Check for errors.
	if err != nil {
		// Compose the error message
		return fmt.Errorf("Error re-encoding infrastructure to JSON: %w", err)
	}

	infraString := string(infraBytes)

	// Check correctness of infrastructure file
	fmt.Println("🔄 Checking if infrastructure is correct.")
	if !infrastructure.IsInfrastructureJsonCorrect(infra) {
		fmt.Println("❌ The infrastructure JSON is NOT correct.")
		return
	}
	fmt.Println("✅ The infrastructure JSON is correct.")

	// Get locations to deploy
	fmt.Println("🔄 Getting all locations of infrastructure.")
	listOfLocations := infrastructure.GetAllLocations(infra, inEvery, inAreas, exceptIn)
	if len(listOfLocations) == 0 {
		fmt.Println("❌ The input does not correspond to any location.")
		return
	}

	// Detect system to use correct shell
	shellPreamble := ""
	if runtime.GOOS == "windows" {
		shellPreamble = "cmd.exe /c "
	}

	// Deploy to all locations
	for _, location := range listOfLocations {
		conf := location.MainLocation
		infrastructureBase64 := base64.StdEncoding.EncodeToString([]byte(infraString))
		envVariablesString := fmt.Sprintf("--env=LOCATION_ID=%s --env=EDGE_DEPLOYMENT_IN_EVERY=%s --env=EDGE_INFRASTRUCTURE=%s --env=REDIS_HOST=%s --env=REDIS_PORT=%s --env=REDIS_PASSWORD=%s --env=FUNCTION_NAME=%s",
			location.AreaName, inEvery, infrastructureBase64, conf.RedisHost, conf.RedisPort, conf.RedisPassword, functionName)

		fmt.Printf("📶 Deploying on location: \"%s\", gateway: \"%s\".\n", location.AreaName, conf.OpenFaaSgateway)

		// Execute commands
		executeCommand(shellPreamble, fmt.Sprintf("faas-cli login --username admin --password %s --gateway %s", conf.OpenFaaSPassword, conf.OpenFaaSgateway))
		executeCommand(shellPreamble, fmt.Sprintf("faas-cli deploy --filter %s --gateway %s %s %s", functionName, conf.OpenFaaSgateway, envVariablesString, strings.Join(faasCliArguments, " ")))
	}
}

// Helper function to execute a command and print its output
func executeCommand(shellPreamble, command string) {
	cmd := exec.Command(shellPreamble + command)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		return
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("command out | %s\n", scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		return
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		return
	}
}
