#!/bin/bash

# filepath: /home/raymond/dev/me/go/shutdown-by-ping/build.sh

# Set the output binary name
OUTPUT_BINARY="ping_shutdown"

# Compile the Go program
echo "Compiling the Go program..."
go build -o $OUTPUT_BINARY main.go

# Check if the compilation was successful
if [ $? -eq 0 ]; then
    echo "Compilation successful. Binary created: $OUTPUT_BINARY"
else
    echo "Compilation failed."
    exit 1
fi