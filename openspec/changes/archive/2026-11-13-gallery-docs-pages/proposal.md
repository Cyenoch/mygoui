## Why

A component page shows one of four views at a time (Preview, Code, API,
Tokens), with one hand-written snippet for all examples and a hand-written
parameter table. The reader cannot see an example beside its own source, the
tables miss most option fields, and both drift from the code. The coss docs
page is the reference: one scrolling page with the demo first, then
installation, usage, API reference and examples, and an "On this page" index.

## What Changes

- Component pages become one scrolling document: title, description and
  actions; the first example with a Preview/Code switch; Installation; Usage;
  API Reference; Examples, each with its own Preview/Code switch; Theme
  tokens. Wide windows show an "On this page" index that follows the scroll
  and jumps to a section.
- Example code is the example's own Go source, read from the gallery's
  embedded files, so it cannot drift.
- API Reference lists every exported field of the page's option types with
  its Go type and the type's doc comment, generated from the library source.
  Hand-written rows only add defaults and descriptions.
- **BREAKING** for the gallery only: the page-wide view switch is removed.

## Impact

`internal/gallery` page layout, state and tests; a generator under
`internal/apidoc` with a checked-in `internal/gallery/api_gen.go` and a
freshness test; `make generate`. Library packages are unchanged.
