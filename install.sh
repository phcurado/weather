#!/usr/bin/env sh
# install.sh — download the latest weather release binary into ~/.local/bin.
set -eu

REPO="phcurado/weather"
BIN_DIR="${BIN_DIR:-$HOME/.local/bin}"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

os=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)
case "$arch" in
  x86_64|amd64) arch=amd64 ;;
  aarch64|arm64) arch=arm64 ;;
  *) echo "unsupported arch: $arch" >&2; exit 1 ;;
esac
case "$os" in
  linux|darwin) ;;
  *) echo "unsupported os: $os" >&2; exit 1 ;;
esac

tag=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" \
      | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' | head -n1)
if [ -z "$tag" ]; then
  echo "could not resolve latest tag" >&2
  exit 1
fi

asset="weather_${tag#v}_${os}_${arch}.tar.gz"
url="https://github.com/$REPO/releases/download/$tag/$asset"

echo "downloading $url"
curl -fsSL "$url" -o "$TMP/$asset"
tar -xzf "$TMP/$asset" -C "$TMP"

mkdir -p "$BIN_DIR"
install -m 0755 "$TMP/weather" "$BIN_DIR/weather"

echo "installed: $BIN_DIR/weather ($tag)"
case ":$PATH:" in
  *":$BIN_DIR:"*) ;;
  *) echo "note: $BIN_DIR is not in PATH" >&2 ;;
esac
