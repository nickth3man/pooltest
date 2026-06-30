# pool-ts

A real-time 2D pool/billiards scaffold for the browser, written in TypeScript on top of Phaser 4.2 and Vite 8. The repo currently ships a runnable hello-world — a green felt stage with a single white cue ball centered — so you can boot the dev server and confirm the toolchain end-to-end before layering on real pool mechanics. Physics, rails, pockets, and the cue stick are deliberately next steps, not yet implemented.

## Stack

- [Phaser](https://phaser.io/) ^4.2.0 — 2D game framework
- [Vite](https://vitejs.dev/) ^8.1.0 — dev server and bundler
- [TypeScript](https://www.typescriptlang.org/) ~5.7.2

## Prerequisites

- **Node.js 20.19+** or **22.12+** — Vite 8 requires a modern Node runtime. Check with:

  ```bash
  node -v
  ```

  v24.x is fully supported.

## Install

```bash
npm install
```

## Run dev server

```bash
npm run dev
```

Opens <http://localhost:5173> with HMR enabled. Edits to `src/main.ts` reload instantly.

## Build for production

```bash
npm run build
```

Runs `tsc --noEmit` (type-check only) and then `vite build`. Static output lands in `dist/`. To preview the production build locally:

```bash
npm run preview
```

## What's next

The scaffold is intentionally minimal. The next phases will add:

- Table outline (rails + cushions) and the six pockets
- Full rack of 15 object balls with stripes/solids and the 8-ball
- Arcade physics for ball–ball and ball–rail collisions, plus rolling friction
- Cue stick with aim line, power meter, and strike
- Turn manager, foul detection, and 8-ball rules (group clearance, 8-on-the-break, etc.)

## Windows notes

- npm, Node, and Vite all run natively on Windows — no WSL required.
- The dev server binds to `0.0.0.0:5173`; open <http://localhost:5173> in your browser.
- If port 5173 is already in use, Vite automatically picks the next free port and prints it in the terminal.
- Use forward slashes in paths — they work fine in the bash shell that ships with modern Windows.

## Project layout

```
pool-ts/
├── .gitignore
├── README.md
├── index.html
├── package.json
├── tsconfig.json
├── tsconfig.node.json
├── vite.config.ts
└── src/
    └── main.ts
```
