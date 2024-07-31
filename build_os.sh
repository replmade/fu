#!/bin/zsh

# Set the program name
PROGRAM_NAME="fu"

# Mac Silicon (ARM64)
GOOS=darwin GOARCH=arm64 go build -o ${PROGRAM_NAME}_macos_arm64

# Mac Intel (AMD64)
GOOS=darwin GOARCH=amd64 go build -o ${PROGRAM_NAME}_macos_amd64

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o ${PROGRAM_NAME}_linux_arm64

# Linux Intel 64 (AMD64)
GOOS=linux GOARCH=amd64 go build -o ${PROGRAM_NAME}_linux_amd64

echo "Compilation completed for all target architectures."
