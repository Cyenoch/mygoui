# mygoui

Native [MyGo](https://github.com/egoist/mygo) components with a [shadcn/ui](https://ui.shadcn.com) inspired visual style. Typed Go options, caller-owned state, semantic themes and native keyboard, editing and accessibility behavior.

The gallery contains **66 component specimens**, from buttons and forms to dialogs, tables, charts, dates, code blocks and conversation surfaces. Shadcn informs the appearance; the APIs and implementation follow Go and MyGo.

Requires **Go 1.27.1+** and **MyGo v0.6.0+**. Icons use the published `github.com/Cyenoch/mygoicons v0.1.0` module. Builds need no JavaScript, asset server or sibling checkout.

## Try it

From this checkout:

```sh
go run ./examples/gallery
go run ./examples/basic
```

The gallery provides search, independent interactive specimens, actual compiled Go source, Neutral/Blue presets, light/dark appearance, high contrast and 200% text previews. Its sample state stays in memory.

The source repository is [github.com/Cyenoch/mygoui](https://github.com/Cyenoch/mygoui). **No version has been released yet.** Once a maintainer publishes a release, consumers can install it with `go get github.com/Cyenoch/mygoui@<version>`. For local development, use a Go workspace:

```sh
# In a directory containing your application's module:
go work init . ../mygoui
```

If you already have a workspace, use `go work use ../mygoui` instead. A workspace is a consumer development choice; the library itself has no replacements or workspace requirements.

## Use components

Apply the theme once at the root of your MyGo view, then build native controls:

```go
import (
    "github.com/egoist/mygo/ui"
    "github.com/Cyenoch/mygoui"
    "github.com/Cyenoch/mygoui/form"
    "github.com/Cyenoch/mygoui/theme"
)

type state struct {
    Name string
    Saved string
}

func (s *state) View(c *ui.Context) {
    theme.Apply(c)
    colors := theme.Colors(c)
    ui.Column(c).Fill().Padding(24).Gap(12).
        Background(colors.Background).TextColor(colors.Foreground).Children(func() {
        form.Input(c.Key("name"), &s.Name, form.InputOptions{
            ControlOptions: form.ControlOptions{Label: "Display name"},
            Placeholder: "Your name",
        })
        mygoui.Button(c, mygoui.ButtonOptions{
            Label: "Save", OnClick: func() { s.Saved = s.Name },
        })
        ui.Text(c, s.Saved)
    })
}
```

Keep model data, widget state and `ui.Handle` across builds, separately for each window. Elements returned by constructors last for one build pass. Run slow work outside the view and publish its results with `Window.Update`.

## Packages

| Import suffix | Components |
| --- | --- |
| root `mygoui` | Button, Badge, Card, Alert, Avatar, Toggle, Progress, Spinner, Skeleton, typography and basic layout surfaces |
| `form` | Inputs, fields, labels, input groups, OTP, checkbox, switch, radio group, sliders and native select |
| `navigation` | Accordion, collapsible, breadcrumb, tabs, select, combobox, command palette, pagination, navigation menu and sidebar |
| `overlay` | Dialog, alert dialog, sheet, drawer, popover, tooltip, hover card, menus and toast |
| `data` | Virtualized Table and DataTable with local filtering, sorting, pagination and column visibility |
| `calendar` | Gregorian civil dates, single/multiple/range calendar and date picker |
| `chart` | Native bar, line, area, pie, radar and radial charts |
| `conversation` | Message, bubble, marker, attachment, message scroller and questionnaire |
| `codeblock`, `trace` | Exact source presentation/copy, syntax highlighting, line treatments and retained disclosures |
| `theme`, `locale` | Semantic colors, metrics, typography, accessibility preferences and namespaced translations |

See [the component catalog and limitations](docs/components.md) and [the design contract](docs/design.md).

## Theme and language

`theme.Apply(c)` installs Neutral using the window's current appearance and OS accessibility preferences. `theme.ApplyPreset(c, theme.Blue)` selects Blue. The theme honors OS accent, high contrast, text scale and reduced motion. Set application-wide font policy with `theme.SetPreferences`, then invalidate affected windows; retain widget state per window.

Library messages include English and Simplified Chinese and use MyGo's native library catalog. The host chooses and registers application languages. Override component strings in the host catalog with keys such as `mygoui.componentClose`. Importing the library does not register application languages or initialize an application.

Install `overlay.Toaster(c)` once per window when using notifications. Emit toasts from actions and job completions, with the caller retaining `overlay.ToastState`; rendering a view repeatedly must not enqueue notifications repeatedly.

## MyGo 0.6

Spinner uses the maintained `plugins/loading` arc, including visibility-aware paint scheduling and reduced-motion poses. Tables expose native cell selection through `state.List.CellSelect`, row dragging through `state.List.DragSelect`, and row headers through `Column.RowHeader`. The gallery's table examples exercise selection modes. Native keyed controls retain MyGo's incremental layout, partial redraw, wheel scrolling and platform menu behavior.

The [upgrade notes](docs/mygo-0.6.md) record the reviewed release changes, integration choices and checks.

## Development

```sh
make check
make race
golangci-lint run ./...
make wasm
make workflows # requires actionlint
```

Ordinary tests use MyGo's headless renderer and native event routing. They do not open windows. Linux needs Pango (`libpangocairo-1.0-0`) for software text rendering. Race tests require a C toolchain. CI configures these prerequisites and checks Linux, macOS and Windows, plus WebAssembly builds, MyGo lifetime analysis, lint and vulnerability scanning.

See [CONTRIBUTING.md](CONTRIBUTING.md) for validation, public API changes, dependency updates and release steps. Headless evidence and outstanding platform checks are recorded in [VALIDATION.md](VALIDATION.md).

## License

Code is [MIT](LICENSE). Shadcn is a visual reference; this project is independently maintained. Dependency and reference notices are retained in [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md).
