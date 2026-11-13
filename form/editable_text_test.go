package form

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func editableTester(name *string, changes *int) *ui.Tester {
	other := ""
	return ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(8).Children(func() {
			EditableText(c, name, EditableTextOptions{
				AriaLabel: "Name",
				OnChange:  func() { *changes++ },
			})
			Input(c, &other, InputOptions{AriaLabel: "Other"})
		})
	}, 360, 200)
}

func TestEditableTextCommitsOnEnter(t *testing.T) {
	name, changes := "Draft", 0
	tt := editableTester(&name, &changes)
	if err := tt.Click("Name"); err != nil {
		t.Fatal(err)
	}
	tt.Type(" two")
	if name != "Draft" {
		t.Fatalf("typing wrote %q before the commit", name)
	}
	tt.Key(0, ui.KeyEnter)
	if name != "Draft two" || changes != 1 {
		t.Fatalf("Enter left %q with %d changes, want \"Draft two\" and 1", name, changes)
	}
	// The text is shown again and keeps the focus, so F2 edits it anew.
	tt.Key(0, ui.KeyF2)
	tt.Type("!")
	tt.Key(0, ui.KeyEnter)
	if name != "Draft two!" {
		t.Fatalf("F2 edit left %q, want \"Draft two!\"", name)
	}
}

func TestEditableTextCommitsOnBlur(t *testing.T) {
	name, changes := "Draft", 0
	tt := editableTester(&name, &changes)
	if err := tt.Click("Name"); err != nil {
		t.Fatal(err)
	}
	tt.Type("ed")
	if err := tt.Click("Other"); err != nil {
		t.Fatal(err)
	}
	if name != "Drafted" || changes != 1 {
		t.Fatalf("blur left %q with %d changes, want \"Drafted\" and 1", name, changes)
	}
}

func TestEditableTextEscapeCancels(t *testing.T) {
	name, changes := "Draft", 0
	tt := editableTester(&name, &changes)
	if err := tt.Click("Name"); err != nil {
		t.Fatal(err)
	}
	tt.Type(" discarded")
	tt.Key(0, ui.KeyEscape)
	if name != "Draft" || changes != 0 {
		t.Fatalf("Escape left %q with %d changes, want the value unchanged", name, changes)
	}
	if !tt.HasText("Draft") {
		t.Fatal("the text is not shown again after Escape")
	}
}
