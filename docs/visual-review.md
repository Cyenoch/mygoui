# Coss visual review

Review date: 2026-10-27. This records an actual visual inspection, not a
pixel-diff certification. The source reference is coss ui commit
`3d8ca28e7d90123e9e0f44665944dc240b026ccb`; the live website may differ.

## Evidence and method

- Opened the official documentation in a browser, rendered each of the 56
  available component pages, and captured light/dark full-page PNGs. Opened
  relevant menus, calendars, dialogs, sheets, command palette, popovers,
  preview cards, tooltips and toast examples. Inspected the saved images,
  rather than treating a successful browser action as evidence of visibility.
- The Sidebar documentation URL returns 404. Its reference is the live
  documentation sidebar plus the pinned MIT registry source, not an invented
  component demo. The catalog union therefore remains 57 subjects.
- Native captures use the actual Go gallery examples through `ui.NewTester`,
  a 720 x 720 DIP canvas, explicit Cal Sans registration, both appearances,
  and named interaction states. The first collection produced 542 captures
  and 99 contact sheets. Those counts describe evidence, not acceptance gates.
- Inspected all native default contact-sheet rows against official rendered
  samples. Compared component chrome, density, typography, state surfaces,
  geometry and hierarchy; ignored different application/sample text and the
  gallery's specimen frame. Native screenshots are 1x; browser images are 2x.
- Inspected selected native state/motion sheets, including the segmented
  indicator and Sheet progression. Headless timed frames establish native
  progression, not browser/native frame-perfect timing or GPU equivalence.

Evidence is local and ignored by Git:

```text
artifacts/coss/reference-<slug>-<light|dark>[-open]-full.png
artifacts/coss/reference-sidebar-<light|dark>-docs-shell-full.png
artifacts/coss/reference-review/
artifacts/coss/native-audit/manifest.json
artifacts/coss/native-audit/contactsheets/defaults-01.png ... defaults-15.png
artifacts/coss/native-audit/contactsheets/states/
artifacts/coss/native-audit-final/manifest.json
artifacts/coss/native-audit-final/contactsheets/
```

The manifest records each native example, mode, target, capture bounds and
state coverage, including states that could not be exercised. Cropped contact
sheets are navigation aids: overlays can lie outside the first specimen crop,
so consult the original full-page image. Some official animations and hover
panels settle between captures; a filename containing `open` alone is not
proof that its panel is visible. Default screenshots do not cover every
possible option or every state permutation.

## Component-by-component inspection

All rows received a light/dark baseline visual inspection. "Native" in the
notes means an intentional adaptation, not a claim of identical web behavior.
Detailed functional limitations remain in [components.md](components.md).

