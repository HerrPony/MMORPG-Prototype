#!/bin/sh
dest="${ERLENWALD_DIR:-${XDG_DATA_HOME:-$HOME/.local/share}/ErlenwaldOnline}"
rm -rf "$dest" "$HOME/.local/bin/erlenwald" "${XDG_DATA_HOME:-$HOME/.local/share}/applications/erlenwald-online.desktop"
echo "Erlenwald Online wurde entfernt (Einstellungen in ~/.local/share/ErlenwaldOnline bleiben nicht erhalten)."
