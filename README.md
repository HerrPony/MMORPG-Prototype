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