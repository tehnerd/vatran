#!/usr/bin/env bash
set -euo pipefail

# setup_submodules.sh
# Initializes git submodules and applies necessary patches for the build.

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PATCH_DIR="${ROOT_DIR}/upstream_patch"
SUBMODULE_DIR="${ROOT_DIR}/external/katran"

echo "Initializing git submodules..."
git -C "${ROOT_DIR}" submodule update --init --recursive

if [[ -d "${SUBMODULE_DIR}" ]]; then
  echo "Applying patches to external/katran..."

  for patch_file in "${PATCH_DIR}"/*.patch; do
    if [[ -f "${patch_file}" ]]; then
      patch_name="$(basename "${patch_file}")"
      # Check if patch is already applied by doing a dry-run
      if git -C "${SUBMODULE_DIR}" apply --check "${patch_file}" 2>/dev/null; then
        echo "  Applying: ${patch_name}"
        git -C "${SUBMODULE_DIR}" apply "${patch_file}"
      else
        echo "  Skipping (already applied or conflicts): ${patch_name}"
      fi
    fi
  done
else
  echo "error: submodule directory not found at ${SUBMODULE_DIR}" >&2
  exit 1
fi

echo "Setup complete."
