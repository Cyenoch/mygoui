# Design

## Context

The gallery is a consumer of the library and a native MyGo application, so
its shell uses library components and native window chrome rather than a web
layout. Before this change it had a 52 DIP brand bar above a separate toolbar
row (about 104 DIPs of chrome), collapsible navigation groups, inline
appearance and compact-navigation rows that pushed the page down, and
"Copied" labels that changed button widths.

## Decisions

### One title bar row

The sidebar runs the full window height. Its top row holds the brand and
reserves `TitleBar().Left` (macOS traffic lights). The page column's title bar
holds the navigation toggle, the breadcrumb, GitHub and Appearance, and
reserves `TitleBar().Right` (Windows/Linux captions), plus `Left` when the
sidebar is closed or the window is compact. Both rows share `headerHeight`, so
their bottom edges line up.

The row is `max(52 / zoom, native band, control + 16)`. When the three
controls cannot fit between the window controls (very large text in a narrow
window), they move below the native band, as the brand previously did. A
fixed 52 DIP row could not hold zoomed or enlarged controls, so the spec now
promises clearance from the native controls rather than a fixed height.

The navigation toggle always lives in the page title bar. Moving it between
the sidebar and the page title bar when the sidebar closes lost keyboard
focus.

`mygoui.Button` aligns itself to the start of its row (so it does not stretch
in columns); title-bar buttons set `AlignSelf(Center)`.

### Flat navigation, overlays for transient controls

Groups are headings, not toggles; every page is visible and the list never
changes height while browsing. `OpenGroup` is gone.

Appearance opens in `overlay.Popover` (mode, font, radius as
`navigation.SegmentedControl`, primary swatches). Compact navigation opens in
a left `overlay.Sheet` whose header holds the brand and search and whose body
scrolls the list. The sheet covers the title band, so its brand row reserves
`TitleBar().Left`. Neither moves the page.

### Page anatomy

The content column is left-aligned with one gutter (32 DIPs, 16 compact).
Every page except home starts with its title (H2) and description, then the
component view switch, so the title, tabs and examples share one left edge.
Demo panels have a 160 DIP minimum height and center their content on both
axes.

Copy actions keep their label slot at the wider of the label and "Copied";
`commandBar` (mono text plus copy button) serves the usage import and the
install command.

### Landing page

Home has a hero (badge, H1, lead, Browse and Customize actions, the install
command), download cards, a two-card live preview built from library
components, the catalog and a footer. Download cards open
`releases/latest/download/<stable name>`; Web opens the Pages URL.

### Distribution

`fonts.Register` fails in browsers by design, so `examples/gallery` treats
that failure as fatal only outside `js`. `State.SystemFontOnly` then hides the
font choice.

Pages builds `go tool mygo build -platform web` on Ubuntu and deploys
`dist/web` with the official Pages actions. Its output uses relative URLs,
so it works under `/mygoui/`.

Release runs on a macOS runner because the disk image needs macOS; MyGo
cross-compiles Windows and Linux without cgo. It writes the tag version into
`mygo.json`, builds, copies outputs to stable names
(`coss-ui-gallery-{macos-arm64.dmg,windows-amd64.exe,linux-amd64.tar.gz,linux-amd64.deb,linux-amd64.rpm,web.zip}`)
and uploads with `gh release upload --clobber`, creating the release if the
tag has none. Builds are ad-hoc signed and not notarized; no signing secret is
introduced.

## Risks

- Unsigned macOS and Windows builds trigger Gatekeeper and SmartScreen
  warnings.
- Pages must be enabled once with "GitHub Actions" as the source.
- The workflows were linted with actionlint and the asset naming was
  dry-run against a local build; neither has run on GitHub.
