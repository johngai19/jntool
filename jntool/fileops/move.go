package fileops

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var moveCmd = &cobra.Command{
	Use:   "move",
	Short: "Scan a folder and move files by extension and/or substring match",
	Long: `Scan the source folder, filter files by extension and/or substring condition,
and move matching files to one or more destination folders. Optionally rename during move.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		src, _ := cmd.Flags().GetString("source")
		dests, _ := cmd.Flags().GetStringSlice("dest")
		exts, _ := cmd.Flags().GetStringSlice("ext")
		substr, _ := cmd.Flags().GetString("substr")
		cond, _ := cmd.Flags().GetString("cond")
		renameFrom, _ := cmd.Flags().GetString("rename-from")
		renameTo, _ := cmd.Flags().GetString("rename-to")
		doRename, _ := cmd.Flags().GetBool("rename")

		entries, err := os.ReadDir(src)
		if err != nil {
			return fmt.Errorf("reading source folder: %w", err)
		}

		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			name := e.Name()
			fullPath := filepath.Join(src, name)

			extMatch := len(exts) == 0
			for _, ex := range exts {
				if strings.EqualFold(filepath.Ext(name), ex) {
					extMatch = true
					break
				}
			}

			strMatch := substr == "" || strings.Contains(strings.ToLower(name), strings.ToLower(substr))

			ok := false
			if cond == "and" {
				ok = extMatch && strMatch
			} else {
				ok = extMatch || strMatch
			}
			if !ok {
				continue
			}

			newName := name
			if doRename && renameFrom != "" {
				newName = strings.ReplaceAll(newName, renameFrom, renameTo)
			}

			for _, d := range dests {
				if err := os.MkdirAll(d, 0755); err != nil {
					return fmt.Errorf("creating dest %s: %w", d, err)
				}
				dstPath := filepath.Join(d, newName)
				if err := moveFile(fullPath, dstPath); err != nil {
					return fmt.Errorf("moving %s -> %s: %w", fullPath, dstPath, err)
				}
			}
		}
		fmt.Println("Move operation completed.")
		return nil
	},
}

func init() {
	FileopsCmd.AddCommand(moveCmd)
	moveCmd.Flags().StringP("source", "s", ".", "source folder to scan")
	moveCmd.Flags().StringSliceP("dest", "d", nil, "destination folders (comma-separated)")
	moveCmd.Flags().StringSlice("ext", []string{}, "file extensions to match (e.g. .pdf,.epub)")
	moveCmd.Flags().String("substr", "", "substring to match in file name")
	moveCmd.Flags().String("cond", "or", "filter condition: and|or")
	moveCmd.Flags().BoolP("rename", "r", false, "enable rename replace")
	moveCmd.Flags().String("rename-from", "", "substring to replace in file name")
	moveCmd.Flags().String("rename-to", "", "replacement substring")
}

// moveFile tries to rename, falling back to copy+remove on cross-device
func moveFile(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	// fallback copy+remove
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
	return os.Remove(src)
}
