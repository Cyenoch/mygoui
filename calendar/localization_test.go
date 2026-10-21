package calendar_test

import (
	"os"
	"os/exec"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/calendar"
	locales "github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

func TestDatePickerLocaleChangePreservesInvalidDrafts(t *testing.T) {
	const helper = "MYGOUI_DATE_PICKER_LOCALE_TEST"
	if os.Getenv(helper) != "1" {
		// SetLocaleFunc has no getter. Isolate its override so the callback
		// installed by MyGo, including its pre-start behavior, stays untouched.
		executable, err := os.Executable()
		if err != nil {
			t.Fatal(err)
		}
		cmd := exec.CommandContext(t.Context(), executable, "-test.run=^TestDatePickerLocaleChangePreservesInvalidDrafts$", "-test.count=1")
		cmd.Env = append(os.Environ(), helper+"=1")
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("isolated locale regression: %v\n%s", err, output)
		}
		return
	}

	previous := i18n.Locale()
	var language atomic.Value
	language.Store("en")
	i18n.SetLocaleFunc(func() string { return language.Load().(string) })
	t.Cleanup(func() { i18n.SetLocaleFunc(func() string { return previous }) })
	selected := calendar.Range{
		Start: calendar.Date{Year: 2026, Month: time.October, Day: 12},
		End:   calendar.Date{Year: 2026, Month: time.October, Day: 20},
	}
	state := calendar.DatePickerState{Calendar: calendar.State{Range: selected}}
	changes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		calendar.DatePicker(c, &state, calendar.DatePickerOptions{
			Label: "Cancel", Range: true, OnChange: func() { changes++ },
		})
	}, 480, 300)
	tt.SetLocale("en")
	for _, input := range []struct {
		label string
		draft string
	}{
		{label: locales.In("en").ComponentDateRangeStart("Cancel"), draft: "Cancel"},
		{label: locales.In("en").ComponentDateRangeEnd("Cancel"), draft: "原始日期"},
	} {
		if err := tt.Click(input.label); err != nil {
			t.Fatal(err)
		}
		tt.Command("selectAll")
		tt.Type(input.draft)
	}
	english := locales.In("en").ComponentEnterISODate()
	chinese := locales.In("zh-Hans").ComponentEnterISODate()
	if english == chinese {
		t.Fatal("fixture needs distinct validation translations")
	}
	for _, locale := range []string{"en", "zh-Hans", "en"} {
		language.Store(locale)
		tt.SetLocale(locale)
		wantError := locales.In(locale).ComponentEnterISODate()
		if state.Invalid != wantError || state.EndInvalid != wantError || !tt.HasText(wantError) {
			t.Fatalf("%s validation was not refreshed: start=%q end=%q texts=%v", locale, state.Invalid, state.EndInvalid, tt.Texts())
		}
		if state.Text != "Cancel" || state.EndText != "原始日期" || state.Calendar.Range != selected || changes != 0 {
			t.Fatalf("%s refresh replaced an invalid draft or valid selection: %+v callbacks=%d", locale, state, changes)
		}
		if !tt.Focused(locales.In(locale).ComponentDateRangeEnd("Cancel")) {
			t.Fatalf("%s refresh lost the active end editor", locale)
		}
	}
	// The retained editor remains usable after both language changes.
	tt.Command("selectAll")
	tt.Type("2026-10-21")
	if state.EndInvalid != "" || state.EndText != "2026-10-21" || state.Calendar.Range.End.Day != 21 || changes != 1 {
		t.Fatalf("correcting retained draft failed: %+v callbacks=%d", state, changes)
	}
	if state.Text != "Cancel" || state.Invalid != english {
		t.Fatal("correcting the end editor discarded the other invalid draft")
	}
}
