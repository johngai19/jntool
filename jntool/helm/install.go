package helm

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
)

var installCmd = &cobra.Command{
	Use:   "install [chartDir] [releaseName]",
	Short: "Install umbrella chart using latest backup values",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		chartDir, release := args[0], args[1]
		defer restoreDefaults(chartDir)
		ns, _ := cmd.Flags().GetString("namespace")

		// check for dry-run flag
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		// 1) find & restore latest backups
		valsBak, err := findLatestBackup(chartDir, "values.yaml")
		if err != nil {
			return err
		}
		tagBak, err := findLatestBackup(chartDir, "values-tag.yaml")
		if err != nil {
			return err
		}
		if err := copyFile(filepath.Join(chartDir, valsBak), filepath.Join(chartDir, "values.yaml")); err != nil {
			return err
		}
		if err := copyFile(filepath.Join(chartDir, tagBak), filepath.Join(chartDir, "values-tag.yaml")); err != nil {
			return err
		}

		// 2) merge the two YAMLs
		baseBytes, err := os.ReadFile(filepath.Join(chartDir, "values.yaml"))
		if err != nil {
			return fmt.Errorf("read base values.yaml: %w", err)
		}
		tagBytes, err := os.ReadFile(filepath.Join(chartDir, "values-tag.yaml"))
		if err != nil {
			return fmt.Errorf("read values-tag.yaml: %w", err)
		}
		var baseMap, tagMap map[string]interface{}
		if err := yaml.Unmarshal(baseBytes, &baseMap); err != nil {
			return fmt.Errorf("unmarshal base values.yaml: %w", err)
		}
		if err := yaml.Unmarshal(tagBytes, &tagMap); err != nil {
			return fmt.Errorf("unmarshal values-tag.yaml: %w", err)
		}
		for k, v := range tagMap {
			baseMap[k] = v
		}
		mergedBytes, err := yaml.Marshal(baseMap)
		if err != nil {
			return fmt.Errorf("marshal merged values: %w", err)
		}

		// 3) debug output (you can remove/comment this later)
		fmt.Printf("Merged values content:\n%s\n", string(mergedBytes))

		// 4) overwrite the *default* values.yaml in-place
		if err := os.WriteFile(filepath.Join(chartDir, "values.yaml"), mergedBytes, 0o644); err != nil {
			return fmt.Errorf("write merged values.yaml: %w", err)
		}

		// 5) Helm client config
		settings := cli.New()
		cfg := new(action.Configuration)
		if err := cfg.Init(
			settings.RESTClientGetter(),
			ns,
			os.Getenv("HELM_DRIVER"),
			func(_ string, _ ...interface{}) {},
		); err != nil {
			return err
		}

		installer := action.NewInstall(cfg)
		installer.ReleaseName = release
		installer.Namespace = ns
		installer.CreateNamespace = true
		installer.DryRun = dryRun // simulate if requested

		// 6) locate & load chart
		chartPath, err := installer.ChartPathOptions.LocateChart(chartDir, settings)
		if err != nil {
			return err
		}
		chart, err := loader.Load(chartPath)
		if err != nil {
			return err
		}

		// 7) run install with the merged map
		if _, err := installer.Run(chart, baseMap); err != nil {
			return fmt.Errorf("install failed: %w", err)
		}

		fmt.Printf("Installed %s in namespace %s\n", release, ns)
		return nil
	},
}

func init() {
	installCmd.Flags().StringP("namespace", "n", "default", "Kubernetes namespace")
	installCmd.Flags().Bool("dry-run", false, "simulate installation (dry-run)")
	HelmCmd.AddCommand(installCmd)
}

// findLatestBackup returns the name of the highest-numbered .bak for base.
func findLatestBackup(dir, base string) (string, error) {
	re := regexp.MustCompile("^" + regexp.QuoteMeta(base) + `\.(\d+)\.bak$`)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	var nums []int
	for _, e := range entries {
		if m := re.FindStringSubmatch(e.Name()); m != nil {
			if n, err := strconv.Atoi(m[1]); err == nil {
				nums = append(nums, n)
			}
		}
	}
	if len(nums) == 0 {
		return "", fmt.Errorf("no backups found for %s", base)
	}
	sort.Ints(nums)
	return fmt.Sprintf("%s.%d.bak", base, nums[len(nums)-1]), nil
}

// restoreDefaults copies *.default back to their working filenames.
func restoreDefaults(dir string) {
	files := []struct{ src, dst string }{
		{"values.default", "values.yaml"},
		{"values-tag.default", "values-tag.yaml"},
		{"variables.default", "variables.json"},
	}
	for _, f := range files {
		_ = copyFile(filepath.Join(dir, f.src), filepath.Join(dir, f.dst))
	}
}
