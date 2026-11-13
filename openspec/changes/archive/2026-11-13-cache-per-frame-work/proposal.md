## Why

A review of every component left eight per-frame costs open: ComboBox,
Autocomplete and Command refilter all their options every frame; Select walks
its options and allocates 26 key handlers per frame; modal panels rasterize a
rounded clip on the CPU while they move; menus build every row; shader
effects are interpreted per pixel where there is no GPU; and the gallery
sidebar builds all its rows.

## What Changes

- ComboBox, Autocomplete and Command keep their filtered rows until the
  query, the option slices or a new caller-owned `Revision` change.
  **BREAKING**: an in-place edit of an option slice shows only after the
  caller bumps `Revision`, as `data.TableState.Revision` already requires.
- A Select row without a key is never the chosen row.
- Menu panels with 32 or more rows build only the rows in view, at a fixed
  width; the path is chosen once per opening.
- Modal panels round their scroll viewport instead of clipping the panel, so
  sliding and fading panels need no CPU clip mask.
- Without a GPU, SkeletonGroup rests static and a mask ScrollFade shows the
  shadow treatment.
- The gallery sidebar builds the rows near its viewport and the focused row.

## Impact

`form` (ComboBox, Autocomplete, Select), `navigation` (Command, menus),
`overlay` panel chrome, root `SkeletonGroup` and `ScrollFade`,
`internal/gallery` navigation; `docs/components.md`, `docs/scrollfade.md`.
