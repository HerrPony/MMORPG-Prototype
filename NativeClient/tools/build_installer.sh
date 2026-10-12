#!/usr/bin/env bash
# Builds the Windows installer from the exported launcher (and optionally the exported game).
#   VERSION=1.0.0 [BUNDLE_GAME=1] [UPDATE_URL=...] [GAME_SERVER=...] NativeClient/tools/build_installer.sh
set -eu
here="$(cd "$(dirname "$0")/.." && pwd)"
VERSION="${VERSION:-1.0.0}"
out="$here/build/dist"
mkdir -p "$out"
file="$out/ErlenwaldOnline-Setup.exe"
args=(-DVERSION="$VERSION" -DLAUNCHER_DIR="$here/build/windows/launcher")
[ -n "${UPDATE_URL:-}" ] && args+=(-DUPDATE_URL="$UPDATE_URL")
[ -n "${GAME_SERVER:-}" ] && args+=(-DGAME_SERVER="$GAME_SERVER")
if [ "${BUNDLE_GAME:-0}" = 1 ]; then
  args+=(-DGAME_DIR="$here/build/windows/game")
  file="$out/ErlenwaldOnline-Setup-Full.exe"
fi
makensis -V2 "${args[@]}" -DOUTFILE="$file" "$here/installer/windows/installer.nsi"
ls -la "$out"
