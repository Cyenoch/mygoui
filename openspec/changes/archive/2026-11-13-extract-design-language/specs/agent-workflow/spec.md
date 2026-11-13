## MODIFIED Requirements

### Requirement: Design work loads durable native context

Design work SHALL use Impeccable with the product record, current visual
system and applicable surface brief. Native desktop constraints and the
design language recorded in root `DESIGN.md` SHALL govern refinement;
upstream references SHALL serve as provenance only. The project SHALL keep
engineering contracts distinct from recorded visual rules.

#### Scenario: An agent refines the gallery

- **WHEN** the agent loads the gallery's design context
- **THEN** it receives the product record, visual system and Operate surface brief
- **AND** it preserves desktop native interaction and the recorded design language

#### Scenario: A visual rule is missing from the record

- **WHEN** a design decision needs a rule that `DESIGN.md` does not state
- **THEN** the contributor records the rule in `DESIGN.md` through a change
  instead of deferring to an upstream reference
