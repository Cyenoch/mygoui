// Command shell demonstrates the layout family: a custom title bar, a
// resizable sidebar, an inspector and a status strip, built with
// layout.Scaffold.
package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/layout"
	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/navigation"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

// titleBarHeight is the native band the window is created with and the
// title bar covers.
const titleBarHeight = 44

// state is the window-owned state. An application that restores its layout
// sets the pane fields before the first build and saves them when it likes;
// the library stores nothing.
type state struct {
	mailboxes list.ListBoxState
	messages  ui.ScrollState

	sidebarWidth, inspectorWidth float32
	sidebarCollapsed             bool
	sidebarOpen                  bool
	inspectorCollapsed           bool
	inspectorOpen                bool
}

var mailboxes = []list.ListBoxItem{
	{Key: "inbox", Label: "Inbox"},
	{Key: "drafts", Label: "Drafts"},
	{Key: "sent", Label: "Sent"},
	{Key: "archive", Label: "Archive"},
}

func main() {
	mygo.App.SetName("MyGo UI Shell")
	mygo.App.OnWindowAllClosed(func() { mygo.App.Quit() })
	mygo.App.WhenReady(func() {
		s := &state{}
		s.mailboxes.SelectedKey = "inbox"
		o := mygo.WindowOptions{
			Title: "MyGo UI Shell", Width: 1040, Height: 680,
			TitleBarStyle:  mygo.TitleBarHiddenInset,
			TitleBarHeight: titleBarHeight,
			Content:        ui.View(func(c *ui.Context) { view(c, s) }),
		}
		if runtime.GOOS == "darwin" {
			// A standard 14 DIP button centered in the band.
			o.TrafficLightPosition = &mygo.Point{X: 14, Y: (titleBarHeight - 14) / 2}
		}
		mygo.NewWindow(o)
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}

func view(c *ui.Context, s *state) {
	theme.Apply(c)
	mailbox := selected(s)
	layout.Scaffold(c, layout.ScaffoldOptions{
		StartBar: layout.TitleBarOptions{Title: "Mail"},
		Start: &layout.PaneOptions{
			Label: "Mailboxes", Size: &s.sidebarWidth, Width: 220, Max: 360,
			Collapsed: &s.sidebarCollapsed, Open: &s.sidebarOpen,
			Body: func() {
				list.ListBox(c, list.ListBoxOptions{
					Label: "Mailboxes", Items: mailboxes, State: &s.mailboxes,
					Selection: list.SelectionSingle, Padding: theme.Space8,
					// Choosing a mailbox closes the compact sheet.
					OnAction: func(string) { s.sidebarOpen = false },
				}).Fill()
			},
		},
		TitleBar: layout.TitleBarOptions{
			Title: mailbox, MinHeight: titleBarHeight,
			Leading: func() {
				layout.PaneToggle(c, layout.PaneToggleOptions{
					Collapsed: &s.sidebarCollapsed, Open: &s.sidebarOpen,
				})
			},
			Trailing: func() {
				mygoui.Button(c, mygoui.ButtonOptions{
					Label: "New message", Icon: lucide.SquarePen(), IconOnly: true,
					Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall,
				})
				layout.PaneToggle(c, layout.PaneToggleOptions{
					Collapsed: &s.inspectorCollapsed, Open: &s.inspectorOpen, End: true,
				})
			},
		},
		Content: func() { messages(c, s, mailbox) },
		End: &layout.PaneOptions{
			Label: "Details", Size: &s.inspectorWidth, Width: 260, Max: 400,
			Collapsed: &s.inspectorCollapsed, Open: &s.inspectorOpen,
			Body: func() { details(c, mailbox) },
		},
		StatusBar: func() {
			navigation.StatusBar(c, navigation.StatusBarOptions{
				Message: "Up to date",
				Items:   []navigation.StatusBarItemOptions{{Text: fmt.Sprintf("%d messages", len(subjects))}},
			})
		},
	})
}

func selected(s *state) string {
	for _, m := range mailboxes {
		if m.Key == s.mailboxes.SelectedKey {
			return m.Label
		}
	}
	return ""
}

var subjects = []string{"Release notes for 0.3", "Splitter keyboard review", "Window chrome on Linux"}

func messages(c *ui.Context, s *state, mailbox string) {
	mygoui.ScrollArea(c, &s.messages, mygoui.ScrollAreaOptions{
		Axis: mygoui.ScrollVertical, Label: "Messages",
	}, func() {
		ui.Column(c).FillWidth().MinWidth(0).Padding(theme.Space24).Gap(theme.Space16).Children(func() {
			mygoui.Heading(c, 4, mailbox)
			for _, subject := range subjects {
				mygoui.Card(c.Key(subject), mygoui.CardOptions{}, func() {
					mygoui.CardHeader(c, func() {
						mygoui.CardTitle(c, subject)
						mygoui.CardDescription(c, "Drag the dividers, or focus one and use the arrow keys.")
					})
				})
			}
		})
	}).Fill()
}

func details(c *ui.Context, mailbox string) {
	ui.Column(c).Fill().MinWidth(0).Padding(theme.Space16).Gap(theme.Space12).Children(func() {
		mygoui.Typography(c, "Details", mygoui.TypographyOptions{Type: theme.TypeH6})
		mygoui.Descriptions(c, mygoui.DescriptionsOptions{
			Columns: 1,
			Items: []mygoui.DescriptionsItem{
				{Label: "Mailbox", Value: mailbox},
				{Label: "Messages", Value: fmt.Sprint(len(subjects))},
			},
		})
	})
}
