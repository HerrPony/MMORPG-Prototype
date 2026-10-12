#!/usr/bin/env python3
"""Shrinks the bundled emoji font to the glyphs the game actually uses (item/skill icons, UI symbols).

    pip install fonttools
    python3 NativeClient/tools/fonts/subset_emoji.py <full NotoColorEmoji.ttf> NativeClient/game/assets/fonts/NotoColorEmoji.ttf

The full font is ~10 MB, the subset ~1 MB; since the game pack is re-downloaded on every update this matters.
Re-run it when new icons are added to the server data. Source of the full font: Noto Color Emoji (SIL OFL).
"""
import glob
import os
import sys

from fontTools import subset
from fontTools.ttLib import TTFont

root = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", "..", ".."))
src, dst = sys.argv[1], sys.argv[2]
files = glob.glob(os.path.join(root, "Server", "**", "*.go"), recursive=True) \
    + glob.glob(os.path.join(root, "NativeClient", "game", "scripts", "**", "*.gd"), recursive=True) \
    + glob.glob(os.path.join(root, "Client", "js", "**", "*.js"), recursive=True)
cmap = TTFont(src).getBestCmap()
used = {0x200D, 0xFE0F, 0x20E3}
for f in files:
    try:
        text = open(f, encoding="utf-8").read()
    except (UnicodeDecodeError, OSError):
        continue
    for ch in text:
        if ord(ch) >= 0x203C and ord(ch) in cmap:
            used.add(ord(ch))
# keycap digits, flags and similar sequences need their base characters
used |= {ord(c) for c in "0123456789#*"} & set(cmap)
opts = subset.Options()
opts.layout_features = ["*"]
opts.hinting = False
opts.notdef_outline = True
opts.glyph_names = False
opts.name_IDs = [1, 2, 3, 4, 6]
font = TTFont(src)
sub = subset.Subsetter(opts)
sub.populate(unicodes=sorted(used))
sub.subset(font)
font.save(dst)
print("%d code points, %.2f MB -> %.2f MB" % (len(used), os.path.getsize(src) / 1e6, os.path.getsize(dst) / 1e6))
