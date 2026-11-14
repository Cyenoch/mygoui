## Context

mygoicons v0.3.0 exports `catalog`, which lists a collection's accessors
with alias and color information. Listing a collection keeps it in the
executable.

## Decisions

- **One documentation page, Foundations group.** The page is a
  `ComponentPage` like every other: a browser example first, then usage and
  sizing examples, so it carries the same Preview / Code views.
- **A virtual list of rows of cells.** `ui.List` builds only visible rows;
  each row holds as many fixed-size cells as the measured width allows. A
  grid of two thousand cells built at once is the alternative and is not
  affordable per frame.
- **Aliases resolve to their drawing.** Search matches an alias's name and
  ID but shows the original, so one drawing never appears twice.
- **Lucide only.** Measured in the browser build: 27.7 MB before, 30.3 MB
  with the Lucide catalog, 43.7 MB with Twemoji as well. Twemoji is also
  CC-BY-4.0, which the gallery's distributed builds would have to attribute;
  the library imports no Twemoji today and the gallery keeps to that.

## Risks

- The page keeps all of Lucide in the gallery executable (+2.6 MB in the
  browser build). Applications are unaffected: they link only the icons
  they name.
