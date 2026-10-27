package form

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestInputGroupTypesAndRunsSuffixAction(t *testing.T) {
	var value string
	cleared := 0
	tester := formTextTester(t, func(c *ui.Context) {
		InputGroup(c, &value, InputGroupOptions{
			Label:       "Amount",
			Placeholder: "0.00",
			Prefix: func() {
				ui.Text(c, "$")
			},
			Suffix: func() {
				ui.ButtonBase(c).Label("Clear").Padding(8).OnClick(func() {
					value = ""
					cleared++
				})
			},
		})
	})
	formTextClick(t, tester, "Amount")
	tester.Type("42")
	if value != "42" {
		t.Fatalf("input group committed %q", value)
	}
	formTextClick(t, tester, "Clear")
	if value != "" || cleared != 1 {
		t.Fatalf("suffix action: value %q, cleared %d", value, cleared)
	}
}

func TestInputGroupBlockPrefixKeepsEditorAsFieldControl(t *testing.T) {
	var value string
	activated := 0
	tester := formTextTester(t, func(c *ui.Context) {
		InputGroup(c, &value, InputGroupOptions{
			Label: "Amount",
			BlockPrefix: func() {
				ui.Text(c, "Block above")
				ui.Button(c, "Add").OnClick(func() { activated++ })
			},
		})
	})
	block, ok := tester.Find("Block above")
	if !ok {
		t.Fatal("input group shows no block prefix")
	}
	// Click the label text at the field's top edge: the label names and
	// focuses the editor, not the Add button that sits above it.
	tester.ClickAt(24, 24)
	caret, focused := tester.TextCaret()
	tester.Type("42")
	if value != "42" || activated != 0 {
		t.Fatalf("label click did not focus the editor: value %q, activated %d", value, activated)
	}
	if !focused || block.Y >= caret.Y {
		t.Fatalf("block prefix is not above the editor: block %v, caret %v", block, caret)
	}
}

func TestInputGroupAddonFocusesEditor(t *testing.T) {
	var value string
	tester := formTextTester(t, func(c *ui.Context) {
		InputGroup(c, &value, InputGroupOptions{
			Label:  "Amount",
			Prefix: func() { ui.Text(c, "$") },
		})
	})
	formTextClick(t, tester, "$")
	tester.Type("42")
	if value != "42" {
		t.Fatalf("typing after addon click = %q", value)
	}
}
