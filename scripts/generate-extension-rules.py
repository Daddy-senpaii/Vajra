#!/usr/bin/env python3
"""Generate a Chrome DNR ruleset bounded by its guaranteed static rule budget."""

import hashlib
import json
import pathlib
import sys
import urllib.parse

ROOT = pathlib.Path(__file__).resolve().parent.parent
CORE_LIST = ROOT / "packages/blocking-engine/blocklist/data/adult-domains.txt"
PRIORITY_LIST = ROOT / "apps/extension/data/default-priority-domains.txt"
OUTPUT = ROOT / "apps/extension/public/rules.json"
RULE_LIMIT = 30_000


def read_domains(path: pathlib.Path) -> list[str]:
    return [
        line.strip().lower()
        for line in path.read_text(encoding="utf-8").splitlines()
        if line.strip() and not line.lstrip().startswith("#")
    ]


def main() -> int:
    domains = set(read_domains(CORE_LIST))
    priority = list(dict.fromkeys(read_domains(PRIORITY_LIST)))
    missing_priority = [domain for domain in priority if domain not in domains]
    if missing_priority:
        print(
            "Priority domains missing from core data: " + ", ".join(missing_priority),
            file=sys.stderr,
        )

    selected = [domain for domain in priority if domain in domains]
    selected_set = set(selected)
    remainder = domains - selected_set
    # A deterministic hash sample spreads the static rules across the complete
    # alphabetically sorted upstream set instead of covering only its prefix.
    remaining_count = RULE_LIMIT - len(selected)
    selected.extend(
        sorted(remainder, key=lambda item: hashlib.sha256(item.encode()).digest())[
            :remaining_count
        ]
    )

    rules = [
        {
            "id": rule_id,
            "priority": 1,
            "action": {
                "type": "redirect",
                "redirect": {
                    "extensionPath": "/blocked.html?domain="
                    + urllib.parse.quote(domain, safe=".-")
                },
            },
            "condition": {
                "urlFilter": f"||{domain}^",
                "resourceTypes": ["main_frame", "sub_frame"],
            },
        }
        for rule_id, domain in enumerate(selected, start=1)
    ]

    OUTPUT.write_text(json.dumps(rules, separators=(",", ":")) + "\n", encoding="utf-8")
    print(f"Generated {len(rules):,} browser rules from {len(domains):,} core domains")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