| Component | Inspected result and boundary |
| --- | --- |
| Accordion | Neutral divided rows, compact labels, chevrons and expanded content. Native layout animation; independent expansion state. |
| Alert | Rounded bordered status callouts, subdued fill, status icon and muted description. Titles stay neutral, matching the registry rather than taking the status color. |
| Alert Dialog | Rounded explicit-choice popup, description and separated action footer. Re-captured the official open panel after discovering an earlier trigger-only image. Native focus ownership and dismissal policy. |
| Autocomplete | Inline editor and filtered list rather than a separate search-dialog trigger. Native list placement/virtualization; clear affordance is configurable. |
| Avatar | Circular size progression and neutral initials fallback. Gallery uses initials rather than copying the site's profile photography. |
| Badge | Compact default, secondary, outline, destructive and semantic faces. Independent badge variants, not inherited chip geometry. |
| Breadcrumb | Muted ancestors, chevrons and stronger current destination. Application text/overflow policy differ. |
| Button | Compact rounded default/secondary/outline/ghost/link/destructive faces and loading/disabled examples. Native disabled opacity is 50%, reference 64%; no press-scale effect. |
| Calendar | Centered month grid, neutral navigation and selected-day surface. Gregorian native calendar; year picker is an explicit example rather than part of the default demo. |
| Card | Bordered 2xl card, compact title/description and section-owned padding. Secondary/tertiary/transparent are native additions. Sample content must provide final-section spacing. |
| Checkbox | Small outlined square, neutral selected fill and check glyph, aligned label. Native mixed/disabled/focus semantics. |
| Checkbox Group | Vertically aligned choices and group label/description. Group content and selection are caller-owned. |
| Collapsible | Compact disclosure row, rotating chevron and revealed content. Native disclosure body instead of React parts. |
| Combobox | Inline editable value, disclosure/clear actions and filtered popup. Gallery's clear option and labels are consumer choices. |
| Composer | Bordered multiline draft, fixed submit-action area and disabled empty submit. Plain text only; no rich structured tokens/mentions. |
| Command | Compact search row and action list, muted disabled items and shortcut text. Official dialog/grouped palette also opened; native inline Command composes with Dialog. |
| Context Menu | Official popup chrome inspected. Native implementation intentionally uses the OS context menu, so web icon/shortcut styling is not reproduced. Headless capture cannot prove OS menu presentation. |
| Date Picker | Native date editor plus anchored calendar; official calendar popup inspected. ISO editing is an intentional difference from the web's button-only default example. |
| Dialog | Rounded popup, close affordance, muted description, dim/blur backdrop and shaded footer. Default scrim corrected to black 32% with 4 DIP blur sigma. |
| Drawer | Attached edge panel, close affordance and inward-edge drag handle. Official bottom placement and native right placement are different configured examples. No snap-point/nested-gesture parity. |
| Empty | Centered icon tile, strong heading and restrained muted description. Recovery actions are caller-provided. |
| Field | Label/control/help hierarchy inspected. Found inherited medium-weight, oversized help text; repair must retain native naming, focus and description association. |
| Fieldset | Semibold legend, muted description and grouped controls. Native grouping rather than browser fieldset styling. |
| Form | Vertical field/action composition inspected. Validation, submission and networking belong to the caller; no browser constraint-validation engine. |
| Frame | Muted rounded wrapper, independently bordered inner panel, header and footer. Different sample content does not affect the framing contract. |
| Group | Shared control boundaries and outer corners. Native explicit segment corners replace CSS child selectors. |
| Input | Compact height, thin border, subtle light shadow, dark input fill and muted placeholder. Native editing/selection retained. |
| Input Group | Shared input chrome around prefix/suffix content. Reuses Field metadata; help-text repair applies here too. |
| Kbd | Small subdued keycap surface and compact shortcut typography. System mono fallback is not the site's Paper Mono font. |
| Label | Compact medium label with optional required/invalid marker. Click-to-focus uses a persistent native handle. |
| Menu | Group headings, check indicators, separator and disabled rows. Standalone Menu is content; the anchored popup owns outer border/shadow. |
| Meter | Square 8 DIP quantity track, aligned label/value and configurable indicator color. Colored gallery meters are optional semantic examples; the default remains neutral. |
| Number Field | Shared input frame and native decrement/increment actions. Found left-aligned value versus centered official value; requires native editor alignment correction. No scrub-area behavior. |
| OTP Field | Compact separated slots, label and help. One native editor paints slots instead of six web inputs; no OS autofill claim. |
| Pagination | Compact previous/next, current-page outline and ellipsis. Result summary is an optional application composition. |
| Popover | Rounded anchored panel, close action and subtle light shadow. Optional arrow and size are native options, not mandatory defaults. |
| Preview Card | Hover/focus-triggered bordered preview, handoff to interactive body and leave grace. Website metadata/profile content is application-owned. |
| Progress | Thin rounded neutral track/indicator, optional label/value. Determinate and indeterminate timing are native adaptations. |
| Radio Group | Small outlined circles, neutral selected fill and aligned choices. Native roving focus and mutually exclusive state. |
| Scroll Area | Clipped bordered viewport and overflow content. Native scrollbars/gestures replace the website's custom scrollbar parts. |
| Select | Compact trigger and selected/disabled popup rows. Native list placement differs from web selected-item alignment and browser form integration. |
| Segmented Control | Muted tray, elevated selected face and disabled choice. Inspected 50/120/600 ms native indicator frames; selected surface moves and settles without becoming a tab panel. |
| Separator | Thin neutral horizontal/vertical rules. Additional emphasis levels are native extensions. |
| Sheet | Flush edge panel, native close/focus, separated footer and dim/blur scrim. Inspected intermediate slide/fade frames as well as settled panel. |
| Sidebar | Live documentation group/selection hierarchy and native expanded navigation inspected. Dedicated official page is unavailable; collapsed/mobile/persistence policies are native options/limits. |
| Skeleton | Neutral circular and text-like blocks, rounded geometry and reduced-motion behavior. Site's timed content-replacement demo is not itself a skeleton primitive. |
| Slider | Thin track, neutral active range, round thumb and optional output/ticks. Native range interaction and vertical limits are documented. |
| Spinner | Native loading arc at multiple sizes inspected. Geometry/animation comes from the MyGo loading plugin, not the identical web SVG. |
| Switch | Small pill track and circular thumb, neutral selected colors and aligned description. Native binary control. |
| Table | Neutral separators, readable body/status cells and compact paging. Native data-grid header/resize grips/virtualization visibly differ from the web HTML table. |
| Tabs | Muted tray, selected surface, disabled entries and panel hierarchy. Native focus/selection lifecycle, including removal of selected items, is regression-tested separately. |
| Textarea | Rounded multiline editing surface, muted metadata and native selection. Native line limits/scrolling replace browser resize-handle behavior. |
| Toast | Rounded popover-like card, title/description and semantic icon. Native queue/deck/timing rather than browser portal behavior. |
| Toggle | Quiet off face and subdued selected fill without press scale. Explicit native pressed state. |
| Toggle Group | Found joined square middle item in default appearance; official default keeps individual rounded toggles with 2 DIP gaps. Joined corners belong to outline appearance. |
| Toolbar | Compact enclosing surface, grouped actions and thin separators. Native arrow-key focus group; application supplies controls. |
| Tooltip | Small bordered light/dark hint; rich native example and official "Helpful hint" popup inspected. Native rest/warm timing differs from the web provider API. |

