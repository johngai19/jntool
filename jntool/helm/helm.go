package helm

import (
	"fmt"

	"github.com/spf13/cobra"
)

var HelmCmd = &cobra.Command{
	Use:   "helm",
	Short: "Extract variables from Helm values files",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Executing Helm variable extraction functionality")
	},
}
