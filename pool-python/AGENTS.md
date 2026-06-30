# AGENTS.md — pool-python

2D pool/billiards game in Python on pygame-ce, managed by uv, with strictest ruff + ty.

## Stack

- **Language:** Python >=3.11 (CPython; managed by uv)
- **Project manager:** uv >=0.6 (tested 0.11.8)
  - PEP 517 build backend: `uv_build` (`requires = ["uv_build>=0.11.25,<0.12"]`)
  - PEP 735 dev groups: `[dependency-groups] dev = [...]`
  - Lockfile: `uv.lock` (committed for reproducible installs)
- **Runtime dependency:**
  - `pygame-ce >=2.5.7,<3` — PyPI distribution name `pygame-ce`; **import name is `pygame`** (API-compatible with legacy pygame 2)
- **Dev dependencies (group `dev`):**
  - `ruff >=0.15.20,<0.16` — lint + format
  - `ty >=0.0.55,<0.1` — Astral type checker
- **Layout:** `src/` layout; import package `pool_python`; console script `pool-python` → `pool_python.main:main`
- **Ruff config (strictest):** `select = ["ALL"]`, `line-length = 100`, `target-version = "py311"`; ignores are only the formatter-conflict set + `D213` (mutually exclusive with `D212`)
- **ty config (strictest):** `[tool.ty.rules] all = "error"`, `[tool.ty.terminal] error-on-warning = true`, `[tool.ty.environment] python-version = "3.11"`
- **Platform:** Windows ships pre-built pygame-ce wheels — **no C compiler required**

## Dev commands

Run from the `pool-python/` folder.

### Dependencies
- `uv sync` — install/sync deps, create `.venv`, resolve `uv.lock`
- `uv lock` — refresh lockfile without installing
- `uv lock --upgrade-package <name>` — upgrade one dep
- `uv lock --upgrade` — upgrade all deps
- `uv add <pkg>` — add a runtime dependency
- `uv add --dev <pkg>` (or `--group dev`) — add a dev dependency
- `uv remove <pkg>` — remove a dependency
- `uv python pin 3.11` — pin project Python; `uv python list` — available interpreters

### Run
- `uv run pool-python` — run via console script
- `uv run python -m pool_python` — run as module
- `uv run python src/pool_python/main.py` — run source file directly

### Quality gate (strict)
- `uv run ruff check` — lint (all rules); `--fix` auto-fix; `--watch` watch mode
- `uv run ruff format` — format; `--check` check only; `--diff` show diff
- `uv run ty check` — strict type check; alternative without install: `uvx ty check`
- Full gate: `uv run ruff check && uv run ruff format --check && uv run ty check`

### Tests / cleanup
- `uv run pytest` — once a test runner is added
- Clean: delete `.venv/` then `uv sync` to recreate

`uv run` auto-syncs `.venv` and `uv.lock`; no manual venv activation is needed. Works identically in cmd, PowerShell, and bash.
