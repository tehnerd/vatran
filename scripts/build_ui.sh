#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
THIRDPARTY_DIR="${ROOT_DIR}/thirdparty"
UI_DIR="${ROOT_DIR}/ui"
DIST_DIR="${UI_DIR}/dist"

ESBUILD_BIN="${THIRDPARTY_DIR}/esbuild"
if [[ ! -x "${ESBUILD_BIN}" ]]; then
  if command -v esbuild >/dev/null 2>&1; then
    ESBUILD_BIN="$(command -v esbuild)"
  else
    echo "error: esbuild not found. Run scripts/get_ui_deps.sh or install esbuild." >&2
    exit 1
  fi
fi

mkdir -p "${DIST_DIR}"

"${ESBUILD_BIN}" \
  "${UI_DIR}/src/index.js" \
  --bundle \
  --platform=browser \
  --format=iife \
  --outfile="${DIST_DIR}/app.js" \
  --minify \
  --define:process.env.NODE_ENV='"production"'

cp "${UI_DIR}/src/index.html" "${DIST_DIR}/index.html"
cp "${UI_DIR}/src/styles.css" "${DIST_DIR}/styles.css"
