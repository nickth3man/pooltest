import { rmSync } from "node:fs";

const paths = ["dist", "coverage"];

for (const path of paths) {
  rmSync(path, { recursive: true, force: true });
}
