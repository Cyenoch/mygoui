package overlay

import (
	"slices"
	"testing"
	"time"

	myui "github.com/egoist/mygo/ui"
)

func TestToastAnnouncementOverrideKeepsExactVisualPayloadAndNativeLifecycle(t *testing.T) {
	f := newNotificationFixture(t, 640, 700)
	f.tt.SetPreferences(myui.Preferences{TextScale: 1, ReduceMotion: true})
	announcement := "Application notice available"
	closed, actions := 0, 0
	f.add(ToastOptions{
		ID: "private", Title: "Private title", Description: "Exact private command and file text",
		Announcement: &announcement, Timeout: -1, Action: "Inspect notice",
		OnAction: func() { actions++ }, OnClose: func() { closed++ },
	})
	if got := f.tt.Announcements(); !slices.Equal(got, []string{announcement}) {
		t.Fatalf("automatic speech exposed source text: %q", got)
	}
	announcement = "mutated by caller"
	f.tt.Frame()
	if !f.tt.HasText("Private title") || !f.tt.HasText("Exact private command and file text") {
		t.Fatalf("exact visual payload lost: %q", f.tt.Texts())
	}
	if got := f.tt.Announcements(); len(got) != 0 {
		t.Fatalf("rebuild repeated speech: %q", got)
	}
	update := "Application notice updated"
	f.add(ToastOptions{
		ID: "private", Title: "Updated private title", Description: "Updated exact private detail",
		Announcement: &update, Timeout: -1, Action: "Inspect notice",
		OnAction: func() { actions++ }, OnClose: func() { closed++ },
	})
	if got := f.tt.Announcements(); !slices.Equal(got, []string{update}) {
		t.Fatalf("update speech = %q", got)
	}
	if !f.tt.HasText("Updated private title") || !f.tt.HasText("Updated exact private detail") || f.tt.HasText("Private title") {
		t.Fatalf("native same-ID update lost current visual payload: %q", f.tt.Texts())
	}
	notificationFocus(t, f.tt, "Inspect notice")
	f.tt.Key(0, myui.KeyEnter)
	if actions != 1 || closed != 1 {
		t.Fatalf("native action/close = %d/%d", actions, closed)
	}
	expiry := "Notice available"
	f.add(ToastOptions{ID: "expires", Title: "Private expires", Description: "Private expiring body", Announcement: &expiry, Timeout: time.Second, OnClose: func() { closed++ }})
	f.tt.Announcements()
	notificationFocus(t, f.tt, "Page action")
	f.tt.Advance(2 * time.Second)
	f.tt.Frame()
	if closed != 2 || f.tt.HasText("Private expiring body") {
		t.Fatal("safe speech override broke native expiry")
	}
}
