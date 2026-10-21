package form

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestInputSelectionPasteUndoAndIME(t *testing.T) {
	value := "A😀café"
	var handle ui.Handle
	readonly, disabled := false, false
	composing := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		input := Input(c, &value, InputOptions{ControlOptions: ControlOptions{
			Label: "Name", Handle: &handle, ReadOnly: readonly, Disabled: disabled,
		}}).Width(280)
		composing = input.Composing()
	}, 320, 120)
	tt.Click("Name")
	tt.Command("selectAll")
	tt.SetClipboard("Ada\nLovelace")
	tt.Command("paste")
	if value != "Ada Lovelace" {
		t.Fatalf("single-line native paste: %q", value)
	}
	tt.Command("undo")
	if value != "A😀café" {
		t.Fatalf("native undo lost original unicode: %q", value)
	}
	tt.Command("selectAll")
	tt.Compose("ni", 2)
	// Native preedit replaces the selection immediately but keeps marked text
	// outside the committed value until the input method accepts a candidate.
	if !composing || value != "" {
		t.Fatalf("preedit was committed or selection was not replaced: %q, composing=%v", value, composing)
	}
	tt.Type("你")
	if value != "你" {
		t.Fatalf("IME commit: %q", value)
	}
	tt.Command("undo")
	if value != "A😀café" {
		t.Fatalf("IME replacement lost its one-step undo transaction: %q", value)
	}
	tt.Command("redo")
	if value != "你" {
		t.Fatalf("IME replacement redo: %q", value)
	}
	readonly = true
	tt.Frame()
	tt.Type("wrong")
	tt.Command("selectAll")
	tt.Command("copy")
	if value != "你" || tt.Clipboard() != "你" {
		t.Fatalf("read-only must preserve selection/copy without editing: %q / %q", value, tt.Clipboard())
	}
	disabled = true
	tt.Frame()
	tt.Click("Name")
	tt.Type("wrong")
	if value != "你" || tt.Focused("Name") {
		t.Fatal("disabled input accepted focus or changed value")
	}
}

func TestPasswordNeverCopiesSecret(t *testing.T) {
	value := "secret😀"
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Input(c, &value, InputOptions{ControlOptions: ControlOptions{Label: "Password"}, Password: true})
	}, 320, 100)
	tt.Click("Password")
	tt.Command("selectAll")
	tt.SetClipboard("sentinel")
	tt.Command("copy")
	if tt.Clipboard() != "sentinel" {
		t.Fatalf("password exposed clipboard data: %q", tt.Clipboard())
	}
	tt.Command("cut")
	if tt.Clipboard() != "sentinel" || value != "secret😀" {
		t.Fatal("password cut exposed or removed the secret")
	}
}

func TestFieldLabelTargetsGroupedEditorBeforeAddonActions(t *testing.T) {
	value := ""
	clicks := 0
	disabled := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Field(c, FieldOptions{
			Label: "Account", Description: "Use your account name", Error: "Account is required", Disabled: disabled,
		}, func() {
			InputGroup(c, &value, InputGroupOptions{
				InputOptions: InputOptions{},
				BlockStart: func() {
					InputGroupButton(c, mygoui.ButtonOptions{Label: "Help", OnClick: func() { clicks++ }})
				},
				InlineStart: func() { InputGroupText(c, "@") },
			})
		})
	}, 420, 200)
	tt.Click("Account")
	tt.Type("ada")
	if value != "ada" || clicks != 0 {
		t.Fatalf("field label targeted an addon instead of editor: %q, %d", value, clicks)
	}
	tt.Click("Help")
	if clicks != 1 {
		t.Fatal("addon button did not invoke real callback")
	}
	disabled = true
	tt.Frame()
	tt.Advance(time.Second)
	tt.Click("Account")
	tt.Type("wrong")
	tt.Click("Help")
	if value != "ada" || clicks != 1 {
		t.Fatal("disabled field left input or addon active")
	}
}

func TestTextareaKeepsMultilinePaste(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Textarea(c, &value, TextareaOptions{InputOptions: InputOptions{
			ControlOptions: ControlOptions{Label: "Notes"},
		}, MinLines: 2, MaxLines: 4}).Width(220)
	}, 260, 180)
	tt.Click("Notes")
	tt.SetClipboard("first\nsecond\n")
	tt.Command("paste")
	if value != "first\nsecond\n" {
		t.Fatalf("textarea paste changed content: %q", value)
	}
	tt.Command("undo")
	if value != "" {
		t.Fatalf("textarea paste not one native undo step: %q", value)
	}
}
