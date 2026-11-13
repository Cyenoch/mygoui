## Decisions

- **Revision instead of hashing options.** Detecting in-place edits would
  cost the walk the cache removes. The slice identity and length catch a
  replaced slice; `Revision` covers edits in place, the convention
  `data.TableState` already set.
- **Refilter into fresh slices.** A copied state, or a list built from the
  previous result, may still hold the old rows.
- **One path per menu opening.** Switching between the plain and the
  virtualized panel rebuilds the rows and drops focus, so the row count is
  read when the menu opens.
- **Fixed width for long menus.** A virtualized panel cannot measure its
  widest row; a width that changed while scrolling would be worse.
- **Round the viewport, not the panel.** MyGo rasterizes a mask when a square
  scroll clip meets a rounded clip under a changing transform. Rounding the
  viewport itself leaves one analytic clip. The scale-in of a desktop Dialog
  still moves one cached bitmap.
- **Shadow layer inside the viewport for the non-GPU mask fade**, so it
  follows the viewport's clip and sits under the scrollbars.

## Review findings (independent agent) and their fixes

- Command cleared its shortcut drawer after building; rows a list builds
  during layout then called nil. The drawer now stays set. Regression test.
- An open menu crossing 32 rows lost focus: path fixed per opening. Test.
- The sidebar culled a focused row scrolled out of view: it stays built. Test.
- Refiltering reused the backing array shared with state copies: fresh slices.
- The non-GPU mask fallback painted over rounded corners and scrollbars:
  moved into the viewport.

## Evidence

Browser (WebAssembly, local build): Bottom Drawer enter max 6.6 ms and exit
max 13.8 ms per frame, AlertDialog enter 11.1 ms and exit 11.5 ms, no frame
over 16 ms; rounded corners, close control and footer intact. Sidebar frame:
192.9 KB/694 allocations to 49.3 KB/327. Long menus in a flyout, RTL long
menus, non-GPU rendering on real hardware and screen readers are not
verified.
