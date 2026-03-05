# Contributing to Pool Game

Thank you for your interest in contributing!

## Code Architecture

- `src/core/` - Core systems like `EventBus`.
- `src/game/` - Game loop and state logic (`GameState`, `GameRules`, `InputHandler`).
- `src/models/` - Data structures (`Ball`, `Table`, `Vector2`).
- `src/physics/` - The `PhysicsEngine` handling movement and collisions.
- `src/rendering/` - Canvas drawing logic (`Renderer`, `BallRenderer`, `TableRenderer`).

## Pull Request Process

1. Fork the repo and create your branch from `main`.
2. Ensure you have installed the dependencies and your code compiles (`npm run build`).
3. Run the test suite (`npm run test:run`) and ensure all tests pass. If you added new features, add corresponding tests.
4. Issue that pull request! Please fill out the PR template completely.

## Style Guide

We strive for strictly typed, modular code. Please ensure:

- The TypeScript compiler passes with `strict: true`.
- You leverage the `EventBus` for cross-system communication rather than tight coupling.
