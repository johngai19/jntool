package main

import (
	"log"
	"path/filepath"

	"github.com/johngai19/jntool/jntool/root"
	"github.com/spf13/cobra/doc"
)

func main() {
	outDir := filepath.Join("docs", "commands")
	log.Printf("🔧 Generating CLI docs into %s...\n", outDir)
	if err := doc.GenMarkdownTree(root.RootCmd, outDir); err != nil {
		log.Fatalf("failed to generate docs: %v", err)
	}
	log.Printf("✅ Documentation generated at %s\n", outDir)
}
