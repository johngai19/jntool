#!/bin/bash

cd ../

# Automatically create project directory structure
mkdir -p cmd/jntool
mkdir -p internal/helm
mkdir -p internal/fileops
mkdir -p internal/utils
mkdir -p build
mkdir -p pkg
mkdir -p tests
mkdir -p docs
mkdir -p scripts

# Create sample Go files in folders
touch cmd/jntool/main.go
touch internal/helm/helm.go
touch internal/fileops/fileops.go
touch internal/utils/string_utils.go
touch tests/helm_test.go
touch tests/fileops_test.go

go mod init jntool
go mod tidy
go mod vendor

# Modify README.md file
echo "# jntool" > README.md
echo "## Description" >> README.md
echo "This is a sample project for jntool." >> README.md

# Modify .gitignore file
echo "vendor/" > .gitignore
echo "build/" >> .gitignore
echo "*.log" >> .gitignore
echo "*.tmp" >> .gitignore
echo "*.swp" >> .gitignore
echo "*.swo" >> .gitignore
echo "*.bak" >> .gitignore
echo "*.exe" >> .gitignore
echo "*.dll" >> .gitignore
echo "*.so" >> .gitignore


echo "Project directory structure generated"