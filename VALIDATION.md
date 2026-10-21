# Validation

Local validation used macOS arm64, Go 1.27.2, MyGo v0.6.0 and mygoicons v0.1.0. Module metadata targets Go 1.27.1. The public repository is `github.com/Cyenoch/mygoui`; no release has been tagged.

## Completed checks

- The complete ordinary suite passed with `CGO_ENABLED=0 go test -shuffle=on ./... -count=1 -timeout=180s`: 198 named tests plus subtests and an executable button example, across 17 packages (14 with tests).
- `go build ./...`, `go vet ./...`, `go tool mygo vet .`, `golangci-lint run ./...`, Go formatting, module tidy and `go mod verify` passed. Lint reported zero issues. The checked-elements migration preview found no required changes.
- CGO-disabled builds passed for darwin/arm64, linux/amd64, windows/amd64 and js/wasm. The first public CI run passed macOS and Windows execution; Linux setup corrections are described below.
- An independent application installed a locally generated Go module archive through a temporary module proxy with `GOWORK=off`, no replace directives and published dependencies. A native headless input accepted `Ada`, and the Save button delivered the exact value to caller state. This proves archive/consumer isolation, not public registry availability.
- Regression checks passed for native table rectangular cell selection, row-header whole-row selection, select-all, wheel scrolling with a fixed header, Spinner animation without rebuilding, reduced-motion settling and namespaced host translation overrides.
- Six headless gallery captures were inspected: dark Button, light Sidebar, dark Message, light Chart, dark Data Table and a 390-DIP light Input at font size 32. The temporary capture harness was removed; the permanent gallery suite exercises all 66 defaults at three widths plus independent light/dark, enlarged text, source and action cases.
- Actionlint accepted the workflow. Repository scans found no application-specific source dependencies or local replacements. Optional development skills retain their references and upstream metadata.
- `govulncheck v1.9.0 -show verbose ./...` found **zero reachable vulnerabilities**. It also reported package-level standard-library advisory [GO-2026-6809](https://pkg.go.dev/vuln/GO-2026-6809) in the local Go 1.27.2 crypto/x509 implementation, fixed in Go 1.27.3; no affected call path was reachable in this library. Application maintainers should use a patched Go toolchain.

## Race checks

The first race run timed out in a pixel assertion that acquired a new screenshot for every pixel. The assertion now acquires one frame before scanning. The affected conversation package then passed under race detection in about two seconds; its ordinary runtime fell from about 16 seconds to under one second. This correction changes test acquisition cost, not component behavior.

The final complete `CGO_ENABLED=1 go test -race -shuffle=on ./... -count=1 -timeout=300s` passed in every tested package. Conversation completed in 1.951 seconds and the largest overlay suite in 143.648 seconds, with no data-race reports.

## Cleanup and CI validation

- Removed the gallery's duplicate implementation inventory and unreachable planned-component branches, including their English/Chinese messages. All 66 specimens still pass the existing three-width gallery tests.
- The two scroll-viewport pixel tests now capture one image per state, preserving every assertion. Their ordinary package time fell from 2.022 seconds to 0.032 seconds on this machine; command wall time fell from 2.66 seconds to 0.44 seconds. These are local samples with warm caches, not GitHub runner measurements.
- `make check` passed in 37.49 seconds. The subsequent `make race` passed in 134.63 seconds, with no race reports. Lint, actionlint v1.7.12 with ShellCheck, WebAssembly, module tidy/integrity and the reachable-vulnerability scan also passed; module metadata remained unchanged.
- CI retains all three OS test/build legs, runs Linux race checks independently and downloads the pinned golangci-lint binary with analysis caching. Module metadata, MyGo lifetimes, WebAssembly, workflows and vulnerabilities are checked in the Linux analysis job. The additional runner setup may increase total runner minutes while shortening the test/race critical path. There is no GitHub baseline for the previous workflow, so local speedups cannot establish total runner-minute savings.

- The [first public CI run](https://github.com/Cyenoch/mygoui/actions/runs/39572820116) exposed Go 1.27.1 standard-library vulnerabilities and missing Linux CJK/emoji fonts. CI now resolves the latest Go 1.27 patch, and Linux test/race jobs install the fixture fonts. Go 1.27.3 vulnerability scanning reports no vulnerabilities. The Linux Unicode pixel failures were reproduced in a Debian container and passed after adding the fonts; font sizing now uses ordinary text to avoid narrow-glyph hinting plateaus, with width/height and family assertions retained.

## Platform limits

No visible application window was opened and no personal application configuration was read or changed. Headless captures verify renderer output and native event routing. Actual OS menu appearance, screen-reader operation, platform IME/gesture behavior and GPU partial-presentation performance remain separate manual platform checks. Component feature limits, including Drawer gestures, arbitrary slider thumbs and alternate calendar arithmetic, are documented in `docs/components.md`.
