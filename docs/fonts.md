# Fonts and native fidelity

The library's default body and heading faces are native system fallbacks;
mono uses native monospace fallbacks. This is a visible fidelity limit:
line breaks, control label widths and heading shapes can differ from the
coss docs site even when color and geometry match.

The pinned MIT
[styling guide](https://github.com/cosscom/ui/blob/3d8ca28e7d90123e9e0f44665944dc240b026ccb/apps/ui/content/docs/%28root%29/styling.mdx)
describes Cal Sans 2.0 for the site's body/headings. The downloaded site
snapshot links `CalSansVF` upright/italic and `PaperMono_Regular` WOFF2 files.
Those observed web asset names do not establish independent redistribution
rights or byte equivalence to another release.

Do not extract fonts from coss `packages/*`. Independent font projects
provide OFL 1.1 sources.

## Cal Sans

Official source:
[`calcom/sans`](https://github.com/calcom/sans/tree/9fe03d44f00806aadeb6d3653e37dce358554de0)
at `9fe03d44f00806aadeb6d3653e37dce358554de0`.

| File under that revision | Tree size | Git blob |
| --- | --- | --- |
| `fonts/calsans-cossui/CalSansVF.ttf` | 529,412 bytes | `0b2fba03354208dd01a4f86c17fc6ee51aa25d3f` |
| `fonts/calsans-cossui/CalSansVF-Italic.ttf` | 560,860 bytes | `829c0ed068a3f2fe8ba06fa9ac9e772d911a1a2b` |
| `OFL.txt` | 4,381 bytes | `52135553ce36d47c868d0dd93de3accaf60fea76` |

The official
[OFL.txt](https://github.com/calcom/sans/blob/9fe03d44f00806aadeb6d3653e37dce358554de0/OFL.txt)
declares SIL Open Font License 1.1 and credits the Cal Sans Project Authors,
2021. The repository README identifies Cal Sans v2. At this exact revision
`package.json` says npm version `2.1.1`; these two font blobs also occur at
tag `v2.011` (`706884a0420ad801f3c21ebe129059e97ae7180d`).
Use the commit/path when recording provenance instead of calling the files
an exact copy of the site's "2.0" assets.

The optional `github.com/Cyenoch/mygoui/fonts` package now embeds these two
TrueType faces from the independent source. Its original notice is retained
at [fonts/OFL.txt](../fonts/OFL.txt). The font package is separate from the
default theme; importing the library does not register it.

Tree sizes, blob IDs and license text were inspected on 2026-10-27. The two
embedded binaries were subsequently downloaded from that independent source;
tests exercise explicit registration and compare regular/italic glyph advances
at weights 400 and 700 with independently registered original faces.
The gallery's headless captures and macOS app use this opt-in family. This does
not establish byte equivalence with the site's WOFF2 files or identical shaping
on every renderer.

## Paper Mono

Official independent source:
[`paper-design/paper-mono`](https://github.com/paper-design/paper-mono/tree/e6eaeceaef02e77e3db997711e07a16378de2bd7)
at `e6eaeceaef02e77e3db997711e07a16378de2bd7`.
`fonts/ttf/PaperMono-Regular.ttf` is 106,088 bytes in that tree;
`fonts/variable/PaperMono[wght].ttf` is 124,024 bytes.

The upstream
[LICENSE.txt](https://github.com/paper-design/paper-mono/blob/e6eaeceaef02e77e3db997711e07a16378de2bd7/LICENSE.txt)
declares OFL 1.1, credits Lost Coast Labs, Inc. (Paper Design), 2025,
and preserves its Geist Mono attribution to Vercel and basement.studio, 2023.
The separate
[OFL.txt](https://github.com/paper-design/paper-mono/blob/e6eaeceaef02e77e3db997711e07a16378de2bd7/OFL.txt)
credits the Paper-Mono.Git Project Authors, 2025. Preserve both notices with
any future application distribution. Paper Mono remains a host-selected
font; this investigation did not add it to the library.

## Registration and distribution

For embedded Cal Sans, call the optional package explicitly during startup:

```go
import (
    "fmt"

    "github.com/Cyenoch/mygoui/fonts"
    "github.com/Cyenoch/mygoui/theme"
)

func configureFonts() error {
    if err := fonts.Register(); err != nil {
        return fmt.Errorf("configure application fonts: %w", err)
    }
    theme.UpdatePreferences(func(p *theme.Preferences) {
        p.FontFamily = fonts.Family
        p.HeadingFontFamily = fonts.Family
    })
    return nil
}
```

`fonts.Family` is `"Cal Sans"`. `Register() error` registers regular and italic
faces once; subsequent calls return the same result. Failure to register the
italic face leaves the regular face registered. Propagate that error explicitly.
Browser builds return an error because MyGo's browser registration path cannot
select these styles or report loading failures reliably. A successful wasm
build is not successful browser font registration.

The basic and gallery application entry points explicitly call `Register`;
the gallery configures body/headings with `fonts.Family`. That application
choice does not change the theme's native-system default for other consumers.

For other host fonts, MyGo v0.7.1 exposes
`ui.RegisterFont(data []byte, family string) error` for TrueType/OpenType data.
Load/embed the selected files and register them explicitly; there is no
import-side font I/O or registration.

Once faces are registered under the chosen family names:

```go
theme.Configure(theme.Preferences{
    FontFamily: "Cal Sans",
    HeadingFontFamily: "Cal Sans",
    MonoFontFamily: "Paper Mono",
})
```

This second snippet configures family resolution; it does not register either font.
Use the family alias actually supplied to `ui.RegisterFont`. If only Cal Sans
is registered, leave `MonoFontFamily` empty to keep native monospace fallbacks.
Changing the configured family cannot supply missing weights or glyphs.

OFL permits embedding/bundling fonts with applications while retaining the
font copyright and license. The font stays OFL; the surrounding application
does not thereby become OFL. Preserve applicable reserved-name restrictions
when modifying a font, and do not sell the font by itself. These terms come
from the linked upstream licenses; the project's MIT license does not replace
them.

Native registration, fallback, weight/italic selection and line metrics need
separate checks on each target platform. Choosing the same family alone does
not establish matching optical sizing, custom variation axes or browser
rasterization.
