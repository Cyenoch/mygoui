# Proposal

## Why

The project needs a durable requirements/change workflow and a complete
design toolkit that preserve native MyGo behavior across agent sessions.
The user explicitly selected OpenSpec and Impeccable and retired the
Matt Pocock skill workflow for this repository.

## What Changes

- Install the complete OpenSpec skill set and a project-scoped profile.
- Pin development CLIs, document setup and validate specs in CI.
- Share canonical skills and rules between Claude Code and Codex, register
  official agent profiles and hooks, and verify the project engine at setup.
- Record native project context, per-artifact rules and operation guidance.
- Capture the existing visual system and gallery surface context for Impeccable.
- Route design review and verification to native evidence with explicit
  browser-tool limitations.
- Seed native-library, semantic-theme and gallery-workbench main specs from
  existing documented contracts; these establish a baseline, not new behavior.
- Replace the retired skill workflow in agent/contributor guidance.

## Capabilities

### New Capabilities

- `agent-workflow`: Reproducible, project-scoped requirements and design work
  with current contracts, tracked changes and honest verification.

### Modified Capabilities

None. The native library's public behavior is unchanged by this adoption.

## Impact

Development configuration, generated skills, documentation, CI and
third-party skill notices. No Go runtime dependency or public API changes.
Existing uncommitted library changes remain outside this adoption's ownership.
