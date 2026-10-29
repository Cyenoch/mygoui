# Agent Workflow Delta

## Purpose

Keep requirements, design decisions, implementation tracking and verification
usable across contributors and agent sessions without changing library runtime.

## ADDED Requirements

### Requirement: Development tooling is reproducible and project scoped

Contributors SHALL have pinned development CLI dependencies and committed
agent integrations. Setup SHALL make the complete OpenSpec workflow set
available without changing the user's global workflow profile or requiring
these tools to import or build the Go library.

#### Scenario: A contributor sets up a checkout

- **WHEN** the contributor runs the documented agent setup command
- **THEN** all twelve selected OpenSpec skills and the Impeccable project
  integration are available
- **AND** the contributor's global OpenSpec profile is unchanged

#### Scenario: A contributor refreshes skills through npm with a global core profile

- **WHEN** the contributor runs `npm run openspec -- update --force`
  while their global profile selects only core workflows
- **THEN** the project still exposes all twelve selected workflows
- **AND** their global configuration remains unchanged

#### Scenario: A contributor already has an older global design engine

- **WHEN** the contributor runs the documented setup and then a design command
- **THEN** the launcher uses the npm-installed engine verified against the
  committed version and checksum ahead of the older global engine
- **AND** a release mismatch fails setup explicitly

### Requirement: Claude and Codex share maintained project integrations

Claude Code and Codex SHALL load the same canonical project skills and
project instructions. Each harness SHALL have its native hook configuration
and official Impeccable agent profiles. Setup and documented refresh commands
SHALL preserve shared directory links without replacing conflicting user content.

#### Scenario: Either harness starts in the project

- **WHEN** Claude Code or Codex loads the configured project integrations
- **THEN** its skill paths resolve to `.agents/skills`
- **AND** it receives the same AGENTS.md rules and the twelve OpenSpec workflows
- **AND** its hooks resolve to the same verified project engine

#### Scenario: A contributor refreshes the shared integrations

- **WHEN** the contributor runs `make agents-update`
- **THEN** both harnesses retain their skill and agent directory links
- **AND** the shared skills retain invocation guidance usable by both harnesses

#### Scenario: A conflicting Claude skills directory already exists

- **WHEN** setup encounters an existing directory at the intended shared link
- **THEN** setup fails with the conflicting path
- **AND** it leaves the existing content untouched

### Requirement: Substantive changes have one current behavior record

New capabilities, public behavior/API changes, architectural changes and
substantial fixes SHALL use scoped OpenSpec change artifacts. Delta scenarios,
decisions and tasks SHALL remain aligned with implementation. Completed work
SHALL be verified and archived with synchronized current specs.

#### Scenario: Implementation changes the planned behavior

- **WHEN** implementation reveals a required behavior change
- **THEN** the active change's affected requirements and tasks are updated
  before the work is reported complete

#### Scenario: A correction preserves the existing contract

- **WHEN** a contributor makes a small correction without changing the contract
- **THEN** the contributor can apply a focused check without inventing a
  new feature spec or parallel PRD workflow

### Requirement: Design work loads durable native context

Design work SHALL use Impeccable with the product record, current visual
system and applicable surface brief. Native desktop constraints and the
pinned coss reference SHALL govern refinement. The project SHALL keep
engineering contracts distinct from recorded visual rules.

#### Scenario: An agent refines the gallery

- **WHEN** the agent loads the gallery's design context
- **THEN** it receives the product record, visual system and Operate surface brief
- **AND** it preserves desktop native interaction and the coss design authority

### Requirement: Verification reports reflect the evidence actually obtained

Contributors SHALL distinguish artifact validity, headless behavior/paint and
actual OS/GPU/assistive-technology results. Web-only tools MUST NOT be reported
as native Go UI validation. Unavailable checks SHALL remain explicitly unverified.

#### Scenario: The web detector does not support a changed Go file

- **WHEN** the change is reviewed for native UI quality
- **THEN** the contributor uses native behavior and rendering evidence
- **AND** does not report an unsupported detector scan as a passed native audit

### Requirement: Requirements and archived completion are mechanically checked

The documented spec check and CI SHALL validate current specs, active deltas
and archived task completion. Structural checks SHALL remain distinct from
implementation verification and SHALL NOT introduce aesthetic score gates.

#### Scenario: A malformed active requirement is submitted

- **WHEN** strict spec validation encounters a requirement with no scenario
- **THEN** the spec check fails and identifies the invalid artifact

#### Scenario: An archive contains unfinished tracked tasks

- **WHEN** the archive completion check encounters an unchecked task
- **THEN** the check fails instead of reporting the archived work complete
