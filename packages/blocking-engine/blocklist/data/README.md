# Bundled adult-domain data

`adult-domains.txt` is a snapshot of the Bon-Appetit `porn-domains` generated
blocklist. Snapshot timestamp, line count, size, SHA-256, and pinned source URL
are recorded in [`metadata.json`](metadata.json).

- Source project: <https://github.com/Bon-Appetit/porn-domains>
- Pinned raw snapshot: <https://raw.githubusercontent.com/Bon-Appetit/porn-domains/16611f68d447783dd6a7b7e6a32ca656a17c326a/block.3b0f12a2c2.997b06.txt>
- Upstream metadata: <https://raw.githubusercontent.com/Bon-Appetit/porn-domains/main/meta.json>
- License: [CC BY-SA 4.0](https://github.com/Bon-Appetit/porn-domains/blob/main/LICENSE)

The bundled file is a point-in-time copy so the engine works offline. Refresh
it with `scripts/update-adult-blocklist.sh`, review the resulting data diff, and
retain this attribution. The source maintainers note that their list can be
incomplete or contain false positives; domain-level blocking can also affect
unrelated services hosted on a listed domain.

`curated-domains.txt` contains additional candidate domains supplied for this
project, including labeled site domains and source-listed alternate domains.
These are loaded alongside the upstream snapshot and prioritized in extension
rules. Pattern-generated guesses from the companion research data are excluded.
