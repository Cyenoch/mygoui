import { createHash } from "node:crypto";
import {
  chmodSync,
  existsSync,
  lstatSync,
  mkdirSync,
  readFileSync,
  realpathSync,
  symlinkSync,
  writeFileSync,
} from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const root = fileURLToPath(new URL("../", import.meta.url));
const platform = process.platform === "win32" ? "windows" : process.platform;
const target = `${platform}-${process.arch}`;
const executable = platform === "windows" ? "impeccable.exe" : "impeccable";
const scripts = path.join(root, ".agents/skills/impeccable/scripts");
const version = readFileSync(path.join(scripts, "VERSION"), "utf8").trim();
const source = path.join(
  root, "node_modules/@impeccable", `cli-${target}`, "bin", executable,
);
const binary = readFileSync(source);
const asset = `impeccable-${target}${platform === "windows" ? ".exe" : ""}`;
const expected = readFileSync(path.join(scripts, "engine-checksums.txt"), "utf8")
  .split(/\r?\n/)
  .map((line) => line.trim().split(/\s+/))
  .find(([release, name]) => release === version && name === asset)?.[2];

if (createHash("sha256").update(binary).digest("hex") !== expected) {
  throw new Error(
    `Installed ${asset} does not match skill engine ${version}; align the npm lockfile and skill release before setup.`,
  );
}

for (const [name, destination] of [
  [".claude/skills", "../.agents/skills"],
  [".claude/agents", "../.agents/claude-agents"],
  [".codex/agents", "../.agents/skills/impeccable/agents"],
]) {
  const link = path.join(root, name);
  mkdirSync(path.dirname(link), { recursive: true });
  const existing = lstatSync(link, { throwIfNoEntry: false });
  if (!existing) {
    symlinkSync(destination, link, "dir");
  } else if (
    !existing.isSymbolicLink() ||
    !existsSync(link) ||
    realpathSync(link) !== realpathSync(path.resolve(path.dirname(link), destination))
  ) {
    throw new Error(`Expected ${link} to link to ${destination}; existing content was left untouched.`);
  }
}

// Both harnesses' hooks use this location before the unversioned home engine.
const destination = path.join(scripts, "bin", target, executable);
mkdirSync(path.dirname(destination), { recursive: true });
writeFileSync(destination, binary);
chmodSync(destination, 0o755);
console.log(`Installed verified Impeccable engine ${version}; Claude and Codex share .agents skills.`);
