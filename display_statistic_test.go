package mygoui_test

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestStatisticShowsValueAffixesAndNamedTrend pins the metric: the caller's
// strings show as given around the value, and the direction of a change is
// named for assistive technology rather than carried by color alone.
func TestStatisticShowsValueAffixesAndNamedTrend(t *testing.T) {
	o := mygoui.StatisticOptions{Label: "Revenue", Value: "12,480", Prefix: "$", Suffix: "USD"}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Statistic(c, o)
	}, 360, 160)
	prefix, value, suffix := findRect(t, tt, "$"), findRect(t, tt, "12,480"), findRect(t, tt, "USD")
	if prefix.X >= value.X || value.X >= suffix.X {
		t.Fatalf("affixes are not around the value: %+v %+v %+v", prefix, value, suffix)
	}
	if _, ok := tt.AccessByRole(ui.RoleGroup, "Revenue"); !ok {
		t.Fatal("the label does not name the metric")
	}
	if _, ok := tt.Access(locale.M().Increase); ok {
		t.Fatal("a metric without a delta names a trend")
	}

	o.Delta, o.Trend = "12.5%", mygoui.TrendUp
	tt.Frame()
	if delta := findRect(t, tt, "12.5%"); delta.Y <= value.Y {
		t.Fatalf("the delta is not under the value: %+v", delta)
	}
	if _, ok := tt.Access(locale.M().Increase); !ok {
		t.Fatal("a rising delta does not name its direction")
	}
	o.Trend = mygoui.TrendDown
	tt.Frame()
	if _, ok := tt.Access(locale.M().Decrease); !ok {
		t.Fatal("a falling delta does not name its direction")
	}
}

// TestStatisticDeltaStatusReplacesTheTrendColor pins the override: a rise
// that is bad news paints in the danger ink, not the success one.
func TestStatisticDeltaStatusReplacesTheTrendColor(t *testing.T) {
	o := mygoui.StatisticOptions{Value: "248", Delta: "31", Trend: mygoui.TrendUp}
	var good, bad ui.Color
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		good = mygoui.StatusColors(c, mygoui.ColorSuccess).Foreground
		bad = mygoui.StatusColors(c, mygoui.ColorDanger).Foreground
		mygoui.Statistic(c, o)
	}, 240, 120)
	if inkCount(tt.Image(), good) == 0 || inkCount(tt.Image(), bad) != 0 {
		t.Fatal("a rising delta is not painted in the success ink")
	}
	o.DeltaStatus = mygoui.ColorDanger
	tt.Frame()
	tt.Advance(300 * time.Millisecond)
	if inkCount(tt.Image(), bad) == 0 || inkCount(tt.Image(), good) != 0 {
		t.Fatal("DeltaStatus did not replace the trend's ink")
	}
}

// TestStatisticLoadingHoldsTheValueHeight pins the placeholder: the label
// stays, the value and delta wait, and the metric keeps the height of a
// loaded one without a delta, so a row of metrics does not jump on arrival.
func TestStatisticLoadingHoldsTheValueHeight(t *testing.T) {
	o := mygoui.StatisticOptions{Label: "Revenue", Value: "12,480", Delta: "12.5%", Loading: true}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Statistic(c, o)
	}, 360, 160)
	group, ok := tt.AccessByRole(ui.RoleGroup, "Revenue")
	if !ok || !group.Busy {
		t.Fatalf("a loading metric is not busy: %+v", group)
	}
	if _, ok := tt.Find("12,480"); ok {
		t.Fatal("the value shows while loading")
	}
	if _, ok := tt.Find("12.5%"); ok {
		t.Fatal("the delta shows while loading")
	}
	findRect(t, tt, "Revenue")
	loading := group.Bounds.H

	o.Loading, o.Delta = false, ""
	tt.Frame()
	group, _ = tt.AccessByRole(ui.RoleGroup, "Revenue")
	if group.Bounds.H != loading {
		t.Fatalf("height jumped on arrival: %v while loading, %v loaded", loading, group.Bounds.H)
	}
}
