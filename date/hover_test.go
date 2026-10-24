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

// TestCalendarDayHoverPaint checks a plain and a today cell take their hover
// faces and revert; a selected cell takes the accent pressed face.
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
		wantDateHoverPixel(t, tt.Image(), x, y, colors.Default, "day hover fill")

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
		normal := colors.AccentSoft.Over(colors.Background)
		hover := colors.AccentSoftHover.Over(colors.Background)
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
		wantDateHoverPixel(t, tt.Image(), x, y, colors.Accent, "selected day fill")

		tt.Press(cell.X+cell.W/2, cell.Y+cell.H/2)
		tt.Frame()
		wantDateHoverPixel(t, tt.Image(), x, y, colors.AccentHover, "selected day pressed fill")
		tt.Release(cell.X+cell.W/2, cell.Y+cell.H/2)
		tt.Frame()
		wantDateHoverPixel(t, tt.Image(), x, y, colors.Accent, "selected day fill after release")
	})
}

// TestDateFieldHoverPaint checks the date input group's field surface follows
// hover and reverts.
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
	// The group sits one Space4 below the label and is 36 DIP high.
	y := int(label.Y+label.H) + theme.Space4 + 18
	x := int(label.X+label.W) - 16
	wantDateHoverPixel(t, tt.Image(), x, y, colors.FieldBackground, "date field normal fill")

	tt.Move(label.X+label.W/2, float32(y))
	tt.Frame()
	wantDateHoverPixel(t, tt.Image(), x, y, colors.FieldHover, "date field hover fill")

	tt.Move(2, 2)
	tt.Frame()
	wantDateHoverPixel(t, tt.Image(), x, y, colors.FieldBackground, "date field fill after leaving")
}
