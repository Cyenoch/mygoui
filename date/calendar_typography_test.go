package date_test

import (
	"image"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/date"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestCalendarTypography(t *testing.T) {
	previous := theme.PreferencesSnapshot()
	t.Cleanup(func() { theme.Configure(previous) })
	muted := ui.RGBA(40, 90, 140, 0.5)
	weekdayInk := ui.RGBA(40, 90, 140, 0.36)
	patch := &theme.PalettePatch{MutedForeground: &muted}
	prefs := previous
	prefs.Palette = theme.PaletteOverrides{Light: patch, Dark: patch}
	theme.Configure(prefs)

	for _, dark := range []bool{false, true} {
		name := "light"
		if dark {
			name = "dark"
		}
		t.Run(name, func(t *testing.T) {
			state := &date.CalendarState{}
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, dark)
				colors := theme.Colors(c)
				ui.Column(c).Fill().Padding(20).Background(colors.Background).Children(func() {
					date.Calendar(c, date.CalendarOptions{
						Today: time.Date(2026, time.October, 15, 12, 0, 0, 0, time.UTC),
					}, state)
					ui.Row(c).Children(func() {
						ui.Box(c).Width(144)
						ui.Box(c).Size(36, 36).Center().Children(func() {
							style.Text(ui.Text(c, "1"), c, theme.SM()).
								TextColor(colors.Foreground).SingleLine().FontFeatures("tnum").
								Label("regular day reference")
						})
					})
					ui.Box(c).Width(36).Height(24).Padding(0, 0, theme.Space8, 0).Center().Children(func() {
						style.Text(ui.Text(c, "Su"), c, theme.XS()).
							FontWeight(theme.WeightMedium).TextColor(weekdayInk).SingleLine().
							Label("weekday reference")
					})
				})
			}, 360, 480)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true})
			// The first frame's entrance transition can outlive the preference change.
			tt.Advance(time.Second)
			img := tt.Image()
			t.Run("day", func(t *testing.T) {
				assertSameTextPixels(t, tt, img, "1", "regular day reference")
			})
			t.Run("weekday", func(t *testing.T) {
				assertSameTextPixels(t, tt, img, "Su", "weekday reference")
			})
		})
	}
}

func assertSameTextPixels(t *testing.T, tt *ui.Tester, img *image.RGBA, actual, reference string) {
	t.Helper()
	a, ok := tt.Find(actual)
	if !ok {
		t.Fatalf("text %q missing", actual)
	}
	b, ok := tt.Find(reference)
	if !ok {
		t.Fatalf("text %q missing", reference)
	}
	if a.W != b.W || a.H != b.H {
		t.Fatalf("%q bounds %+v, %q bounds %+v", actual, a, reference, b)
	}
	for y := range int(a.H) {
		for x := range int(a.W) {
			got := img.RGBAAt(int(a.X)+x, int(a.Y)+y)
			want := img.RGBAAt(int(b.X)+x, int(b.Y)+y)
			if got != want {
				t.Fatalf("%q differs from %q at offset %d,%d: %v != %v", actual, reference, x, y, got, want)
			}
		}
	}
}
