# Code Audit Report

**Project**: pool-game
**Date**: 2026-03-04
**Auditor**: OpenCode (automated static analysis)
**Scope**: Full codebase - 6-dimension audit

---

## Executive Summary

| Metric | Value |
|---|---|
| **Overall Health Score** | 7.6/10 |
| **Critical Issues** | 0 |
| **High Priority Issues** | 3 |
| **Medium Priority Issues** | 4 |
| **Low Priority Issues** | 3 |

### Top 3 Priorities
1. **God-object risk in game orchestrator** - Architecture - `Game` centralizes loop, state machine, event wiring, and UI update concerns in one class.
2. **Coverage gap in highest-risk modules** - Testing - Core runtime modules (`Game`, `Renderer`, `TableRenderer`, `UIRenderer`) are not directly unit tested.
3. **Tooling guardrail gap (lint baseline missing)** - Quality - ESLint execution fails because no flat config exists, reducing automated static protection.

### What's Working Well
- Strong TypeScript safety baseline (`strict`, `noImplicitAny`, `noUnusedLocals`, `noUnusedParameters`) in `tsconfig.json:6`.
- Security sweeps found no hardcoded secrets, key material, eval/exec injection patterns, XSS sinks, or insecure TLS overrides.
- CI posture is solid: build/test pipeline and CodeQL scanning are active (`.github/workflows/ci.yml`, `.github/workflows/codeql.yml`).

---

## Findings

### 🔴 Critical

No critical findings identified in this audit.

### 🟠 High Priority

#### [H-01] Monolithic Game Orchestration
- **File**: `src/Game.ts:20`
- **Category**: Architecture
- **Issue**: `Game` owns object composition, event subscriptions, main loop timing, simulation control, UI callback fan-out, and lifecycle teardown.
- **Impact**: High change-surface area and coupling increase regression risk; feature additions will likely create brittle cross-effects.
- **Evidence**:
  ```
  export class Game {
    private stateManager: GameStateManager;
    private physicsEngine: PhysicsEngine;
    private renderer: Renderer;
    private audioManager: AudioManager;
  ```
- **Recommendation**: Split into `GameLoop`, `GameSession` (composition), and `GamePresenter`/UI bridge; keep `Game` as thin facade.
- **Effort**: Large (1+ days)

#### [H-02] Core Runtime Modules Lack Direct Tests
- **File**: `src/Game.ts:1`
- **Category**: Testing
- **Issue**: Existing tests cover `Ball`, `GameRules`, `PhysicsEngine`, `EventBus`, `InputHandler`, `AudioManager`, but not `Game` and render pipeline modules.
- **Impact**: End-to-end behavior can regress while unit suite still passes; orchestration and rendering faults may escape CI.
- **Evidence**:
  ```
  src/__tests__/AudioManager.test.ts
  src/__tests__/Ball.test.ts
  src/__tests__/EventBus.test.ts
  src/__tests__/GameRules.test.ts
  src/__tests__/PhysicsEngine.test.ts
  ```
- **Recommendation**: Add integration tests around `Game` state transitions and renderer contract tests for `Renderer`/`TableRenderer`/`UIRenderer`.
- **Effort**: Medium (1-4 hrs)

#### [H-03] Linting Baseline Missing
- **File**: `package.json:6`
- **Category**: Quality
- **Issue**: Static analysis invocation (`npx eslint .`) fails because no ESLint flat config exists.
- **Impact**: Code style/safety regressions are less likely to be caught early in local dev and CI.
- **Evidence**:
  ```
  "scripts": {
    "test": "vitest",
    "test:run": "vitest run",
    "build": "tsc"
  }
  ```
- **Recommendation**: Add `eslint.config.js` and a `lint` script; enforce in CI before build/test.
- **Effort**: Small (<30 min)

### 🟡 Medium Priority

#### [M-01] Win Condition Uses Count Instead of Unique Balls
- **File**: `src/game/GameRules.ts:74`
- **Category**: Quality
- **Issue**: Win state is computed by `sunkBalls.length === 7`; duplicate sink events can false-trigger victory.
- **Impact**: Incorrect game-end behavior under repeated/duplicate event emission.
- **Evidence**:
  ```
  checkWinCondition(): boolean {
    return this.sunkBalls.length === 7;
  }
  ```
