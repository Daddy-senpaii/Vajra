#!/bin/sh
set -eu

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
data_dir="$repo_root/packages/blocking-engine/blocklist/data"
metadata_url="https://raw.githubusercontent.com/Bon-Appetit/porn-domains/main/meta.json"
metadata_tmp=$(mktemp)
list_tmp=$(mktemp "$data_dir/.adult-domains.XXXXXX")
metadata_out_tmp=$(mktemp "$data_dir/.metadata.XXXXXX")
trap 'rm -f "$metadata_tmp" "$list_tmp" "$metadata_out_tmp"' EXIT HUP INT TERM

curl --fail --silent --show-error --location --retry 2 "$metadata_url" -o "$metadata_tmp"
list_url=$(python3 -c 'import json,sys; print(json.load(open(sys.argv[1]))["blocklist"]["raw_url"])' "$metadata_tmp")
curl --fail --silent --show-error --location --retry 2 "$list_url" -o "$list_tmp"

python3 - "$metadata_tmp" "$list_tmp" "$metadata_out_tmp" <<'PY'
import hashlib
import json
import pathlib
import sys

upstream = json.loads(pathlib.Path(sys.argv[1]).read_text())['blocklist']
content = pathlib.Path(sys.argv[2]).read_bytes()
lines = len(content.splitlines())
if len(content) != upstream['size'] or lines != upstream['lines']:
    raise SystemExit(
        f"download validation failed: got {lines} lines/{len(content)} bytes, "
        f"expected {upstream['lines']} lines/{upstream['size']} bytes"
    )

metadata = {
    'source': 'Bon-Appetit/porn-domains',
    'source_url': upstream['raw_url'],
    'updated': upstream['updated'],
    'lines': lines,
    'size': len(content),
    'sha256': hashlib.sha256(content).hexdigest(),
    'license': 'CC BY-SA 4.0',
}
pathlib.Path(sys.argv[3]).write_text(json.dumps(metadata, indent=2) + '\n')
PY

mv "$list_tmp" "$data_dir/adult-domains.txt"
mv "$metadata_out_tmp" "$data_dir/metadata.json"
python3 "$repo_root/scripts/generate-extension-rules.py"
printf 'Updated adult-domain snapshot from %s (%s lines)\n' "$list_url" "$(python3 -c 'import json,sys; print(json.load(open(sys.argv[1]))["blocklist"]["lines"])' "$metadata_tmp")"
