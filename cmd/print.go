package cmd

import (
	"log"
	"os"

	"github.com/ermes-labs/api-go/infrastructure"
	"github.com/spf13/cobra"
)

var PrintCmd = &cobra.Command{
	Use:     "print <path to infrastructure file>",
	Short:   "Print the infrastructure",
	Long:    "Print the infrastructure defined in the specified file.",
	Example: `  ermes-cli print my-infra.json`,
	Args:    cobra.ExactArgs(1), // Expecting exactly 1 argument
	Run: func(cmd *cobra.Command, args []string) {
		// Log the print.
		log.Println("Printing infrastructure defined in", args[0])

		// Get the infrastructure file name.
		infrastructureFileName := args[0]

		// Read and print the infrastructure file.
		infraBytes, err := os.ReadFile(infrastructureFileName)
		// Check for errors.
		if err != nil {
			log.Fatal("Error reading infrastructure file:", err)
			return
		}

		// Unmarshal and print the infrastructure JSON.
		infra, _, err := infrastructure.UnmarshalInfrastructure(infraBytes)
		// Check for errors.
		if err != nil {
			log.Fatal("Error parsing infrastructure JSON:", err)
			return
		}

		// Print the infrastructure.
		log.Println(infra.String())
	},
}

func init() {
	// Add the print command to the root command.
	RootCmd.AddCommand(PrintCmd)
}
