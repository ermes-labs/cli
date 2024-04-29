package cmd

import (
	"log"
	"os/exec"

	"github.com/spf13/cobra"
)

// templatePullCmd allows the user to fetch a template from a repository
var TemplatePullCmd = &cobra.Command{
	Use:   `pull [REPOSITORY_URL]`,
	Short: `Downloads templates from the specified git repo`,
	Long: `Downloads templates from the specified git repo specified by [REPOSITORY_URL], and copies the 'template'
directory from the root of the repo, if it exists (The default repo is https://github.com/ermes-labs/templates).

[REPOSITORY_URL] may specify a specific branch or tag to copy by adding a URL fragment with the branch or tag name.
	`,
	Example: `
  ermes-cli template pull https://github.com/ermes-labs/templates
`,
	Args: cobra.RangeArgs(0, 1),
	Run: func(_ *cobra.Command, args []string) {
		repositoryURL := "https://github.com/ermes-labs/templates"
		if len(args) > 0 {
			repositoryURL = args[0]
		}
		// Exec using the underlying ermes-cli.
		err := exec.Command("faas-cli", "template", "pull", repositoryURL).Run()
		// Check for errors.
		if err != nil {
			// Log the error.
			log.Fatal("Error pulling template:", err)
		}
	},
}

func init() {
	// Add the template pull command to the template command.
	RootCmd.AddCommand(TemplatePullCmd)
}
