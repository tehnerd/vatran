#!/bin/bash
#
# Script to copy build artifacts from _build/ to _build_go/ directory.
# These artifacts are required for Go bindings compilation.
#

set -e

# Get the script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

BUILD_DIR="${PROJECT_DIR}/_build"
BUILD_GO_DIR="${PROJECT_DIR}/_build_go"

# List of files to copy (excluding config.yaml)
FILES=(
    "balancer.bpf.o"
    "healthchecking.bpf.o"
    "healthchecking_ipip.o"
    "libbpf.a"
    "libbpfadapter.a"
    "libchhelpers.a"
    "libfolly.a"
    "libiphelpers.a"
    "libkatran_capi_static.a"
    "libkatranlb.a"
    "libkatransimulator.a"
    "libmurmur3.a"
    "libpcapwriter.a"
    "xdp_root.o"
)

# Create _build_go directory if it doesn't exist
mkdir -p "$BUILD_GO_DIR"

echo "Copying build artifacts from $BUILD_DIR to $BUILD_GO_DIR"

# Find and copy each file
for file in "${FILES[@]}"; do
    found=$(find "$BUILD_DIR" -name "$file" -type f 2>/dev/null | head -n 1)
    if [ -n "$found" ]; then
        cp "$found" "$BUILD_GO_DIR/"
        echo "  Copied: $file"
    else
        echo "  Warning: $file not found in $BUILD_DIR"
    fi
done

echo "Done."
