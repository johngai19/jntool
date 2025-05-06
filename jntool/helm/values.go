package helm

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const defaultOutputFormat = "json"

var valuesCmd = &cobra.Command{
	Use:   "values [file]",
	Short: "Extract all placeholder variables (starting with @) from a Helm values.yaml file",
	Long:  `Read a Helm values.yaml, scan for placeholders in the form @{VAR_NAME}, and emit a mapping VAR_NAME -> "" in JSON or YAML for later filling.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]
		outFmt, _ := cmd.Flags().GetString("output")
		if outFmt == "" {
			outFmt = defaultOutputFormat
		}

		vars, err := extractValues(input)
		if err != nil {
			return err
		}

		switch outFmt {
		case "json":
			b, err := json.MarshalIndent(vars, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(b))
		case "yaml":
			y, err := yaml.Marshal(vars)
			if err != nil {
				return fmt.Errorf("failed to marshal YAML: %w", err)
			}
			fmt.Println(string(y))
		default:
			return fmt.Errorf("unsupported output format: %s", outFmt)
		}
		return nil
	},
}

func init() {
	HelmCmd.AddCommand(valuesCmd)
	valuesCmd.Flags().StringP("output", "o", defaultOutputFormat, "output format (json or yaml)")
}

// extractValues reads the file at path, scans for all @{VAR} placeholders,
// and returns a map VAR->"".
func extractValues(path string) (map[string]interface{}, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read file %s: %w", path, err)
	}

	re := regexp.MustCompile(`@\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(string(content), -1)

	vars := make(map[string]interface{}, len(matches))
	for _, m := range matches {
		if len(m) > 1 {
			vars[m[1]] = ""
		}
	}
	return vars, nil
}
