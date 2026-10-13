# Arbeitsregeln

## Native-Client-Releases
Bei jedem Versions-Update des Native Clients (z.B. `native-v1.0.6`):
1. Die Version im Code erhöhen: Spiel in `NativeClient/game/{project.godot,export_presets.cfg}` und `NativeClient/installer/windows/installer.nsi`. **Den Launcher** (`NativeClient/launcher/{project.godot,export_presets.cfg,scripts/launcher.gd}`) **nur dann mit hochzählen, wenn sich der Launcher selbst ändert.** Ein reines Client-Update darf keine Launcher-Version ändern, sonst lädt sich der Launcher bei jedem Update selbst neu herunter.
2. Den Launcher-Changelog in `NativeClient/release-notes/X.Y.Z.txt` schreiben (siehe unten).
3. Committen und pushen.
4. Den Tag `native-vX.Y.Z` setzen und pushen, damit die Action `.github/workflows/native-client.yml` läuft.

## Launcher-Changelog
Zu jedem Native-Client-Release schreibe ich einen vollständigen Changelog für das Neuigkeiten-Panel des Launchers (BBCode, deutsch), im bisherigen Format:
```
[b]Version X.Y.Z[/b]
[b]Behoben[/b]            (oder [b]Neu[/b], [b]Geändert[/b] ...)
[list]
[*]Nutzerverständliche Beschreibung der Änderung.
[/list]
[b]Technik[/b]
[list]
[*]Technische Details.
[/list]
```
Der Text steht in `NativeClient/release-notes/X.Y.Z.txt` (wird mit der Versionserhöhung committet). Der Workflow liest ihn beim Tag-Push und übernimmt ihn als `notes` ins Update-Manifest; beim manuellen Start (`workflow_dispatch`) hat die Eingabe `notes` Vorrang. Ohne Datei und ohne Eingabe bleibt nur `Version X.Y.Z`. Den Changelog zeige ich dem Nutzer zusätzlich im Chat.
