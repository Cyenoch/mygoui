# Tasks

## 1. Shell

- [x] 1.1 Full-height sidebar with brand row reserving the leading window controls
- [x] 1.2 Page title bar with toggle, breadcrumb, GitHub and Appearance, clear of window controls; stacked below the band when they cannot fit
- [x] 1.3 Flat grouped navigation and search with shortcut hint
- [x] 1.4 Appearance popover and compact navigation sheet

## 2. Pages

- [x] 2.1 Page title, description and view switch in the content column
- [x] 2.2 Width-stable copy feedback and consistent demo panels
- [x] 2.3 Landing page with install command, downloads, live preview and catalog

## 3. Distribution

- [x] 3.1 System font fallback for the browser build
- [x] 3.2 Pages workflow
- [x] 3.3 Release workflow with stable asset names

## 4. Verification and docs

- [x] 4.1 Update gallery tests to the new contracts
- [x] 4.2 Update README, CONTRIBUTING, DESIGN.md, surface brief and CHANGELOG
- [x] 4.3 Run make check, golangci-lint, make wasm, make spec-check and actionlint; record results

## 5. Visual acceptance

Findings from a review of wide/compact, light/dark, native-inset captures
and user reports; each is covered by a test where the behavior is observable.

- [x] 5.1 Center the search shortcut hint in the field (Kbd set `AlignSelf(Start)`)
- [x] 5.2 Date and range picker triggers sit 4 DIP from the field edge, as coss input-group buttons
- [x] 5.3 Date demo readouts keep a fixed width so a changing range no longer moves the centered calendar (`TestRangeCalendarSummaryDoesNotShiftGrid` now renders inside the example panel); drop the duplicated range hint
- [x] 5.4 Calendar year-picker toggle uses a down chevron instead of a second right chevron
- [x] 5.5 Landing preview cards share one height; Link centers its icon; footer link uses an external-link icon
- [x] 5.6 Sidebar content fades at its scroll edges
- [x] 5.7 Slider mark labels sit under their ticks (`TestSliderMarksSitInTrackRow`)
- [x] 5.8 ComboBox and Autocomplete stay closed after an outside click (`TestComboboxesOutsideClickClosesAndStaysClosed`)
- [x] 5.9 Overlapping AvatarGroup members are opaque (`TestAvatarGroupOverlapDoesNotMixFills`)
