# Erlenwald Online (MMORPG-Prototyp)

Go-Server (`Server/`) + HTML/Three.js-Client (`Client/`, wird vom Server ausgeliefert).

## Starten
```
cd Server
go run ./cmd/server      (oder start.bat)
```
Dann im Browser http://localhost:8080 öffnen (mehrere Tabs/Browser = mehrere Spieler).
Mit `go run ./cmd/server -dev` gibt es Testbefehle im Chat: `/dev tp x z`, `/dev level n`, `/dev gold n`, `/dev unlock`, `/dev heal`.

## Starten mit Docker
Zwei Container: `server` (Go, WebSocket) und `client` (nginx mit den statischen Dateien, leitet `/ws` an den Server weiter).
```
docker compose up --build
```
Dann http://localhost:8080 öffnen. Accounts/Charaktere liegen im Volume `gamedata`.
Fertige Images (werden per GitHub Actions gebaut): `ghcr.io/herrpony/mmorpg-server` und `ghcr.io/herrpony/mmorpg-client` – `docker compose pull && docker compose up`.

## Welt
Zwei Kontinente (je ca. 600 × 600 m), vom Server als Höhen- und Biomkarte erzeugt (`gamedata/world.go`):
- **Erlenmark** (Stufe 1–20): Erlenwald, Goldene Ebenen, Düstermoor, Eisgipfel, Glutwüste. Hauptstadt Sonnenburg, Dörfer Eichenhain, Moorbruck, Frostwacht, Sandfurt, Hafen Sturmhafen. Dungeon: Schwarzfelsmine (ab Stufe 15).
- **Jadereich** (Stufe 20–40): Smaragddschungel, Kirschblütenhain, Rotkluft-Canyon, Aschenlande (mit Lava). Hauptstadt Himmelspforte (Pagoden-Stadt), Dörfer Lotushafen (Hafen), Bambushain, Kirschbach, Rotfels, Aschenwacht – jeweils mit Greifenmeister.
  Dungeons: Versunkener Nebeltempel (Stufe 26–30) und Feuerschlund (Stufe 36–40), je 2 Minibosse, Trash und Endboss für 5 Spieler.
- **Fähre:** Ein Schiff pendelt im festen Takt (45 s Liegezeit) zwischen Sturmhafen und Lotushafen. Es liegt längsseits am Steg:
  einfach hinüberlaufen, an Deck frei bewegen und im Zielhafen wieder über den Steg von Bord gehen.
- Orte lösen beim Betreten ein Namensbanner aus und werden auf der Weltkarte (M, Reiter pro Kontinent) aufgedeckt.
  Aktive Quests markieren ihr Zielgebiet golden auf Weltkarte und Minimap.
- Gegner spawnen nie auf Straßen. Maximalstufe 40, alle 2 Stufen eine neue Fertigkeit beim Klassenlehrer.
## Aufbau Server (Server/internal)
- `config` Flags · `protocol` Nachrichten · `storage` Store-Interface + JSON-Dateien (`Server/data`)
- `auth` Registrierung/Login · `gateway` WebSocket, Account-/Charakter-Verwaltung
- `gamedata` Welt, Orte, Klassen, Skills, Items/Ausrüstung, Mobs, NPCs, Quests, Dungeon (hier Inhalte erweitern)
- `quest` Questlogik · `loot` persönliche Beute · `progression` XP (Stufe 20)
- `game` Simulation (20 Hz): Kampf, Bedrohung, Elite/Boss-Fähigkeiten, Begleiter, Gruppen, Instanzen, Greifenflug, Zonen

## Steuerung
WASD bewegen · Leertaste springen · Rechtsklick halten: Kamera drehen · Linksklick: auswählen ·
Rechtsklick: NPC/Objekt ansprechen, Gegner angreifen, Leiche looten, Spieler → Gruppenmenü · Tab: Ziel wechseln ·
1–0, -, =, [, ]: Fertigkeiten · B Taschen, L Questlog, C Charakter/Ausrüstung, P Zauberbuch, **M Weltkarte**, O Gruppe ·
R Autorun · F mit Ziel interagieren · Enter Chat (`/p` Gruppenchat, `/einladen Name`, `/verlassen`) · N Musik an/aus · F11 Vollbild · Esc Fenster/Ziel/Menü

## Spielablauf
- Start nur mit Basisangriff. Fertigkeiten alle 2 Stufen beim Klassenlehrer (2–8, 12–20); Stufe 10: Klassenquest in Sonnenburg
  (Jäger: Tier zähmen/rufen · Krieger: Verteidigungshaltung + Spott · Magier: Wasserelementar).
