package form

import (
	"testing"

	"github.com/Cyenoch/mygoui/locale"
	"github.com/egoist/mygo/ui"
)

func TestPathFieldBrowseCallsCallerAndKeepsValue(t *testing.T) {
	value, browsed := "/tmp/report.pdf", 0
	tester := formTextTester(t, func(c *ui.Context) {
		PathField(c, &value, PathFieldOptions{Label: "Output", OnBrowse: func() { browsed++ }})
	})
	formTextClick(t, tester, locale.M().Browse)
	if browsed != 1 || value != "/tmp/report.pdf" {
		t.Fatalf("browse ran %d times and left %q", browsed, value)
	}
}

// The resting label that truncates a long path must not take the editor's
// input or its accessible value.
func TestPathFieldEditsAndReadsFullPathWhileTruncated(t *testing.T) {
	value := "/Users/ada/Projects/mygoui/internal/gallery/component_pages_form_inputs"
	changes := 0
	tester := formTextTester(t, func(c *ui.Context) {
		PathField(c, &value, PathFieldOptions{Label: "Output", OnChange: func() { changes++ }})
	})
	info, ok := tester.AccessByRole(ui.RoleTextField, "Output")
	if !ok || info.Value != value {
		t.Fatalf("resting field reads %q (found %v), want the whole path", info.Value, ok)
	}
	formTextClick(t, tester, "Output")
	tester.Key(0, ui.KeyEnd)
	tester.Type(".go")
	if want := "/Users/ada/Projects/mygoui/internal/gallery/component_pages_form_inputs.go"; value != want || changes == 0 {
		t.Fatalf("typed path = %q after %d changes", value, changes)
	}
}

func TestPathFieldReadOnlyStillBrowses(t *testing.T) {
	value, browsed := "/tmp/a", 0
	tester := formTextTester(t, func(c *ui.Context) {
		PathField(c, &value, PathFieldOptions{Label: "Output", ReadOnly: true, OnBrowse: func() { browsed++ }})
	})
	formTextClick(t, tester, "Output")
	tester.Type("x")
	formTextClick(t, tester, locale.M().Browse)
	if value != "/tmp/a" || browsed != 1 {
		t.Fatalf("read-only field holds %q after %d browses", value, browsed)
	}
}

func TestPathFieldDisabledDoesNotBrowse(t *testing.T) {
	value, browsed := "/tmp/a", 0
	tester := formTextTester(t, func(c *ui.Context) {
		PathField(c, &value, PathFieldOptions{Label: "Output", Disabled: true, OnBrowse: func() { browsed++ }})
	})
	_ = tester.Click(locale.M().Browse)
	tester.Frame()
	if browsed != 0 {
		t.Fatal("a disabled path field browsed")
	}
}

func TestPathFieldDropSetsPathAndReports(t *testing.T) {
	value, changes := "/tmp/old", 0
	readOnly := false
	tester := formTextTester(t, func(c *ui.Context) {
		PathField(c, &value, PathFieldOptions{Label: "Output", ReadOnly: readOnly, OnChange: func() { changes++ }})
	})
	rect, ok := tester.Find("/tmp/old")
	if !ok {
		t.Fatal("the path is not shown")
	}
	x, y := rect.X+rect.W/2, rect.Y+rect.H/2
	if !tester.DropFiles(x, y, "/data/a.csv", "/data/b.csv") {
		t.Fatal("the field did not take a dropped file")
	}
	tester.Frame()
	if value != "/data/a.csv" || changes != 1 {
		t.Fatalf("drop left %q after %d changes, want the first file", value, changes)
	}
	readOnly = true
	tester.Frame()
	if tester.DropFiles(x, y, "/data/c.csv") {
		t.Fatal("a read-only field took a drop")
	}
}
