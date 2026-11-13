# Spec Delta

## Purpose

Give an agent that builds or maintains a mygoui application the library's
component, theme, design and behavior, motion and shader-effect rules as
installable skills, with one source per topic for both audiences.

## ADDED Requirements

### Requirement: Published skills install from the repository

The repository SHALL keep the mygoui skills, each named with the `mygoui-`
prefix, in its root `skills/` directory, where the `skills` CLI discovers
them. Using that directory as the source SHALL offer only the mygoui skills.

#### Scenario: An author lists a checkout's skills

- **WHEN** the author runs `npx skills add <checkout> --list`
- **THEN** the listing includes every `mygoui-` skill with its description

#### Scenario: An author installs only the published skills

- **WHEN** the author runs `npx skills add <checkout>/skills --list`
- **THEN** the CLI offers the `mygoui-` skills and no maintainer or generated
  upstream skill

### Requirement: The project's skill directory links the published skills

`.agents/skills` SHALL hold a relative link to each published skill in
`skills/`, so a harness reading `.agents/skills`, directly or through
`.claude/skills`, reads the published files. Setup and the OpenSpec refresh
SHALL leave those links resolving.

#### Scenario: A harness resolves a project skill

- **WHEN** `.agents/skills/mygoui-<topic>` or `.claude/skills/mygoui-<topic>` is read
- **THEN** it resolves to `skills/mygoui-<topic>`

#### Scenario: A contributor runs setup or the OpenSpec refresh

- **WHEN** the contributor runs the setup script (`npm run agents:setup`) or `npm run openspec -- update --force`
- **THEN** the `mygoui-` entries in `.agents/skills` still resolve to the published copies

### Requirement: Application authors load only the public contract

Each skill's main document SHALL state only behavior and APIs an application
can rely on. Maintainer-only material SHALL live in a sibling document that
the main document reaches through one pointer naming work on the mygoui
repository as its condition.

#### Scenario: An application author's agent loads a skill

- **WHEN** the agent reads a `mygoui-` skill's main document for application work
- **THEN** it finds no internal package, OpenSpec or Impeccable instruction
- **AND** the maintainer document is named once, conditioned on repository work

### Requirement: Skills defer to the owning records

A skill SHALL point at the record that owns a rule or value instead of
restating values the record owns, and SHALL reach that record in the module
version the reader's application uses.

#### Scenario: A recorded value changes

- **WHEN** a change alters a value owned by `DESIGN.md` or a `docs/` record
- **THEN** a skill that points at that record needs no edit unless the rule it states changed

#### Scenario: An installed skill is read outside this repository

- **WHEN** an application author's agent follows a skill's record pointer
- **THEN** the pointer resolves through the module source of the version the application requires, or the repository at that version when the application vendors its dependencies

### Requirement: Names and examples in skills are real

Every API a skill names SHALL exist in the module or MyGo version the
repository pins, and every Go example SHALL compile against them. A change
that renames or removes such an API, or alters a rule a skill states, SHALL
update the affected skill in the same change.

#### Scenario: A skill names an API

- **WHEN** a reviewer looks up a name the skill gives with `go doc`
- **THEN** the pinned module or MyGo version defines it

#### Scenario: A later change alters a stated rule

- **WHEN** a change modifies a public API or rule that a `mygoui-` skill states
- **THEN** that change's tasks include updating the skill
