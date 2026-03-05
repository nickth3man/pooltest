import { execSync } from "node:child_process";

const commands = [
  "npm run lint",
  "npm run typecheck",
  "npm run test:run",
  "npm run build",
];

for (const command of commands) {
  try {
    execSync(command, { stdio: "inherit" });
  } catch (error) {
    if (typeof error === "object" && error !== null && "status" in error) {
      const status = error.status;
      if (typeof status === "number") {
        process.exit(status);
      }
    }

    process.exit(1);
  }
}
