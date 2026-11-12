# Tasks

## 1. Shared foundations

- [x] 1.1 Overlay elevation uses shadow-lg/5 plus appearance edge; add `TooltipShadows`
- [x] 1.2 Move popup surface and motion into `internal/anchored`; remove per-package copies and offset defaults
- [x] 1.3 Share overlay warm-up between overlay and navigation

## 2. Duplicate-caused defects

- [x] 2.1 Slider value text uses the number field's decimal count
- [x] 2.2 Error message implies invalid in shared field state; one error-text helper
- [x] 2.3 InputOTP uses the shared field face
- [x] 2.4 One status glyph mapping for alert, alert dialog and toast

## 3. Native replacements

- [x] 3.1 `Painter.FillRadii` for color slider and calendar bands
- [x] 3.2 `Color.Mix` for table row tint
- [x] 3.3 Context menu items map Description, Shortcut and danger natively

## 4. Documentation and verification

- [x] 4.1 Update docs/theme.md, docs/components.md, DESIGN.md and CHANGELOG
- [x] 4.2 Focused `ui.NewTester` regressions for the changed contracts
- [x] 4.3 `make check`, `golangci-lint run ./...`, `make wasm`, `make spec-check`; record in VALIDATION.md
