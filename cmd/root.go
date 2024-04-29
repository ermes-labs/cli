package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:     "ermes-cli",
	Short:   "CLI for Ermes",
	Version: "0.0.1",
}
