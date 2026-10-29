# Design

## Context

See proposal.md. MyGo 0.7.2 already supplies `ui/chart` scales, axes, marks,
selection, curves, stacking and sectors. There is no existing mygoui chart
family. Existing `theme.Palette` already has Chart1 through Chart5.
Bklit is pinned at ec5a9705ee973f95fd0ef2d974e9d00ae3fda011.

## Goals / Non-Goals

**Goals:** Cover every documented chart and utility subject with useful typed
native APIs, real geometry, input and inspectable Gallery examples.

**Non-Goals:** Port React providers/hooks, CSS/DOM APIs, proprietary Studio,
export/video tools, or every upstream decorative rendering option.

## Decisions

- Use coherent Cartesian, polar, matrix, geographic, hierarchy and flow APIs,
  not an untyped map or a second graphics framework. Reuse MyGo geometry and
  native focus/selection before implementing missing layout algorithms.
- Errors belong at public data boundaries. No catch-all/default-success recovery.
- Caller-owned state contains values and handles, never contexts or elements.
  Live data is supplied by the host; maps are supplied polygons, not downloads.
- Extend the existing workbench in Operate mode. The established world and
  `.impeccable/surfaces/internal-gallery.md` remain authoritative. Keep its
  Preview/Code/API/Tokens structure; add chart and utility navigation groups.
- Coss chart tokens carry series identity. Border and MutedForeground carry
  axes; native coss controls carry legends, brushes and breadcrumbs. No new
  typography, palette or persistent explanatory UI is needed.
- Layout fills available native width; labels and height honor text scaling.
  No decorative entrance choreography; static rendering inherently respects
  reduced motion. Native interactions are checked through Tester.

## Risks / Trade-offs

- Broad catalog can hide shallow implementations -> publish a per-subject
  matrix and verify representative geometry and meaningful interactions.
- Custom geography and flow shapes lack native named marks -> expose chart
  selection descriptions and keyboard traversal; do not claim screen-reader
  platform verification without observing it.
- Upstream includes proprietary Studio -> reference only MIT chart sources and
  docs, preserve attribution, never copy Studio implementation.
- Headless output differs from GPU -> inspect native Tester captures and record
  actual OS/GPU validation separately; do not infer it from PNGs.
