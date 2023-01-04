#!/bin/bash
set -e

readonly output_dir="$1"

swag init -g "main.go" -o "$output_dir/api/openapi/" --dir "$output_dir/,internal/common/" --ot go,yaml
