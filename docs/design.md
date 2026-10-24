# Design and public contract

`github.com/Cyenoch/mygoui` is one native Go library covering the union of
[HeroUI v3](https://heroui.com/docs/react/components)'s web and native
component catalogs, built for [MyGo](https://github.com/egoist/mygo). It is
not a React port, a browser compatibility layer or a WebView component kit.

## Ownership

A constructor receives a build-scoped `*ui.Context` and returns a build-scoped
`ui.Element`. Callers retain values, selection, expansion, validation and
`ui.Handle` focus targets per window. MyGo owns editing, input, focus,
selection, scrolling, animation scheduling and accessibility. Views perform no
network, filesystem, clipboard or background work.

Composition is Go-native. Roots take typed options and either a body callback
for freely composed content or a typed slice for coordinated collections
(accordion sections, tabs, menus). There is no provider tree, context
propagation, ambient build scope or slot registry. Collection components use
stable keys and native virtualized state rather than rebuilding every row.

The package layout follows component families:

| Package | Responsibility |
| --- | --- |
| `mygoui` | Surfaces, typography, actions, status, progress, avatars, cards, lists |
| `theme` | HeroUI semantic tokens, geometry, typography, host preferences |
| `form` | Field primitives and form controls |
| `list` | ListBox selection and rendering |
| `navigation` | Accordion, disclosure, tabs, breadcrumbs, pagination, menus |
| `overlay` | Modal family, popovers, tooltips and toasts |
| `date` | Calendars, date fields/pickers and time fields |
| `color` | Color area, channel sliders, fields, swatches and pickers |
| `data` | Table |
| `locale` | Library-level accessible messages |

## Appearance

Colors, radii, spacing and type come from HeroUI's published tokens, translated
from `references/heroui/packages/styles/themes/default/variables.css`.
Components read `theme.Colors(c)` live; palettes are immutable snapshots rebuilt
only when `theme.Configure` changes accent, radius, vibrancy or typography.

The default palette is HeroUI's light/dark pair, including surface, overlay,
field, default, accent, success, warning and danger roles with their derived
hover, soft and separator levels. HeroUI's desktop metrics are authoritative:
4 DIP spacing, 8 DIP radius base, 14 DIP control text, 36 DIP control height,
16 DIP card padding and 24 DIP large radius.

`theme.Apply(c)` installs the matching native widget theme for the current
appearance. `theme.ApplyMode(c, dark)` exists for previews and tests.
`theme.Configure` is process-wide host policy; widget state remains
window-local.

## API rules

- Options are typed Go structs. Unsupported combinations fail visibly rather
  than silently degrading.
- Callbacks describe user intent; a nil callback never fabricates success.
- Disabled controls cannot activate; pending controls remain readable but
  cannot submit.
- Icon-only controls keep the caller's label as their accessible name and
  tooltip.
- Field controls associate their label, description and error through MyGo's
  native field mechanism.
- Package upgrades remove obsolete code at the same time; no deprecated
  aliases, adapters or compatibility shims are kept.

## Upstream fidelity and limits

`references/heroui` and `references/heroui-native` are pinned submodules. Their
CSS, variant definitions and compound APIs are read for every component; see
`docs/references.md` for the update procedure. Web-only behavior is translated
only where native semantics allow it. Known limits are recorded in
`docs/components.md` instead of being implied by a component's presence, and
the concrete MyGo/CSS differences are in `docs/mygo-vs-css.md`.

The module depends on MyGo v0.6.3 and mygoicons v0.2.1. Module metadata contains
no local replacements or workspace requirements.

## Verification

Consumers must build using only this module and resolvable dependencies.
Meaningful tests cover editing, focus, disabled and invalid states, selection,
dates, filtering, overlays, motion and narrow or enlarged layouts. Headless
rendering establishes paint and native event routing; OS gestures, GPU
presentation and screen-reader integration remain separate platform checks.

The gallery devotes one page to each component and follows HeroUI's docs
anatomy: Usage, Examples, API Reference and Theme tokens, with a right-hand
"On this page" rail. `ComponentPage` carries one full `Code` snippet, an
optional `Import` override for the compact Usage bar (empty derives the block
from the snippet's package qualifiers, always naming the module root), and
named `ComponentExample` frames; there is no page-level demo fallback. Each
named example renders its title outside its own bordered panel; the page
snippet renders once after the examples (inside the panel for a single
example, or in a shared Code panel for several).
