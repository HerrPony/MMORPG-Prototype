#!/usr/bin/env bash
# Compiles every GDScript of the game project (godot binary from $GODOT or PATH) and prints the errors.
set -u
GODOT="${GODOT:-godot}"
cd "$(dirname "$0")/../game"
"$GODOT" --headless --path . --import >/dev/null 2>&1
out=$("$GODOT" --headless --path . res://tests/load_all.tscn 2>&1 | grep -E "SCRIPT ERROR|Parse Error|Compile Error|cannot load|load_all:" -A1 | grep -v "^--")
echo "$out"
echo "$out" | grep -q "load_all: ok"
