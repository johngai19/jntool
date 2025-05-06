#!/usr/bin/env bash
set -euo pipefail

# Script to generate cobra command docs
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
#echo $"REPO_ROOT: $REPO_ROOT"

DOCS_DIR="$REPO_ROOT/docs/commands"

# ensure output directory exists
mkdir -p "$DOCS_DIR"

echo "🔧 Generating CLI docs into $DOCS_DIR…"

# invoke the cobra “docs” subcommand via the main entrypoint
go run "$REPO_ROOT/scripts/gen_docs.go" -output "$DOCS_DIR"

echo "✅ Docs generated at $DOCS_DIR"
