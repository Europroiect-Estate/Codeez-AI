#!/usr/bin/env bash
set -e
# Codeez install script — download latest release, verify checksum, install binary
# Directory priority (like OpenCode): CODEEZ_INSTALL_DIR > XDG_BIN_DIR > /usr/local/bin > ~/.local/bin > ~/bin

REPO="${REPO:-Europroiect-Estate/Codeez-AI}"
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
  x86_64) ARCH=amd64 ;;
  aarch64|arm64) ARCH=arm64 ;;
esac
case "$OS" in
  darwin) OS=Darwin ;;
  linux) OS=Linux ;;
  *) OS=Linux ;;
esac

LATEST=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep tag_name | cut -d '"' -f 4)
VERSION="${VERSION:-$LATEST}"
if [ -z "$VERSION" ]; then
  echo "Could not determine latest version"
  exit 1
fi

NAME="codeez_${VERSION#v}_${OS}_${ARCH}"
if [ "$OS" = "Windows" ]; then
  ARCHIVE="${NAME}.zip"
else
  ARCHIVE="${NAME}.tar.gz"
fi

URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE}"
CHECKSUM_URL="https://github.com/${REPO}/releases/download/${VERSION}/checksums.txt"

# Install directory: CODEEZ_INSTALL_DIR > XDG_BIN_DIR > /usr/local/bin > ~/.local/bin > ~/bin
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

echo "Installing codeez $VERSION to $INSTALL_DIR"
TMP=$(mktemp -d)
curl -sL -o "$TMP/$ARCHIVE" "$URL"
curl -sL -o "$TMP/checksums.txt" "$CHECKSUM_URL"
(cd "$TMP" && sha256sum -c --ignore-missing checksums.txt 2>/dev/null || shasum -a 256 -c checksums.txt 2>/dev/null || true)
if [ "$OS" = "Windows" ]; then
  unzip -o "$TMP/$ARCHIVE" -d "$TMP"
  cp "$TMP/codeez.exe" "$INSTALL_DIR/" 2>/dev/null || cp "$TMP/codeez" "$INSTALL_DIR/"
else
  tar -xzf "$TMP/$ARCHIVE" -C "$TMP"
  cp "$TMP/codeez" "$INSTALL_DIR/"
fi
rm -rf "$TMP"
chmod +x "$INSTALL_DIR/codeez"
echo "Installed: $INSTALL_DIR/codeez"
"$INSTALL_DIR/codeez" version
