package form

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestTextFieldTypesAndSubmits(t *testing.T) {
	var value string
	submitted, changed := 0, 0
	tester := formTextTester(t, func(c *ui.Context) {
		TextField(c, TextFieldOptions{
			Value:       &value,
			Label:       "City",
			Description: "Town or city",
			Required:    true,
			FullWidth:   true,
			OnSubmit:    func() { submitted++ },
			OnChange:    func() { changed++ },
		})
	})
	if !tester.HasText("Town or city") {
		t.Fatal("text field shows no description")
	}
	formTextClick(t, tester, "City")
	tester.Type("Paris")
	if value != "Paris" || changed != 1 {
		t.Fatalf("typed %q after %d changes", value, changed)
	}
	tester.Key(0, ui.KeyEnter)
	if submitted != 1 {
		t.Fatalf("submitted %d times", submitted)
	}
}

func TestTextFieldInvalidHidesDescription(t *testing.T) {
	var value string
	tester := formTextTester(t, func(c *ui.Context) {
		TextField(c, TextFieldOptions{
			Value:       &value,
			Label:       "Email",
			Description: "We never share it",
			Error:       "Required",
			Invalid:     true,
		})
	})
	if !tester.HasText("Required") {
		t.Fatal("text field shows no error")
	}
	if tester.HasText("We never share it") {
		t.Fatal("invalid text field kept its description")
	}
}
