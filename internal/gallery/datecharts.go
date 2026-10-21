package gallery

import (
	"strconv"
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/calendar"
	"github.com/Cyenoch/mygoui/chart"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

type DateChartState struct {
	Calendar                                                  calendar.State
	Picker                                                    calendar.DatePickerState
	Chart                                                     chart.State
	Mode                                                      calendar.Mode
	Kind                                                      chart.ChartKind
	Stack                                                     chart.StackMode
	Weekends, Monday, ReadOnly, Range, German, Empty, Missing bool
}

func NewDateChartState() DateChartState {
	date := calendar.Date{Year: 2026, Month: time.October, Day: 20}
	return DateChartState{Calendar: calendar.State{Date: date, Month: date}, Picker: calendar.DatePickerState{Calendar: calendar.State{Date: date, Month: date}}, Chart: chart.State{Selected: -1}}
}

var sampleChartData = chart.Data{Labels: []string{"May", "Jun", "Jul", "Aug", "Sep", "Oct"}, Series: []chart.Series{
	{ID: "desktop", Label: "Desktop", Color: chart.Chart1, Values: []chart.Value{{Number: 24}, {Number: 32}, {Number: 28}, {Number: 40}, {Number: 36}, {Number: 48}}},
	{ID: "mobile", Label: "Mobile", Color: chart.Chart2, Values: []chart.Value{{Number: 12}, {Number: 18}, {Number: 16}, {Number: 22}, {Number: 24}, {Number: 30}}},
}}

func PreviewDateChart(c *ui.Context, s *DateChartState, id ComponentID) {
	ui.Column(c).FillWidth().MinWidth(0).Gap(16).Children(func() {
		ui.Row(c).FillWidth().Wrap().Gap(8).Children(func() {
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Reset sample", Variant: mygoui.ButtonOutline, OnClick: func() { *s = NewDateChartState() }})
			if id != "chart" {
				mygoui.Toggle(c, &s.ReadOnly, mygoui.ToggleOptions{Label: "Read only"})
				mygoui.Toggle(c, &s.Weekends, mygoui.ToggleOptions{Label: "Disable weekends"})
				mygoui.Toggle(c, &s.Monday, mygoui.ToggleOptions{Label: "Monday first"})
				mygoui.Toggle(c, &s.German, mygoui.ToggleOptions{Label: "German locale"})
			}
		})
		options := calendar.Options{Label: "Sample calendar", ReadOnly: s.ReadOnly, MonthChoice: true, StartYear: 1900, EndYear: 2026, Mode: s.Mode}
		if s.German {
			options.Locale = "de"
		}
		if s.Monday {
			day := time.Monday
			options.FirstWeekday = &day
		}
		if s.Weekends {
			options.Disabled = func(d calendar.Date) bool {
				day := d.In(time.UTC).Weekday()
				return day == time.Saturday || day == time.Sunday
			}
		}
		switch id {
		case "calendar":
			ui.Row(c).FillWidth().Wrap().Gap(8).Children(func() {
				for _, item := range []struct {
					label string
					mode  calendar.Mode
				}{{"Single", calendar.Single}, {"Multiple", calendar.Multiple}, {"Range", calendar.RangeMode}} {
					mygoui.Button(c, mygoui.ButtonOptions{Label: item.label, Variant: mygoui.ButtonOutline, OnClick: func() { s.Mode = item.mode }})
				}
			})
			ui.ScrollHorizontal(c).FillWidth().Height(theme.Scale(c, 360)).Children(func() { calendar.Calendar(c, &s.Calendar, options) })
			selection := s.Calendar.Date.String()
			switch s.Mode {
			case calendar.Multiple:
				selection = strconv.Itoa(len(s.Calendar.Dates)) + " selected dates"
			case calendar.RangeMode:
				selection = s.Calendar.Range.Start.String() + " – " + s.Calendar.Range.End.String()
			}
			ui.Text(c, "Gregorian calendar · "+selection).MinHeight(theme.Scale(c, 24))
		case "date-picker":
			mygoui.Toggle(c, &s.Range, mygoui.ToggleOptions{Label: "Date range"})
			calendar.DatePicker(c, &s.Picker, calendar.DatePickerOptions{Label: "Choose sample date", Calendar: options, Range: s.Range, Clearable: true,
				Presets: []calendar.Preset{{Label: "Sample week", Range: calendar.Range{Start: calendar.Date{Year: 2026, Month: time.October, Day: 19}, End: calendar.Date{Year: 2026, Month: time.October, Day: 25}}}}})
		case "chart":
			ui.Row(c).FillWidth().Wrap().Gap(8).Children(func() {
				for _, item := range []struct {
					label string
					kind  chart.ChartKind
				}{{"Bar", chart.Bar}, {"Line", chart.Line}, {"Area", chart.Area}, {"Pie", chart.Pie}, {"Radar", chart.Radar}, {"Radial", chart.Radial}} {
					mygoui.Button(c, mygoui.ButtonOptions{Label: item.label, Variant: mygoui.ButtonOutline, OnClick: func() { s.Kind = item.kind }})
				}
			})
			ui.Row(c).FillWidth().Wrap().Gap(8).Children(func() {
				for _, item := range []struct {
					label string
					stack chart.StackMode
				}{{"Unstacked", chart.Unstacked}, {"Stacked", chart.Stacked}, {"Normalized", chart.Normalized}, {"Centered", chart.Centered}} {
					mygoui.Button(c, mygoui.ButtonOptions{Label: item.label, Variant: mygoui.ButtonOutline, OnClick: func() { s.Stack = item.stack }})
				}
			})
			ui.Row(c).FillWidth().Wrap().Gap(8).Children(func() {
				mygoui.Toggle(c, &s.Empty, mygoui.ToggleOptions{Label: "Empty data"})
				mygoui.Toggle(c, &s.Missing, mygoui.ToggleOptions{Label: "Missing month"})
			})
			chartData := sampleChartData
			if s.Empty {
				chartData = chart.Data{}
			}
			if s.Missing && !s.Empty {
				chartData.Series = append([]chart.Series(nil), sampleChartData.Series...)
				for i := range chartData.Series {
					chartData.Series[i].Values = append([]chart.Value(nil), chartData.Series[i].Values...)
					chartData.Series[i].Values[2].Missing = true
				}
			}
			ui.Text(c, "Local monthly samples. Pie and radial show the first month's Desktop/Mobile values.")
			chart.Chart(c, &s.Chart, chart.Options{Kind: s.Kind, Stack: s.Stack, Data: chartData, Label: "Monthly sample usage", Height: 240})
		}
	})
}
