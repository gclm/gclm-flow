#!/usr/bin/env python3
"""Merge hooks.json into settings.json.
Usage: python3 inject-hooks.py <hooks.json> <settings.json>
"""
import json
import sys
from pathlib import Path

src = Path(sys.argv[1])
dest = Path(sys.argv[2])

with open(src) as f:
    hooks = json.load(f)

settings: dict = {}
if dest.exists():
    with open(dest) as f:
        settings = json.load(f)

settings['hooks'] = hooks

with open(dest, 'w') as f:
    json.dump(settings, f, indent=2, ensure_ascii=False)
    f.write('\n')
