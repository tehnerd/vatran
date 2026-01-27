#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
THIRDPARTY_DIR="${ROOT_DIR}/thirdparty"

mkdir -p "${THIRDPARTY_DIR}"

fetch() {
  local url="$1"
  local out="$2"
  if command -v curl >/dev/null 2>&1; then
    curl -v -fsSL "$url" -o "$out"
  elif command -v wget >/dev/null 2>&1; then
    wget -q "$url" -O "$out"
  else
    echo "error: curl or wget required" >&2
    exit 1
  fi
}

# Pinned versions
REACT_VERSION="18.2.0"
ROUTER_VERSION="6.22.3"
REMIX_ROUTER_VERSION="1.15.3"
CHART_VERSION="4.4.1"
HTM_VERSION="3.1.1"
ESBUILD_VERSION="0.19.12"

fetch "https://unpkg.com/react@${REACT_VERSION}/umd/react.production.min.js" \
  "${THIRDPARTY_DIR}/react.production.min.js"
fetch "https://unpkg.com/react-dom@${REACT_VERSION}/umd/react-dom.production.min.js" \
  "${THIRDPARTY_DIR}/react-dom.production.min.js"
fetch "https://unpkg.com/@remix-run/router@${REMIX_ROUTER_VERSION}/dist/router.umd.min.js" \
  "${THIRDPARTY_DIR}/remix-router.umd.min.js"
fetch "https://unpkg.com/react-router@${ROUTER_VERSION}/dist/umd/react-router.production.min.js" \
  "${THIRDPARTY_DIR}/react-router.production.min.js"
fetch "https://unpkg.com/react-router-dom@${ROUTER_VERSION}/dist/umd/react-router-dom.production.min.js" \
  "${THIRDPARTY_DIR}/react-router-dom.production.min.js"
fetch "https://unpkg.com/chart.js@${CHART_VERSION}/dist/chart.umd.js" \
  "${THIRDPARTY_DIR}/chart.umd.min.js"
fetch "https://unpkg.com/htm@${HTM_VERSION}/dist/htm.umd.js" \
  "${THIRDPARTY_DIR}/htm.umd.js"

ESBUILD_PLATFORM="$(uname -s | tr '[:upper:]' '[:lower:]')"
ESBUILD_ARCH="$(uname -m)"
case "${ESBUILD_ARCH}" in
  x86_64|amd64) ESBUILD_ARCH="x64" ;;
  aarch64|arm64) ESBUILD_ARCH="arm64" ;;
  *)
    echo "warning: unsupported arch for esbuild binary: ${ESBUILD_ARCH}" >&2
    ESBUILD_ARCH=""
    ;;
esac

if [[ -n "${ESBUILD_ARCH}" ]]; then
  ESBUILD_NAME="${ESBUILD_PLATFORM}-${ESBUILD_ARCH}"
  ESBUILD_URL="https://registry.npmjs.org/@esbuild/${ESBUILD_NAME}/-/${ESBUILD_NAME}-${ESBUILD_VERSION}.tgz"
  #ESBUILD_URL="https://github.com/evanw/esbuild/archive/refs/tags/v0.19.12.tar.gz"
  ESBUILD_TGZ="${THIRDPARTY_DIR}/${ESBUILD_NAME}.tgz"
  fetch "${ESBUILD_URL}" "${ESBUILD_TGZ}"
  tar -xzf "${ESBUILD_TGZ}" -C "${THIRDPARTY_DIR}"
  mv "${THIRDPARTY_DIR}/package/bin/esbuild" "${THIRDPARTY_DIR}/esbuild"
  rm -rf "${THIRDPARTY_DIR}/package" "${ESBUILD_TGZ}"
fi

cat > "${THIRDPARTY_DIR}/VERSIONS.json" <<JSON
{
  "react": "${REACT_VERSION}",
  "react-dom": "${REACT_VERSION}",
  "react-router": "${ROUTER_VERSION}",
  "@remix-run/router": "${REMIX_ROUTER_VERSION}",
  "react-router-dom": "${ROUTER_VERSION}",
  "chart.js": "${CHART_VERSION}",
  "htm": "${HTM_VERSION}",
  "esbuild": "${ESBUILD_VERSION}"
}
JSON
