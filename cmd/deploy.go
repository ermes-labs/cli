package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var DeployCmd = &cobra.Command{
	Use:   "deploy <function name> <path to infrastructure file>",
	Short: "Deploy a function to specified infrastructure.",
	Long:  "Deploy a function to specified infrastructure.",
	Args:  cobra.ExactArgs(2), // Expecting exactly 2 arguments
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implementation of the deploy command
		log.Println("Deploying function", args[0], "to infrastructure", args[1])
	},
}

func init() {
	DeployCmd.Flags().String("inEvery", "", "Area type identifier for deployment")
	DeployCmd.Flags().String("inAreas", "", "Name of the areas to deploy the function")
	DeployCmd.Flags().String("exceptIn", "", "Name of the areas to NOT deploy the function")
	DeployCmd.Flags().StringArray("faas-cli", []string{}, "faas-cli deploy compatible parameters")
}
