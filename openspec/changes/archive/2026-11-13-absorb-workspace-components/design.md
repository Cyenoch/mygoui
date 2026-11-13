## Decisions

- **Loading is one state.** A region that has content dims it and shows a
  line; one without content shows a skeleton or spinner. An indicator waits
  150 ms before it appears, so short work shows nothing. `internal/busy`
  owns this; components expose `Loading bool`.
- **State motion.** 150 ms for a control's state, 200 ms for a reveal, 300 ms
  for a number roll or a tour's travel. The first build never animates,
  reduced motion settles movement while fades run, and nothing requests a
  frame at rest.
- **No weight change on state.** Heavier text is wider and moves what follows
  it; selection uses fill, ink or an indicator.
- **Tour.** A scrim with a spotlight cut out says a tour is running; one
  panel travels with the spotlight while its text crossfades, instead of a
  new bubble per step. A dark window adds an outline to the spotlight.
- **Log levels.** Pressed level toggles narrow the log to those levels; none
  pressed shows every line. A line without a level follows Info.
- **Window shortcuts do not move focus.** `DocumentTabs` switches tabs from
  Ctrl+Tab without taking focus from the document; Cmd/Ctrl+W closes the
  active tab and never a pinned one.
- **A pushed page is not a modal.** `NavigationStack` goes back by button,
  Alt+Left and the mouse back button, not Escape, which belongs to editors.
- **Unloaded lazy branches check as a unit**; children that load later start
  unchecked.

## Findings and verification

An independent review reported 25 defects; 23 were fixed with regression
tests where the tester could reproduce them. Left: the spoiler opens on any
focus inside it (MyGo 0.9 does not report the focused descendant's box), and
tone ink on a status bar is black or white rather than a theme token.
Rendered images at 3x in both appearances were inspected for every component;
real-window, GPU, screen-reader and trackpad checks were not run.
