#!/usr/bin/env bash
set -eux

CACHE_DIR="$1"; shift
mkdir -p "$CACHE_DIR"
cd "$CACHE_DIR"

VERSION='v0.11.0'
GLIDE="glide-${VERSION}-linux-amd64.tar.gz"

wget -nc "https://github.com/Masterminds/glide/releases/download/$VERSION/$GLIDE"
mkdir -p glide
tar -xzf "$GLIDE" -C ./glide
sudo cp ./glide/linux-amd64/glide /usr/local/bin/glide
