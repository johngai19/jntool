package helm

import (
	"fmt"

	"github.com/spf13/cobra"
)

var HelmCmd = &cobra.Command{
	Use:   "helm",
	Short: "helm related commands",
	Long:  `This command contains subcommands related to Helm, a package manager for Kubernetes.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Executing Helm variable extraction functionality")
	},
}
