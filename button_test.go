package mygoui_test

import (
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	myui "github.com/egoist/mygo/ui"
)

func TestButtonActivationAndUnavailableStates(t *testing.T) {
	calls := 0
	loading, disabled := false, false
	var handle myui.Handle
	tt := myui.NewTester(func(c *myui.Context) {
		theme.Apply(c)
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Save", Handle: &handle, Loading: loading, Disabled: disabled,
			OnClick: func() { calls++ },
		})
	}, 320, 160)
	if err := tt.Click("Save"); err != nil {
		t.Fatal(err)
	}
	if calls != 1 || !tt.Focused("Save") {
		t.Fatalf("pointer activation calls=%d focus=%v", calls, tt.Focused("Save"))
	}
	tt.Key(0, myui.KeyEnter)
	tt.Key(0, myui.KeySpace)
	if calls != 3 {
		t.Fatalf("keyboard activations=%d, want3", calls)
	}
	loading = true
	tt.Frame()
	if err := tt.Click("Save"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, myui.KeyEnter)
	if calls != 3 {
		t.Fatalf("loading control activated: %d", calls)
	}
	loading, disabled = false, true
	tt.Frame()
	if err := tt.Click("Save"); err != nil {
		t.Fatal(err)
	}
	if calls != 3 {
		t.Fatalf("disabled control activated: %d", calls)
	}
	disabled = false
	tt.Frame()
	if err := tt.Click("Save"); err != nil {
		t.Fatal(err)
	}
	if calls != 4 {
		t.Fatalf("reenabled control did not activate once: %d", calls)
	}
}

func TestButtonWithoutActionIsNotATabStop(t *testing.T) {
	tt := myui.NewTester(func(c *myui.Context) {
		theme.Apply(c)
		myui.Row(c).Children(func() {
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Unavailable"})
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Continue", OnClick: func() {}})
		})
	}, 320, 160)
	tt.Key(0, myui.KeyTab)
	if !tt.Focused("Continue") || tt.Focused("Unavailable") {
		t.Fatal("missing action remained an enabled keyboard target")
	}
}

func TestButtonScalesToKeepLargeTextInsideTarget(t *testing.T) {
	tt := myui.NewTester(func(c *myui.Context) {
		theme.Apply(c)
		mygoui.Button(c, mygoui.ButtonOptions{Label: "Save", Size: mygoui.ButtonXS, OnClick: func() {}})
	}, 320, 160)
	before, ok := tt.Find("Save")
	if !ok {
		t.Fatal("missing button")
	}
	tt.SetPreferences(myui.Preferences{TextScale: 2})
	after, ok := tt.Find("Save")
	if !ok || after.H <= before.H {
		t.Fatalf("large text did not expand target: before=%v after=%v", before, after)
	}
}

func TestButtonComposedContentKeepsNativeActivationAndLoading(t *testing.T) {
	calls := 0
	loading := false
	tt := myui.NewTester(func(c *myui.Context) {
		theme.Apply(c)
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Accessible action", Loading: loading, OnClick: func() { calls++ },
			Content: func() { myui.Text(c, "Custom body") },
		})
	}, 320, 160)
	if err := tt.Click("Accessible action"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, myui.KeyEnter)
	tt.Key(0, myui.KeySpace)
	if calls != 3 || !tt.Focused("Accessible action") {
		t.Fatalf("custom button lost native activation/focus: %d", calls)
	}
	loading = true
	tt.Frame()
	if _, ok := tt.Find("Loading"); !ok || !tt.HasText("Custom body") {
		t.Fatal("loading custom button lost spinner or content")
	}
	if err := tt.Click("Accessible action"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, myui.KeyEnter)
	if calls != 3 {
		t.Fatal("loading custom button activated")
	}
}

func TestButtonWrapUsesAvailableWidthAndCompleteLabel(t *testing.T) {
	const label = "Keep this complete action readable across several lines"
	calls := 0
	tt := myui.NewTester(func(c *myui.Context) {
		theme.Apply(c)
		myui.Column(c).Width(150).Children(func() {
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: label, Wrap: true, Size: mygoui.ButtonXS, OnClick: func() { calls++ },
			}).MaxWidthPercent(100).PaddingY(6)
		})
	}, 200, 220)
	r, ok := tt.Find(label)
	if !ok || r.W > 150 || r.H <= 32 || !tt.HasText(label) {
		t.Fatalf("wrapping action lost width/height/full label: %+v", r)
	}
	if err := tt.Click(label); err != nil || calls != 1 {
		t.Fatalf("wrapping action failed: calls=%d err=%v", calls, err)
	}
}
