#!/bin/bash

# TARGET_DIR="gen"
# TARGET_DIR="ultimate"
# TARGET_DIR="sync_table"
TARGET_DIR="without_db"

PROJECT_DIR=$(dirname "$0")
GENERATE_DIR="$PROJECT_DIR/cmd/$TARGET_DIR"

cd "$GENERATE_DIR" || exit

echo "Start Generating"
go run .
