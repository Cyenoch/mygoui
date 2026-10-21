package gallery

import (
	"fmt"
	"strconv"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/conversation"
	"github.com/egoist/mygo/ui"
)

type SampleMessage struct {
	ID   conversation.MessageID
	Text string
}

type ConversationState struct {
	Text, Result                                        string
	End, Selected, Disabled, Removed, Uploading, Failed bool
	Progress                                            float64
	Messages                                            []SampleMessage
	Next                                                int
	Scroller                                            conversation.ScrollerState
	Questionnaire                                       conversation.QuestionnaireState
}

func NewConversationState() ConversationState {
	messages := make([]SampleMessage, 12)
	for i := range messages {
		messages[i] = SampleMessage{ID: conversation.MessageID("sample-" + strconv.Itoa(i)),
			Text: fmt.Sprintf("Local sample row %d. This is fixture content, not recorded agent activity.", i+1)}
	}
	return ConversationState{Text: "An editable local conversation sample.", Progress: 0.64, Messages: messages, Next: 12}
}

func PreviewConversation(c *ui.Context, s *ConversationState, id ComponentID) {
	ui.Column(c).FillWidth().Gap(12).Children(func() {
		if id != "questionnaire" && id != "message" {
			ui.TextArea(c, &s.Text).Label("Conversation sample text").Lines(2, 4).FillWidth()
		}
		ui.Row(c).Wrap().Gap(8).Children(func() {
			if id == "bubble" {
				mygoui.Toggle(c, &s.End, mygoui.ToggleOptions{Label: "Align end"})
			}
			if id == "message" || id == "bubble" || id == "attachment" || id == "message-scroller" || id == "questionnaire" {
				mygoui.Toggle(c, &s.Disabled, mygoui.ToggleOptions{Label: "Disable actions"})
			}
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Reset conversation fixture", Variant: mygoui.ButtonOutline,
				OnClick: func() { *s = NewConversationState() }})
		})
		align := conversation.Start
		if s.End {
			align = conversation.End
		}
		body := func() { ui.Text(c, s.Text).Selectable() }
		switch id {
		case "message":
			messageConversation(c, s)
		case "bubble":
			for _, variant := range []conversation.BubbleVariant{conversation.BubbleDefault, conversation.BubbleSecondary,
				conversation.BubbleMuted, conversation.BubbleTinted, conversation.BubbleOutline, conversation.BubbleGhost, conversation.BubbleDestructive} {
				conversation.Bubble(c, conversation.BubbleOptions{Variant: variant, Align: align,
					Reactions: []conversation.Reaction{{Key: "sample-reaction", Label: "Helpful", Count: 1, Selected: s.Selected,
						Disabled: s.Disabled, OnClick: func() { s.Selected = !s.Selected }}}}, body)
			}
		case "marker":
			for _, variant := range []conversation.MarkerVariant{conversation.MarkerDefault, conversation.MarkerSeparator, conversation.MarkerBorder} {
				conversation.Marker(c, conversation.MarkerOptions{Text: s.Text, Icon: lucide.Info(), Variant: variant})
			}
			conversation.Marker(c, conversation.MarkerOptions{Text: "Record sample marker action", Icon: lucide.Check(),
				OnClick: func() { s.Result = "Marker action activated"; c.AddToast(ui.Toast{Title: s.Result}) }})
		case "attachment":
			ui.Row(c).Wrap().Gap(8).Children(func() {
				mygoui.Toggle(c, &s.Uploading, mygoui.ToggleOptions{Label: "Uploading"})
				mygoui.Toggle(c, &s.Failed, mygoui.ToggleOptions{Label: "Upload error"})
			})
			if s.Removed {
				ui.Text(c, "Sample attachment removed. Reset the fixture to restore it.")
				break
			}
			state := conversation.AttachmentDone
			if s.Uploading {
				state = conversation.AttachmentUploading
			}
			if s.Failed {
				state = conversation.AttachmentError
			}
			conversation.Attachment(c, conversation.AttachmentOptions{Title: "sample.go", Description: s.Text, State: state,
				Icon: lucide.FileText(), Disabled: s.Disabled, Progress: &s.Progress,
				OnOpen:   func() { s.Result = "Opened sample metadata: sample.go"; c.AddToast(ui.Toast{Title: s.Result}) },
				OnRemove: func() { s.Removed = true }, OnRetry: func() { s.Failed = false; s.Uploading = true; s.Result = "Retry fixture selected" },
			})
		case "message-scroller":
			ui.Row(c).Wrap().Gap(8).Children(func() {
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Append sample message", Disabled: s.Disabled, OnClick: func() {
					s.Messages = append(s.Messages, SampleMessage{ID: conversation.MessageID("sample-" + strconv.Itoa(s.Next)), Text: s.Text})
					s.Next++
				}})
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Grow latest sample", Disabled: s.Disabled, OnClick: func() {
					if len(s.Messages) > 0 {
						s.Messages[len(s.Messages)-1].Text += "\n" + s.Text
					}
				}})
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Jump to first sample", Variant: mygoui.ButtonOutline,
					OnClick: func() {
						if len(s.Messages) > 0 {
							s.Scroller.JumpTo(s.Messages[0].ID, ui.Start)
						}
					}})
			})
			ui.Box(c).FillWidth().Height(360).Children(func() {
				conversation.Scroller(c, &s.Scroller, conversation.ScrollerOptions{Scope: "gallery-samples", Label: "Local sample messages",
					Count: len(s.Messages), Key: func(i int) conversation.MessageID { return s.Messages[i].ID },
					RowLabel:    func(i int) string { return "Sample " + string(s.Messages[i].ID) },
					Row:         func(row *ui.Context, i int) { ui.Text(row, s.Messages[i].Text).Selectable() },
					MoreHistory: true, OnLoadOlder: func() {
						s.Messages = append([]SampleMessage{{ID: conversation.MessageID("older-" + strconv.Itoa(s.Next)), Text: "Prepended local sample history"}}, s.Messages...)
						s.Next++
					},
				})
			})
		case "questionnaire":
			conversation.Questionnaire(c, &s.Questionnaire, []conversation.Question{
				{Name: "direction", Title: "Choose a sample direction", Description: "A required single choice or freeform answer.", AllowText: true,
					TextLabel: "Another sample direction", Choices: []conversation.Choice{{Value: "compact", Label: "Compact"}, {Value: "complete", Label: "Complete"}}},
				{Name: "signals", Title: "Choose sample details", Description: "Optional multiple choice.", Multiple: true, Optional: true,
					Choices: []conversation.Choice{{Value: "status", Label: "Status"}, {Value: "source", Label: "Source"}, {Value: "disabled", Label: "Unavailable", Disabled: true}}},
			}, conversation.QuestionnaireOptions{Label: "Local sample questionnaire", SubmitLabel: "Save sample answers", Busy: s.Disabled,
				LetterShortcuts: true, OnSubmit: func(answers map[conversation.QuestionID]conversation.Answer) error {
					s.Result = fmt.Sprintf("Saved %d sample answers", len(answers))
					c.AddToast(ui.Toast{Title: s.Result})
					return nil
				}})
		}
	})
}
