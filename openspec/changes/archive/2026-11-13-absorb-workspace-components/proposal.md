## Why

Desktop applications built on the library still had to write their own tree,
document tabs, menu bar, status bar, viewers for JSON, logs and diffs, masked
and tagged inputs, settings rows, guided tours and bulk-action bars. Loading
was shown differently by each component, and most state changes snapped.
A hands-on pass over the first versions also found them thin: missing
filters, gestures, keyboard paths and exit motion.

## What Changes

- Root: `Descriptions`, `Banner`, `SettingsGroup`/`SettingsRow`, `Spoiler`,
  `LoadingOverlay`, `Statistic`, `StatusDot`, `Timeline`, `NumberFlip`.
- `list.Tree`; `navigation`: `DocumentTabs`, `StatusBar`, `OverflowList`,
  `Menubar`, `NavigationStack`, `SplitButton`, `Carousel`.
- `form`: `TagsInput`, `KeyValueEditor`, `MaskInput`, `PathField`,
  `PasswordField`, card choices, range slider painting.
- `content`: `JSONViewer`, `LogViewer`, `DiffViewer`.
- `overlay`: `Tour`, `TeachingTip`, `ActionBar`, `Popconfirm`,
  `ShortcutsDialog`, `AboutDialog`.
- One loading presentation (`internal/busy`) and shared state motion
  (`internal/motion`) across existing and new components.
- A state never changes a text's font weight.
- Gallery: finished compositions move to "Viewers" and "Patterns" groups;
  the document column is centred.

## Impact

New files in the packages above; `overlay` now imports `navigation`, `form`
and `content`, and `navigation` imports `form`. `Button` pending no longer
dims, the range slider paints its own track, `NavigationStack` has no Escape
path. No new module dependency.
