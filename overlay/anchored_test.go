package overlay

import (
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	myui "github.com/egoist/mygo/ui"
	"testing"
	"time"
)

func TestTooltipDelayKeyboardAndEscapeSuppression(t *testing.T) {
	tt := myui.NewTester(func(c *myui.Context) {
		theme.Apply(c)
		myui.Column(c).Padding(30).Children(func() {
			b := mygoui.Button(c, mygoui.ButtonOptions{Label: "Share", OnClick: func() {}})
			Tooltip(c, b, TooltipOptions{Text: "Share with others"})
		})
	}, 400, 300)
	r, _ := tt.Find("Share")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Advance(500 * time.Millisecond)
	if tt.HasText("Share with others") {
		t.Fatal("tooltip opened before native delay")
	}
	tt.Advance(200 * time.Millisecond)
	if !tt.HasText("Share with others") {
		t.Fatal("tooltip did not open after hover delay")
	}
	tt.Key(0, myui.KeyEscape)
	if tt.HasText("Share with others") {
		t.Fatal("Escape did not hide tooltip")
	}
	tt.Advance(time.Second)
	if tt.HasText("Share with others") {
		t.Fatal("tooltip reopened while trigger remained hovered")
	}
	tt.Move(0, 0)
	tt.Key(0, myui.KeyTab)
	if !tt.HasText("Share with others") {
		t.Fatal("keyboard focus should show tooltip immediately")
	}
}

func TestHoverCardGraceCrossingAndFocusedContent(t *testing.T) {
	var s HoverCardState
	tt := myui.NewTester(func(c *myui.Context) {
		theme.Apply(c)
		myui.Column(c).Padding(30).Children(func() {
			b := mygoui.Button(c, mygoui.ButtonOptions{Label: "Author", OnClick: func() {}})
			HoverCard(c, b, &s, HoverCardOptions{PopoverOptions: PopoverOptions{Title: "Author preview"}, OpenDelay: 100 * time.Millisecond, CloseDelay: 300 * time.Millisecond}, func() { myui.Button(c, "View author") })
		})
	}, 500, 400)
	r, _ := tt.Find("Author")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Advance(110 * time.Millisecond)
	if !tt.HasText("Author preview") {
		t.Fatal("hover card did not open")
	}
	tt.Move(0, 0)
	tt.Advance(100 * time.Millisecond)
	if !s.Open {
		t.Fatal("card closed before pointer could cross gap")
	}
	p, _ := tt.Find("View author")
	tt.Move(p.X+p.W/2, p.Y+p.H/2)
	tt.Advance(400 * time.Millisecond)
	if !s.Open {
		t.Fatal("card did not retain hover while crossing to panel")
	}
	if err := tt.Click("View author"); err != nil {
		t.Fatal(err)
	}
	tt.Move(0, 0)
	tt.Advance(time.Second)
	if !s.Open {
		t.Fatal("focused preview content must remain open")
	}
	tt.Key(0, myui.KeyEscape)
	if s.Open {
		t.Fatal("Escape did not close hover card")
	}
}

func TestPopoverCollidesWithinWindow(t *testing.T) {
	open := true
	var panel myui.Handle
	tt := myui.NewTester(func(c *myui.Context) {
		theme.Apply(c)
		b := myui.Button(c, "Anchor").Absolute().Right(0).Bottom(0)
		Popover(c, b, &open, PopoverOptions{Title: "Collision preview", Placement: Placement{Side: SideRight}}, func() { myui.Button(c, "Child").Bind(&panel) })
	}, 320, 240)
	tt.Frame()
	r, ok := tt.Find("Collision preview")
	if !ok {
		t.Fatal("popover absent")
	}
	if r.X < 0 || r.Y < 0 || r.X+r.W > 320 || r.Y+r.H > 240 {
		t.Fatalf("popover content escaped viewport: %+v", r)
	}
}
