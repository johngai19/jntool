package fileops

import (
	"fmt"

	"github.com/spf13/cobra"
)

var FileopsCmd = &cobra.Command{
	Use:   "fileops",
	Short: "File operations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Executing file operations functionality")
	},
}
