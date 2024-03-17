package core

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/ermes-labs/api-go/infrastructure"
)

// Deploy function in Go
func Deploy(functionName string, openFaasCliArguments []string, areas []*infrastructure.Area) error {
	for _, area := range areas {
		nodeBytes, err := infrastructure.MarshalNode(area.Node)

		if err != nil {
			log.Fatal("Error marshalling node:", err)
		}

		envVariablesString := fmt.Sprintf(""+
			"--env=LOCATION_ID= --env=EDGE_DEPLOYMENT_IN_EVERY= --env=EDGE_INFRASTRUCTURE= --env=REDIS_HOST= --env=REDIS_PORT= --env=REDIS_PASSWORD= --env=FUNCTION_NAME="+
			"--env=NODE=%s"+
			"--env=AREA_NAME=%s"+
			"--env=HOST=%s",
			string(nodeBytes), area.AreaName, area.Host)

		log.Println("Deploying node:", area.Node)

		// Execute commands
		// exec.Command("faas-cli", "login", "--username", "admin", "--password", "admin", "--gateway", area.Host)
		cmd := exec.Command("faas-cli", "deploy", "--filter", "functionName", "--gateway", area.Host, envVariablesString, strings.Join(openFaasCliArguments, " "))
		cmdReader, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal("Error creating StdoutPipe for Cmd", err)
		}

		scanner := bufio.NewScanner(cmdReader)
		go func() {
			for scanner.Scan() {
				log.Printf("command out | %s\n", scanner.Text())
			}
		}()

		err = cmd.Start()
		if err != nil {
			log.Fatal("Error starting Cmd", err)
		}

		err = cmd.Wait()
		if err != nil {
			log.Fatal("Error waiting for Cmd", err)
		}
	}

	return nil
}
