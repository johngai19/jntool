package helm

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var applyVarsCmd = &cobra.Command{
	Use:   "apply-vars [chartDir]",
	Short: "Backup and replace placeholders in Helm chart values files",
	Long: `Backs up values.yaml, values-tag.yaml and variables.json (as .default),
then replaces @{VAR} placeholders in the values files using variables.json,
validates the result is legal YAML, and writes numbered backups like
values.yaml.1.bak, values-tag.yaml.1.bak, variables.json.1.bak.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		chartDir := args[0]

		// Step 1: Backup default copies
		defs := []struct{ src, dst string }{
			{"values.yaml", "values.default"},
			{"values-tag.yaml", "values-tag.default"},
			{"variables.json", "variables.default"},
		}
		for _, d := range defs {
			if err := copyFile(
				filepath.Join(chartDir, d.src),
				filepath.Join(chartDir, d.dst),
			); err != nil {
				return fmt.Errorf("backup %s -> %s failed: %w", d.src, d.dst, err)
			}
		}

		// Step 2: Load variables.json
		varsPath := filepath.Join(chartDir, "variables.json")
		raw, err := os.ReadFile(varsPath)
		if err != nil {
			return fmt.Errorf("read variables.json: %w", err)
		}
		var vars map[string]string
		if err := json.Unmarshal(raw, &vars); err != nil {
			return fmt.Errorf("unmarshal variables.json: %w", err)
		}

		// Step 3: Replace placeholders and validate YAML
		re := regexp.MustCompile(`@\{([^}]+)\}`)
		targets := []string{"values.yaml", "values-tag.yaml"}
		for _, name := range targets {
			path := filepath.Join(chartDir, name)
			orig, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("read %s: %w", name, err)
			}
			replaced := re.ReplaceAllStringFunc(string(orig), func(m string) string {
				key := strings.TrimSuffix(strings.TrimPrefix(m, "@{"), "}")
				if v, ok := vars[key]; ok {
					return v
				}
				return m
			})
			// validate YAML
			var tmp interface{}
			if err := yaml.Unmarshal([]byte(replaced), &tmp); err != nil {
				return fmt.Errorf("invalid YAML in %s after replacement: %w", name, err)
			}
			bak := getNextBackupName(chartDir, name)
			if err := os.WriteFile(filepath.Join(chartDir, bak), []byte(replaced), 0644); err != nil {
				return fmt.Errorf("write %s: %w", bak, err)
			}
		}

		// Step 4: Numbered backup of variables.json
		varsBak := getNextBackupName(chartDir, "variables.json")
		if err := copyFile(varsPath, filepath.Join(chartDir, varsBak)); err != nil {
			return fmt.Errorf("backup variables.json -> %s failed: %w", varsBak, err)
		}

		fmt.Printf("Completed apply-vars on %s\n", chartDir)
		return nil
	},
}

func init() {
	HelmCmd.AddCommand(applyVarsCmd)
}

// copyFile copies a file from src to dst.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}

// getNextBackupName returns the next available backup filename like "name.1.bak".
func getNextBackupName(dir, name string) string {
	for i := 1; ; i++ {
		candidate := fmt.Sprintf("%s.%d.bak", name, i)
		if _, err := os.Stat(filepath.Join(dir, candidate)); os.IsNotExist(err) {
			return candidate
		}
	}
}
