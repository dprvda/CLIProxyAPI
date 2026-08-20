# MIGRATION-MANIFEST - CLIProxyAPI

KEPT by owner ruling — the LLM router returns on Linux; this is a REVIVE, not an archive.
Born on Windows 11; frozen 2026-08. This machine migrated to Ubuntu Linux; everything else
(growth-os crons, n8n, cc-wrapper, hermes, panel, mining) was decommissioned, but this repo's
service comes back. Read `docs/planning/linux-migration-2026-08-20.md` +
`linux-system-setup-2026-08-20.md` in dprvda/pravda-automations-page before reviving.

## Data (gitignored / external)
- `release-v7.2.89/cli-proxy-api.exe` (+ `.previous.exe`, `.pre-keepalive.exe`) — gitignored
  (`*.exe`), REBUILD via `go build` (go.mod: go 1.26.0). Packaged zip/checksums/example config in
  that same dir ARE tracked (committed 2026-08-20 in the freeze commit).
- `auths/` in repo — empty placeholder, `.gitkeep` only. Not where real credentials live.
- Real runtime state lives OUTSIDE this repo at `C:\Users\dprvd\.cli-proxy-api\`
  (planned mirror dest: `E:\migration-mirror\C\Users\dprvd\.cli-proxy-api\` — mirror copy was
  still running at freeze time, treat as PLANNED not verified-present):
  - `auth/` — 143M, per-account OAuth credential JSONs (Claude + Codex, several Google accounts,
    some `.disabled`). Listed in `growth-os-ops/migration-exports/secrets-manifest.csv` as
    `agent-creds`.
  - `config.yaml` — 687B, router config: model routing table + account/auth mapping (not just
    path config — this is live credentials-adjacent state).
  - `claude-routing-verdict.json` — 1.5K, a one-time billing-plan verification record.

## Services + scheduled tasks this repo ran
- Scheduled task `GrowthOS-llm-router` | every 1m | ENABLED at freeze | in PROTECTED_CRONS | via
  `C:\Users\dprvd\growth-os-ops\run-hidden.vbs` -> `powershell -File
  C:\Users\dprvd\growth-os-ops\llm-router-ensure.ps1`. That ensure script and its sibling
  `llm-router-supervise.ps1` live OUTSIDE this repo, in growth-os-ops (mirrored at
  `E:\migration-mirror\C\Users\dprvd\growth-os-ops\`).
- Ensure script presence-checks the supervisor and, if absent, spawns
  `llm-router-supervise.ps1` DETACHED via `Win32_Process.Create` (dodges Task Scheduler's
  `ExecutionTimeLimit` kill of descendant processes — a plain child process got killed ~2min in).
- Supervise script is the actual launcher: `release-v7.2.89\cli-proxy-api.exe --config
  C:\Users\dprvd\.cli-proxy-api\config.yaml`, restarts it the instant it exits, logs exit
  code/lifetime to `growth-os-ops\llm-router-exits.log`.
- Listens on `127.0.0.1:17999`. Serves BOTH providers through one endpoint (Claude OAuth models
  + Codex/GPT OAuth models) — every ADE session's `ANTHROPIC_BASE_URL` on this box points here.
- Known unfixed bug carried into the freeze: the binary exits by itself under some condition (no
  panic, no crash event, empty stderr) — contained by the supervisor's ~1s relaunch, never
  root-caused. Worth fixing before porting the ensure/supervise pattern forward.
- Linux target: a systemd user service with `Restart=always` replaces the ensure-loop +
  supervisor + VBS wrapper entirely — no Task Scheduler analog needed on Linux.

## Windows coupling
- `docker-build.ps1` is tracked in the repo (Windows Docker build helper) — not blocking, use
  `go build` / the Linux/Docker path directly.
- No hardcoded `C:\Users\dprvd` paths in tracked repo source (`git grep` clean) — all the
  Windows-specific coupling (ensure loop, supervisor, VBS, Task Scheduler, the
  `.cli-proxy-api` config location) lives OUTSIDE this repo, in growth-os-ops.

## Secrets
- OAuth credential JSONs + `config.yaml` at `C:\Users\dprvd\.cli-proxy-api\` are NOT in this
  repo's git history (repo's own `auths/` stays an empty placeholder).
- Per plan, harvest target is 1Password vault `btc-bot-dev`, doc `migration/C/Users/dprvd/
  .cli-proxy-api/<file>`, indexed by `migration/INDEX` — harvest status NOT verified by this
  pass; confirm before treating those Google-account OAuth tokens as recoverable.

## Remote
- `fork` = `https://github.com/dprvda/CLIProxyAPI.git` (push target, owner's fork) — checked
  out branch `wip/migration-freeze`, tracking `fork/wip/migration-freeze`, up to date at freeze.
- `origin` = `https://github.com/router-for-me/CLIProxyAPI.git` (upstream) — pushes here 403'd
  during the freeze. Push to `fork`, not `origin`, until upstream access is re-verified.
