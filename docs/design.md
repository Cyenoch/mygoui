# Design and public contract

`github.com/Cyenoch/mygoui` provides reusable native MyGo components with a shadcn-inspired visual style. The root package is `mygoui`. Family packages group coherent APIs: `form`, `navigation`, `overlay`, `data`, `calendar`, `chart`, `conversation`, `codeblock`, and `trace`. `theme` owns semantic appearance and `locale` owns library messages.

## Ownership

A constructor receives `*ui.Context` and returns a build-scoped `ui.Element`. Callers retain model data, component state and `ui.Handle` per window. MyGo owns editing, input, focus, selection, animation scheduling, virtualization and accessibility. Event callbacks change caller state. Views perform no network or filesystem work.

The theme defaults to Neutral, with an optional Blue preset. Apply it at the view root. Semantic tokens and text scaling have one authority in `theme`. Host preferences are explicit process-wide policy; widget state is window-local. Library messages use the `mygoui` namespace and do not register application languages. Applications choose languages and may override namespaced messages.

The module has two direct dependencies: MyGo v0.6.0 and the published mygoicons v0.1.0. There are no local replacements, workspace requirements or sibling checkout dependencies. The gallery is an internal example consumer and owns its own English and Simplified Chinese application catalog.

## MyGo 0.6 integration

- Spinner delegates paint scheduling, reduced-motion poses and progress semantics to `plugins/loading.Arc`.
- Native table cell selection and row dragging use `TableState.List.CellSelect` and `DragSelect`. `Column.RowHeader` maps directly to the toolkit's row-header behavior. Custom cell painters consult `List.CellSelected` when overriding selection colors. Row reordering is mutually exclusive with native cell/drag selection.
- Keyed elements and retained widget state allow the toolkit's incremental layout and partial redraw to work. The library adds no frame loop or parallel rendering cache.
- Native lists own follow-end and scroll anchors. Native menus retain platform appearance, including Windows 11 presentation.
- Headless tests use the 0.6 Tester contract: off-window clicks return errors, and wheel input can exercise smooth native scrolling.

## Scope and compatibility

Shadcn is a visual reference, not a React API or browser compatibility promise. Native Go options, comparable identities, explicit errors and small callbacks are the public interface. No CSS, WebView, JavaScript, React or generated runtime registry is required.

Initial source includes 66 component specimens. Coverage of a named specimen does not imply every web variant is supported. Current limits are documented in `components.md`. While the module is v0, incompatible API changes require a minor version and an entry in the changelog; patch versions preserve API behavior.

## Verification

Consumers must be able to build using only this module and resolvable dependencies. Meaningful tests cover actual editing, focus, disabled controls, selection, dates, filtering, source fidelity, overlays, motion and narrow/enlarged layouts. Headless rendering establishes pixels and native event routing; OS gestures, GPU presentation and screen-reader integration require separate platform verification.
