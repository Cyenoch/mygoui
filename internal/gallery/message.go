package gallery

import (
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/conversation"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func messageConversation(c *ui.Context, s *ConversationState) {
	ui.Column(c).FillWidth().MaxWidth(440).MinWidth(0).Margin(0, ui.Auto).Gap(24).Children(func() {
		avatar := func(name string) func() {
			return func() { mygoui.Avatar(c, mygoui.AvatarOptions{Name: name, Size: mygoui.AvatarSM}) }
		}
		turn := func(key, text string, align conversation.Align, footer string, name string) {
			conversation.Message(c.Key(key), conversation.MessageOptions{Align: align, Label: "Sample conversation message", Footer: footer, Avatar: avatar(name)}, func() {
				variant := conversation.BubbleDefault
				if align == conversation.Start {
					variant = conversation.BubbleSecondary
				}
				conversation.Bubble(c, conversation.BubbleOptions{Align: align, Variant: variant}, func() { ui.Text(c, text).Selectable() })
			})
		}
		turn("outgoing-one", "Deploying to prod real quick.", conversation.End, "", "Ava")
		turn("incoming-one", "It's 4:55 PM. On a Friday.", conversation.Start, "", "Oliver")
		turn("outgoing-two", "It's a one-line change.", conversation.End, "Delivered", "Ava")
		conversation.Message(c.Key("incoming-group"), conversation.MessageOptions{Align: conversation.Start, Label: "Grouped sample replies", Avatar: avatar("Oliver")}, func() {
			conversation.MessageGroup(c, func() {
				conversation.Bubble(c.Key("first-reply"), conversation.BubbleOptions{Align: conversation.Start, Variant: conversation.BubbleSecondary}, func() { ui.Text(c, "It's always a one-line change 😭.").Selectable() })
				count := 1
				if s.Selected {
					count = 2
				}
				conversation.Bubble(c.Key("last-reply"), conversation.BubbleOptions{Align: conversation.Start, Variant: conversation.BubbleSecondary,
					Reactions: []conversation.Reaction{{Key: "thumbs-up", Label: "👍", AccessibilityLabel: "Thumbs up", Count: count, Selected: s.Selected, Disabled: s.Disabled, OnClick: func() { s.Selected = !s.Selected }}}}, func() { ui.Text(c, "Alright, let me take a look.").Selectable() })
			})
		})
		conversation.Message(c.Key("typing"), conversation.MessageOptions{Align: conversation.Start, Label: "Oliver is typing", Avatar: avatar("Oliver")}, func() {
			conversation.Bubble(c, conversation.BubbleOptions{Align: conversation.Start, Variant: conversation.BubbleSecondary, Label: "Typing indicator"}, func() {
				ui.Row(c).Gap(4).Padding(5, 0).Children(func() {
					for range 3 {
						ui.Column(c).Size(6, 6).Radius(3).Background(theme.Colors(c).MutedForeground)
					}
				})
			})
		})
		ui.Text(c, "Oliver is typing...").Role(ui.RoleStatus).FontSize(theme.Scale(c, 13)).TextColor(theme.Colors(c).MutedForeground)
	})
}
