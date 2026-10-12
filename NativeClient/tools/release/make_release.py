#!/usr/bin/env python3
"""Publishes a build into an update tree that the launcher understands.

    make_release.py --platform windows --version 1.2.0 --game-dir build/windows/game --out updates
                    [--launcher-dir build/windows/launcher --launcher-version 1.0.0] [--notes notes.txt]

The tree is plain static files (any web server, GitHub Pages, the Go server's /updates/ endpoint):
    manifest-<platform>.json          what a platform needs: version, notes, files with SHA-256
    objects/<h2>/<sha256>             every file, stored under its hash (identical files exist once, old versions stay)
"""
import argparse
import hashlib
import json
import os
import shutil
import sys
import time


def sha256_of(path):
    h = hashlib.sha256()
    with open(path, "rb") as f:
        for chunk in iter(lambda: f.read(1 << 20), b""):
            h.update(chunk)
    return h.hexdigest()


def collect(directory, skip=()):
    files = []
    for dp, dn, fn in os.walk(directory):
        dn.sort()
        for name in sorted(fn):
            full = os.path.join(dp, name)
            rel = os.path.relpath(full, directory).replace(os.sep, "/")
            if rel in skip or rel.endswith(".old") or name in ("installed.json",):
                continue
            files.append((rel, full))
    return files


def entry(rel, full, out):
    digest = sha256_of(full)
    dest = os.path.join(out, "objects", digest[:2], digest)
    if not os.path.exists(dest):
        os.makedirs(os.path.dirname(dest), exist_ok=True)
        shutil.copyfile(full, dest)
    exe = rel.endswith(".x86_64") or rel.endswith(".sh") or ("." not in os.path.basename(rel))
    e = {"path": rel, "sha256": digest, "size": os.path.getsize(full)}
    if exe:
        e["exec"] = True
    return e


def prune(out):
    """Deletes objects that no manifest in the tree refers to anymore."""
    used = set()
    for name in os.listdir(out):
        if name.startswith("manifest-") and name.endswith(".json"):
            with open(os.path.join(out, name), encoding="utf-8") as f:
                m = json.load(f)
            for e in m.get("files", []) + m.get("launcher_files", []):
                used.add(e["sha256"])
    removed = 0
    obj = os.path.join(out, "objects")
    for dp, _, fn in os.walk(obj):
        for name in fn:
            if name not in used:
                os.remove(os.path.join(dp, name))
                removed += 1
    return removed


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--platform", required=True, choices=["windows", "linux", "macos"])
    ap.add_argument("--version", required=True)
    ap.add_argument("--game-dir", required=True, help="exported game (engine binary, .pck, music/ ...)")
    ap.add_argument("--out", required=True)
    ap.add_argument("--launch", help="executable inside game-dir (default: ErlenwaldOnline.exe / .x86_64)")
    ap.add_argument("--launcher-dir", help="folder with the launcher files (Launcher.exe + Launcher.pck)")
    ap.add_argument("--launcher-version", default="1.0.0")
    ap.add_argument("--notes", help="text/BBCode file shown in the launcher")
    ap.add_argument("--prune", action="store_true", help="remove unreferenced objects afterwards")
    a = ap.parse_args()

    os.makedirs(a.out, exist_ok=True)
    launch = a.launch or ("ErlenwaldOnline.exe" if a.platform == "windows" else "ErlenwaldOnline.x86_64")
    if not os.path.exists(os.path.join(a.game_dir, launch)):
        sys.exit("launch file %s not found in %s" % (launch, a.game_dir))
    notes = ""
    if a.notes:
        with open(a.notes, encoding="utf-8") as f:
            notes = f.read()
    files = [entry(rel, full, a.out) for rel, full in collect(a.game_dir)]
    manifest = {
        "format": 1, "product": "erlenwald", "platform": a.platform, "version": a.version,
        "released": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()), "launch": launch,
        "objects": "objects/{h2}/{h}", "notes": notes, "files": files,
    }
    if a.launcher_dir:
        manifest["launcher_version"] = a.launcher_version
        manifest["launcher_files"] = [entry(rel, full, a.out) for rel, full in collect(a.launcher_dir)]
    path = os.path.join(a.out, "manifest-%s.json" % a.platform)
    tmp = path + ".tmp"
    with open(tmp, "w", encoding="utf-8") as f:
        json.dump(manifest, f, indent=1, ensure_ascii=False)
    os.replace(tmp, path)  # the manifest appears last, so clients never see a half published release
    total = sum(e["size"] for e in files)
    print("published %s %s: %d files, %.1f MB -> %s" % (a.platform, a.version, len(files), total / 1e6, path))
    if a.prune:
        print("pruned %d old objects" % prune(a.out))


if __name__ == "__main__":
    main()
