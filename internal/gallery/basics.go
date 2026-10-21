package gallery

import (
	"strconv"
	"strings"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// BasicState is caller-owned fixture data for the independent basic previews.
type BasicState struct {
	Label, Result                                        string
	Disabled, Loading, Invalid, Pressed, Vertical, Small bool
	Count                                                int
	Group                                                mygoui.ToggleGroupState
}

func NewBasicState() BasicState {
	return BasicState{Label: "Native sample", Group: mygoui.ToggleGroupState{Values: []string{"left"}}}
}

// PreviewBasic can be embedded in another native window with persistent state.
func PreviewBasic(c *ui.Context, s *BasicState, id ComponentID) {
	ui.Column(c).FillWidth().MinWidth(0).Justify(ui.Start).AlignItems(ui.Start).Gap(16).Children(func() {
		if id != "skeleton" && id != "button-group" && id != "toggle-group" {
			ui.TextInput(c, &s.Label).Label("Sample label").Placeholder("Enter sample text").FillWidth().MinWidth(0)
		}
		ui.Row(c).FillWidth().Wrap().Justify(ui.Start).AlignContent(ui.Start).Gap(8).Children(func() {
			if id == "button" || id == "badge" || id == "item" || id == "button-group" || id == "toggle" || id == "toggle-group" {
				mygoui.Toggle(c.Key("disabled"), &s.Disabled, mygoui.ToggleOptions{Label: "Disabled"})
			}
			if id == "button" {
				mygoui.Toggle(c.Key("loading"), &s.Loading, mygoui.ToggleOptions{Label: "Loading"})
				mygoui.Toggle(c.Key("invalid"), &s.Invalid, mygoui.ToggleOptions{Label: "Invalid"})
			}
			if id == "button-group" || id == "toggle-group" {
				mygoui.Toggle(c.Key("vertical"), &s.Vertical, mygoui.ToggleOptions{Label: "Vertical"})
			}
			if id == "card" {
				mygoui.Toggle(c.Key("small"), &s.Small, mygoui.ToggleOptions{Label: "Small"})
			}
			if id != "skeleton" {
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Reset fixture", Variant: mygoui.ButtonOutline,
					OnClick: func() { *s = NewBasicState() }})
			}
		})
		action := func() {
			s.Count++
			s.Result = "Activated " + strconv.Itoa(s.Count) + " times"
			c.AddToast(ui.Toast{Title: s.Result})
		}
		switch id {
		case "button":
			ui.Row(c).FillWidth().Wrap().Justify(ui.Start).AlignContent(ui.Start).Gap(8).Children(func() {
				for _, item := range []struct {
					name    string
					variant mygoui.ButtonVariant
				}{
					{"Default", mygoui.ButtonDefault}, {"Outline", mygoui.ButtonOutline},
					{"Secondary", mygoui.ButtonSecondary}, {"Ghost", mygoui.ButtonGhost},
					{"Destructive", mygoui.ButtonDestructive}, {"Link", mygoui.ButtonLink},
				} {
					mygoui.Button(c.Key(item.name), mygoui.ButtonOptions{Label: item.name,
						Variant: item.variant, Disabled: s.Disabled, Loading: s.Loading, Invalid: s.Invalid, OnClick: action})
				}
			})
			ui.Row(c).FillWidth().Wrap().Justify(ui.Start).AlignContent(ui.Start).Gap(8).Children(func() {
				for _, size := range []mygoui.ButtonSize{mygoui.ButtonXS, mygoui.ButtonSM, mygoui.ButtonSizeDefault, mygoui.ButtonLG} {
					mygoui.Button(c, mygoui.ButtonOptions{Label: s.Label, Size: size, OnClick: action,
						Disabled: s.Disabled, Loading: s.Loading, Icon: lucide.Plus()})
				}
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Increment sample", Size: mygoui.ButtonIcon,
					Icon: lucide.Plus(), Disabled: s.Disabled, OnClick: action})
			})
		case "badge":
			ui.Row(c).FillWidth().Wrap().Justify(ui.Start).AlignContent(ui.Start).Gap(8).Children(func() {
				for _, variant := range []mygoui.BadgeVariant{mygoui.BadgeDefault, mygoui.BadgeSecondary, mygoui.BadgeDestructive,
					mygoui.BadgeOutline, mygoui.BadgeGhost, mygoui.BadgeLink} {
					mygoui.Badge(c, mygoui.BadgeOptions{Label: s.Label, Variant: variant, Icon: lucide.Check()})
				}
				mygoui.Badge(c, mygoui.BadgeOptions{Label: "Action badge", OnClick: action, Disabled: s.Disabled})
			})
		case "card":
			mygoui.Card(c, mygoui.CardOptions{Small: s.Small, Label: "Sample card"}, func() {
				mygoui.CardHeader(c, mygoui.CardHeaderOptions{Title: s.Label, Description: "A native header, body and footer."})
				mygoui.CardContent(c, func() { ui.Text(c, "This content wraps as the available width changes.").Selectable() })
				mygoui.CardFooter(c, func() {
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Card action", OnClick: action, Disabled: s.Disabled})
				})
			})
		case "separator":
			ui.Row(c).FillWidth().Wrap().Justify(ui.Start).AlignContent(ui.Start).Gap(12).Children(func() {
				ui.Text(c, "Before")
				mygoui.Separator(c, mygoui.SeparatorOptions{Vertical: true, Length: 40})
				ui.Text(c, "After")
			})
			mygoui.Separator(c, mygoui.SeparatorOptions{Label: s.Label})
		case "skeleton":
			mygoui.Skeleton(c, mygoui.SkeletonOptions{Width: 48, Height: 48, Round: true})
			mygoui.Skeleton(c, mygoui.SkeletonOptions{Width: 220, Height: 16})
			mygoui.Skeleton(c, mygoui.SkeletonOptions{Width: 160, Height: 16})
		case "spinner":
			ui.Row(c).FillWidth().Wrap().Justify(ui.Start).AlignContent(ui.Start).Gap(12).Children(func() {
				mygoui.Spinner(c, mygoui.SpinnerOptions{Size: 16, Label: "Sample loading"})
				mygoui.Spinner(c, mygoui.SpinnerOptions{Size: 28, Label: "Larger sample loading"})
				mygoui.Spinner(c, mygoui.SpinnerOptions{Size: 16, Label: s.Label})
			})
		case "kbd":
			mygoui.KbdGroup(c, func() { mygoui.Kbd(c, "⌘"); mygoui.Kbd(c, "K") })
			mygoui.Kbd(c, s.Label)
		case "typography":
			for _, kind := range []mygoui.TypographyKind{mygoui.Heading1, mygoui.Heading2, mygoui.Heading3,
				mygoui.Heading4, mygoui.Paragraph, mygoui.Lead, mygoui.Large, mygoui.Small,
				mygoui.Muted, mygoui.InlineCode, mygoui.Blockquote} {
				mygoui.Typography(c, s.Label, kind)
			}
		case "aspect-ratio":
			mygoui.AspectRatio(c, 16.0/9, func() {
				ui.Column(c).Fill().Center().Background(theme.Colors(c).Muted).Children(func() { ui.Text(c, s.Label+" · 16:9") })
			}).MaxWidth(480)
		case "empty":
			mygoui.Empty(c, mygoui.EmptyOptions{Title: s.Label, Description: "Add a local sample to exercise the recovery action.",
				Icon: lucide.Folder(), Actions: func() { mygoui.Button(c, mygoui.ButtonOptions{Label: "Add sample", OnClick: action}) }})
		case "item":
			for _, variant := range []mygoui.ItemVariant{mygoui.ItemDefault, mygoui.ItemOutline, mygoui.ItemMuted} {
				mygoui.Item(c, mygoui.ItemOptions{Label: s.Label, Description: "Selectable source information with an independent action.",
					Icon: lucide.FileText(), Variant: variant, Disabled: s.Disabled, OnClick: action,
					Actions: func() {
						mygoui.Button(c, mygoui.ButtonOptions{Label: "Item action", Variant: mygoui.ButtonOutline, OnClick: action})
					}})
			}
		case "avatar":
			ui.Row(c).FillWidth().Wrap().Justify(ui.Start).AlignContent(ui.Start).Gap(12).Children(func() {
				mygoui.Avatar(c, mygoui.AvatarOptions{Name: s.Label, Size: mygoui.AvatarSM})
				mygoui.Avatar(c, mygoui.AvatarOptions{Name: s.Label})
				mygoui.Avatar(c, mygoui.AvatarOptions{Name: s.Label, Size: mygoui.AvatarLG,
					Badge: func() { ui.Icon(c, lucide.Check()).Size(12, 12) }})
			})
			mygoui.AvatarGroup(c, func() {
				mygoui.Avatar(c, mygoui.AvatarOptions{Name: "Ada Lovelace"})
				mygoui.Avatar(c, mygoui.AvatarOptions{Name: "Grace Hopper"})
				mygoui.AvatarGroupCount(c, 3)
			})
		case "button-group":
			mygoui.ButtonGroup(c, mygoui.ButtonGroupOptions{Vertical: s.Vertical, Label: "Sample actions"},
				mygoui.ButtonOptions{Label: "Back", Variant: mygoui.ButtonOutline, OnClick: func() { s.Result = "Back selected"; c.AddToast(ui.Toast{Title: s.Result}) }, Disabled: s.Disabled},
				mygoui.ButtonOptions{Label: "Forward", Variant: mygoui.ButtonOutline, OnClick: func() { s.Result = "Forward selected"; c.AddToast(ui.Toast{Title: s.Result}) }, Disabled: s.Disabled})
		case "toggle":
			mygoui.Toggle(c, &s.Pressed, mygoui.ToggleOptions{Label: s.Label, Variant: mygoui.ToggleOutline, Disabled: s.Disabled})
			ui.Text(c, "Pressed: "+strconv.FormatBool(s.Pressed))
		case "toggle-group":
			mygoui.ToggleGroup(c, &s.Group, mygoui.ToggleGroupOptions{Label: "Text alignment", Vertical: s.Vertical,
				Disabled: s.Disabled, Variant: mygoui.ToggleOutline}, []mygoui.ToggleGroupItem{
				{Value: "left", Label: "Left"}, {Value: "center", Label: "Center"}, {Value: "right", Label: "Right"}, {Value: "disabled", Label: "Unavailable", Disabled: true},
			})
			ui.Text(c, "Selected: "+strings.Join(s.Group.Values, ", "))
		}
	})
}
