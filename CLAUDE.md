# Arbeitsregeln

## Native-Client-Releases
Bei jedem Versions-Update des Native Clients (z.B. `native-v1.0.5`):
1. Die Version im Code erhöhen (aktuell in `NativeClient/installer/windows/installer.nsi`, `NativeClient/game/{project.godot,export_presets.cfg}`, `NativeClient/launcher/{project.godot,export_presets.cfg,scripts/launcher.gd}`).
2. Committen und pushen.
3. Den Tag `native-vX.Y.Z` setzen und pushen, damit die Action `.github/workflows/native-client.yml` läuft.
