package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// DefaultTemplateRepository contains the Git repo for the official templates
const DefaultTemplateRepository = "https://github.com/openfaas/templates.git"

// TemplateDirectory contains the directory where templates are stored
const TemplateDirectory = "./template/"

var (
	defaultGateway = "http://127.0.0.1:8080"
	language       string
	gateway        string
	handlerDir     string
	imagePrefix    string
	appendFile     string
	list           bool
	quiet          bool
	memoryLimit    string
	cpuLimit       string
	memoryRequest  string
	cpuRequest     string
)

func init() {
	NewFunctionCmd.Flags().StringVar(&language, "lang", "", "Language or template to use")
	NewFunctionCmd.Flags().StringVarP(&gateway, "gateway", "g", defaultGateway, "Gateway URL to store in YAML stack file")
	NewFunctionCmd.Flags().StringVar(&handlerDir, "handler", "", "directory the handler will be written to")
	NewFunctionCmd.Flags().StringVarP(&imagePrefix, "prefix", "p", "", "Set prefix for the function image")

	NewFunctionCmd.Flags().StringVar(&memoryLimit, "memory-limit", "", "Set a limit for the memory")
	NewFunctionCmd.Flags().StringVar(&cpuLimit, "cpu-limit", "", "Set a limit for the CPU")

	NewFunctionCmd.Flags().StringVar(&memoryRequest, "memory-request", "", "Set a request or the memory")
	NewFunctionCmd.Flags().StringVar(&cpuRequest, "cpu-request", "", "Set a request value for the CPU")

	NewFunctionCmd.Flags().BoolVar(&list, "list", false, "List available languages")
	NewFunctionCmd.Flags().StringVarP(&appendFile, "append", "a", "", "Append to existing YAML file")
	NewFunctionCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Skip template notes")

	RootCmd.AddCommand(NewFunctionCmd)
}

// NewFunctionCmd displays newFunction information
var NewFunctionCmd = &cobra.Command{
	Use:   "new FUNCTION_NAME --lang=FUNCTION_LANGUAGE [--gateway=http://host:port] | --list | --append=STACK_FILE)",
	Short: "Create a new template in the current folder with the name given as name",
	Long: `The new command creates a new function based upon hello-world in the given
language or type in --list for a list of languages available.`,
	Example: `  ermes-cli new chatbot --lang node
  ermes-cli new chatbot --lang node --append stack.yml
  ermes-cli new text-parser --lang python --quiet
  ermes-cli new text-parser --lang python --gateway http://mydomain:8080
  ermes-cli new --list`,
	Args: cobra.MinimumNArgs(1),
	RunE: runNewFunction,
}

func runNewFunction(cmd *cobra.Command, args []string) error {
	if list {
		_, err := os.ReadDir(TemplateDirectory)
		if err != nil {
			return fmt.Errorf(`no language templates were found.

Download templates:
  ermes-cli template pull           download the default templates`)
		}
	}

	arg := []string{"new", args[0]}
	// Get the raw flags string
	cmd.Flags().Visit(func(f *pflag.Flag) {
		if f.Changed {
			arg = append(arg, f.Name, f.Value.String())
		}
	})

	// Pass everything to the faas-cli new command
	return exec.Command("faas-cli", arg...).Run()
}
