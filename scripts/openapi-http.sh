#!/bin/bash
set -e

readonly output_dir="$1"

swag init -g "$output_dir/main.go" -o "$output_dir/api/openapi/" --ot go,yaml
