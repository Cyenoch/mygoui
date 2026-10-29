# MyGo UI

A native Go component library for [MyGo](https://github.com/egoist/mygo),
reconstructing the component catalog and visual language of
[coss ui](https://coss.com/ui/docs).

Typed options, explicit application state, semantic light/dark colors, and
native editing, selection, focus and scrolling. Components compose through Go
callbacks and typed collections. Applications run without a JavaScript runtime
or WebView.

The `coss` branch replaces the previous design and tokens. The
[capability matrix](docs/components.md) maps every component found across the
pinned coss catalogs, plus their two documented hooks. A mapping records the
native implementation and its limits; it does not certify visual or behavioral
parity. Current verification belongs in [VALIDATION.md](VALIDATION.md).

The module targets **Go 1.27.3** and pins **MyGo v0.7.2** and
**mygoicons v0.2.1**.

Development uses [OpenSpec and Impeccable](docs/agent-workflow.md) for
requirements, scoped changes and native design review. Start with
`make agents`; these tools are optional development dependencies.

## Try it

```sh
go run ./examples/basic
go run ./examples/gallery

make dev
make app
```

The gallery is a native component workbench with a custom title bar, grouped
sidebar, Cmd/Ctrl+K search, appearance controls, and Preview / Go / API / Tokens
views. It reserves the platform's actual window-control insets rather than
drawing replacement window buttons. Examples and snippets are consumers of the
public API; the capability matrix records native adaptations and limits.

For external native-app tests, build the Gallery with `-tags e2e` and set
`MYGO_AUTOMATION` to a socket path in a private directory. Normal builds do not
import the automation plugin. See [MyGo integration](docs/mygo-0.7.md).

## Use

Apply the theme at the start of each view build. Keep values and widget state
with the window's model:

```go
import (
    "github.com/Cyenoch/mygoui"
    "github.com/Cyenoch/mygoui/form"
    "github.com/Cyenoch/mygoui/theme"
    "github.com/egoist/mygo/ui"
)

type settings struct {
    name  string
    saved bool
}

func (s *settings) View(c *ui.Context) {
    theme.Apply(c)
    mygoui.Card(c, mygoui.CardOptions{}, func() {
        mygoui.CardHeader(c, func() {
            mygoui.CardTitle(c, "Project")
        })
        mygoui.CardContent(c, func() {
            form.Input(c, &s.name, form.InputOptions{
                Label: "Name", FullWidth: true,
            })
        })
        mygoui.CardFooter(c, func() {
            mygoui.Button(c, mygoui.ButtonOptions{
                Label: "Save",
                OnClick: func() { s.saved = true },
            })
        })
    })
}
```

`ui.Element`, custom parts and `*ui.Context` belong to the current build.
Retain model values, component state and `ui.Handle` instead. Use native
services for clipboard and URL actions; publish worker results through
`Window.Update`. See [ownership rules](docs/design.md).

## Customize

The canonical colors are coss roles such as `Primary`, `Secondary`,
`MutedBackground`, `MutedForeground`, `Accent`, `Card`, `Popover`,
`Destructive`, `Input` and `Ring`. Each has an independent override in
`theme.PalettePatch`; light and dark appearances have separate patches.

```go
primary := ui.RGB(30, 90, 180)
onPrimary := ui.RGB(255, 255, 255)
theme.Configure(theme.Preferences{
    Radius: 10,
    Palette: theme.PaletteOverrides{
        Light: &theme.PalettePatch{
            Primary: &primary, PrimaryForeground: &onPrimary,
        },
    },
})
```

Call `Configure` before building views, or use `UpdatePreferences` on the UI
thread and invalidate the affected windows. Read colors live with
`theme.Colors(c)`. See [theme tokens and provenance](docs/theme.md).

Body, heading and mono font families are independently configurable. Defaults
use native system fonts, so text metrics differ from the coss site. Applications
can opt into the embedded, independently OFL-licensed Cal Sans faces:

```go
if err := fonts.Register(); err != nil {
    return err
}
theme.UpdatePreferences(func(p *theme.Preferences) {
    p.FontFamily = fonts.Family
    p.HeadingFontFamily = fonts.Family
})
```

Import `github.com/Cyenoch/mygoui/fonts` and call this during application
startup in a function returning an error. Registration is explicit and
idempotent; browser builds return an error. Mono keeps native fallbacks unless
the host independently registers another face. See
[font sources, registration and limits](docs/fonts.md). Configuring a family
name alone does not load that font.

## Scrolling

`ScrollArea` supports `ScrollVertical`, `ScrollHorizontal` and `ScrollBoth`.
`ScrollFade` can also decorate an existing viewport and its `ui.ScrollState`.
`MaskFade` fades the native subtree to transparency through a shader effect.
`ShadowFade` paints an edge shadow; `SurfaceFade` blends into the actual opaque
surface color. Effect layers have an 8192-pixel side limit and exclude embedded
webviews; GPU/platform checks remain separate from headless rendering.
See [scroll examples and limits](docs/scrollfade.md).

## References and license

The visual reference is coss ui commit
`3d8ca28e7d90123e9e0f44665944dc240b026ccb`.
Run `make refs` to check out its source at `references/coss`.
[Reference provenance](docs/references.md) identifies the MIT `apps/ui`
registry and documentation. The upstream monorepo has mixed licensing;
`packages/*` is outside that MIT scope, and its source and font assets must
not be copied as MIT material.

Project code retains its [MIT license](LICENSE). See
[third-party notices](THIRD_PARTY_NOTICES.md) for reference attribution,
dependencies and optional font licenses. This is an independent native
implementation.
