# Tasks

## 1. Product record

- [x] 1.1 Update `PRODUCT.md` through Impeccable's `init` playbook: Purpose names the four pillars (theme tokens, component library, design and behavior reference, agent skills for application authors and maintainers); each pillar that is not delivered yet is marked planned with its follow-up change; verify the file keeps the `impeccable:product-schema` marker and its existing sections
- [x] 1.2 Rewrite the Visual Authority section of `PRODUCT.md`: coss remains the authority until `extract-design-language`, after which `DESIGN.md` is the authority and coss is provenance; verify the section states both the current rule and the end state
- [x] 1.3 Run `make design-doctor` and read its findings array; verify no new finding refers to `PRODUCT.md`

## 2. Entry points

- [x] 2.1 Replace the opening paragraph of `README.md` with the product statement from `PRODUCT.md`, keeping the coss reference as the current visual source; verify every link in the paragraph resolves
- [x] 2.2 Replace the first line of `AGENTS.md` ("generic native MyGo component library with a shadcn-inspired style") with the same statement in one sentence; verify `CLAUDE.md` still imports it unchanged

## 3. Integration check

- [x] 3.1 Run `make spec-check`; verify the change validates with `skip_specs: true` and no main spec changed
- [x] 3.2 Run `git diff --stat`; verify only `PRODUCT.md`, `README.md`, `AGENTS.md` and this change directory are touched

## Workflow follow-up

- Archive this change after review.
- Propose `publish-agent-skills` next, then `extract-design-language`, `open-theme-tokens`, `component-behavior-reference` and `native-effects` in the order given in `design.md`.
