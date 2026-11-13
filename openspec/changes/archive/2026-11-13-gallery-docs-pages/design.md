## Context

Reference: the coss docs component page (screenshots supplied by the user:
Alert Dialog). The surface stays the native workbench shell; only the page
inside the content column changes. Mode: Read inside an Operate shell.

## Decisions

- **One document, not four views.** Sections in coss order. Tokens stay, as
  the last section, because theme roles are this library's customization
  surface. The first example is the lead demo; the rest sit under Examples.
- **Preview/Code per example.** A two-segment switch above each demo, coss
  style. Example state is window-owned and survives switching, as today.
- **Generated, by cost:**
  - *Example source*: generated at run time. The gallery embeds its own page
    files and parses them once; an example's code is its `Demo` body, or the
    body of the one page function it calls, wraps or is. No build step, no
    drift.
    Pages built in loops fall back to the page snippet.
  - *API reference*: generated at build time by `internal/apidoc` from the
    library source (stdlib `go/parser`, `go/doc`) into `api_gen.go`: for each
    exported struct its doc comment and fields with types and comments. A
    test regenerates and compares, so a stale file fails CI. A page lists
    the structs its usage snippet names and describes a field of (every
    options struct it names when it describes none), so no page keeps a
    type list. A `Params` row adds the default and description of the field
    it names; rows naming no field (positional arguments) lead the section.
  - *Import line, install command, Markdown copy, on-page index*: derived
    from page data.
  - *Hand-maintained*: titles, descriptions, the usage snippet, defaults and
    field descriptions (Go has no default-value syntax, and most fields are
    documented in the struct comment, which the page shows above the table),
    and token lists.
  - *Rejected*: reflection (no doc comments, needs a value per type);
    generating usage snippets (they are editorial); parsing the library at
    run time (the gallery cannot embed parent directories).
- **On this page.** A fixed column right of the scroll area when the area
  beside the sidebar is at least 900 DIP; hidden otherwise. The current
  section is the last one whose heading has passed the top of the viewport,
  or the last section at the end of the page. Sections record where the
  previous frame laid them out; choosing one sets the content offset, with
  no animation, as an anchor does. The content column stays left-aligned,
  so on very wide windows the index sits away from it.
- **Code blocks** are one surface everywhere (install, import, usage,
  example source): no header, a copy control floating in the top corner,
  the fold kept since the usage snippet is always on the page. Example
  source is numbered; usage is not, as in coss. Go is coloured by keyword,
  string, number, comment, called function and package qualifier from theme
  roles; spans are cached per listing until the code or palette changes.
- **An example without page source has no Code tab** (chart pages built in a
  loop), rather than repeating the usage snippet.
- **API tables are compact**: 8 DIP row padding inside a 4 DIP panel inset.
- **Home page.** It owns the window: no sidebar, a centered column up to
  1120 DIP and a title bar with the brand, Components, Theming and search.
  Navigation and Cmd/Ctrl+K go through the sheet compact layouts already
  use. Persuade mode: a centered display headline over a fading
  grid, a stage of live components, catalog numbers counted from the
  registered pages, four statements of what the library is built on, the
  family grid and the downloads. Copy states only what PRODUCT.md and the
  engineering contract claim.

## Risks

- `go/parser` adds to the browser build's size.
- Example source shows gallery state helpers as written; examples that lean
  on them read less cleanly than a hand-edited snippet.
