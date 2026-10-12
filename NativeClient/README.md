# Erlenwald Online – Native Desktop-Client

Der native Client ist eine vollwertige Alternative zum Web-Client (`../Client`). Er spricht dasselbe WebSocket-/JSON-Protokoll
mit demselben Go-Server (`../Server`) – **die Serverlogik ist unverändert**. Alle Funktionen des Web-Clients sind vorhanden,
dazu kommen bessere Grafik, mehr Einstellungen und eine flüssigere Steuerung.

* **Engine:** Godot 4.5 (Forward+, Vulkan/D3D12), Sprache GDScript
* **Optik:** Low-Poly im Minecraft-Stil (Blockmodelle, flache Schattierung), aber mit Echtzeitschatten, SSAO, Bloom, Nebel,
  Himmels-Shader mit Sonne/Mond/Sternen/Wolken/Aurora, Wasser-Shader mit Tiefenfärbung und Schaum, Tag-/Nachtzyklus (6 h)
* **Launcher + Installer:** findet Updates selbst, lädt nur geänderte Dateien, repariert beschädigte Installationen

```
NativeClient/
  game/        Godot-Projekt des Spiels (Client)
  launcher/    Godot-Projekt des Launchers (Updates, Start)
  installer/   NSIS-Skript (Windows) und Installationsskript (Linux)
  tools/       Musik-Renderer, Release-Werkzeuge, Tests
  build/       Ausgaben (nicht im Repository)
```

## Spielen

1. Server starten (siehe Haupt-README), z. B. `cd Server && go run ./cmd/server`
2. Client starten – entweder über den Installer/Launcher (siehe unten) oder direkt mit Godot:
   `godot --path NativeClient/game` (Godot 4.5.1 genügt, keine weiteren Abhängigkeiten)
3. Serveradresse eintragen (Standard `localhost:8080`), Account anlegen, Charakter erstellen, spielen.

### Steuerung (alles in den Einstellungen umbelegbar)

| Taste | Aktion |
| --- | --- |
| `W A S D` | Bewegen · `Leertaste` springen · `Shift` gehen · `R` automatisch laufen |
| Rechte Maustaste halten | Kamera drehen, Charakter blickt mit · beide Tasten: vorwärts |
| Linksklick / Rechtsklick | auswählen / ansprechen, angreifen, looten · Spieler: Menü (Gruppe) |
| `Tab` | nächster Gegner · `F` mit Ziel interagieren · `Mausrad` Zoom |
| `1`–`0`, `-`, `=`, `[`, `]` | Aktionsleiste |
| `B L C P M O` | Taschen, Questlog, Charakter, Zauberbuch, Weltkarte, Gruppe |
| `Enter` / `/` | Chat / Befehl · `Esc` Fenster schließen, Ziel abwählen, Menü |
| `F11` / `F12` | Vollbild / Screenshot (`user://screenshots`) |
| Gamepad | linker Stick laufen, rechter Stick Kamera, A springen, X interagieren, B/Y/LB und Steuerkreuz: Fertigkeiten |

### Einstellungen

Grafik-Voreinstellungen (Niedrig bis Ultra) und einzeln: Fenstermodus, V-Sync, FPS-Limit, Renderauflösung, FSR, MSAA/FXAA/TAA,
Schattenqualität und -reichweite, Sichtweite, Vegetationsdichte, Wasser, Partikel, SSAO, SSIL, Bloom, Nebel, volumetrischer Nebel,
Tonemapping, Helligkeit/Kontrast/Sättigung, Sichtfeld, Oberflächengröße, Namensschilder, Schadenszahlen, FPS-Anzeige.
Audio (Gesamt/Musik/Effekte/Umgebung), Maus (Empfindlichkeit, Invertierung, weiche Kamera) und alle Tasten.
Gespeichert wird in `settings.cfg` im Benutzerordner (`%APPDATA%\ErlenwaldOnline\game` bzw. `~/.local/share/ErlenwaldOnline/game`).

## Aufbau des Spiels (`game/`)

| Pfad | Inhalt |
| --- | --- |
| `scripts/autoload/` | `S` (Zustand), `Net` (WebSocket), `Settings`, `Sound` |
| `scripts/world/` | Weltaufbau aus den Serverdaten: Gelände (Vertex-Shader aus Höhenmap), Wasser, Himmel, Props (MultiMesh), Gebäude, Modelle mit prozeduraler Animation, Entities, Steuerung, Effekte, Dungeons |
| `scripts/ui/` | UI-Theme, HUD, Fenster (Taschen, Charakter, Zauberbuch, Quests, Dialog, Händler, Beute, Gruppe, Karte, Menü, Einstellungen), Login/Charakterwahl |
| `assets/shaders/` | Gelände, Wasser, Himmel, Laub, Karte |
| `assets/sfx/` | Soundeffekte (OGG, aus der WebAudio-Engine gerendert) |

Weltkoordinaten, Kollision (Bäume, Gebäude, Hänge, Wasser), Platzierung der Props und alle Regeln sind identisch zum Web-Client:
Der Client baut die Welt beim Start aus der `world`-Nachricht des Servers (Höhen, Biome, Props, Straßen, Gebäude …) –
Änderungen im Karten-Editor des Servers wirken also auch hier ohne zweite Pflege.

## Musik

Die Musik des Web-Clients (`Client/js/audio.js`, prozedurale orchestrale Themen je Biome, Tag/Nacht, Kampf) wird **offline gerendert**:

```
node NativeClient/tools/music/render_music.mjs NativeClient/build/audio       # Musik + Effekte (Chromium, ffmpeg, playwright)
```

