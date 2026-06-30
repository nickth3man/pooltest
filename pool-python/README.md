# Pool - Python

A real-time 2D pool/billiards game in Python, built on **pygame-ce** (the actively maintained fork of pygame). The current state is a runnable hello-world scaffold: an 800x600 green-felt window running at 60 FPS, with a clean exit on ESC or window close. Physics, balls, cushions, collisions, game state, and AI are **not yet implemented** — this is the foundation those pieces will be built on.

## Stack

- **Python** 3.11+
- **pygame-ce** `>=2.5.7,<3` — imported as `import pygame` (pygame-ce is API-compatible with pygame 2)
- **uv** — package & project manager (PEP 517 backend: `uv_build`; PEP 735 dev groups)
- **Ruff** `>=0.15.20,<0.16` — linting (`select = ["ALL"]`) and formatting (Black-compatible)
- **ty** `>=0.0.55,<0.1` — Astral's type checker, strictest mode (`all = "error"`)

## Prerequisites

- **uv** ≥ 0.6 — install from https://docs.astral.sh/uv/getting-started/installation/. uv installs and manages the Python interpreter automatically, so no separate `python`/`py` install is required.

## Install

Clone the repo and sync dependencies (this creates `.venv` and resolves `uv.lock`):

```bash
cd pool-python
uv sync
```

That single command is all you need — no separate `pip install`, no manual venv creation, no platform-specific steps.

## Run

Launch the game with the installed console script:

```bash
uv run pool-python
```

Or run the package as a module:

```bash
uv run python -m pool_python
```

Press **ESC** or close the window to exit.

## Windows notes

- Use `uv run` instead of activating a venv manually — uv handles the per-platform venv path for you.
- pygame-ce ships **pre-built wheels** for Windows, so no C compiler is required.
- All `uv run` / `uv sync` / `uv add` commands work identically on `cmd`, PowerShell, and POSIX shells.

## Project layout

```
pool-python/
├── .gitignore
├── .python-version
├── README.md
├── pyproject.toml
└── src/
    └── pool_python/
        ├── __init__.py
        └── main.py
```

## Quality gates

The project is configured for the strictest stable settings of Ruff and ty. Run them (in this order) before committing:

```bash
uv run ruff check --fix   # lint with all rules; --fix applies safe auto-fixes
uv run ruff format        # format in place (Black-compatible)
uv run ty check           # strictest type checking — all rules set to error
```

`uvx ty check` is equivalent if you prefer not to install ty into the project venv. `uv run` keeps `.venv` and `uv.lock` in sync automatically.

## What's next

- Vector math & integration (semi-implicit Euler with friction)
- Ball entity (position, velocity, radius, sprite)
- Cushions (table rails, circle-vs-segment collision)
- Ball-ball collisions (elastic, equal-mass)
- Game state (rack, break, turn, fouls, win)
- AI opponent (basic shot selection + difficulty)
- Tests, CI, and packaging polish

## License

MIT
