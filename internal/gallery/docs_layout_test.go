package gallery

import (
	"strings"
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/egoist/mygo/ui"
)

func exampleFramesFixture() ComponentPage {
	return ComponentPage{
		ID: "examples-fixture", Title: "Examples fixture", Group: "Foundations",
		Examples: []ComponentExample{
			{Title: "Button variants", Demo: func(c *ui.Context, s *State) {
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Primary"})
			}},
			{Title: "Second variant", Demo: func(c *ui.Context, s *State) {
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Ghost", Variant: mygoui.ButtonGhost})
			}},
		},
		Code: "package main\n\nimport \"github.com/Cyenoch/mygoui\"\n\n" +
			"func main() {\n\tprintln(\"hi\")\n}\n",
		Tokens: []string{"Accent"},
	}
}

func TestWorkbenchCopyActions(t *testing.T) {
	state := New()
	state.Selected, state.Panel = "button", 1
	tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, 1180, 900)
	tt.Frame()
	if err := tt.Click("Copy Markdown"); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"# Button", "## Usage", "```go", "## API Reference", "## Theme tokens"} {
		if !strings.Contains(tt.Clipboard(), want) {
			t.Fatalf("copied Markdown is missing %q", want)
		}
	}
	if !tt.HasText("Copied") {
		t.Fatal("copy action did not report success")
	}
	tt.Advance(codeCopyFeedback)
	if err := tt.Click("Copy import"); err != nil {
		t.Fatal(err)
	}
	p, _ := componentPage("button")
	if got := tt.Clipboard(); got != p.importStatement() {
		t.Fatalf("copied import = %q, want %q", got, p.importStatement())
	}
}

func TestComponentExampleFrames(t *testing.T) {
	p := exampleFramesFixture()
	state := New()
	tt := ui.NewTester(func(c *ui.Context) { buildComponentPage(c, state, p) }, 1180, 1400)
	tt.Frame()
	for _, example := range p.Examples {
		title, ok := tt.Find(example.Title)
		if !ok {
			t.Fatalf("missing example title %q", example.Title)
		}
		panel, ok := tt.Find("Example panel: " + example.Title)
		if !ok || title.Y+title.H > panel.Y {
			t.Fatalf("%q title and preview panel: title %+v, panel %+v", example.Title, title, panel)
		}
	}
	if !tt.HasText("Primary") || !tt.HasText("Ghost") {
		t.Fatal("preview omitted an interactive specimen")
	}
	state.Panel = 1
	tt.Frame()
	if !tt.HasText("Go code") || !tt.HasText("package main") {
		t.Fatal("source view omitted the snippet")
	}
	if tt.HasText("Example panel: Button variants") {
		t.Fatal("preview remains mounted in source view")
	}
}

func TestUsageBarCopiesDerivedImport(t *testing.T) {
	p, ok := componentPage("input")
	if !ok {
		t.Fatal("input page is not registered")
	}
	state := New()
	state.Panel = 1
	tt := ui.NewTester(func(c *ui.Context) { buildComponentPage(c, state, p) }, 1180, 1400)
	tt.Frame()
	bar, ok := tt.Find("Usage import")
	if !ok || bar.H <= 48 || !tt.HasText(p.importStatement()) {
		t.Fatalf("multi-line import bar is incomplete: %+v", bar)
	}
	if err := tt.Click("Copy import line"); err != nil {
		t.Fatal(err)
	}
	if got := tt.Clipboard(); got != p.importStatement() {
		t.Fatalf("copied import = %q, want %q", got, p.importStatement())
	}
}

func TestComponentPageImportStatement(t *testing.T) {
	for _, tc := range []struct{ name, in, want string }{
		{"default", "", `import "github.com/Cyenoch/mygoui"`},
		{"path", "github.com/Cyenoch/mygoui/form", `import "github.com/Cyenoch/mygoui/form"`},
		{"statement", `import "github.com/Cyenoch/mygoui/list"`, `import "github.com/Cyenoch/mygoui/list"`},
	} {
		if got := (ComponentPage{Import: tc.in}).importStatement(); got != tc.want {
			t.Fatalf("%s: importStatement() = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestComponentPageImportStatementDerivesFromCode(t *testing.T) {
	for _, tc := range []struct {
		id   string
		want []string
	}{
		{"input", []string{`"github.com/Cyenoch/mygoui"`, `"github.com/Cyenoch/mygoui/form"`}},
		{"dialog", []string{`"github.com/Cyenoch/mygoui"`, `"github.com/Cyenoch/mygoui/overlay"`}},
		{"typography", []string{`"github.com/Cyenoch/mygoui"`, `"github.com/Cyenoch/mygoui/theme"`}},
		{"link", []string{`"github.com/Cyenoch/mygoui"`, `"github.com/Cyenoch/mygoicons/lucide"`}},
	} {
		p, ok := componentPage(tc.id)
		if !ok {
			t.Fatalf("%s: page is not registered", tc.id)
		}
		got := p.importStatement()
		if !strings.HasPrefix(got, "import (\n") || !strings.HasSuffix(got, "\n)") {
			t.Fatalf("%s: invalid import block:\n%s", tc.id, got)
		}
		for _, want := range tc.want {
			if count := strings.Count(got, want); count != 1 {
				t.Fatalf("%s: import block names %s %d times", tc.id, want, count)
			}
		}
	}
}

func TestCopyControlsFitTheirRows(t *testing.T) {
	for _, width := range []int{1180, 640} {
		state := New()
		state.Panel = 1
		p := exampleFramesFixture()
		tt := ui.NewTester(func(c *ui.Context) { buildComponentPage(c, state, p) }, width, 1400)
		tt.Frame()
		for _, item := range []struct{ row, control string }{
			{"Usage import", "Copy import line"},
			{"Go code header", "Copy code"},
		} {
			row, rowOK := tt.Find(item.row)
			button, buttonOK := tt.Find(item.control)
			if !rowOK || !buttonOK || button.W <= 0 || button.H <= 0 ||
				button.X < row.X || button.X+button.W > row.X+row.W ||
				button.Y < row.Y || button.Y+button.H > row.Y+row.H {
				t.Fatalf("width %d: %s outside %s: button %+v row %+v", width, item.control, item.row, button, row)
			}
		}
	}
}
