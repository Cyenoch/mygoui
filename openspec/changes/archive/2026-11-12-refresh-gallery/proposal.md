# Proposal

## Why

The gallery stacks two header bars, keeps components behind collapsed groups,
misaligns titles, view tabs and content, and inserts appearance and compact
navigation rows that push the page around. Its first page is a plain category
grid. The browser build exits at startup because bundled font registration is
unavailable there, and no workflow publishes the gallery to GitHub Pages or
packages it for download.

## What Changes

- One title bar row beside a full-height sidebar; the sidebar's top row
  holds the native window controls and brand. The row is at least 52 DIPs and
  grows with text scale and zoom instead of letting controls overflow it. Desktop title-bar drag regions
  and control insets are kept; the browser build has none.
- Flat grouped navigation with every component visible, search with a
  shortcut hint, and a compact navigation sheet.
- Component pages lead with title, description and the view switch inside the
  content column; all chrome shares one gutter.
- Appearance controls move into a popover; copy actions keep their width.
- A landing page as the home view: hero, install command, live component
  showcase, catalog cards and download links for each platform.
- The browser build falls back to system fonts instead of exiting.
- A Pages workflow publishes the WebAssembly gallery; a release workflow
  builds macOS, Windows, Linux and web packages with stable asset names and
  attaches them to the tag's GitHub release.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `gallery-workbench`: navigation layout, landing page, layout stability,
  browser build and distribution.

## Impact

`internal/gallery`, `examples/gallery`, `.github/workflows`, `mygo.json`,
README, CONTRIBUTING, the gallery surface brief and DESIGN.md. No library API
change.