- **Recommendation**: Track sunk balls in a `Set<number>` and assert exactly `{1..7}`.
- **Effort**: Small (<30 min)

#### [M-02] Overlapping-Ball Edge Case Not Resolved
- **File**: `src/physics/PhysicsEngine.ts:137`
- **Category**: Correctness
- **Issue**: Collision resolver exits on `distSq === 0`, skipping separation when balls occupy identical position.
- **Impact**: Rare stuck/overlap states can persist and destabilize simulation.
- **Evidence**:
  ```
  // No collision
  if (distSq === 0 || distSq > minDist * minDist) return;
  const dist = Math.sqrt(distSq)
  ```
- **Recommendation**: Inject deterministic fallback normal for zero-distance collisions and separate balls before impulse.
- **Effort**: Medium (1-4 hrs)

#### [M-03] Global DOM Listeners Not Explicitly Unbound
- **File**: `src/main.ts:42`
- **Category**: Maintainability
- **Issue**: `document` and `window` listeners are attached without stored references/removal path.
- **Impact**: Re-initialization/hot-reload scenarios can accumulate handlers and duplicate lifecycle actions.
- **Evidence**:
  ```
  document.addEventListener("visibilitychange", () => {
  });
  window.addEventListener("beforeunload", () => {
    game.destroy();
  });
  ```
- **Recommendation**: Hoist listener functions and unregister on teardown or use one-time handler where appropriate.
- **Effort**: Small (<30 min)

#### [M-04] Event Handler Failures Are Swallowed
- **File**: `src/core/EventBus.ts:167`
- **Category**: Maintainability
- **Issue**: `emit` catches and logs handler exceptions but does not surface failure state to caller/telemetry.
- **Impact**: Production faults can hide behind console noise; difficult root-cause analysis.
- **Evidence**:
  ```
  try {
    handler(event);
  } catch (error) {
    console.error(`Error in event handler for ${event.type}:`, error);
  }
  ```
- **Recommendation**: Add configurable error hook/reporter and optionally fail-fast in non-production/test modes.
- **Effort**: Medium (1-4 hrs)

### 🟢 Low Priority / Improvements

#### [L-01] Dev Dependency Vulnerabilities (Moderate)
- **File**: `package.json:11`
- **Category**: Security
- **Issue**: `npm audit` reports moderate vulnerabilities in `esbuild` via `vite`/`vitest` chain.
- **Impact**: Primarily affects developer tooling/dev server attack surface, not runtime game bundle.
- **Evidence**:
  ```
  "devDependencies": {
    "typescript": "^5.0.0",
    "vitest": "^1.0.0",
    "@vitest/ui": "^1.0.0"
  }
  ```
- **Recommendation**: Upgrade Vitest/Vite chain and re-run `npm audit`; pin secure versions in lockfile.
- **Effort**: Medium (1-4 hrs)

#### [L-02] Widespread Event Type Assertions
- **File**: `src/Game.ts:74`
- **Category**: Quality
- **Issue**: Event payloads are cast (`as BallSunkEvent`, `as ScratchEvent`, etc.) instead of type-narrowed by event type.
- **Impact**: Reduces compile-time guarantees; easier to ship mismatched payload handling.
- **Evidence**:
  ```
  this.eventBus.on(EventType.BALL_SUNK, (event) => {
    const { ball } = event as BallSunkEvent;
    const result = this.gameRules.handleBallSunk(ball);
  })
  ```
- **Recommendation**: Introduce typed per-event subscriptions (`onBallSunk`, `onScratch`) or generic mapped event signatures.
- **Effort**: Medium (1-4 hrs)

#### [L-03] Placeholder Reset Path in Physics Engine
- **File**: `src/physics/PhysicsEngine.ts:265`
- **Category**: Maintainability
- **Issue**: `reset()` exists but has no logic and no assertion around intended invariants.
- **Impact**: API appears complete but gives no guarantees; future callers may assume behavior that does not exist.
- **Evidence**:
  ```
  /** Reset the physics engine */
  reset(): void {
    // Any per-game reset logic
  }
  ```
