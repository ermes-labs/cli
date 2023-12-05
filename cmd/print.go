package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var PrintCmd = &cobra.Command{
	Use:   "print <path to infrastructure file>",
	Short: "Print the infrastructure",
	Long:  "Print the infrastructure defined in the specified file.",
	Args:  cobra.ExactArgs(1), // Expecting exactly 1 argument
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implementation of the print command
		log.Println("Printing infrastructure defined in", args[0])
	},
}
