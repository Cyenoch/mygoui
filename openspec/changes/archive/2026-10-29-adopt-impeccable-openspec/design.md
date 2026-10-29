# Design

## Context

See proposal.md for scope. The repository already has Go/MyGo maintenance
skills, detailed capability/ownership documents and native rendering collectors.
It has PRODUCT.md but no root visual-system record or OpenSpec root.

## Goals / Non-Goals

Use the official integrations with enough project context to be effective
on the next task. Keep runtime and library import behavior unchanged.
Do not redesign controls, uninstall unrelated global tools, add a browser
runtime or manufacture validation of existing uncommitted code.

## Decisions

1. Pin the two npm development CLIs and commit generated skills. A private
   development-only package manifest makes versions reproducible without
   adding Go dependencies. Impeccable's update distribution is a separate
   refresh and its resulting diff must be reviewed.
2. Select all twelve OpenSpec workflows using its official custom profile
   and skills-only delivery. Scope XDG_CONFIG_HOME in one Node entry point
   used by npm and Make to a checked-in project configuration. This avoids mutating global settings
   or generating home-directory prompts.
3. Keep the standard spec-driven schema. Project context, artifact rules
   and operation guidance provide native constraints without a forked workflow.
4. Extract root DESIGN.md from current Go tokens/components and keep
   docs/design.md as the engineering authority. A small sidecar holds motion
   and compact-layout metadata; no fake HTML versions of native components.
5. Use the adaptive selector for native platform guidance, with explicit
   desktop scope. MyGo semantics override phone-specific examples.
   Browser live/generate and web detection remain unsupported for native UI;
   existing Tester collectors and real-window evidence cover that work.
6. Wire strict spec and archive-task validation into a small CI job. Keep
   design scores and unsupported source scans out of blocking gates.
7. Preserve upstream skill files, licenses and provenance. Project rules
   live in AGENTS.md/config.yaml so integrations can be regenerated.
8. Share `.agents/skills` with Claude through `.claude/skills`, import
   AGENTS.md from CLAUDE.md, and link each harness's official agent profiles
   from canonical files under `.agents`. Keep native hooks in each harness's
   committed manifest. Refresh Claude profiles before restoring the Codex
   variant as the common skill so the updater cannot leave a Claude-only
   skill in the shared directory.
9. Verify the npm platform engine against the committed skill version and
   SHA256 manifest, then install it into the ignored sibling binary directory.
   This keeps hooks and CLI commands ahead of stale global installations
   without forking the upstream launcher.

## Risks / Trade-offs

- Web-oriented tooling misses native sources: document the boundary, use
  native review playbooks and actual rendering evidence.
- Claude and Codex hooks need harness trust: record that configuration alone is not
  execution, and do not claim native detector coverage.
- Generated integrations can drift: keep the full profile in the repository,
  provide explicit refresh commands and inspect doctor/dry-run output.
- Summarized baseline specs can drift from detailed docs: link their source
  records and update requirements alongside public changes.

## Migration Plan

Install and validate the tooling, seed current contracts, record this change,
review and archive it after checks. Existing working-tree changes are
preserved. Reverting this development integration does not require changing
Go module metadata or application behavior.
