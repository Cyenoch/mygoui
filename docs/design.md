# Design and public contract

`github.com/Cyenoch/mygoui` reconstructs suitable
[coss ui](https://coss.com/ui/docs) components with native MyGo controls.
The coss component catalog, MIT registry and published styling documentation
are the design reference. Previous appearance and tokens impose no
constraint on the `coss` branch.

The [component matrix](components.md) records the complete catalog mapping,
implemented subsets and gaps. Names in that matrix establish correspondence,
not pixel equivalence or completed platform testing.

## Ownership

A constructor receives a build-scoped `*ui.Context`. Returned `ui.Element`
values and custom parts expire before the next build pass, even when two passes
happen in one frame. Retain model values, widget state and `ui.Handle` per
window. Never save an element or context in application state or a worker.

MyGo bases own input, editing, caret and selection, focus routing, scroll
behavior and the accessible control tree. Style those bases rather than
replacing them with pointer-only drawings. Use stable identities before
construction: `c.Key(key)` for controls and collection entries, and native
handles for persistent focus and committed geometry.

Composition uses typed Go options, body callbacks for independent content,
and typed slices for coordinated items. There are no ambient component
providers or child introspection. `GroupItem.Build` receives its segment's
corners explicitly; it must apply them to the control it returns.
Vertical groups own their members' width, including icon-only and explicitly
sized controls, and match their widest member without filling the parent.
`FullWidth` fills the parent instead. Explicit separators are ordinary
`GroupItem` builders, matching the upstream Group examples. Explicit `Corners` changes geometry, not a
button's elevation; specialized ButtonGroup owns its joined-button styling.

Native input is applied as constructors return. Use handled actions such as
`OnClick`, `OnChange` and `OnSubmit` for side effects; a view may rebuild
several times for one event. Capture `c.Services()` for persistent clipboard
and URL callbacks. Slow I/O belongs in application workers, with results
published through `Window.Update`; invalidation alone does not synchronize
model writes.

Importing the library must not initialize an application, register application
languages or fonts, read files, use persistent storage, or start background
jobs. Application startup owns those policies.

## Packages

| Package | Responsibility |
| --- | --- |
| `mygoui` | Actions, surfaces, text, feedback, Frame, Group, Sidebar, ScrollArea and ScrollFade |
| `theme` | Semantic colors, radius/spacing/type metrics and process-wide host preferences |
| `form` | Native fields, choices, editors, search and plain-text Composer |
| `list` | Virtualized keyed ListBox and selection state |
| `navigation` | Accordion/Disclosure, tabs, segments, breadcrumbs, pagination, toolbar, menus and commands |
| `overlay` | Dialogs, sheets, drawers, popovers, previews, tooltips and toast state |
| `date` | Gregorian calendars, date/range pickers and time editing |
| `data` | Native virtualized table and column/selection state |
| `color` | Additional native color editing components |
| `locale` | Library messages; application language registration remains outside the library |
| `fonts` | Opt-in embedded Cal Sans and explicit startup registration; no import-side registration |

Public APIs live at the module root or in coherent family packages.
`internal/gallery` and `examples/` are consumers, not an alternate API layer.
Additional library families can remain useful without becoming coss guarantees.

The gallery is a native workbench consumer, not a web documentation clone. Its
52 DIP custom title bar reads `c.TitleBar()` for actual left/right system
control insets, keeps draggable content regions, and uses native macOS traffic
lights. A grouped sidebar, Cmd/Ctrl+K search and Preview/Code/API/Tokens views
share window-owned component state. Compact layouts preserve focus rather
than rebuilding a second application shell. Dialog owns the Modal size
examples; Drawer owns the native floating BottomSheet examples. Their Go APIs
remain available, without duplicate top-level Gallery pages.

Cmd/Ctrl `+`, `-` and `0` use native content zoom, not text-only scaling.
The system-button strip stays 52 window DIPs high while content breakpoints
and insets follow the zoom. Windows caption buttons follow the explicit
Gallery appearance. A separate native Flyout example lives on the Popover
page; it is not a duplicate overlay catalog entry.

Install `ToastState.View` once at the window root, after theme application and
before modal views; page navigation must not remove its renderer. Overlay
footers use a full-bleed muted face and top rule. Footer callbacks build their
actions directly; set `ButtonOptions.FullWidth` on narrow layouts rather than
expecting a parent to override a control's explicit alignment. Body scrolling
keeps actions fixed normally; an outer native viewport makes header and footer
reachable when short windows or enlarged text cannot fit the chrome.

## Appearance

Read current semantic roles with `theme.Colors(c)`. `Primary` is the action
color; `Accent` is the subtle interaction surface. `MutedBackground` and
`MutedForeground` have distinct meanings. Preserve translucent border/input
tokens so they composite over the real surface.

`theme.Apply(c)` installs the native widget theme for the current appearance.
`ApplyMode(c, dark)` supports explicit previews. Preferences are process-wide;
widget values and focus state remain window-local. `Configure` replaces the
whole preference value, while `UpdatePreferences` preserves fields the host
does not edit. Each appearance has a `PalettePatch`, including canonical
sidebar and chart roles.

The default radius base is 10 DIP, giving SM 6, MD 8, LG 10 and XL 14.
Desktop buttons use 32 DIP by default, 28/36 for small/large and 24/40 for
extra small/extra large. Button press feedback uses color and shadow changes,
without a pressed scale. Native disabled painting uses 50% opacity, whereas
the coss button reference specifies 64%.

Default body/heading/mono families are native fallbacks. Configurable family
names do not load a font. See [theme.md](theme.md) for role customization and
the documented/runtime dark palette discrepancy, and [fonts.md](fonts.md)
for independently licensed host font sources.

## Native behavior and boundaries

The GPU path was observed to drop an ancestor's rounded clip beneath a
rectangular descendant clip on MyGo 0.7.1; the 0.7.2 upgrade does not establish
that it is fixed. Headless rendering does not expose that backend limitation.
State backgrounds use native per-corner fills; progress and table row painting
apply their final clip directly around the fill. Arbitrary nested content
still requires a toolkit fix. Do not flatten every container into an effect:
that would hide embedded WebViews, change fade behavior and impose texture limits.

- Caller-owned values and callbacks express selection, validation, submission
  and external actions. A nil callback does not fabricate a successful action.
- Labels and help/errors must reach native field and control associations.
  Icon-only actions still need a meaningful accessible name.
- Disabled actions must remain inert. Read-only editors preserve selection.
  Loading/submission policy must be explicit in the owning component.
- Anchors and overlay parts are consumed in the same build. Keep overlay
  constructors installed while closed when native retained exits require it.
  Exit snapshots are inert and must not rebuild consumer content.
- `navigation.ContextMenu` keeps the OS menu by default; `PopoverMenu` opts
  into MyGo's themed in-window menu with toolkit-owned focus, keyboard and
  selection semantics. Neither rebuilds custom pointer-only menu routing.
- Tooltip timing, clipboard services, scrollbars and input methods keep
  toolkit/platform behavior. The matrix names adaptations such
  as plain-text Composer and Drawer without snap points.
- `ScrollFade` offers true subtree alpha masking through `MaskFade`, alongside
  shadow and opaque-surface edge treatments. Native effect layers have an
  8192-pixel side limit and omit embedded webviews.

The library translates appearance and behavior using the toolkit's contracts.
Tailwind selectors, DOM polymorphism, portals, contenteditable tokens and
React hooks do not become parallel native infrastructure.

## Verification and reference updates

The coss snapshot and licensing boundary are recorded in
[references.md](references.md). Read affected MIT registry files, their MDX
examples and published styling before changing behavior. The pinned
`references/coss` checkout is the only component design reference.

Meaningful tests protect observable editing, activation, focus, selection,
validation, dismissal and retained paint. `ui.NewTester` can exercise native
event routing, semantic metadata and headless paint; it does not prove GPU
presentation, OS gestures, IME integration or screen-reader operation.
Record those separately in [VALIDATION.md](../VALIDATION.md).
