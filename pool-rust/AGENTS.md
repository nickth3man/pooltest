# AGENTS.md — pool-rust

2D pool/billiards game — Bevy ECS + Rapier2D physics (Rust).

## Stack

- **Language:** Rust stable >=1.82 (edition 2021); Windows toolchain `x86_64-pc-windows-msvc`
- **Dependencies:**
  - `bevy 0.18.1` (crate `bevy`; `default-features = false, features = ["2d"]`) — ECS + 2D renderer
  - `bevy_rapier2d 0.34.0` (crate; feature `debug-render-2d`) — Rapier physics bridge
  - `rapier2d 0.32` — transitive, pulled by bevy_rapier2d
- **Version lock (important):** `bevy_rapier2d 0.34.0` hard-requires `bevy ^0.18.1`. Do NOT bump `bevy` to 0.19 — it will not compile. Wait for a bevy_rapier2d release that supports 0.19.
- **Crate type:** binary `pool-rust`; entry `src/main.rs`
- **Build profiles:** `profile.dev` opt-level=1; `profile.dev.package."*"` opt-level=3 (fast deps in dev); `profile.release` lto="thin", codegen-units=1
- **Optional dev feature:** `dev-dynamic = ["bevy/dynamic_linking"]` for faster rebuilds — opt-in only (unconditional dynamic_linking breaks release builds)
- **Windows prerequisite:** MSVC toolchain + Visual Studio Build Tools "Desktop development with C++" workload (provides `link.exe` + Windows SDK). Default `rustup` host `x86_64-pc-windows-msvc` is correct.

## Dev commands

Run from the `pool-rust/` folder.

### Build / run
- `cargo run` — build + run (first build ~7–10 min; Bevy compiles a lot)
- `cargo check` — fast type/check (no codegen)
- `cargo build` — build debug; `cargo build --release` — optimized
- `cargo run --features dev-dynamic` — opt-in fast-iteration (dynamic linking)

### Quality
- `cargo clippy` — lints; strict: `cargo clippy -- -D warnings`
- `cargo fmt` — format; `cargo fmt --check` — check only

### Tests / deps / docs
- `cargo test` — run tests; `cargo test -- --nocapture` — show println output
- `cargo update` — update deps within Cargo.toml ranges
- `cargo add <crate>` — add dependency; `cargo tree` — dependency graph
- `cargo doc --open` — generate + open docs

### Toolchain / cleanup
- `rustup show` — current toolchain; `rustup update` — update Rust
- `cargo clean` — remove `target/`

Ensure the VS Build Tools C++ workload is installed before the first `cargo run` (linking fails otherwise with `linker 'link.exe' not found`).

## Development tooling

This section documents the project-quality tooling configured in this repo.

### Required Rust components

```sh
rustup component add rustfmt clippy rust-analyzer llvm-tools-preview
```

### Recommended Cargo plugins

Install via `cargo-binstall` if you have it (faster):

```sh
# Core quality tools
cargo install --locked cargo-deny cargo-nextest cargo-machete

# Local background runner (better than cargo-watch for Bevy's long compile)
cargo install --locked bacon

# Task runner (optional)
cargo install --locked just
```

### Git hooks

This repo uses `cargo-husky` (configured in `Cargo.toml`) to install a pre-commit hook that runs `cargo fmt --check` and `cargo clippy`. The hooks are generated the first time you build or test:

```sh
cargo test --no-run
```

To bypass hooks in an emergency: `git commit --no-verify`.

### Running checks locally

```sh
cargo fmt --check                 # formatting
cargo clippy -- -D warnings       # lints
cargo test                        # tests
cargo deny check                  # licenses + advisories (requires cargo-deny)
```

### Bevy lints (`bevy_lint`)

`bevy_lint` is the official Bevy-aware linter. It requires a specific nightly toolchain.

**One-time install:**

```powershell
# Windows (PowerShell)
rustup toolchain install nightly-2026-01-22 --component rustc-dev --component llvm-tools
rustup run nightly-2026-01-22 cargo install --git https://github.com/TheBevyFlock/bevy_cli.git --tag lint-v0.6.0 --locked bevy_lint
```

```sh
# Linux / macOS
rustup toolchain install nightly-2026-01-22 --component rustc-dev --component llvm-tools
rustup run nightly-2026-01-22 cargo install --git https://github.com/TheBevyFlock/bevy_cli.git --tag lint-v0.6.0 --locked bevy_lint
```

**Run:**

```sh
bevy_lint --workspace --all-targets --all-features
```

Keep the `nightly-2026-01-22` toolchain installed; `bevy_lint` loads `librustc_driver` from it at runtime.

### VS Code

Open the workspace in VS Code and install the recommended extensions (`.vscode/extensions.json`). The shared settings in `.vscode/settings.json` keep rust-analyzer's build artifacts separate from `cargo run`, which is essential for Bevy-sized projects.
