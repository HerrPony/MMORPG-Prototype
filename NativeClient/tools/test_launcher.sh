#!/usr/bin/env bash
# End-to-end test of the updater: publishes fake releases into a temp tree, serves them over HTTP and lets the
# launcher (command line mode) install, update and repair a game folder.
#   GODOT=/path/to/godot NativeClient/tools/test_launcher.sh
set -eu
GODOT="${GODOT:-godot}"
here="$(cd "$(dirname "$0")/.." && pwd)"
tmp="$(mktemp -d)"
trap 'kill $srv 2>/dev/null || true; rm -rf "$tmp"' EXIT
port=$((20000 + RANDOM % 20000))
mkdir -p "$tmp/g1" "$tmp/upd"
head -c 3000000 /dev/urandom > "$tmp/g1/ErlenwaldOnline.x86_64"; chmod +x "$tmp/g1/ErlenwaldOnline.x86_64"
head -c 1000000 /dev/urandom > "$tmp/g1/ErlenwaldOnline.pck"
mkdir -p "$tmp/g1/music"; head -c 500000 /dev/urandom > "$tmp/g1/music/forest_day.ogg"
printf 'Version 1' > "$tmp/notes.txt"
pub() { python3 "$here/tools/release/make_release.py" --platform linux --version "$1" --game-dir "$2" --out "$tmp/upd" --notes "$tmp/notes.txt" "${@:3}"; }
(cd "$tmp/upd" && python3 -m http.server $port --bind 127.0.0.1 >/dev/null 2>&1) & srv=$!
pub 1.0.0 "$tmp/g1" >/dev/null
sleep 1
run() { "$GODOT" --headless --path "$here/launcher" -- --cli-update --url "http://127.0.0.1:$port" --root "$tmp/inst" "$@" 2>&1 | grep -E "RESULT|Lade|launcher update" || true; }
fail() { echo "FAIL: $*"; exit 1; }
echo "== first install"
out=$(run); echo "$out"; echo "$out" | grep -q "RESULT: ok version=1.0.0" || fail "install"
[ "$(sha256sum < "$tmp/inst/game/ErlenwaldOnline.pck")" = "$(sha256sum < "$tmp/g1/ErlenwaldOnline.pck")" ] || fail "pck differs"
[ -x "$tmp/inst/game/ErlenwaldOnline.x86_64" ] || fail "exec bit"
echo "== nothing to do"
out=$(run); echo "$out"; echo "$out" | grep -q "Lade" && fail "downloaded again"; echo "$out" | grep -q "RESULT: ok" || fail "second run"
echo "== update: only the pack changes"
cp -r "$tmp/g1" "$tmp/g2"; head -c 1200000 /dev/urandom > "$tmp/g2/ErlenwaldOnline.pck"
pub 1.0.1 "$tmp/g2" >/dev/null
out=$(run); echo "$out"; echo "$out" | grep -q "RESULT: ok version=1.0.1" || fail "update"
[ "$(echo "$out" | grep -c Lade)" = "1" ] || fail "more than the pack downloaded"
[ "$(sha256sum < "$tmp/inst/game/ErlenwaldOnline.pck")" = "$(sha256sum < "$tmp/g2/ErlenwaldOnline.pck")" ] || fail "pck not updated"
echo "== removed file disappears"
rm "$tmp/g2/music/forest_day.ogg"; pub 1.0.2 "$tmp/g2" >/dev/null
out=$(run); echo "$out" | grep -q "RESULT: ok version=1.0.2" || fail "update 3"
[ ! -e "$tmp/inst/game/music/forest_day.ogg" ] || fail "stale file kept"
echo "== repair a damaged file"
echo broken > "$tmp/inst/game/ErlenwaldOnline.pck"
out=$(run --repair); echo "$out"; echo "$out" | grep -q "RESULT: ok" || fail "repair"
[ "$(sha256sum < "$tmp/inst/game/ErlenwaldOnline.pck")" = "$(sha256sum < "$tmp/g2/ErlenwaldOnline.pck")" ] || fail "pck not repaired"
echo "== corrupted download is rejected"
python3 - "$tmp/upd" <<'PY'
import json, sys, os
d = sys.argv[1]
m = json.load(open(os.path.join(d, "manifest-linux.json")))
e = [f for f in m["files"] if f["path"].endswith(".pck")][0]
h = e["sha256"]; p = os.path.join(d, "objects", h[:2], h)
open(p, "wb").write(b"x" * e["size"])
PY
rm "$tmp/inst/game/ErlenwaldOnline.pck"
out=$(run); echo "$out"; echo "$out" | grep -q "RESULT: error" || fail "bad hash accepted"
echo "== launcher updates itself (decided by content)"
mkdir -p "$tmp/launcher"; echo "new launcher pck" > "$tmp/launcher/Launcher.pck"; head -c 100000 /dev/urandom > "$tmp/launcher/Launcher.x86_64"
head -c 1200000 /dev/urandom > "$tmp/g2/ErlenwaldOnline.pck"
pub 1.0.3 "$tmp/g2" --launcher-dir "$tmp/launcher" --launcher-version 9.9.9 >/dev/null
echo "old launcher pck" > "$tmp/inst/Launcher.pck"
out=$(run --self-update); echo "$out"; echo "$out" | grep -q "RESULT: launcher update staged" || fail "launcher not staged"
for i in 1 2 3 4 5 6 7 8; do [ -e "$tmp/inst/update_launcher.sh" ] || break; sleep 1; done
[ "$(cat "$tmp/inst/Launcher.pck")" = "new launcher pck" ] || fail "launcher pack not swapped in"
[ -x "$tmp/inst/Launcher.x86_64" ] || fail "launcher binary not installed"
[ ! -d "$tmp/inst/.update" ] || fail "staging folder left behind"
out=$(run --self-update); echo "$out" | grep -q "launcher update available" && fail "a current launcher wants to update again (restart loop)"
echo "launcher is current: no further update"
echo "== a release with the launcher version of the running launcher never replaces it"
lv="$(grep -oP 'config/version="\K[^"]+' "$here/launcher/project.godot")"
echo "launcher files that differ" > "$tmp/launcher/Launcher.pck"
head -c 1300000 /dev/urandom > "$tmp/g2/ErlenwaldOnline.pck"
pub 1.0.4 "$tmp/g2" --launcher-dir "$tmp/launcher" --launcher-version "$lv" >/dev/null
out=$(run --self-update); echo "$out"; echo "$out" | grep -q "launcher update available" && fail "a game-only release replaced the launcher"
echo "game-only release: launcher untouched"
echo "ALL LAUNCHER TESTS PASSED"