- **Recommendation**: Either remove dead API or implement concrete reset contract with tests.
- **Effort**: Small (<30 min)

---

## Category Deep Dives

### 1. Architecture & Design
Architecture is cleanly modular at folder level (`models`, `physics`, `rendering`, `game`, `core`) and uses decoupling via event bus (`src/core/EventBus.ts`). The main risk is concentration of orchestration concerns in `src/Game.ts:20` ([H-01]) plus broad dependency fan-in (`src/Game.ts` has 11 imports). Renderer decomposition is good (`src/rendering/Renderer.ts` delegates to specialized renderers), but orchestration boundaries should be tightened.

### 2. Code Quality
Type system discipline is strong (`tsconfig.json:6-25`) and no obvious dead debug statements were found. Medium quality risks come from cast-heavy event handling (`src/Game.ts:74`, [L-02]) and edge-case correctness in physics (`src/physics/PhysicsEngine.ts:137`, [M-02]).

### 3. Security
No direct OWASP-style red flags were found in this codebase: no hardcoded secrets, no unsafe eval/exec, no XSS sink usage, no insecure crypto/TLS toggles. Current security debt is dependency-level in dev tooling ([L-01]) plus lack of lint baseline ([H-03]) that can otherwise catch unsafe patterns earlier.

### 4. Performance
Physics complexity is expected O(n^2) pairwise collision (`src/physics/PhysicsEngine.ts:116`) and acceptable for 8 balls. No synchronous filesystem/process blocking patterns are present. Minor perf/leak risk exists around global listeners (`src/main.ts:42`, [M-03]) in non-traditional runtime flows (hot reload/re-init).

### 5. Testing
Test suite is healthy and passing (6 files, 28 tests), with strong assertion density in `Ball`, `GameRules`, and `EventBus`. Primary gap is that high-impact orchestration and rendering modules are untested ([H-02]), leaving regression space where unit tests may not detect frame-loop or UI contract breakage.

### 6. Maintainability
Naming and file organization are clear, and class-level responsibilities are mostly coherent. Maintainability concerns focus on error observability strategy in event handling ([M-04]), API clarity around no-op methods ([L-03]), and concentrated complexity in orchestration ([H-01]).

---

## Prioritized Action Plan

### Quick Wins (< 1 day each)
- [ ] **H-03** `package.json:6` - Add ESLint flat config and `npm run lint`; wire into CI.
- [ ] **M-01** `src/game/GameRules.ts:74` - Replace array-length win check with unique-number set logic.
- [ ] **M-03** `src/main.ts:42` - Track/remove global DOM listeners on teardown.
- [ ] **L-03** `src/physics/PhysicsEngine.ts:265` - Remove or implement/reset contract and test it.

### Medium-term (1-5 days each)
- [ ] **H-02** Add integration tests for `Game` lifecycle/state transitions and renderer contracts.
- [ ] **M-02** Add robust zero-distance collision fallback and test with forced overlap fixture.
- [ ] **M-04** Add EventBus error-report hook + structured telemetry path.
- [ ] **L-01** Upgrade Vitest/Vite dependency chain and re-validate with `npm audit`.

### Strategic Initiatives (> 5 days)
- [ ] Decompose `Game` into loop/session/presentation components and reduce fan-in coupling.
- [ ] Establish full quality gate pipeline: lint + typecheck + tests + security audit as required PR checks.

---

## Metrics Dashboard

| Metric | Value |
|---|---|
| Files Analyzed | 39 |
| Total Lines of Code | 4,120 (excluding `package-lock.json`) |
| Languages Detected | TypeScript, HTML, CSS, JSON, Markdown, YAML |
| Test-to-Source File Ratio | 6:19 |
| Complexity Hotspots (files) | 5 |
| Security Findings | 🔴 0  🟠 0  🟡 0  🟢 1 |
| TODO / FIXME / HACK Count | 0 / 0 / 0 |
| Direct Dependencies | 5 (dev), 0 runtime |
| Avg File Length (LOC) | 106 |
| Longest File | `src/Game.ts` (285 lines) |
