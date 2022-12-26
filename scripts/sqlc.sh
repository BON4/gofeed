#!/bin/bash
set -e

readonly config="$1"

sqlc generate -f "$1" 
