package cmd

import (
	"log"
	"os"

	"github.com/ermes-labs/api-go/infrastructure"
	"github.com/spf13/cobra"
)

var CheckCmd = &cobra.Command{
	Use:     "check <path to infrastructure file>",
	Short:   "Check the infrastructure",
	Long:    "Check the infrastructure defined in the specified file.",
	Example: `  ermes-cli check my-infra.json`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Log the check.
		log.Println("Checking infrastructure defined in", args[0])

		// Get the infrastructure file name.
		infrastructureFileName := args[0]

		// Read and parse the infrastructure file.
		infraBytes, err := os.ReadFile(infrastructureFileName)
		// Check for errors.
		if err != nil {
			log.Fatal("Error reading infrastructure file:", err)
			return
		}

		// Parse the infrastructure JSON.
		_, _, err = infrastructure.UnmarshalInfrastructure(infraBytes)
		// Check for errors.
		if err != nil {
			log.Fatal("Error parsing infrastructure JSON:", err)
			return
		}

		log.Println("The infrastructure is valid.")
	},
}

func init() {
	// Add the check command to the root command.
	RootCmd.AddCommand(CheckCmd)
}
