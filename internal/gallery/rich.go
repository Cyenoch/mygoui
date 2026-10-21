package gallery

import (
	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/codeblock"
	"github.com/Cyenoch/mygoui/trace"
	"github.com/egoist/mygo/ui"
)

// CodeState keeps editable exact source separate from its native presentation state.
type CodeState struct {
	Source                      string
	Block                       codeblock.State
	Numbers, Destructive, Small bool
}

func NewCodeState() CodeState {
	return CodeState{Source: "package sample\n\nfunc greeting() string {\n\treturn \"Hello from native Go\"\n}\n", Numbers: true}
}

func PreviewCode(c *ui.Context, s *CodeState) {
	ui.Column(c).FillWidth().MinWidth(0).Justify(ui.Start).AlignItems(ui.Stretch).Gap(12).Children(func() {
		ui.Text(c, "Exact source").Bold()
		ui.TextArea(c.Key("code-source"), &s.Source).Label("Code source").Lines(4, 8).FillWidth()
		ui.Row(c).FillWidth().Wrap().Justify(ui.Start).AlignContent(ui.Start).Gap(8).Children(func() {
			mygoui.Toggle(c, &s.Numbers, mygoui.ToggleOptions{Label: "Line numbers"})
			mygoui.Toggle(c, &s.Destructive, mygoui.ToggleOptions{Label: "Destructive"})
			mygoui.Toggle(c, &s.Small, mygoui.ToggleOptions{Label: "Small output"})
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Append sample line", Variant: mygoui.ButtonOutline,
				OnClick: func() { s.Source += "// appended sample line\n" }})
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Reset code", Variant: mygoui.ButtonOutline,
				OnClick: func() { *s = NewCodeState() }})
		})
		variant := codeblock.Default
		if s.Destructive {
			variant = codeblock.Destructive
		}
		size := codeblock.SizeDefault
		if s.Small {
			size = codeblock.SizeSmall
		}
		codeblock.View(c.Key("source-preview"), s.Source, &s.Block, codeblock.Options{
			Title: "Local sample source", Language: "go", Footer: "Copy uses the exact editable source, including its final newline.",
			Variant: variant, Size: size, LineNumbers: s.Numbers, MaxHeight: 280,
		})
		ui.Text(c, "Line treatments").Bold()
		codeblock.View(c.Key("line-preview"), "unchanged\nremoved\nadded\nhighlighted\n", nil, codeblock.Options{
			Title: "Explicit presentation variants", LineNumbers: true,
			Lines: []codeblock.Line{
				{Text: "unchanged", Variant: codeblock.LineDefault}, {Text: "removed", Variant: codeblock.Removed},
				{Text: "added", Variant: codeblock.Added}, {Text: "highlighted", Variant: codeblock.Highlighted},
			},
		})
	})
}

// TraceState owns both nesting levels independently so hiding a parent retains child expansion.
type TraceState struct {
	Source, Output, Result    string
	Running, Failed, Disabled bool
	Outer, Nested             trace.State
	Code                      codeblock.State
}

func NewTraceState() TraceState {
	return TraceState{Source: "This is an editable gallery reasoning sample, not observed agent activity.",
		Output: "sample output\n\nfinal line\n"}
}

func PreviewTrace(c *ui.Context, s *TraceState) {
	ui.Column(c).FillWidth().MinWidth(0).Justify(ui.Start).AlignItems(ui.Stretch).Gap(12).Children(func() {
		ui.TextArea(c, &s.Source).Label("Trace sample text").Lines(2, 5).FillWidth()
		ui.TextArea(c, &s.Output).Label("Trace sample output").Lines(2, 5).FillWidth()
		ui.Row(c).FillWidth().Wrap().Justify(ui.Start).AlignContent(ui.Start).Gap(8).Children(func() {
			mygoui.Toggle(c, &s.Running, mygoui.ToggleOptions{Label: "Running"})
			mygoui.Toggle(c, &s.Failed, mygoui.ToggleOptions{Label: "Error"})
			mygoui.Toggle(c, &s.Disabled, mygoui.ToggleOptions{Label: "Disabled disclosure"})
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Append output", Variant: mygoui.ButtonOutline,
				OnClick: func() { s.Output += "additional fixture output\n" }})
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Reset trace", Variant: mygoui.ButtonOutline,
				OnClick: func() { *s = NewTraceState() }})
		})
		status, title := trace.Default, "Completed sample activity"
		if s.Running {
			status, title = trace.Running, "Running sample activity"
		}
		if s.Failed {
			status, title = trace.Error, "Error in sample activity"
		}
		trace.Trace(c).Children(func() {
			trace.Collapsible(c.Key("sample-activity"), &s.Outer, trace.Options{Title: title, Status: status,
				Icon: lucide.Brain(), DefaultOpen: true, Disabled: s.Disabled}, func() {
				ui.Text(c, s.Source).Selectable()
				trace.Collapsible(c.Key("sample-output"), &s.Nested, trace.Options{Title: "Nested sample output", Icon: lucide.Terminal()}, func() {
					variant := codeblock.Default
					if s.Failed {
						variant = codeblock.Destructive
					}
					codeblock.View(c.Key("output-body"), s.Output, &s.Code, codeblock.Options{Size: codeblock.SizeSmall,
						Variant: variant, MaxHeight: 220})
				})
			})
			trace.Line(c, trace.Pending).Children(func() {
				trace.Icon(c, lucide.Info())
				trace.Title(c, "Pending sample activity", trace.Pending)
			})
			trace.Action(c, trace.Options{Title: "Record sample action", Icon: lucide.Check()}, func() { s.Result = "Sample action activated"; c.AddToast(ui.Toast{Title: s.Result}) })
			trace.Link(c, trace.Options{Title: "Official Trace documentation", Icon: lucide.FileText()}, "https://ui.shadcn.com/docs/components/trace")
		})
	})
}
