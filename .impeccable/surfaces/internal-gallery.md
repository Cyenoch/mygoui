---
version: 1
slug: "internal-gallery"
primary_target: "internal/gallery"
related_targets: ["examples/gallery"]
---

# Native gallery workbench

## Scope

Targets: internal/gallery and examples/gallery. Captured from the existing
PRODUCT.md, docs/design.md and gallery source; no new visual direction is selected.

## Mode

Operate

## Audience and job

Go application authors discover a component, exercise its native behavior,
and inspect usable public Go API examples in one workspace.

## Workflow

The home view is the landing page (Persuade) and owns the window, with no
sidebar and its own title-bar links and search: a centered hero with the
install command, a stage of live components, catalog numbers, what the
library is built on, the family grid and downloads. Flat grouped navigation and Cmd/Ctrl+K search select
a component. Its page is one document: title, description and page actions,
the lead example, Installation, Usage, API Reference, further Examples and
Theme tokens. Each example switches between Preview and its own source and
retains window-owned state. Wide windows keep an "On this page" index beside
the content that marks the section being read and scrolls to a chosen one.
Example source and the reference's fields and types come from the code; the
page supplies titles, descriptions, the usage snippet and defaults. Appearance (popover) and content zoom remain immediately
available; neither moves the page.

## Layout and native constraints

Preserve the custom native title bar, real OS control insets and draggable
regions. The sidebar's top row and the page title bar share one height. The
browser build has no window controls and uses the system font. Use the existing compact layout and scaling thresholds; keep navigation,
content and overlay actions reachable at compact widths and enlarged text.
Use the design language recorded in root DESIGN.md. Keep native
focus, selection, window commands and accessible control naming.

## Review evidence

Use focused Tester interactions and the existing light/dark, compact/wide
specimen collectors. Record actual window/GPU observations separately.
Browser live mode and HTML/CSS detectors do not analyze this surface.

## Open decisions

None introduced by workflow adoption. Future behavior or visual changes belong
in their scoped OpenSpec change and update this brief only when needed.
