# Product

<!-- impeccable:product-schema 1 -->

## Platform

adaptive

Native desktop applications on macOS, Windows and Linux, using MyGo's
platform window controls, menus, focus and services. This selector enables
Impeccable's native guidance; it does not add iOS or Android targets.
Desktop MyGo behavior takes precedence over phone-specific platform examples.
The WebAssembly build is a compatibility check, not a separate web product.

## Purpose

A complete native UI solution for MyGo applications, without a browser or
JavaScript runtime. It has four pillars, each with one owning record:

| Pillar | Owner | State |
| --- | --- | --- |
| Theme tokens | `theme/`, `docs/theme.md` | Colour, radius, type, elevation, the spacing unit, motion and effect tokens are public, and `theme.Theme` supplies a whole theme as one value (`theme.Use`). Timings peculiar to one component stay in the component |
| Component library | `docs/components.md` | Delivered; the matrix records subsets and gaps |
| Design and behavior reference | `DESIGN.md`, `docs/design.md`, `docs/behavior/` | Visual rules are recorded. The per-component behavior contract covers every family, `layout` included (`component-behavior-reference`) |
| Agent skills | `skills/`, `.agents/skills/` | Five published `mygoui-*` skills cover components, theme, design and behavior, motion and shader effects (`publish-agent-skills`). The effects skill leads with the public `effect` family: Mask, Shimmer and Reveal (`native-effects`) |

## Users And Stack

Go application authors using MyGo v0.9.0. Caller-owned state, typed options and
native focus, editing, selection and accessibility are binding constraints.

## Visual Authority

Root `DESIGN.md` is the authority. coss ui (https://coss.com/ui/docs and
the corresponding cosscom/ui source) is provenance: where the design language
came from, recorded in `docs/references.md`. A later design language can
replace the record without component changes.

## Capabilities

Native charts use Bklit's public chart catalog as capability reference and the
design language recorded in `DESIGN.md`. The chart matrix records supported native subsets;
the library does not embed Studio, a browser runtime or background data services.

Customizable semantic theme tokens in both appearances; independently composed
components; a generic scroll fade using shadow, surface blending or a native
subtree alpha mask; a public shader effect family (`effect`: Mask, Shimmer,
Reveal) with a designed fallback where no GPU runs it. Web-only
or unavailable toolkit features must be documented rather than simulated as
fully supported. GPU, OS gestures and screen-reader results require separate
evidence from headless testing.

## Native Gallery

The gallery is a native component workbench with a completely redesigned
application shell, custom titlebar and sidebar. Component discovery, live
examples, Go source and theme controls form one coherent workflow. Use real
platform window controls and draggable titlebar regions, preserve keyboard
navigation and native focus, and adapt the workspace to compact windows.
Its appearance follows `DESIGN.md`.

## Development Context

OpenSpec records current requirements and scoped change proposals.
Impeccable supplies design discovery, critique, audit, refinement and finish
playbooks. Go tokens and native controls implement the design language
that root `DESIGN.md` records; `DESIGN.md` wins on a mismatch, which a change
settles. Browser-only live iteration
and HTML/CSS detector results do not establish native UI quality.
