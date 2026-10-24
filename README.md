# MyGo UI

Native [HeroUI v3](https://heroui.com/docs/react/components) components for
[MyGo](https://github.com/egoist/mygo), written in Go.

Typed Go options, caller-owned state, HeroUI's semantic color system and
desktop metrics, and native MyGo keyboard, editing, selection, scrolling and
accessibility behavior. No JavaScript, WebView, CSS runtime or browser
compatibility layer is involved.

The library unifies the HeroUI component catalogs into one native Go component
set; shared components are implemented once, and native-only components are
included where they fit. It is not a React port: components compose through
body callbacks and explicit state, not providers, contexts or slots.

Requires **Go 1.27.1+** and **MyGo v0.6.3+**. Icons use the published
`github.com/Cyenoch/mygoicons` module.

## Try it

```sh
go run ./examples/basic
go run ./examples/gallery

# Live-reload workbench through the MyGo CLI (uses mygo.json at the root,
# and runs go run github.com/egoist/mygo/cmd/mygo@latest dev)
make dev

# Production app bundle of the gallery
make app
```

The MyGo UI gallery follows the HeroUI docs layout: a main column with a
sticky "On this page" rail (hidden below 1000 DIP), a page title and
description with Copy Markdown and Copy import actions, and a compact Usage
import bar whose block is derived from the packages the page's snippet names
(one line, or a parenthesized block for several). Each component page has an
Examples section of 2-5 named examples with live previews and no page-level
fallback; the page's full Go snippet renders once below them in a code area
with line numbers, syntax highlighting, copy, and 6-line collapse/expand with
a surface fade. Pages end with a wrapping API Reference table and theme-token
cards. The sidebar groups every component and uses ScrollShadow for overflow;
light/dark appearance and accent/radius customization stay live, and data
stays in memory.

## Use

Apply the theme once at the view root:

```go
import (
    "github.com/Cyenoch/mygoui"
    "github.com/Cyenoch/mygoui/theme"
    "github.com/egoist/mygo/ui"
)

func (s *state) View(c *ui.Context) {
    theme.Apply(c)
    mygoui.Card(c, mygoui.CardOptions{}, func() {
        mygoui.CardHeader(c, func() {
            mygoui.CardTitle(c, "Deploy")
            mygoui.CardDescription(c, "Ship the current commit to production.")
        })
        mygoui.CardFooter(c, func() {
            mygoui.Button(c, mygoui.ButtonOptions{
                Label: "Deploy", OnClick: func() { s.deploy() },
            })
            mygoui.Button(c, mygoui.ButtonOptions{
                Label: "Cancel", Variant: mygoui.ButtonGhost,
                OnClick: func() { s.cancelled = true },
            })
        })
    })
}
```

Customize the design system process-wide before views are built:

```go
theme.Configure(theme.Preferences{
    Accent:   ui.Oklch(0.62, 0.20, 300),
    Radius:   10,
    FontSize: 16,
})
```

See `docs/components.md` for the inventory and native limits,
`docs/design.md` for ownership, `docs/mygo-vs-css.md` for the layout and paint
differences from CSS/Tailwind, and `docs/references.md` for the pinned
upstream sources used to translate each component.

## Status

This checkout is a complete HeroUI-oriented rewrite; the previous
shadcn-inspired implementation is not part of the API. No version has been
released yet.
