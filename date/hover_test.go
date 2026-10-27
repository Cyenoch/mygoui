package date_test

import (
	"image"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/date"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func dateHoverPixel(img *image.RGBA, x, y int) (uint8, uint8, uint8) {
	c := img.RGBAAt(x, y)
	return c.R, c.G, c.B
}

func dateHoverNear(got, want uint8) bool {
	d := int(got) - int(want)
	return d >= -3 && d <= 3
}

func wantDateHoverPixel(t *testing.T, img *image.RGBA, x, y int, c ui.Color, what string) {
	t.Helper()
	r, g, b := dateHoverPixel(img, x, y)
	if !dateHoverNear(r, c.R) || !dateHoverNear(g, c.G) || !dateHoverNear(b, c.B) {
		t.Fatalf("%s at (%d, %d) = %d,%d,%d, want %d,%d,%d", what, x, y, r, g, b, c.R, c.G, c.B)
	}
}

// TestCalendarDayHoverPaint checks plain and today cells take their hover
// faces and revert, while a selected cell keeps its primary face.
func TestCalendarDayHoverPaint(t *testing.T) {
	today := time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC)

	t.Run("plain day", func(t *testing.T) {
		var colors theme.Palette
		var state date.CalendarState
		tt := ui.NewTester(func(c *ui.Context) {
			theme.ApplyMode(c, false)
			colors = theme.Colors(c)
			ui.Column(c).FillWidth().FillHeight().Padding(24).
				Background(colors.Background).Children(func() {
				date.Calendar(c, date.CalendarOptions{Label: "Cal", Today: today}, &state)
			})
		}, 460, 520)

		cell, ok := tt.Find("Tuesday, January 6, 2026")
		if !ok {
			t.Fatal("day cell missing")
		}
		x, y := int(cell.X)+6, int(cell.Y)+6
		wantDateHoverPixel(t, tt.Image(), x, y, colors.Background, "day normal fill")

		tt.Move(cell.X+cell.W/2, cell.Y+cell.H/2)
		tt.Frame()
		wantDateHoverPixel(t, tt.Image(), x, y, colors.Accent.Over(colors.Background), "day hover fill")

		tt.Move(2, 2)
		tt.Frame()
		wantDateHoverPixel(t, tt.Image(), x, y, colors.Background, "day fill after leaving")
	})

	t.Run("today", func(t *testing.T) {
		var colors theme.Palette
		var state date.CalendarState
		tt := ui.NewTester(func(c *ui.Context) {
			theme.ApplyMode(c, false)
			colors = theme.Colors(c)
			ui.Column(c).FillWidth().FillHeight().Padding(24).
				Background(colors.Background).Children(func() {
				date.Calendar(c, date.CalendarOptions{Label: "Cal", Today: today}, &state)
			})
		}, 460, 520)

		cell, ok := tt.Find("Monday, January 5, 2026")
		if !ok {
			t.Fatal("today cell missing")
		}
		normal := colors.Background
		hover := colors.Accent.Over(colors.Background)
		x, y := int(cell.X)+6, int(cell.Y)+6
		wantDateHoverPixel(t, tt.Image(), x, y, normal, "today normal fill")

		tt.Move(cell.X+cell.W/2, cell.Y+cell.H/2)
		tt.Frame()
		wantDateHoverPixel(t, tt.Image(), x, y, hover, "today hover fill")

		tt.Move(2, 2)
		tt.Frame()
		wantDateHoverPixel(t, tt.Image(), x, y, normal, "today fill after leaving")
	})

	t.Run("selected pressed", func(t *testing.T) {
		var colors theme.Palette
		state := date.CalendarState{Value: today}
		tt := ui.NewTester(func(c *ui.Context) {
			theme.ApplyMode(c, false)
			colors = theme.Colors(c)
			ui.Column(c).FillWidth().FillHeight().Padding(24).
				Background(colors.Background).Children(func() {
				date.Calendar(c, date.CalendarOptions{Label: "Cal", Today: today}, &state)
			})
		}, 460, 520)

		cell, ok := tt.Find("Monday, January 5, 2026")
		if !ok {
			t.Fatal("selected cell missing")
		}
		x, y := int(cell.X)+6, int(cell.Y)+6
		wantDateHoverPixel(t, tt.Image(), x, y, colors.Primary, "selected day fill")

		tt.Press(cell.X+cell.W/2, cell.Y+cell.H/2)
		tt.Frame()
		wantDateHoverPixel(t, tt.Image(), x, y, colors.Primary, "selected day pressed fill")
		tt.Release(cell.X+cell.W/2, cell.Y+cell.H/2)
		tt.Frame()
		wantDateHoverPixel(t, tt.Image(), x, y, colors.Primary, "selected day fill after release")
	})
}

// TestDateFieldHoverPaint checks the date input surface stays quiet on hover.
func TestDateFieldHoverPaint(t *testing.T) {
	var colors theme.Palette
	value := time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC)
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			date.DateField(c, date.DateFieldOptions{
				Label: "When", Value: &value, FullWidth: true,
			})
		})
	}, 460, 180)

	label, ok := tt.Find("When")
	if !ok {
		t.Fatal("date field label missing")
	}
	// The group sits below the label.
	y := int(label.Y+label.H) + theme.Space4 + 18
	x := int(label.X+label.W) - 16
	wantDateHoverPixel(t, tt.Image(), x, y, colors.Background, "date field normal fill")

	tt.Move(label.X+label.W/2, float32(y))
	tt.Frame()
	wantDateHoverPixel(t, tt.Image(), x, y, colors.Background, "date field hover fill")

	tt.Move(2, 2)
	tt.Frame()
	wantDateHoverPixel(t, tt.Image(), x, y, colors.Background, "date field fill after leaving")
}

func TestCalendarUsesPrimaryAndAccentTokens(t *testing.T) {
	previous := theme.PreferencesSnapshot()
	defer theme.Configure(previous)
	primary := ui.RGB(21, 92, 170)
	accent := ui.RGBA(213, 55, 32, 0.5)
	theme.Configure(theme.Preferences{Palette: theme.PaletteOverrides{
		Light: &theme.PalettePatch{Primary: &primary, Accent: &accent},
	}})
	selected := time.Date(2026, time.March, 18, 0, 0, 0, 0, time.UTC)
	st := &date.CalendarState{Value: selected}
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).Padding(20).Background(colors.Background).Children(func() {
			date.Calendar(c, date.CalendarOptions{Today: march, DayLabel: label}, st)
		})
	}, 360, 480)
	day, _ := tt.Find("2026-03-18")
	wantDateHoverPixel(t, tt.Image(), int(day.X+6), int(day.Y+6), primary, "selected day")
	hover, _ := tt.Find("2026-03-19")
	tt.Move(hover.X+hover.W/2, hover.Y+hover.H/2)
	wantDateHoverPixel(t, tt.Image(), int(hover.X+6), int(hover.Y+6),
		accent.Over(colors.Background), "hovered day")
}
