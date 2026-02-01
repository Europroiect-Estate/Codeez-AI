#!/usr/bin/env bash
set -e
# Codeez install script — download latest release or build from source
# Usage: curl -sSL https://raw.githubusercontent.com/Europroiect-Estate/Codeez-AI/main/scripts/install.sh | bash
# Directory priority: CODEEZ_INSTALL_DIR > XDG_BIN_DIR > /usr/local/bin > ~/.local/bin > ~/bin

REPO="${REPO:-Europroiect-Estate/Codeez-AI}"
REPO_URL="https://github.com/${REPO}.git"
RAW_BASE="https://raw.githubusercontent.com/${REPO}/main"

# Detect OS and arch
OS_RAW=$(uname -s)
ARCH_RAW=$(uname -m)
case "$ARCH_RAW" in
  x86_64) ARCH=amd64 ;;
  aarch64|arm64) ARCH=arm64 ;;
  *) ARCH=amd64 ;;
esac
case "$OS_RAW" in
  Darwin) OS=Darwin ;;
  Linux) OS=Linux ;;
  MINGW*|MSYS*|CYGWIN*) OS=Windows ;;
  *) OS=Linux ;;
esac

# Install directory
if [ -n "$CODEEZ_INSTALL_DIR" ]; then
  INSTALL_DIR="$CODEEZ_INSTALL_DIR"
elif [ -n "$XDG_BIN_DIR" ]; then
  INSTALL_DIR="$XDG_BIN_DIR"
elif [ -w /usr/local/bin ] 2>/dev/null; then
  INSTALL_DIR=/usr/local/bin
elif [ -n "$HOME" ]; then
  if [ -d "$HOME/.local/bin" ] || mkdir -p "$HOME/.local/bin" 2>/dev/null; then
    INSTALL_DIR="${HOME}/.local/bin"
  elif [ -d "$HOME/bin" ] || mkdir -p "$HOME/bin" 2>/dev/null; then
    INSTALL_DIR="${HOME}/bin"
  else
    INSTALL_DIR="${HOME}/.local/bin"
    mkdir -p "$INSTALL_DIR"
  fi
else
  INSTALL_DIR="."
fi

install_from_release() {
  local version="$1"
  local ver="${version#v}"
  local name="codeez_${ver}_${OS}_${ARCH}"
  local archive
  if [ "$OS" = "Windows" ]; then
    archive="${name}.zip"
  else
    archive="${name}.tar.gz"
  fi
  local url="https://github.com/${REPO}/releases/download/${version}/${archive}"
  local checksum_url="https://github.com/${REPO}/releases/download/${version}/checksums.txt"
  echo "Installing codeez ${version} from release to $INSTALL_DIR"
  TMP=$(mktemp -d)
  if ! curl -sfL -o "$TMP/$archive" "$url"; then
    echo "Download failed: $url"
    rm -rf "$TMP"
    return 1
  fi
  curl -sfL -o "$TMP/checksums.txt" "$checksum_url" 2>/dev/null || true
  (cd "$TMP" && sha256sum -c --ignore-missing checksums.txt 2>/dev/null || shasum -a 256 -c checksums.txt 2>/dev/null || true)
  if [ "$OS" = "Windows" ]; then
    unzip -o -q "$TMP/$archive" -d "$TMP"
    cp "$TMP/codeez.exe" "$INSTALL_DIR/" 2>/dev/null || cp "$TMP/codeez" "$INSTALL_DIR/"
  else
    tar -xzf "$TMP/$archive" -C "$TMP"
    cp "$TMP/codeez" "$INSTALL_DIR/"
  fi
  rm -rf "$TMP"
  chmod +x "$INSTALL_DIR/codeez" 2>/dev/null || true
  echo "Installed: $INSTALL_DIR/codeez"
  "$INSTALL_DIR/codeez" version
}

install_from_source() {
  echo "No release found. Building from source..."
  if ! command -v git >/dev/null 2>&1; then
    echo "Error: git is required to install from source. Install git and retry."
    exit 1
  fi
  if ! command -v go >/dev/null 2>&1; then
    echo "Error: Go is required to build from source. Install Go 1.23+ from https://go.dev and retry."
    exit 1
  fi
  TMP=$(mktemp -d)
  git clone --depth 1 "$REPO_URL" "$TMP/repo"
  (cd "$TMP/repo" && go build -ldflags "-s -w" -o codeez ./cmd/codeez)
  cp "$TMP/repo/codeez" "$INSTALL_DIR/"
  rm -rf "$TMP"
  chmod +x "$INSTALL_DIR/codeez"
  echo "Installed (from source): $INSTALL_DIR/codeez"
  "$INSTALL_DIR/codeez" version
}

# Try releases/latest first
VERSION="${VERSION:-}"
if [ -z "$VERSION" ]; then
  LATEST=$(curl -sfL "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null | grep -o '"tag_name": *"[^"]*"' | cut -d'"' -f4)
  VERSION="$LATEST"
fi

if [ -n "$VERSION" ]; then
  # Normalize: v0.1.0 -> 0.1.0 for archive name
  VERSION_STR="${VERSION#v}"
  if install_from_release "$VERSION"; then
    exit 0
  fi
fi

# Fallback: build from source
install_from_source
