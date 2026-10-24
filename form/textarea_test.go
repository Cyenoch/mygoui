package form

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestTextAreaTypesMultipleLines(t *testing.T) {
	var value string
	tester := formTextTester(t, func(c *ui.Context) {
		TextArea(c, &value, TextAreaOptions{Label: "Notes", MinLines: 2, MaxLines: 5})
	})
	formTextClick(t, tester, "Notes")
	tester.Type("one\ntwo")
	if value != "one\ntwo" {
		t.Fatalf("textarea committed %q", value)
	}
}

func TestTextAreaReadOnlyKeepsValue(t *testing.T) {
	value := "kept"
	tester := formTextTester(t, func(c *ui.Context) {
		TextArea(c, &value, TextAreaOptions{Label: "Notes", ReadOnly: true})
	})
	formTextClick(t, tester, "Notes")
	tester.Type("more")
	if value != "kept" {
		t.Fatalf("read-only textarea committed %q", value)
	}
}
