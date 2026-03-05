# TypeScript HTML5 Canvas Pool Game

[![CI](https://github.com/nicolas/pooltest/actions/workflows/ci.yml/badge.svg)](https://github.com/nicolas/pooltest/actions/workflows/ci.yml)
[![GitHub Pages](https://github.com/nicolas/pooltest/actions/workflows/pages.yml/badge.svg)](https://github.com/nicolas/pooltest/actions/workflows/pages.yml)

A pure TypeScript and HTML5 Canvas implementation of a Billiards/Pool game. Built using modern web standards, strict TypeScript, and continuous integration.

## Play

**[Play the game here](https://nicolas.github.io/pooltest/)** *(Update with GitHub pages link after deployment)*

## Features

- Custom Physics Engine (Collisions, friction, rebounds)
- HTML5 Canvas rendering
- Strict Type-Safety (`tsc`)
- Zero-bundle development environment leveraging ES Modules
- Component-based architecture (EventBus, GameState, Renderers)

## Project Structure

- Source code: `src/` (organized by domain: `game/`, `physics/`, `rendering/`, `models/`, `core/`, `audio/`)
- Tests: `src/__tests__/`
- Documentation: `docs/`
- Maintenance scripts: `scripts/`

See `docs/project-structure.md` for a full layout and conventions.

## Local Development

### Prerequisites

- Node.js (v18+)

### Installation

```bash
git clone https://github.com/nicolas/pooltest.git
cd pooltest
npm install
```

### Running Locally

```bash
# Start the TypeScript compiler in watch mode
npm run dev
```

Then, open `index.html` in your browser. (Since we are using native ES modules, you may need a local static web server like `npx serve .` or the VS Code Live Server extension).

### Testing

Tests are written using Vitest.

```bash
# Run tests in watch mode
npm test

# Run tests once (used in CI)
npm run test:run
```

### Quality Check

```bash
# Run lint + typecheck + tests + build
npm run check
```

### Building

```bash
npm run build
```

### Cleanup

```bash
# Remove generated artifacts
npm run clean
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.
