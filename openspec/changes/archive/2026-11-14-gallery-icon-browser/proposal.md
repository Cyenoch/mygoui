## Why

The library draws every glyph from mygoicons, yet the gallery gives no way
to find an icon: a reader has to leave for another tool to learn that the
trash can is `lucide.Trash2`. The gallery's own navigation and page chrome
also use few icons, so it undersells what the pairing looks like.

## What Changes

- Add an Icons page to the gallery: search the Lucide collection by Go
  name or Iconify ID, see each drawing once, and copy its accessor.
- Use more icons in the gallery's own chrome and examples where they name a
  subject (navigation groups, page actions).
- Raise mygoicons to v0.3.0, which adds the `catalog` package the page
  reads and Lucide 1.2.154.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `gallery-workbench`: adds the requirement that icons can be found and
  copied in the gallery.

## Impact

`internal/gallery`, `go.mod` (mygoicons v0.3.0), `README.md`,
`CONTRIBUTING.md` and `THIRD_PARTY_NOTICES.md` where they name the pinned
version, `CHANGELOG.md`. No public
API of the library changes. The gallery executable grows by the collections
it lists.
