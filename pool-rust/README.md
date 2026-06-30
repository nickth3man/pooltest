# Pool - Rust

A 2D pool/billiards game scaffold in Rust, built on [Bevy](https://bevyengine.org/) (ECS game engine) and [Dimforge Rapier2D](https://rapier.rs/) (2D physics). This is a runnable hello-world: a window titled **"Pool - Rust"** opens with a green pool-table-felt background and one white cue-ball circle. There is no game logic yet — the scaffold exists so that subsequent phases (physics, cushions, cue input, pockets) can be added without re-deriving the toolchain.

## Stack

| Crate             | Version  | Release date | Role                        |
|-------------------|----------|--------------|-----------------------------|
| `bevy`            | 0.18.1   | 13 Jan 2026  | ECS game engine, 2D renderer |
| `bevy_rapier2d`   | 0.34.0   | 14 May 2026  | Rapier 2D physics integration |
| `rapier2d`        | ^0.32    | transitive   | Underlying physics engine    |

**Why these versions?** `bevy_rapier2d` 0.34.0 hard-requires `bevy ^0.18.1`, so the two must move together. Bevy 0.19 has no Rapier binding yet — bumping Bevy without a matching `bevy_rapier2d` release would silently break physics.

## Prerequisites

- **Rust** stable **>= 1.82** via [rustup](https://rustup.rs/).
- **Windows**: MSVC toolchain + Windows 10/11 SDK. Install [Visual Studio Build Tools](https://visualstudio.microsoft.com/downloads/) and select the **"Desktop development with C++"** workload. The default rustup target `x86_64-pc-windows-msvc` is the correct one.
- **Linux**: `udev`, `xkbcommon`, and (for Wayland sessions) the Wayland development packages. On Debian/Ubuntu: `sudo apt install libudev-dev libxkbcommon-dev libwayland-dev`.
- **macOS**: Xcode Command Line Tools — `xcode-select --install`.

## Run

```sh
cargo run
```

The first build downloads roughly 1–2k crates and takes several minutes. Subsequent builds are incremental and typically under a second for small edits.

## Fast dev iteration

The default `dev` profile already sets `opt-level = 1` globally and `opt-level = 3` for all dependencies, so most dependency crates are already optimised even in debug builds. To speed up **your own** crate's rebuilds, Bevy supports dynamic linking behind a feature flag. Add this to `Cargo.toml`:

```toml
[features]
dev-dynamic = ["bevy/dynamic_linking"]
```

Then run:

```sh
cargo run --features dev-dynamic
```

This is **opt-in only** — leaving it unconditional would break `cargo build --release` (the `dynamic_linking` feature conflicts with release builds). It is also host-specific and is therefore guarded behind a feature, not enabled by default.

## What's next

The scaffold is intentionally minimal so each step below is a small, isolated change:

1. **Wire Rapier** — register `RapierPhysicsPlugin::<NoUserData>::pixels_per_meter(100.0)`, mark the cue ball as `RigidBody::Dynamic` with `Collider::ball(BALL_RADIUS / 100.0)` and `Restitution::coefficient(0.7)` to approximate billiard bounce.
2. **Cushion walls** — six `Collider::cuboid` colliders with `RigidBody::Fixed` forming the table boundary.
3. **Textured sprites** — load `assets/ball.png` and use it as a `Sprite` (or `ImageMaterial`) on the ball, with the standard numbered-ball texture pack.
4. **Cue stick + input** — a second entity representing the cue; mouse aim + drag-to-shoot applies an impulse to the cue ball.
5. **Pockets + win detection** — six sensor colliders at pocket positions; track which balls have been pocketed and detect the 8-ball foul state.

## License

Licensed under either of [MIT](https://opensource.org/licenses/MIT) or [Apache-2.0](https://www.apache.org/licenses/LICENSE-2.0), at your option.
