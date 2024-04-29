package cmd

import (
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	RootCmd.AddCommand(pushCmd)

	pushCmd.Flags().IntVar(&parallel, "parallel", 1, "Push images in parallel to depth specified.")
	// pushCmd.Flags().Var(&tagFormat, "tag", "Override latest tag on function Docker image, accepts 'digest', 'latest', 'sha', 'branch', 'describe'")
	pushCmd.Flags().BoolVar(&envsubst, "envsubst", true, "Substitute environment variables in stack.yml file")
	pushCmd.Flags().BoolVar(&quietBuild, "quiet", false, "Perform a quiet build, without showing output from Docker")
}

// pushCmd handles pushing function container images to a remote repo
var pushCmd = &cobra.Command{
	Use:   `push -f YAML_FILE [--regex "REGEX"] [--filter "WILDCARD"] [--parallel] [--tag <sha|branch>]`,
	Short: "Push OpenFaaS functions to remote registry (Docker Hub)",
	Long: `Pushes the OpenFaaS function container image(s) defined in the supplied YAML
config to a remote repository.

These container images must already be present in your local image cache.`,

	Example: `  faas-cli push -f https://domain/path/myfunctions.yml
  faas-cli push -f ./stack.yml
  faas-cli push -f ./stack.yml --parallel 4
  faas-cli push -f ./stack.yml --filter "*gif*"
  faas-cli push -f ./stack.yml --regex "fn[0-9]_.*"
  faas-cli push -f ./stack.yml --tag sha
  faas-cli push -f ./stack.yml --tag branch
  faas-cli push -f ./stack.yml --tag describe`,
	RunE: runPush,
}

func runPush(cmd *cobra.Command, args []string) error {
	arg := []string{"push"}

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			arg = append(arg, f.Name, f.Value.String())
		}
	})

	return exec.Command("faas-cli", arg...).Run()
}
