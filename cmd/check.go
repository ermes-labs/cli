package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var CheckCmd = &cobra.Command{
	Use:   "check <path to infrastructure file>",
	Short: "Check the infrastructure",
	Long:  "Check the infrastructure defined in the specified file.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implementation of the check command
		log.Println("Checking infrastructure defined in", args[0])
	},
}
