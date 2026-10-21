# Components and limits

The gallery contains 66 native component specimens. Each entry below routes to a compact default and an interactive examples section. These are real library calls with independent caller-owned state and embedded compiled Go source. The inventory is a capability guide; the upstream visual reference does not promise browser API compatibility or every web variant.

All import suffixes are relative to `github.com/Cyenoch/mygoui`. `root` is the `mygoui` package. Run `go run ./examples/gallery` to inspect the actual constructors and fixture state. See `VALIDATION.md` for completed evidence.

| Component | Package | Native behavior |
| --- | --- | --- |
| [Accordion](https://ui.shadcn.com/docs/components/accordion) | `navigation` | Native disclosure group with single/multiple selection and arrow-key focus. |
| [Alert](https://ui.shadcn.com/docs/components/alert) | `root` | Native callout with semantic status, description and optional real action. |
| [Alert Dialog](https://ui.shadcn.com/docs/components/alert-dialog) | `overlay` | Modal focus, explicit action/cancel, no outside-click acceptance. |
| [Aspect Ratio](https://ui.shadcn.com/docs/components/aspect-ratio) | `root` | Native aspect-ratio layout with caller-owned media. |
| [Attachment](https://ui.shadcn.com/docs/components/attachment) | `conversation` | File/image metadata, lifecycle status and independent open/remove/retry actions. |
| [Avatar](https://ui.shadcn.com/docs/components/avatar) | `root` | Caller-loaded native image, initials fallback, badge and overflow count. |
| [Badge](https://ui.shadcn.com/docs/components/badge) | `root` | Semantic variants; passive text or actual keyboard-accessible action. |
| [Breadcrumb](https://ui.shadcn.com/docs/components/breadcrumb) | `navigation` | Native path actions, current-page semantics and overflow menu. |
| [Bubble](https://ui.shadcn.com/docs/components/bubble) | `conversation` | Selectable content, grouped alignment, reactions and retained expansion. |
| [Button](https://ui.shadcn.com/docs/components/button) | `root` | Native variants, sizes, loading/disabled states and real actions. |
| [Button Group](https://ui.shadcn.com/docs/components/button-group) | `root` | Joined native actions preserving independent keyboard focus. |
| [Calendar](https://ui.shadcn.com/docs/components/calendar) | `calendar` | Localized native date grid, selection and month keyboard navigation. |
| [Card](https://ui.shadcn.com/docs/components/card) | `root` | Composable header, content and footer with shared native spacing. |
| [Carousel](https://ui.shadcn.com/docs/components/carousel) | `layout` | Scrollable items with actual next/previous selection and bounded navigation. |
| [Chart](https://ui.shadcn.com/docs/components/chart) | `chart` | Maintained native charts, semantic series colors and keyboard selection. |
| [Checkbox](https://ui.shadcn.com/docs/components/checkbox) | `form` | Native checked/mixed/disabled states and accessible label. |
| [Code Block](https://ui.shadcn.com/docs/components/code-block) | `codeblock` | Exact source copy, wrap, numbers, line variants and native selection. |
| [Collapsible](https://ui.shadcn.com/docs/components/collapsible) | `navigation` | Native disclosure with persistent caller-owned expansion. |
| [Combobox](https://ui.shadcn.com/docs/components/combobox) | `navigation` | Native filtering, keyboard highlight and single/multiple choices. |
| [Command](https://ui.shadcn.com/docs/components/command) | `navigation` | Native search input and list with real command callbacks. |
| [Context Menu](https://ui.shadcn.com/docs/components/context-menu) | `overlay` | Native right-click and Shift+F10 menu, checked/disabled/submenu actions. |
| [Data Table](https://ui.shadcn.com/docs/components/data-table) | `data` | Native table with real Go sorting, filtering, selection and pagination. |
| [Date Picker](https://ui.shadcn.com/docs/components/date-picker) | `calendar` | Native date editor and calendar popup with caller-owned date/range. |
| [Dialog](https://ui.shadcn.com/docs/components/dialog) | `overlay` | Native modal with editable content, focus trap and focus restoration. |
| [Direction](https://ui.shadcn.com/docs/components/direction) | `layout` | Native logical layout direction; code remains left to right. |
| [Drawer](https://ui.shadcn.com/docs/components/drawer) | `overlay` | Native edge modal; gestures require separate native input evidence. |
| [Dropdown Menu](https://ui.shadcn.com/docs/components/dropdown-menu) | `overlay` | Native menu with actual actions, checks and submenus. |
| [Empty](https://ui.shadcn.com/docs/components/empty) | `root` | Native empty-state media, explanation and real recovery action. |
| [Field](https://ui.shadcn.com/docs/components/field) | `form` | Native associated label, description, errors and field groups. |
| [Hover Card](https://ui.shadcn.com/docs/components/hover-card) | `overlay` | Native anchored preview with keyboard route and delayed hover. |
| [Input](https://ui.shadcn.com/docs/components/input) | `form` | Native editor with IME, selection, undo and actual bound value. |
| [Input Group](https://ui.shadcn.com/docs/components/input-group) | `form` | Native editor with independent inline/block addon actions. |
| [Input OTP](https://ui.shadcn.com/docs/components/input-otp) | `form` | One native editor with slot guides, paste filtering and completion. |
| [Item](https://ui.shadcn.com/docs/components/item) | `root` | Native content/media/actions; static and actionable variants. |
| [Kbd](https://ui.shadcn.com/docs/components/kbd) | `root` | Native keycap labels and chord grouping. |
| [Label](https://ui.shadcn.com/docs/components/label) | `form` | Native focus handle association and accessible input naming. |
| [Marker](https://ui.shadcn.com/docs/components/marker) | `conversation` | Textual status, bordered row and labelled separator. |
| [Menubar](https://ui.shadcn.com/docs/components/menubar) | `overlay` | Native menu actions; macOS application placement is platform-owned. |
| [Message](https://ui.shadcn.com/docs/components/message) | `conversation` | Caller-composed avatar/header/content/footer and alignment. |
| [Message Scroller](https://ui.shadcn.com/docs/components/message-scroller) | `conversation` | Native stable list keys, intentional follow, anchors and history preservation. |
| [Native Select](https://ui.shadcn.com/docs/components/native-select) | `form` | Native platform choice menu with grouped/disabled options. |
| [Navigation Menu](https://ui.shadcn.com/docs/components/navigation-menu) | `navigation` | Native actions and popovers with keyboard traversal. |
| [Pagination](https://ui.shadcn.com/docs/components/pagination) | `navigation` | Actual caller page changes, boundary disabling and current-page semantics. |
| [Popover](https://ui.shadcn.com/docs/components/popover) | `overlay` | Native anchored editable popup, edge flipping and dismissal. |
| [Progress](https://ui.shadcn.com/docs/components/progress) | `root` | Native bounded value and accessible progress text. |
| [Questionnaire](https://ui.shadcn.com/docs/components/questionnaire) | `conversation` | Typed Go questions/answers, validation, previous/skip/submit and scoped shortcuts. |
| [Radio Group](https://ui.shadcn.com/docs/components/radio-group) | `form` | Native single-choice selection and roving keyboard focus. |
| [Resizable](https://ui.shadcn.com/docs/components/resizable) | `layout` | Maintained native Split with pointer and keyboard divider resizing. |
| [Scroll Area](https://ui.shadcn.com/docs/components/scroll-area) | `layout` | Native scroll container with wheel, keyboard and scrollbar input. |
| [Select](https://ui.shadcn.com/docs/components/select) | `navigation` | Native popup choice, typeahead, groups and selected-value semantics. |
| [Separator](https://ui.shadcn.com/docs/components/separator) | `root` | Native horizontal/vertical semantic rule. |
| [Sheet](https://ui.shadcn.com/docs/components/sheet) | `overlay` | Native edge modal preserving focus and editable content. |
| [Sidebar](https://ui.shadcn.com/docs/components/sidebar) | `navigation` | Native source list with collapse, item actions and narrow overlay. |
| [Skeleton](https://ui.shadcn.com/docs/components/skeleton) | `root` | Native paint-only shimmer with stationary reduced-motion/high-contrast shapes. |
| [Slider](https://ui.shadcn.com/docs/components/slider) | `form` | Native drag/keyboard stepped single value and two-thumb range. |
| [Spinner](https://ui.shadcn.com/docs/components/spinner) | `root` | Native busy indicator with reduced-motion-safe text. |
| [Switch](https://ui.shadcn.com/docs/components/switch) | `form` | Native boolean control with actual checked state. |
| [Table](https://ui.shadcn.com/docs/components/table) | `data` | Native virtual rows and columns with selection and row actions. |
| [Tabs](https://ui.shadcn.com/docs/components/tabs) | `navigation` | Native tab selection, roving keyboard focus and actual panels. |
| [Textarea](https://ui.shadcn.com/docs/components/textarea) | `form` | Native multiline editor with caller-owned source and read-only selection. |
| [Toast](https://ui.shadcn.com/docs/components/toast) | `overlay` | Native toast manager with actual action, update and dismissal. |
| [Toggle](https://ui.shadcn.com/docs/components/toggle) | `root` | Native persistent pressed-state button. |
| [Toggle Group](https://ui.shadcn.com/docs/components/toggle-group) | `root` | Native single/multiple selection and disabled-skipping navigation. |
| [Tooltip](https://ui.shadcn.com/docs/components/tooltip) | `overlay` | Native hover/focus description with real anchor and Escape dismissal. |
| [Trace](https://ui.shadcn.com/docs/components/trace) | `trace` | Native status/disclosure composition with persistent source-keyed expansion. |
| [Typography](https://ui.shadcn.com/docs/components/typography) | `root` | Native semantic type roles and append-stable prose rhythm. |

## Current boundaries

- **Styling:** Neutral and Blue presets use semantic colors, native metrics and accessibility preferences. Native menus follow the OS. Typography uses native font roles; arbitrary CSS and DOM slot compositions are outside this API.
- **Media:** avatar/attachment images are loaded by the caller. Network fetching, uploading, downloads and file pickers remain application work. An action requires a real callback.
- **Forms:** OTP uses one native editor with committed ASCII filtering and slot guides. Slider supports one value, including vertical orientation; RangeSlider supports two horizontal thumbs. Arbitrary thumb counts and vertical ranges are not implemented. Button Group joins action buttons; compose mixed controls with ordinary native layout.
- **Dates:** arithmetic is Gregorian, with civil date values and explicit location/constraints. Date Picker editors use strict ISO `YYYY-MM-DD`. Alternate calendar arithmetic is not implemented.
- **Charts:** bar, line, area, pie, radar and radial marks use native chart primitives. Polar values must be nonnegative, radar needs at least three categories, and empty/missing data remain explicit. Pie/radial samples use the first category per series.
- **Tables:** filtering, sorting and pagination operate on supplied local rows. Increment DataTable Revision when changing existing row cells or comparator behavior in place. Keep stable row/column IDs. Native cell/drag selection cannot be combined with row reordering. Selection may retain IDs across filtered/paged views; the caller owns bulk-operation scope. A custom cell color must account for `List.CellSelected`.
- **Navigation:** combobox option labels must be unique because the native base commits by display label. Pagination is one-based; DataTable's page state is zero-based. Tabs use native roving focus and a named group for panel semantics.
- **Overlays:** Drawer is an edge modal. Swipe dismissal, velocity snapping and nested gesture arbitration are not implemented. Tooltips retain native timing. ToastState retains up to 32 identities; MyGo shows the newest three. Loading toasts are sticky unless explicitly timed. All other timer, focus and lifetime behavior belongs to MyGo.
- **Conversation:** callers supply all message data, history and actions. The scroller uses native stable keys, follow-end, prepend anchoring and explicit jumps. Persistence and every possible saved-transcript/turn-anchor policy are application concerns.
- **Source:** Code Block preserves exact input and copies without line numbers; syntax and ANSI support are those of the native source widgets. Trace expansion is separate caller-owned presentation state.

## Evidence vocabulary

Behavior tests exercise native headless rendering and event routing. A passing test establishes only its observed contract. OS gestures, native menus, GPU partial presentation, IME integrations and screen-reader operation require real platform checks. The catalog does not declare every upstream variant verified.
