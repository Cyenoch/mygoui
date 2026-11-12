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

The home view introduces the library: install command, downloads, a live
preview and the catalog. Flat grouped navigation and Cmd/Ctrl+K search select
a component. Its page leads with title, description and the Preview / Code /
API / Tokens switch, which share that selection and retain window-owned
example state. Appearance (popover) and content zoom remain immediately
available; neither moves the page.

## Layout and native constraints

Preserve the custom native title bar, real OS control insets and draggable
regions. The sidebar's top row and the page title bar share one height. The
browser build has no window controls and uses the system font. Use the existing compact layout and scaling thresholds; keep navigation,
content and overlay actions reachable at compact widths and enlarged text.
Use the root visual system and pinned coss component appearance. Keep native
focus, selection, window commands and accessible control naming.

## Review evidence

Use focused Tester interactions and the existing light/dark, compact/wide
specimen collectors. Record actual window/GPU observations separately.
Browser live mode and HTML/CSS detectors do not analyze this surface.

## Open decisions

None introduced by workflow adoption. Future behavior or visual changes belong
in their scoped OpenSpec change and update this brief only when needed.
