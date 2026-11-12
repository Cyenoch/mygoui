# Tasks

## 1. Foundation

- [x] 1.1 `mygoui.FlyoutHost`; the presenter fits content within component bounds (no new `overlay.FlyoutOptions` field)
- [x] 1.2 Shared presenter in `internal/anchored`: open from the action, re-anchor, refresh, close sync, fallback

## 2. Components

- [x] 2.1 Popover
- [x] 2.2 Dropdown (top level)
- [x] 2.3 Select
- [x] 2.4 DatePicker and DateRangePicker
- [x] 2.5 ColorPicker

## 3. Verification and docs

- [x] 3.1 Tester coverage: unchanged in-window behavior without a host and the window-less host fallback; flyout context, close sync and dismissal are verified in a real macOS window (headless Tester creates no windows)
- [x] 3.2 Gallery: Popover example and the appearance panel opt in; other examples stay in-window
- [x] 3.3 docs/flyout.md, docs/components.md, CHANGELOG; record real-window evidence and limits in VALIDATION.md
