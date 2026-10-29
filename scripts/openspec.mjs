import { spawnSync } from "node:child_process";
import { fileURLToPath } from "node:url";

const root = new URL("../", import.meta.url);
const result = spawnSync(
  process.execPath,
  [
    fileURLToPath(new URL("node_modules/@fission-ai/openspec/bin/openspec.js", root)),
    ...process.argv.slice(2),
  ],
  {
    stdio: "inherit",
    env: {
      ...process.env,
      OPENSPEC_TELEMETRY: "0",
      XDG_CONFIG_HOME: fileURLToPath(new URL(".agents/config/", root)),
    },
  },
);
if (result.error) throw result.error;
process.exit(result.status ?? 1);
