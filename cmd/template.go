package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(TemplateCmd)
}

// templateCmd allows access to store and pull commands
var TemplateCmd = &cobra.Command{
	Use:     `template [COMMAND]`,
	Short:   "OpenFaaS template store and pull commands",
	Long:    "Allows pulling custom templates",
	Example: `  ermes-cli template pull https://github.com/custom/template`,
}
