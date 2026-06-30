# AGENTS.md — pool-ts

2D web pool/billiards game — Vite + Phaser 4 + TypeScript (browser-only).

## Stack

- **Runtime target:** browser (Vite dev server / static host)
- **Node.js:** 20.19+ or 22.12+ (Vite 8 floor; tested v24.18.0); npm bundled (tested 11.16)
- **Dependencies:**
  - `phaser ^4.2.0` (npm `phaser`) — HTML5 2D framework, WebGL2 renderer (v4 "Giedi")
- **Dev dependencies:**
  - `vite ^8.1.0` (npm `vite`) — dev server + bundler
  - `typescript ~5.7.2` (npm `typescript`) — strict mode
- **Module system:** ESM (`"type": "module"`)
- **TypeScript config:** strict; `moduleResolution: bundler`; target ES2020; DOM + DOM.Iterable libs; `noEmit: true`; separate `tsconfig.node.json` for `vite.config.ts`
- **HTML entry:** `index.html` mounts `<div id="game-container">`; Phaser canvas attaches via `parent: 'game-container'`
- **Build output:** `dist/` (single JS bundle, ~360 KB gzipped — Phaser is large; the "chunk > 500 kB" warning is expected, not an error)
- **Build gate:** `npm run build` runs `tsc --noEmit` FIRST — type errors fail the build

## Dev commands

Run from the `pool-ts/` folder.

### Dependencies
- `npm install` — install deps (creates `node_modules/`, `package-lock.json`)
- `npm update` — update deps within ranges
- `npm outdated` — list outdated
- `npm install <pkg>` — add runtime dep; `npm install -D <pkg>` (`--save-dev`) — add dev dep
- `npm uninstall <pkg>` — remove
- `npm audit` — vulnerability audit; `npm audit fix`

### Dev / build
- `npm run dev` — Vite dev server with HMR → http://localhost:5173 (bound 0.0.0.0, LAN-accessible; auto next free port if busy)
- `npx tsc --noEmit` — type-check only
- `npm run build` — production build (`tsc --noEmit` + `vite build` → `dist/`)
- `npm run preview` — preview the production build locally

### Cleanup
- `rm -rf node_modules dist` — full clean reinstall

Requires Node 20.19+ or 22.12+ (Vite 8). No C compiler or native build tools needed.
