package chart_test

import (
	"reflect"
	"slices"
	"testing"

	"github.com/Cyenoch/mygoui/chart"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

func TestChartLocaleChangeRefreshesCachedAccessibleReading(t *testing.T) {
	data := chart.Data{
		Labels: []string{"Cancel", "原始分类"},
		Series: []chart.Series{{ID: "raw", Label: "原始系列", Values: []chart.Value{{Number: 1234.5}, {Number: 2}}}},
	}
	state := chart.State{Selected: -1}
	selections := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		chart.Chart(c, &state, chart.Options{
			Label: "Cancel", Data: data,
			OnSelection: func(int) { selections++ },
		})
	}, 420, 310)
	tt.SetLocale("en-US")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyHome)
	reading := func(locale string) string {
		return "Cancel · 原始系列: " + i18n.FormatNumber(locale, 1234.5, i18n.NumberOptions{MaxFractionDigits: 100})
	}
	english := reading("en-US")
	german := reading("de-DE")
	if english == german {
		t.Fatal("fixture needs distinct regional number formats")
	}
	if !tt.HasText(english) || !slices.Contains(tt.Announcements(), english) {
		t.Fatalf("initial keyboard reading missing: %v", tt.Texts())
	}
	beforeSelections := selections
	tt.SetLocale("de-DE")
	if state.Selected != 0 || selections != beforeSelections {
		t.Fatalf("locale refresh changed selection: index=%d callbacks=%d, want 0 and %d", state.Selected, selections, beforeSelections)
	}
	if !tt.HasText(german) || tt.HasText(english) {
		t.Fatalf("unchanged chart retained its cached reading: %v", tt.Texts())
	}
	// The native selection callback uses this same cached description for its
	// accessible value and announces it when keyboard navigation returns.
	tt.Announcements()
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyHome)
	if announcements := tt.Announcements(); !slices.Contains(announcements, german) || slices.Contains(announcements, english) {
		t.Fatalf("keyboard announced a stale description: %v", announcements)
	}
	tt.SetLocale("en-US")
	if !tt.HasText(english) || state.Selected != 0 {
		t.Fatalf("returning to English lost the reading or selection: %v", tt.Texts())
	}
	want := chart.Data{
		Labels: []string{"Cancel", "原始分类"},
		Series: []chart.Series{{ID: "raw", Label: "原始系列", Values: []chart.Value{{Number: 1234.5}, {Number: 2}}}},
	}
	if !reflect.DeepEqual(data, want) {
		t.Fatalf("locale refresh modified caller data: %+v", data)
	}
}
