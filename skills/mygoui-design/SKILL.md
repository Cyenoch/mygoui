---
name: mygoui-design
description: "Design and behavior rules for mygoui screens. Use when designing, reviewing or testing a screen built with mygoui."
---

# mygoui design and behavior

For work on the mygoui repository itself, also read [MAINTAINING.md](MAINTAINING.md).

`DESIGN.md` owns the visual rules and `docs/components.md` each component's behavior and gaps; both ship with the module version your application uses (`go list -m -f '{{.Dir}}' github.com/Cyenoch/mygoui`; when vendored, the repository at that version). Read `DESIGN.md` before a visual decision these rules leave open. This is a native desktop system: follow MyGo's desktop window, keyboard and pointer conventions.

## Visual rules

- **Hierarchy** comes from type roles, tone and position: headings through `mygoui.Heading`/`Typography`, secondary text in `MutedForeground`, compact neutral surfaces with thin borders. Color carries meaning: `Primary` marks the action, status roles mark state.
- **Density** is the desktop scale: the default control sizes, `theme.Space*` gaps on the 4-DIP rhythm, and explicit alignment (MyGo's zero alignment is `ui.Start`, and rows center children by default).
- **Composition** reuses the components: a Card, Frame or settings group for a section, a Dialog or Sheet for a task that interrupts. Transient controls open as overlays.
- **State shows** in fill, ink, an indicator or an icon. Text keeps its weight in every state.
- **Windows change size.** Check compact and wide windows, both appearances and enlarged text; actions stay reachable in each.

## Behavior rules

- **Feedback.** Every state the user causes shows at once or moves into place; the `mygoui-motion` skill holds the timing.
- **Disabled** controls are inert. **Read-only** editors keep selection and copying. **Invalid** fields set `Invalid` and an `Error` that says what to fix, on the field itself.
- **Loading.** Use the component's caller-owned `Loading` flag: actions keep their size and swap content for a spinner, fields stay editable, data regions show skeletons on first load and dim under a refresh. Composed content without a flag goes in `mygoui.LoadingOverlay`. A region of your own uses `mygoui.SkeletonGroup` with `SkeletonNone` placeholders. The loading section of `DESIGN.md` holds the full treatment.
- **Keyboard and focus.** Everything the pointer does, the keyboard does. Keep MyGo and mygoui controls in charge of focus, editing and selection; a pointer-only drawing loses them. Move focus deliberately with a `ui.Handle` (`.Bind(&h)`, `h.Focus(c)`), and declare commands with `h.OnShortcut` or `c.OnShortcut`.
- **Accessible names.** Every control has one: `Label` on buttons (icon-only too), fields and lists; `AriaLabel` on a bare input; `Label` on `ScrollArea` and `LoadingOverlay`. Field help and errors stay attached through the field's `Description` and `Error`.
- **Dismissal.** Dialogs and popovers own their focus and Escape handling; keep their constructors built while closed so the exit can play.

## Copy

Labels name the thing or the action in a word or two. Each fact appears once. Secondary help sits in a field's `Description`, a tooltip or an empty state, shown where needed. Errors state the fix. An explanatory sentence on a screen usually signals a control that should explain itself; fix the control.

## Verify with ui.NewTester

`ui.NewTester(view, width, height)` runs layout, input routing, the accessibility tree and headless paint in a test.

1. Drive it as a user: `Click(name)`, `Type(text)`, `Key(mods, key)`, then `Frame()`.
2. Check what assistive technology sees: `Access(name)` returns role, name, value and `Disabled`, `Invalid`, `Busy`, `Focused`, `Checked`, `Selected`.
3. Check motion settles: `Advance(d)`, then `FramePending()` is false at rest.
4. Repeat under `SetDark(true)`, `SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1.5})` and `SetSize(w, h)` for a compact window.
5. Inspect paint with `Image()` where a pixel is the contract.

```go
func TestSaveShowsBusyThenSettles(t *testing.T) {
	saving := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Save", Loading: saving, OnClick: func() { saving = true },
		})
	}, 320, 200)
	if err := tt.Click("Save"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if info, ok := tt.Access("Save"); !ok || info.Role != ui.RoleButton || !info.Busy {
		t.Fatalf("Save = %+v, want a busy button", info)
	}
	saving = false // the work finished
	tt.Frame()
	tt.Advance(time.Second)
	if tt.FramePending() {
		t.Fatal("the view still asks for frames at rest")
	}
}
```

A Tester paints without a GPU and has no OS window. It does not show GPU shader output, window chrome, IME, OS gestures or screen-reader speech; check those in the running app on each platform and report them separately.