Das erzeugt 33 OGG-Stücke (`music/<biom>_day|_night|_fight.ogg`, `menu`, `dungeon`, `dungeon_fight`, zusammen ca. 17 MB) und die Effekte.
Die Musik liegt neben dem Spiel (`game/music/`, vom Launcher installiert) und wird beim Wechsel von Region, Tageszeit und Kampf überblendet.
Fehlen die Dateien, läuft das Spiel einfach ohne Musik. Die Effekte (`game/assets/sfx`) sind im Repository enthalten.

## Launcher, Installer, Updates

```
Erlenwald Online/              (Windows: %LOCALAPPDATA%\Programs\Erlenwald Online)
  Launcher.exe  Launcher.pck   der Launcher (Godot-Engine + kleines Skript-Paket)
  launcher_defaults.json       Update-Adresse und Spielserver (vom Installer geschrieben)
  game/                        vom Launcher verwaltet
    ErlenwaldOnline.exe        Engine
    ErlenwaldOnline.pck        Spielinhalt (ca. 11 MB, ändert sich bei jedem Update)
    music/*.ogg                Musik
    installed.json             installierte Version und Prüfsummen
```

Der Launcher fragt beim Start (und alle 10 Minuten) `manifest-<plattform>.json` vom Update-Server ab, vergleicht die SHA-256-Prüfsummen mit
den installierten Dateien, lädt **nur Geändertes** (meist nur das `.pck`), prüft es und ersetzt die Dateien. Er aktualisiert sich bei
Bedarf selbst (Skript tauscht Dateien nach dem Beenden). Ist die Engine-Datei des Spiels identisch zu der des Launchers, wird sie kopiert
statt heruntergeladen. „Prüfen“ hasht alle Dateien und repariert beschädigte. Beim Start übergibt der Launcher den Spielserver
(`-- --server host:port`). Ohne Netz startet er das installierte Spiel trotzdem.

### Installer

* **Windows:** `ErlenwaldOnline-Setup.exe` (Web-Installer, ~24 MB, enthält nur den Launcher) und `ErlenwaldOnline-Setup-Full.exe`
  (mit Spiel, ohne Download). Installation pro Benutzer ohne Administratorrechte, Startmenü, Desktop-Verknüpfung, Deinstallation.
  Bauen: `VERSION=1.0.0 [BUNDLE_GAME=1] [UPDATE_URL=…] [GAME_SERVER=…] NativeClient/tools/build_installer.sh` (benötigt `makensis`).
* **Linux:** `ErlenwaldOnline-linux.tar.gz` entpacken, `./install.sh` ausführen (Menüeintrag und Befehl `erlenwald`), `uninstall.sh` entfernt alles.
* macOS ist noch nicht vorbereitet (nur Windows und Linux).
* Die Programme sind nicht signiert: Windows SmartScreen zeigt beim ersten Start „Weitere Informationen → Trotzdem ausführen“.
  Mit einem Zertifikat lässt sich der Installer und die Engine-Datei signieren (Godot-Exportoption `codesign`, NSIS `signtool`).

### Eine neue Version veröffentlichen

Automatisch: Tag `native-v1.2.3` pushen oder den Workflow **Native client release** manuell starten
(`.github/workflows/native-client.yml`). Er exportiert Spiel und Launcher, rendert die Musik, testet den Updater, schreibt den Update-Baum
in den Branch `gh-pages` (Ordner `updates/`, Standard-Adresse des Launchers: `https://herrpony.github.io/MMORPG-Prototype/updates`) und hängt die
Installer an ein GitHub-Release.

Von Hand:

```
godot --headless --path NativeClient/game     --export-release "Windows Desktop" ../build/windows/game/ErlenwaldOnline.exe
godot --headless --path NativeClient/launcher --export-release "Windows Desktop" ../build/windows/launcher/Launcher.exe
cp -r NativeClient/build/audio/music NativeClient/build/windows/game/music
python3 NativeClient/tools/release/make_release.py --platform windows --version 1.2.0 \
    --game-dir NativeClient/build/windows/game --launcher-dir NativeClient/build/windows/launcher --launcher-version 1.2.0 \
    --notes notes.txt --out updates --prune
```

`updates/` ist ein statischer Ordner (`manifest-*.json` und `objects/<hash>`), den jeder Webserver ausliefern kann. Der Go-Server kann das
selbst: `server -updates ./updates` (oder `MMORPG_UPDATES=./updates`) stellt ihn unter `/updates/` bereit – dann im Launcher
(Einstellungen) bzw. per Installer-Option `UPDATE_URL=http://mein-server:8080/updates` eintragen.

## Tests

```
GODOT=godot NativeClient/tools/check_scripts.sh         # alle Skripte kompilieren
GODOT=godot NativeClient/tools/test_launcher.sh         # Updater: installieren, Delta-Update, Reparatur, kaputter Download, Selbst-Update
```

Spieltests laufen headless gegen einen Testserver (`server -dev -client "" -data /tmp/x -addr :8084`):

```
godot --headless --path NativeClient/game -- --autologin localhost:8084 user pass --autoenter --play
```

(meldet sich an, läuft, wählt einen Gegner an, kämpft und prüft Server-Synchronisation, Schaden und Musikwechsel).
`--play2` testet Händler, Ausrüstung, Quests, Trainer, Karten-/Gruppenfenster und den Greifenflug. Mit
`--shot out.png sekunden [--open bags,char] [--chat "/dev level 20;/dev tp 6 12"] [--hour 12] [--interact npc]` entstehen Screenshots.
Lizenzen: Godot (MIT), Noto Color Emoji (OFL, `game/assets/fonts`).
