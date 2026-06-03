#!/bin/bash

# SecretShift Linux Installation Script
# Installs or updates secret-shift to the latest version

set -e

REPO="PapaDanielVi/secret-shift"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64) SUFFIX="x86_64" ;;
    aarch64|arm64) SUFFIX="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Find package manager
if command -v dnf &>/dev/null; then
    PKG_MANAGER="rpm"
elif command -v yum &>/dev/null; then
    PKG_MANAGER="rpm"
elif command -v pacman &>/dev/null; then
    PKG_MANAGER="arch"
elif command -v zypper &>/dev/null; then
    PKG_MANAGER="rpm"
elif command -v apt-get &>/dev/null; then
    PKG_MANAGER="deb"
elif command -v apk &>/dev/null; then
    PKG_MANAGER="apk"
else
    echo "No supported package manager found. Falling back to tarball."
    PKG_MANAGER="tar"
fi

# Get latest release version
LATEST_VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep -oP '"tag_name": "\K[^"]+')
if [ -z "$LATEST_VERSION" ]; then
    echo "Failed to fetch latest release version"
    exit 1
fi

echo "Installing secret-shift ${LATEST_VERSION}..."

# Remove existing installation if present
if command -v secret-shift &>/dev/null; then
    echo "Removing existing installation..."
    SECRET_SHIFT_BIN=$(which secret-shift)
    case "$PKG_MANAGER" in
        rpm) sudo rpm -e secret-shift 2>/dev/null || true ;;
        deb) sudo dpkg -r secret-shift 2>/dev/null || true ;;
        apk) sudo apk del secret-shift 2>/dev/null || true ;;
        arch) sudo pacman -R secret-shift 2>/dev/null || true ;;
        tar) sudo rm -f "$SECRET_SHIFT_BIN" ;;
    esac
fi

# Download and install based on package manager
case "$PKG_MANAGER" in
    rpm)
        URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/secret-shift_${LATEST_VERSION}_linux_${SUFFIX}.rpm"
        echo "Downloading RPM package..."
        curl -Lo /tmp/secret-shift.rpm "$URL"
        sudo rpm -i /tmp/secret-shift.rpm
        rm /tmp/secret-shift.rpm
        ;;
    deb)
        URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/secret-shift_${LATEST_VERSION}_linux_${SUFFIX}.deb"
        echo "Downloading DEB package..."
        curl -Lo /tmp/secret-shift.deb "$URL"
        sudo dpkg -i /tmp/secret-shift.deb || sudo apt-get -f install -y
        rm /tmp/secret-shift.deb
        ;;
    apk)
        URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/secret-shift_${LATEST_VERSION}_linux_${SUFFIX}.apk"
        echo "Downloading APK package..."
        curl -Lo /tmp/secret-shift.apk "$URL"
        sudo apk add --allow-untrusted /tmp/secret-shift.apk
        rm /tmp/secret-shift.apk
        ;;
    arch)
        URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/secret-shift_${LATEST_VERSION}_linux_${SUFFIX}.pkg.tar.zst"
        echo "Downloading Arch package..."
        curl -Lo /tmp/secret-shift.pkg.tar.zst "$URL"
        sudo pacman -U /tmp/secret-shift.pkg.tar.zst --noconfirm
        rm /tmp/secret-shift.pkg.tar.zst
        ;;
    tar)
        URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/secret-shift_Linux_${SUFFIX}.tar.gz"
        echo "Downloading tarball..."
        curl -Lo /tmp/secret-shift.tar.gz "$URL"
        sudo tar -xzf /tmp/secret-shift.tar.gz -C "$INSTALL_DIR"
        rm /tmp/secret-shift.tar.gz
        ;;
esac

echo "secret-shift ${LATEST_VERSION} installed successfully!"
secret-shift version