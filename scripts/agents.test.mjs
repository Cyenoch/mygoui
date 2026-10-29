import assert from "node:assert/strict";
import { spawnSync } from "node:child_process";
import {
  chmodSync, copyFileSync, cpSync, existsSync, mkdirSync, mkdtempSync,
  readFileSync, readdirSync, realpathSync, rmSync, symlinkSync, writeFileSync,
} from "node:fs";
import os from "node:os";
import path from "node:path";
import test from "node:test";
import { fileURLToPath } from "node:url";

const repository = fileURLToPath(new URL("../", import.meta.url));

function fixture(t) {
  const root = mkdtempSync(path.join(os.tmpdir(), "mygoui-agents-test-"));
  t.after(() => rmSync(root, { recursive: true, force: true }));
  for (const directory of [".agents", ".claude", ".codex", "openspec", "scripts"]) {
    cpSync(path.join(repository, directory), path.join(root, directory), {
      recursive: true,
      // Reproduce a checkout without ignored binaries or local consent.
      filter: (source) => !source.includes(`${path.sep}bin${path.sep}`) &&
        !source.endsWith(`${path.sep}bin`) &&
        !source.endsWith("settings.local.json") &&
        !source.endsWith("config.local.json"),
      verbatimSymlinks: true,
    });
  }
  copyFileSync(path.join(repository, "package.json"), path.join(root, "package.json"));
  symlinkSync(path.join(repository, "node_modules"), path.join(root, "node_modules"), "junction");
  const home = path.join(root, "home");
  mkdirSync(home);
  return {
    root,
    env: { ...process.env, HOME: home, USERPROFILE: home, IMPECCABLE_BIN: "" },
  };
}

function nodeRun({ root, env }, args) {
  return spawnSync(process.execPath, args, {
    cwd: root, env, encoding: "utf8", timeout: 30_000,
  });
}

function succeeded(result) {
  assert.equal(result.status, 0, result.error?.message ?? result.stderr ?? result.stdout);
}

test("npm OpenSpec refresh retains all workflows despite a global core profile", (t) => {
  const context = fixture(t);
  const globalConfig = path.join(context.root, "global/openspec/config.json");
  mkdirSync(path.dirname(globalConfig), { recursive: true });
  const original = JSON.stringify({ profile: "core", delivery: "both" });
  writeFileSync(globalConfig, original);
  context.env.XDG_CONFIG_HOME = path.dirname(path.dirname(globalConfig));
  succeeded(nodeRun(context, [
    process.env.npm_execpath, "run", "openspec", "--", "update", "--force",
  ]));
  const skills = path.join(context.root, ".agents/skills");
  const available = readdirSync(skills).filter(
    (name) => name.startsWith("openspec-") && existsSync(path.join(skills, name, "SKILL.md")),
  );
  assert.equal(available.length, 12);
  const shared = readFileSync(path.join(skills, "openspec-apply-change/SKILL.md"), "utf8");
  assert.match(shared, /\$openspec-apply-change \(Codex\) or \/openspec-apply-change/);
  assert.equal(realpathSync(path.join(context.root, ".claude/skills")), realpathSync(skills));
  assert.equal(readFileSync(globalConfig, "utf8"), original);
});

test("both harnesses use the verified project engine ahead of a stale home engine", (t) => {
  const context = fixture(t);
  const stale = path.join(context.env.HOME, ".impeccable/bin/impeccable");
  const used = path.join(context.root, "stale-engine-used");
  mkdirSync(path.dirname(stale), { recursive: true });
  writeFileSync(stale, `#!/bin/sh
if [ "$1" = engine-probe ]; then
  echo "impeccable-engine 0.0.1"
  exit 0
fi
touch "${used}"
exit 99
`);
  chmodSync(stale, 0o755);
  succeeded(nodeRun(context, ["scripts/setup-agents.mjs"]));
  const version = readFileSync(
    path.join(context.root, ".agents/skills/impeccable/scripts/VERSION"), "utf8",
  ).trim();
  for (const harness of [".agents", ".claude"]) {
    const result = spawnSync(
      path.join(context.root, harness, "skills/impeccable/scripts/impeccable"),
      ["engine-probe"],
      { cwd: context.root, env: context.env, encoding: "utf8", timeout: 5_000 },
    );
    succeeded(result);
    assert.equal(result.stdout.trim(), `impeccable-engine ${version}`);
  }
  assert.equal(existsSync(used), false);
  assert.equal(
    realpathSync(path.join(context.root, ".claude/agents")),
    realpathSync(path.join(context.root, ".agents/claude-agents")),
  );
  assert.equal(
    realpathSync(path.join(context.root, ".codex/agents")),
    realpathSync(path.join(context.root, ".agents/skills/impeccable/agents")),
  );
});

test("setup refuses conflicting Claude content without replacing it", (t) => {
  const context = fixture(t);
  const skills = path.join(context.root, ".claude/skills");
  rmSync(skills);
  mkdirSync(skills);
  const existing = path.join(skills, "custom-skill.md");
  writeFileSync(existing, "keep my skill\n");
  const result = nodeRun(context, ["scripts/setup-agents.mjs"]);
  assert.notEqual(result.status, 0);
  assert.match(result.stderr, /existing content was left untouched/);
  assert.equal(readFileSync(existing, "utf8"), "keep my skill\n");
});

test("setup refuses a skill release that does not match the locked engine", (t) => {
  const context = fixture(t);
  const scripts = path.join(context.root, ".agents/skills/impeccable/scripts");
  writeFileSync(path.join(scripts, "VERSION"), "0.0.1\n");
  const result = nodeRun(context, ["scripts/setup-agents.mjs"]);
  assert.notEqual(result.status, 0);
  assert.match(result.stderr, /does not match skill engine 0\.0\.1/);
  assert.equal(existsSync(path.join(scripts, "bin")), false);
});
