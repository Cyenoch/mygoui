package gallery

import (
	"bytes"
	"encoding/json"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

// Snippets use caller-owned state and application actions. Compile them
// against real exports so renamed options and unrelated enum types cannot
// silently turn the gallery into unusable documentation.
func TestComponentSnippetsTypeCheck(t *testing.T) {
	packages := []string{"fmt", "time", "github.com/egoist/mygo/ui"}
	for _, path := range qualifierImports {
		packages = append(packages, path)
	}
	args := append([]string{"list", "-deps", "-export", "-json"}, packages...)
	output, err := exec.Command("go", args...).Output()
	if err != nil {
		t.Fatal(err)
	}
	exports := map[string]string{}
	decoder := json.NewDecoder(bytes.NewReader(output))
	for decoder.More() {
		var item struct{ ImportPath, Export string }
		if err := decoder.Decode(&item); err != nil {
			t.Fatal(err)
		}
		exports[item.ImportPath] = item.Export
	}
	fset := token.NewFileSet()
	imp := importer.ForCompiler(fset, "gc", func(path string) (io.ReadCloser, error) {
		return os.Open(exports[path])
	})
	var imports strings.Builder
	imports.WriteString("import (\n")
	for _, path := range packages {
		imports.WriteString(`"` + path + "\"\n")
	}
	imports.WriteString(")\n")
	for _, page := range Pages() {
		p, ok := componentPage(page.ID)
		if !ok {
			continue
		}
		t.Run(p.ID, func(t *testing.T) {
			source := p.Code
			if !strings.HasPrefix(source, "package ") {
				source = "package snippets\n" + imports.String() + snippetState + snippetStateFor(p.ID) +
					"\nfunc build(c *ui.Context) {\n" + source + "\n}\n"
			}
			file, err := parser.ParseFile(fset, p.ID+".go", source, 0)
			if err != nil {
				t.Fatal(err)
			}
			cfg := types.Config{
				Importer: imp, DisableUnusedImportCheck: true,
				Error: func(err error) { t.Error(err) },
			}
			if _, err := cfg.Check("snippets/"+p.ID, fset, []*ast.File{file}, nil); err != nil {
				t.Log(source)
			}
		})
	}
}

const snippetState = `
var email, name, slug, username, query, bio, price, draft, sent, last, destination string
var invalid, updates, releaseNotes, wifi, analytics, autoSave, day, week, enabled, collapsed bool
var volume float64
var save, back, next, openDocs, showAll, copyLink, copySelection, dismiss, submit, invite, deleteProject, revokeInvite func()
var remove func(string)
var report func([]string)
var selectCrumb func(int)
var verify func(string)
var resetResults func()
var muted func(string)
var maxPriceButton func()
var shortcutRow func(*ui.Context, string, string)
`

func snippetStateFor(id string) string {
	switch id {
	case "listgroup":
		return "var open func(string)"
	case "togglebutton":
		return "var bold bool"
	case "toolbar":
		return "var bold func()"
	case "select":
		return "var fruits []list.ListBoxItem; var fruit form.SelectState"
	case "combobox":
		return "var languages []list.ListBoxItem; var language form.ComboBoxState"
	case "autocomplete":
		return "var cities []list.ListBoxItem; var city form.AutocompleteState"
	case "listbox":
		return "var fruits []list.ListBoxItem; var fruit list.ListBoxState"
	case "taggroup":
		return "var state struct{Tags mygoui.TagGroupState}"
	case "scrollshadow":
		return "var state struct{Feed ui.ScrollState}"
	case "scrollarea", "scrollfade":
		return "var scroll ui.ScrollState"
	case "command":
		return "var commands navigation.CommandState"
	case "segmentedcontrol":
		return "var selected int"
	case "previewcard":
		return "var preview overlay.PreviewCardState"
	case "sheet", "colorpicker":
		return "var open bool; var state color.ColorState"
	case "checkboxgroup":
		return "var checks form.CheckboxGroupState"
	case "radio", "radiogroup":
		return "var plan form.RadioGroupState[string]"
	case "inputotp":
		return "var note func(int); var otp form.OTPState"
	case "numberfield":
		return "var seats form.NumberFieldState"
	case "drawer":
		return "var note string"
	case "calendar", "rangecalendar", "datefield", "datepicker", "daterangepicker":
		return `var today, min, max, launch, deploy time.Time
var cal date.CalendarState
var rangeState date.RangeCalendarState
var window date.DateRange
var isWeekend func(time.Time) bool
var open bool`
	case "timefield":
		return "var start date.TimeFieldState"
	case "colorarea", "colorslider", "colorfield", "colorswatch", "colorswatchpicker":
		return "var state color.ColorState"
	case "table":
		return `var columns []data.TableColumnOptions
var pageRows []struct{ID string}
var state data.TableState
var page, pages int
var summary string`
	}
	return ""
}

// TestCodeBlockCopyControl pins the compact header and the copy feedback: the
// header stays 40 DIP or shorter, Copy writes the untouched snippet, the
// button reports Copied, and it reverts after the feedback window.
func TestCodeBlockCopyControl(t *testing.T) {
	const code = "package main\n\nfunc main() {\n\tprintln(\"hi\")\n}\n"
	p := ComponentPage{ID: "copyblock", Title: "Copy block", Group: "Foundations", Code: code}
	state := New()
	state.Panel = 1
	tt := ui.NewTester(func(c *ui.Context) { buildComponentPage(c, state, p) }, 700, 700)
	tt.Frame()

	header, ok := tt.Find("Go code header")
	if !ok || header.W <= 0 || header.H <= 0 || header.H > 40 {
		t.Fatalf("code header is not compact (rect %+v)", header)
	}
	if _, ok := tt.Find("Copy code"); !ok {
		t.Fatal("code header has no Copy code button")
	}
	if err := tt.Click("Copy code"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if got := tt.Clipboard(); got != code {
		t.Fatalf("clipboard = %q, want the raw snippet %q", got, code)
	}
	if _, ok := tt.Find("Copied"); !ok {
		t.Fatal("Copy button does not report Copied")
	}

	tt.Advance(codeCopyFeedback)
	if _, ok := tt.Find("Copy code"); !ok {
		t.Fatal("Copy code button did not revert after the feedback window")
	}
	if _, ok := tt.Find("Copied"); ok {
		t.Fatal("Copy button still says Copied after the feedback window")
	}
}

// TestCodeBlockHugsSnippet proves the panel adds only its compact header and
// 12 DIP body padding, with no extra blank line under the last source line.
func TestCodeBlockHugsSnippet(t *testing.T) {
	code := strings.Join([]string{
		"mygoui.Button(c, mygoui.ButtonOptions{",
		"    Label: \"Save\",",
		"    Variant: mygoui.ButtonPrimary,",
		"})",
		"",
	}, "\n")
	p := ComponentPage{ID: "hugblock", Title: "Hug block", Group: "Foundations", Code: code}
	state := New()
	state.Panel = 1
	tt := ui.NewTester(func(c *ui.Context) { buildComponentPage(c, state, p) }, 700, 700)
	tt.Frame()
	block, ok := tt.Find("Go code")
	if !ok || block.H <= 0 {
		t.Fatalf("code block is missing (rect %+v)", block)
	}
	lines := float32(len(strings.Split(strings.TrimRight(code, "\n"), "\n")))
	if block.H >= lines*24+60 {
		t.Fatalf("code block wastes space: height %v for %v lines (rect %+v)", block.H, lines, block)
	}
}
