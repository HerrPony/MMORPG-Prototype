#!/bin/sh
# Installs the Erlenwald launcher for the current user (no root needed) and adds a menu entry.
# The launcher downloads and updates the game itself on its first start.
set -e
here="$(cd "$(dirname "$0")" && pwd)"
dest="${ERLENWALD_DIR:-${XDG_DATA_HOME:-$HOME/.local/share}/ErlenwaldOnline}"
mkdir -p "$dest" "$HOME/.local/bin" "${XDG_DATA_HOME:-$HOME/.local/share}/applications"
cp "$here/Launcher.x86_64" "$here/Launcher.pck" "$dest/"
[ -f "$here/launcher_defaults.json" ] && cp "$here/launcher_defaults.json" "$dest/"
[ -d "$here/game" ] && cp -r "$here/game" "$dest/"
chmod +x "$dest/Launcher.x86_64"
cat > "$HOME/.local/bin/erlenwald" <<SH
#!/bin/sh
cd "$dest" && exec ./Launcher.x86_64 "\$@"
SH
chmod +x "$HOME/.local/bin/erlenwald"
cat > "${XDG_DATA_HOME:-$HOME/.local/share}/applications/erlenwald-online.desktop" <<DESK
[Desktop Entry]
Type=Application
Name=Erlenwald Online
Comment=Low-Poly-MMORPG
Exec=$dest/Launcher.x86_64
Path=$dest
Terminal=false
Categories=Game;RolePlaying;
DESK
echo "Erlenwald Online wurde nach $dest installiert. Start: 'erlenwald' oder über das Anwendungsmenü."