## Corrections discovered by visual review

The baseline pass identified three concrete production discrepancies: Field
help/error typography, NumberField value alignment and default ToggleGroup
corner/gap behavior. All three were corrected and re-captured after production
freeze. Inspected final contact sheets `defaults-06.png`, `defaults-09.png` and
`defaults-14.png`: help text is compact and regular, the number is centered,
and the selected middle toggle has rounded corners with separated neighbors.
The parent re-ran the collector successfully: 57 subjects, 114 mode entries,
550 captures and 102 contact sheets under `native-audit-final`. Baseline images
remain evidence of the findings, not passing snapshots.

## Real window and platform limits

Built and launched `Coss UI Gallery.app` on macOS 27.0.1 / arm64. Inspected
actual title-bar traffic lights and content alignment after setting their
position to (16, 19) in a 52 DIP custom bar. Operated Cmd+K search, entered
"Button", navigated to its page, activated Default (counter advanced), opened
Appearance and switched to Dark. The real native window and both appearances
rendered; evidence is in the task's native screenshots. Also opened ScrollFade:
on the macOS Metal renderer, MaskFade reveals the colored gradient beneath its
fading content. Activated a masked child button (the counter changed) and
used Page Down; the content moved from builds 01-04 to builds 05-08 while
both overflowing edges faded. This is a scoped real-GPU visual/input check,
not an effect performance benchmark.

Windows/Linux builds and simulated title-bar inset tests do not establish
actual Windows/Linux window presentation. This review does not certify screen
readers, IME, clipboard, every GPU/backend or mask performance, or every
component's all-state interaction on a physical OS window. See
[VALIDATION.md](../VALIDATION.md).
