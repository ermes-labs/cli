package cmd

import (
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Flags that are to be added to commands.
var (
	image        string
	handler      string
	functionName string
	nocache      bool
	squash       bool
	parallel     int
	shrinkwrap   bool
	buildArgs    []string
	buildArgMap  map[string]string
	buildOptions []string
	copyExtra    []string
	// tagFormat        pflag.Value
	buildLabels      []string
	buildLabelMap    map[string]string
	envsubst         bool
	quietBuild       bool
	disableStackPull bool
)

func init() {
	// Setup flags that are used by multiple commands (variables defined in faas.go)
	BuildCmd.Flags().StringVar(&image, "image", "", "Docker image name to build")
	BuildCmd.Flags().StringVar(&handler, "handler", "", "Directory with handler for function, e.g. handler.js")
	BuildCmd.Flags().StringVar(&functionName, "name", "", "Name of the deployed function")
	BuildCmd.Flags().StringVar(&language, "lang", "", "Programming language template")

	// Setup flags that are used only by this command (variables defined above)
	BuildCmd.Flags().BoolVar(&nocache, "no-cache", false, "Do not use Docker's build cache")
	BuildCmd.Flags().BoolVar(&squash, "squash", false, `Use Docker's squash flag for smaller images [experimental] `)
	BuildCmd.Flags().IntVar(&parallel, "parallel", 1, "Build in parallel to depth specified.")
	BuildCmd.Flags().BoolVar(&shrinkwrap, "shrinkwrap", false, "Just write files to ./build/ folder for shrink-wrapping")
	BuildCmd.Flags().StringArrayVarP(&buildArgs, "build-arg", "b", []string{}, "Add a build-arg for Docker (KEY=VALUE)")
	BuildCmd.Flags().StringArrayVarP(&buildOptions, "build-option", "o", []string{}, "Set a build option, e.g. dev")
	// BuildCmd.Flags().Var(&tagFormat, "tag", "Override latest tag on function Docker image, accepts 'digest', 'sha', 'branch', or 'describe', or 'latest'")
	BuildCmd.Flags().StringArrayVar(&buildLabels, "build-label", []string{}, "Add a label for Docker image (LABEL=VALUE)")
	BuildCmd.Flags().StringArrayVar(&copyExtra, "copy-extra", []string{}, "Extra paths that will be copied into the function build context")
	BuildCmd.Flags().BoolVar(&envsubst, "envsubst", true, "Substitute environment variables in stack.yml file")
	BuildCmd.Flags().BoolVar(&quietBuild, "quiet", false, "Perform a quiet build, without showing output from Docker")
	BuildCmd.Flags().BoolVar(&disableStackPull, "disable-stack-pull", false, "Disables the template configuration in the stack.yml")

	// Set bash-completion.
	_ = BuildCmd.Flags().SetAnnotation("handler", cobra.BashCompSubdirsInDir, []string{})

	RootCmd.AddCommand(BuildCmd)
}

// BuildCmd allows the user to build an OpenFaaS function container
var BuildCmd = &cobra.Command{
	Use: `build -f YAML_FILE [--no-cache] [--squash]
  faas-cli build --image IMAGE_NAME
                 --handler HANDLER_DIR
                 --name FUNCTION_NAME
                 [--lang <ruby|python|python3|node|csharp|dockerfile>]
                 [--no-cache] [--squash]
                 [--regex "REGEX"]
                 [--filter "WILDCARD"]
                 [--parallel PARALLEL_DEPTH]
                 [--build-arg KEY=VALUE]
                 [--build-option VALUE]
                 [--copy-extra PATH]`,
	Short: "Builds OpenFaaS function containers",
	Long: `Builds OpenFaaS function containers either via the supplied YAML config using
the "--yaml" flag (which may contain multiple function definitions), or directly
via flags.`,
	Example: `  faas-cli build -f https://domain/path/myfunctions.yml
  faas-cli build -f ./stack.yml --no-cache --build-arg NPM_VERSION=0.2.2
  faas-cli build -f ./stack.yml --build-option dev
  faas-cli build -f ./stack.yml --filter "*gif*"
  faas-cli build -f ./stack.yml --regex "fn[0-9]_.*"
  faas-cli build --image=my_image --lang=python --handler=/path/to/fn/
                 --name=my_fn --squash
  faas-cli build -f ./stack.yml --build-label org.label-schema.label-name="value"`,
	RunE: runBuild,
}

func runBuild(cmd *cobra.Command, args []string) error {
	arg := []string{"build"}
	// Get the raw flags string
	cmd.Flags().Visit(func(f *pflag.Flag) {
		if f.Changed {
			arg = append(arg, f.Name, f.Value.String())
		}
	})

	// Pass everything to the faas-cli new command
	return exec.Command("faas-cli", arg...).Run()
}
