## Why

The library lacks several generic desktop components that
`references/lobe-mygo` already proves on MyGo: copyable code and content,
inline editing, shortcut capture, file intake, rating, steps, reordering,
resizable and floating panels, and image display. The user selected sixteen
to absorb.

## What Changes

Each component is rebuilt with coss anatomy, this module's theme tokens,
typed options and caller-owned state. Lobe tokens, shadows, motion presets
and brand assets are not copied.

- New `content` family: `CopyButton`, `Snippet`, `Highlighter`, `Markdown`,
  `QRCode`.
- `form`: `EditableText`, `HotkeyInput`, `Upload`, `Rate`, `ImageSelect`.
- `navigation.Steps`, `list.SortableList`.
- `overlay`: `FloatingPanel`, `FloatingSheet`, `ImageViewer`.
- Root: `Image`, `Result`, `GuideCard`, `ResizablePanel`.
- Gallery pages, `docs/components.md` and the changelog cover them.

## Impact

New files in the packages above and a new `content` package. One new module
dependency, `rsc.io/qr` (BSD-3-Clause), used only by `content.QRCode`.
Existing APIs are unchanged.
