package cmd

import (
	"log"
	"os"

	"github.com/ermes-labs/api-go/infrastructure"
	"github.com/ermes-labs/cli/core"
	"github.com/ermes-labs/cli/query"
	"github.com/spf13/cobra"
)

var DeployCmd = &cobra.Command{
	Use:   "deploy <function name> <path to infrastructure file>",
	Short: "Deploy a function to specified infrastructure.",
	Long:  "Deploy a function to specified infrastructure.",
	Example: `  ermes-cli deploy my-function my-infra.json
	ermes-cli deploy my-function my-infra.json --deploy-in "#Milan"
	ermes-cli deploy my-function my-infra.json --ermes-cli "--gateway http://localhost:8080"`,
	Args: cobra.ExactArgs(2), // Expecting exactly 2 arguments
	Run: func(cmd *cobra.Command, args []string) {
		// Log the deployment.
		log.Println("Deploying function", args[0], "to infrastructure", args[1])

		// Get the function name and infrastructure file name.
		functionName, infrastructureFileName := args[0], args[1]

		// Read and parse the infrastructure file.
		infraBytes, err := os.ReadFile(infrastructureFileName)
		// Check for errors.
		if err != nil {
			log.Fatal("Error reading infrastructure file:", err)
			return
		}

		// Parse the infrastructure JSON.
		infra, areasMap, err := infrastructure.UnmarshalInfrastructure(infraBytes)
		// Check for errors.
		if err != nil {
			log.Fatal("Error parsing infrastructure JSON:", err)
			return
		}

		// Read the "in" flag if exists.
		in, _ := cmd.Flags().GetString("deploy-in")
		// If the "in" flag is not empty, use it as the area type identifier.
		if in == "" {
			in = "*"
		}

		// Compute the areas to deploy the function to.
		areas, err := query.CollectAreas(infra, areasMap, in)
		// Check for errors.
		if err != nil {
			log.Fatal("Error collecting areas:", err)
			return
		}

		// Get the faas-cli arguments.
		openFaasCliArguments, err := cmd.Flags().GetStringArray("faas-cli")
		// Check for errors.
		if err != nil {
			log.Fatal("Error getting faas-cli arguments:", err)
			return
		}

		core.Deploy(functionName, openFaasCliArguments, areas)
	},
}

func init() {
	DeployCmd.Flags().String("deploy-in", "", "area type identifier for deployment")
	DeployCmd.Flags().StringArray("faas-cli", []string{}, "faas-cli deploy compatible parameters")
	// Add the deploy command to the root command.
	RootCmd.AddCommand(DeployCmd)
}
