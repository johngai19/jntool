package root

import (
	"fmt"
	"os"

	"github.com/johngai19/jntool/jntool/fileops"
	"github.com/johngai19/jntool/jntool/helm"
	"github.com/johngai19/jntool/jntool/utils"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "jntool",
	Short: "jntool is a versatile command-line toolbox",
	Long: `jntool is a command-line toolbox written in Go,
providing common functionalities such as extracting Helm values,
file search, and file move operations.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add subcommands here
	// For example:
	RootCmd.AddCommand(helm.HelmCmd)
	RootCmd.AddCommand(fileops.FileopsCmd)
	RootCmd.AddCommand(utils.UtilsCmd)
}