- Quest-Arten: Töten, Sammeln, Weltobjekte, Botenquests mit Zeitlimit, Erkunden, Zähmen, Elite-/Gruppen-/Dungeonquests.
- Greifenmeister schalten Flugpunkte frei (Gespräch nötig); Flüge kosten Geld und gehen zu bekannten Städten.
- Ausrüstung ändert Werte und das Aussehen (auch bei anderen Spielern sichtbar); Händler kaufen/verkaufen.
- Gruppen bis 5 Spieler: geteilte XP und Questfortschritt, persönliche Beute, eigene Dungeon-Instanz pro Gruppe.

## Assets
Alle Modelle sind selbst erzeugte Pixel-Modelle ohne fremde Dateien: Figuren und Tiere bestehen aus Quadern, deren Flächen zur Laufzeit
mit kleinen Pixel-Texturen bemalt werden (`Client/js/world/human.js`, `beasts.js`, `creatures2.js`).
Spieler zeigen gewählte Haut-/Haarfarbe und Geschlecht sowie die angelegte Ausrüstung; Skelette, Zombies, Mumien, Yetis, Kobolde,
Elementare usw. haben eigene Stile und Proportionen. Welt, Musik und Effekte werden ebenfalls im Browser erzeugt.
## Mobile / Touch (Querformat)

Auf Smartphones und Tablets erkennt der Client automatisch Touch-Geräte und blendet eine virtuelle Steuerung ein:

- **Linke Bildschirmhälfte** (unten): schwebender Joystick – leicht kippen = gehen, voll = rennen
- **Rechte Seite**: Kamera drehen durch Wischen, Zwei-Finger-Pinch = Zoom, Tippen = Ziel wählen (NPCs/Objekte/Beute werden direkt benutzt, bei Gegnern/Spielern zweites Tippen = Angriff/Menü)
- **Buttons**: ⚔️ Angriff (Auto-Ziel), ⤒ Springen, 🎯 nächstes Ziel, ✋ Interagieren/Looten, 💬 Chat
- Aktionsleiste (2×7) und Menüleiste oben; lange auf einen Taschen-Slot drücken = Rechtsklick (benutzen/anlegen)
- Vollbild + Querformat-Sperre werden beim ersten Tippen angefordert; im Hochformat erscheint ein Hinweis zum Drehen.
  iPhone/Safari kennt keine Fullscreen-API: dort „Zum Home-Bildschirm“ hinzufügen (Web-App-Modus, Vollbild).

## Welt-Hindernisse (Bäume, Felsen, Gebäude)

Wo welche Props stehen, wird **nur an einer Stelle** festgelegt: `Server/internal/gamedata/props.go` (`generateProps`).
Die Liste wird beim Login an die Clients geschickt (`scenery.props`). Der Client zeichnet daraus die Vegetation und
baut sein Kollisionsgitter, der Server nutzt sie für Mob-Bewegung und Sichtlinie (`collide.go`).
Neue Prop-Arten: Namen in `PropKinds` ergänzen und das Modell in `Client/js/world/vegetation.js` (`parts()`) anlegen.
Gebäude kommen aus `towns.go` und gelten für Client und Server gleichermaßen.

## Admin-Modus & Map-Editor

Accounts auf der Admin-Liste dürfen die Welt bearbeiten: `-admins name1,name2` (oder Umgebungsvariable
`MMORPG_ADMINS`, im Docker-Setup in `docker-compose.yml`). **Den Account vorher normal registrieren**, sonst könnte
ihn jemand anderes anlegen. Admins sehen im Spiel oben links den Button 🛠; der Editor ist unter `/editor.html`
erreichbar (Anmeldung mit dem Spiel-Account). Jede Editor-Nachricht wird serverseitig erneut auf Admin-Rechte geprüft.

Im Editor (2D-Karte, Mausrad = Zoom, rechte Maustaste = schieben, Strg+Z/Y, Strg+S = speichern):

- **Spawns, Sammelobjekte, Gebäude** setzen, verschieben (ziehen/Pfeiltasten), löschen; Aggro-Radius einblendbar
- **Props** (Bäume, Felsen …) malen und wegradieren, **Gelände** formen (anheben/absenken/glätten/ebnen) und **Biome** malen
- **Quests, Mobs, NPCs**: Definitionen als JSON bearbeiten, duplizieren/neu anlegen (Tabs oben); NPCs lassen sich auf der Karte verschieben
- **Prüfen**: findet Spawns im Wasser, zu dicht stehende Spawns, Quests mit fehlenden Mobs/NPCs

„Speichern“ schreibt `world_edits.json` im Datenverzeichnis, wendet die Änderungen sofort auf den laufenden Server an
und lädt die Clients der Spieler neu. Die Datei liegt als Schicht über den Standarddaten im Code: „Auf Standard
zurücksetzen“ stellt die Code-Werte wieder her; Änderungen im Code an nicht bearbeiteten Dingen bleiben wirksam.
Nicht editierbar: Dungeons, Items, Fertigkeiten, Straßen/Flüsse/Orte.
