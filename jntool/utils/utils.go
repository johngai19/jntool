package utils

import (
	"fmt"

	"github.com/spf13/cobra"
)

var UtilsCmd = &cobra.Command{
	Use:   "utils",
	Short: "Commands related to utility functions",
	Long:  `This command contains subcommands related to utility functions.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Utils command executed")
	},
}
