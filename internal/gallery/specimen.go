package gallery

import (
	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/calendar"
	"github.com/Cyenoch/mygoui/chart"
	"github.com/Cyenoch/mygoui/codeblock"
	"github.com/Cyenoch/mygoui/conversation"
	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/overlay"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/Cyenoch/mygoui/trace"
	"github.com/egoist/mygo/ui"
	"time"
)

// defaultSpecimen is the compact first example. The existing complete, editable
// fixture builders are rendered independently in the Examples section below it.
func defaultSpecimen(c *ui.Context, s *State, id ComponentID) {
	b := &s.DemoBasic
	f := &s.DemoForms
	o := &s.DemoOverlay
	action := func() { b.Count++; b.Result = "Sample action completed"; c.AddToast(ui.Toast{Title: b.Result}) }
	align := ui.Start
	switch id {
	case "button", "badge", "avatar", "spinner", "kbd", "toggle", "toggle-group", "button-group",
		"dialog", "sheet", "drawer", "alert-dialog", "popover", "tooltip", "hover-card", "dropdown-menu", "context-menu", "menubar", "toast":
		align = ui.Center
	}
	measure := float32(480)
	if id == "sidebar" || id == "message" {
		measure = 800
	}
	ui.Column(c).FillWidth().MaxWidth(measure).MinWidth(0).Margin(0, ui.Auto).Gap(12).
		Justify(ui.Center).AlignItems(align).Children(func() {
		switch id {
		case "accordion", "breadcrumb", "collapsible", "combobox", "command", "navigation-menu", "pagination", "select", "sidebar", "tabs":
			navigationSample(c, &s.DemoNavigation, id, false)
		case "alert", "progress", "direction", "resizable", "scroll-area", "carousel", "table", "data-table":
			dataLayoutSample(c, &s.DemoDataLayout, id, false)
		case "calendar":
			ui.ScrollHorizontal(c).FillWidth().Height(theme.Scale(c, 340)).Children(func() { calendar.Calendar(c, &s.DemoDates.Calendar, calendar.Options{Label: "Sample date"}) })
		case "date-picker":
			calendar.DatePicker(c, &s.DemoDates.Picker, calendar.DatePickerOptions{Label: "Sample date", Clearable: true,
				Calendar: calendar.Options{Min: calendar.Date{Year: 2026, Month: time.January, Day: 1}}})
		case "chart":
			chart.Chart(c, &s.DemoDates.Chart, chart.Options{Kind: chart.Bar, Data: sampleChartData, Label: "Monthly sample usage", Height: 220}).FillWidth()
		case "button":
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Default", OnClick: action})
		case "badge":
			mygoui.Badge(c, mygoui.BadgeOptions{Label: "New release"})
		case "card":
			mygoui.Card(c, mygoui.CardOptions{Label: "Project card"}, func() {
				mygoui.CardHeader(c, mygoui.CardHeaderOptions{Title: "Create a project", Description: "Start with a local sample."})
				mygoui.CardContent(c, func() { ui.Text(c, "Your project settings can be changed later.") })
				mygoui.CardFooter(c, func() { mygoui.Button(c, mygoui.ButtonOptions{Label: "Create project", OnClick: action}) })
			})
		case "separator":
			ui.Text(c, "Components").Bold()
			mygoui.Separator(c, mygoui.SeparatorOptions{})
			ui.Text(c, "Native controls and composable views.")
		case "skeleton":
			mygoui.Skeleton(c, mygoui.SkeletonOptions{Width: 48, Height: 48, Round: true})
			mygoui.Skeleton(c, mygoui.SkeletonOptions{Width: 180, Height: 16})
			mygoui.Skeleton(c, mygoui.SkeletonOptions{Width: 140, Height: 16})
		case "spinner":
			mygoui.Spinner(c, mygoui.SpinnerOptions{Label: "Loading sample"})
		case "kbd":
			mygoui.KbdGroup(c, func() { mygoui.Kbd(c, "⌘"); mygoui.Kbd(c, "K") })
		case "typography":
			mygoui.Typography(c, "The native interface", mygoui.Heading3)
			mygoui.Typography(c, "Clear type makes content easier to read.", mygoui.Paragraph)
		case "aspect-ratio":
			mygoui.AspectRatio(c, 16.0/9, func() {
				ui.Box(c).Fill().Center().Background(theme.Colors(c).Muted).Children(func() { ui.Text(c, "16:9") })
			}).FillWidth()
		case "empty":
			mygoui.Empty(c, mygoui.EmptyOptions{Title: "No projects yet", Description: "Create your first local sample.", Icon: lucide.Folder(),
				Actions: func() { mygoui.Button(c, mygoui.ButtonOptions{Label: "Create project", OnClick: action}) }})
		case "item":
			mygoui.Item(c, mygoui.ItemOptions{Label: "Project notes", Description: "A local sample document", Icon: lucide.FileText(), Variant: mygoui.ItemOutline, OnClick: action})
		case "avatar":
			mygoui.Avatar(c, mygoui.AvatarOptions{Name: "Ada Lovelace"})
		case "button-group":
			mygoui.ButtonGroup(c, mygoui.ButtonGroupOptions{Label: "Navigate samples"},
				mygoui.ButtonOptions{Label: "Back", Variant: mygoui.ButtonOutline, OnClick: action},
				mygoui.ButtonOptions{Label: "Next", Variant: mygoui.ButtonOutline, OnClick: action})
		case "toggle":
			mygoui.Toggle(c, &b.Pressed, mygoui.ToggleOptions{Label: "Bold", Variant: mygoui.ToggleOutline})
		case "toggle-group":
			mygoui.ToggleGroup(c, &b.Group, mygoui.ToggleGroupOptions{Label: "Alignment", Variant: mygoui.ToggleOutline},
				[]mygoui.ToggleGroupItem{{Value: "left", Label: "Left"}, {Value: "center", Label: "Center"}, {Value: "right", Label: "Right"}})
		case "code-block":
			codeblock.View(c, s.Code.Source, &s.DemoCode, codeblock.Options{Language: "go", Title: "sample.go", LineNumbers: true, MaxHeight: 200})
		case "trace":
			trace.Collapsible(c, &s.DemoTrace, trace.Options{Title: "Reviewed the sample source", Icon: lucide.Brain()}, func() { ui.Text(c, "The function returns a greeting from local sample data.") })
		case "input", "field":
			form.Field(c, form.FieldOptions{Label: "Sample input", Description: "Choose a display name."}, func() {
				form.Input(c, &f.Text, form.InputOptions{ControlOptions: form.ControlOptions{Label: "Display name"}})
			})
		case "textarea":
			form.Textarea(c, &f.Notes, form.TextareaOptions{InputOptions: form.InputOptions{ControlOptions: form.ControlOptions{Label: "Notes"}}, MinLines: 3, MaxLines: 5})
		case "label":
			form.Input(c, &f.Text, form.InputOptions{ControlOptions: form.ControlOptions{Label: "Display name", Handle: &f.InputHandle}})
			form.Label(c, "Focus the display name", &f.InputHandle)
		case "input-group":
			form.InputGroup(c, &f.Text, form.InputGroupOptions{InputOptions: form.InputOptions{ControlOptions: form.ControlOptions{Label: "Website"}},
				BlockEnd: func() {
					ui.Row(c).FillWidth().Wrap().Gap(8).Children(func() {
						form.InputGroupText(c, "Website address")
						form.InputGroupButton(c, mygoui.ButtonOptions{Label: "Use", Variant: mygoui.ButtonGhost, Size: mygoui.ButtonXS, OnClick: action})
					})
				}})
		case "input-otp":
			ui.ScrollHorizontal(c).FillWidth().Height(112).Children(func() {
				form.InputOTP(c, &f.OTP, form.OTPOptions{ControlOptions: form.ControlOptions{Label: "Verification code"}, Length: 6, SeparatorAfter: 3})
			})
		case "checkbox":
			form.Checkbox(c, &f.Checked, form.CheckboxOptions{ControlOptions: form.ControlOptions{Label: "Accept sample notifications"}})
		case "radio-group":
			form.RadioGroup(c, &f.Choice, form.RadioGroupOptions{ControlOptions: form.ControlOptions{Label: "Response detail"}}, []form.RadioItem{{Value: "fast", Label: "Fast"}, {Value: "balanced", Label: "Balanced"}})
		case "switch":
			form.Switch(c, &f.Enabled, form.SwitchOptions{ControlOptions: form.ControlOptions{Label: "Airplane mode"}})
		case "slider":
			form.Slider(c, &f.Value, form.SliderOptions{ControlOptions: form.ControlOptions{Label: "Volume"}, Min: 0, Max: 100, Step: 5, Unit: "%"})
		case "native-select":
			form.NativeSelect(c, &f.Choice, form.NativeSelectOptions{ControlOptions: form.ControlOptions{Label: "Response detail"}}, []form.SelectOption{{Value: "fast", Label: "Fast"}, {Value: "balanced", Label: "Balanced"}})
		case "dialog", "sheet", "drawer", "alert-dialog":
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Open " + entryFor(id).Name, Variant: mygoui.ButtonOutline, OnClick: func() { o.Open = true }})
			options := overlay.DialogOptions{Title: "Edit profile", Description: "Make changes to this local sample.", Footer: func() {
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Save changes", OnClick: func() { o.Result = "Profile saved"; o.Open = false; c.AddToast(ui.Toast{Title: o.Result}) }})
			}}
			body := func() { form.Input(c, &o.Text, form.InputOptions{ControlOptions: form.ControlOptions{Label: "Name"}}) }
			switch id {
			case "dialog":
				overlay.Dialog(c, &o.Open, options, body)
			case "sheet":
				overlay.Sheet(c, &o.Open, overlay.SheetOptions{DialogOptions: options, Side: overlay.SideRight}, body)
			case "drawer":
				overlay.Drawer(c, &o.Open, overlay.DrawerOptions{DialogOptions: options, Side: overlay.SideBottom}, body)
			case "alert-dialog":
				overlay.AlertDialog(c, &o.Open, overlay.AlertDialogOptions{Title: "Apply this sample change?", Description: "This changes only the local gallery fixture.", ActionLabel: "Continue", OnAction: func() { o.Result = "Change applied"; c.AddToast(ui.Toast{Title: o.Result}) }}, nil)
			}
		case "popover":
			anchor := mygoui.Button(c, mygoui.ButtonOptions{Label: "Open Popover", Variant: mygoui.ButtonOutline, OnClick: func() { o.Open = !o.Open }})
			overlay.Popover(c, anchor, &o.Open, overlay.PopoverOptions{Title: "Dimensions", Description: "Set the sample display name."}, func() { form.Input(c, &o.Text, form.InputOptions{ControlOptions: form.ControlOptions{Label: "Name"}}) })
		case "tooltip":
			anchor := mygoui.Button(c, mygoui.ButtonOptions{Label: "Hover or focus", Variant: mygoui.ButtonOutline, OnClick: action})
			overlay.Tooltip(c, anchor, overlay.TooltipOptions{Text: "Add to your sample library"})
		case "hover-card":
			anchor := mygoui.Button(c, mygoui.ButtonOptions{Label: "@native", Variant: mygoui.ButtonLink, OnClick: action})
			overlay.HoverCard(c, anchor, &o.Hover, overlay.HoverCardOptions{PopoverOptions: overlay.PopoverOptions{Title: "Native components", Description: "Go views built with MyGo."}}, nil)
		case "dropdown-menu":
			overlay.DropdownMenu(c, &o.Menu, overlay.MenuOptions{Label: "Open menu"}, sampleMenu(c, o))
		case "context-menu":
			anchor := mygoui.Button(c, mygoui.ButtonOptions{Label: "Right click or Shift+F10", Variant: mygoui.ButtonOutline, OnClick: action})
			overlay.ContextMenu(anchor, sampleMenu(c, o))
		case "menubar":
			overlay.Menubar(c, &o.Bar, []overlay.BarMenu{{ID: "file", Label: "File", Items: sampleMenu(c, o)}, {ID: "edit", Label: "Edit", Items: sampleMenu(c, o)}})
		case "toast":
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Show notification", Variant: mygoui.ButtonOutline, OnClick: func() {
				overlay.Toast(c, &o.Toasts, overlay.ToastOptions{Title: "Event created", Description: "The local sample has been saved."})
			}})
		case "message":
			messageConversation(c, &s.DemoConversation)
		case "bubble":
			conversation.Bubble(c, conversation.BubbleOptions{Variant: conversation.BubbleMuted}, func() { ui.Text(c, "This is a native message bubble.") })
		case "marker":
			conversation.Marker(c, conversation.MarkerOptions{Text: "Today", Variant: conversation.MarkerSeparator})
		case "attachment":
			conversation.Attachment(c, conversation.AttachmentOptions{Title: "sample.go", Description: "Go source · local sample", Icon: lucide.FileText(), OnOpen: action})
		case "message-scroller":
			ui.Box(c).FillWidth().Height(180).Children(func() {
				conversation.Scroller(c, &s.DemoConversation.Scroller, conversation.ScrollerOptions{Scope: "default-sample", Count: 3, Key: func(i int) conversation.MessageID { return s.DemoConversation.Messages[i].ID }, Row: func(row *ui.Context, i int) { ui.Text(row, s.DemoConversation.Messages[i].Text) }})
			})
		case "questionnaire":
			conversation.Questionnaire(c, &s.DemoConversation.Questionnaire, []conversation.Question{{Name: "detail", Title: "How much detail?", Choices: []conversation.Choice{{Value: "short", Label: "Short"}, {Value: "full", Label: "Full"}}}}, conversation.QuestionnaireOptions{Label: "Sample preference", SubmitLabel: "Save preference", OnSubmit: func(_ map[conversation.QuestionID]conversation.Answer) error {
				b.Result = "Preference saved"
				c.AddToast(ui.Toast{Title: b.Result})
				return nil
			}})
		}
	})
}
