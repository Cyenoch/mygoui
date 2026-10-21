package gallery

import (
	"strings"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/codeblock"
	"github.com/Cyenoch/mygoui/internal/gallery/locales"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func selectComponent(s *State, id ComponentID) {
	s.Selected = id
	s.previewScroll = ui.ScrollState{}
	s.SourceOpen = false
	s.JumpSection = ""
}

func preview(c *ui.Context, s *State) {
	entry := entryFor(s.Selected)
	colors := theme.Colors(c)
	ui.Scroll(c.Key("preview-"+entry.ID)).Label(locales.ComponentGalleryPreview()).TrackScroll(&s.previewScroll).
		Grow(1).Basis(0).FillHeight().MinHeight(0).MinWidth(0).Padding(24, 16).
		Justify(ui.Start).AlignItems(ui.Stretch).Children(func() {
		ui.Column(c).Label(locales.ComponentGalleryDocumentation()).FillWidth().MaxWidth(800).Margin(0, ui.Auto).
			MinWidth(0).Justify(ui.Start).AlignItems(ui.Stretch).Gap(24).Children(func() {
			ui.Column(c).Gap(12).Children(func() {
				ui.Text(c, locales.ComponentGalleryName(strings.ReplaceAll(string(entry.ID), "-", "_"))).Bold().FontSize(theme.Scale(c, 30))
				ui.Text(c, entry.Plan).FontSize(theme.Scale(c, 16)).LineHeight(1.5).
					TextColor(colors.MutedForeground)
				ui.Row(c).Wrap().Gap(8).Children(func() {
					ui.Text(c, "Go / MyGo").FontSize(theme.Scale(c, 13))
					ui.Text(c, locales.ComponentGalleryImplementedUnverified()).FontSize(theme.Scale(c, 12)).TextColor(colors.MutedForeground)
				})
			})
			section(c, s, "Preview")
			ui.Column(c).FillWidth().MinWidth(0).Border(1, colors.Border).Radius(theme.Radii().LG).Clip().Children(func() {
				ui.Column(c).Label(locales.ComponentGalleryInteractivePreview()).FillWidth().MinWidth(0).MinHeight(300).
					Padding(24).Justify(ui.Center).AlignItems(ui.Stretch).Children(func() { defaultSpecimen(c.Key("default-demo"), s, entry.ID) })
				ui.Column(c).FillWidth().MinWidth(0).Gap(8).Background(colors.Muted).
					BorderWidth(1, 0, 0, 0).BorderColor(colors.Border).Children(func() {
					label := locales.ComponentGalleryViewCode()
					if s.SourceOpen {
						label = locales.ComponentGalleryHideCode()
					}
					mygoui.Button(c.Key("view-code"), mygoui.ButtonOptions{Label: label,
						Variant: mygoui.ButtonOutline, Size: mygoui.ButtonSM,
						OnClick: func() { s.SourceOpen = !s.SourceOpen }}).AlignSelf(ui.Center).Margin(12, 0, 0, 0)
					if s.SourceOpen {
						codeblock.View(c.Key("actual-preview-source"), sourceFor(entry.ID), &s.DocumentationCode,
							codeblock.Options{Language: "go", Title: locales.ComponentGalleryActualSource(), LineNumbers: true, MaxHeight: 440, Embedded: true})
					} else {
						ui.Text(c, "Go · Native MyGo API").FontSize(theme.Scale(c, 12)).TextColor(colors.MutedForeground).AlignSelf(ui.Center).Margin(0, 0, 12, 0)
					}
				})
			})
			section(c, s, "Usage")
			ui.Text(c, usageFor(entry)).LineHeight(1.5)
			section(c, s, "Examples")
			ui.Text(c, locales.ComponentGalleryVariants(locales.ComponentGalleryName(strings.ReplaceAll(string(entry.ID), "-", "_")))).Bold().FontSize(theme.Scale(c, 18))
			ui.Text(c, examplesFor(entry)).LineHeight(1.5)
			ui.Column(c).FillWidth().MinWidth(0).Padding(24).Border(1, colors.Border).Radius(theme.Radii().LG).
				Children(func() { specimen(c.Key("advanced-examples"), s, entry.ID) })
			section(c, s, "Behavior and API")
			ui.Text(c, "State belongs to the caller and survives view rebuilds. Actions above change only local sample fixtures. Native controls retain keyboard focus, selection and accessibility semantics.").LineHeight(1.5)
			reference := "https://ui.shadcn.com/docs/components/" + string(entry.ID)
			ui.Text(c, locales.ComponentGalleryReference(reference)).Selectable().FontSize(theme.Scale(c, 12)).TextColor(colors.MutedForeground)
			adjacent(c, s)
		})
	})
}

func section(c *ui.Context, s *State, name string) {
	heading := ui.Text(c, locales.ComponentGallerySection(strings.ReplaceAll(name, " ", "_"))).Bold().FontSize(theme.Scale(c, 20)).Margin(12, 0, 0, 0)
	if s.JumpSection == name {
		heading.ScrollIntoView()
		s.JumpSection = ""
	}
}

func outline(c *ui.Context, s *State) {
	ui.Column(c).Width(176).Shrink(0).Padding(28, 16).Gap(8).Children(func() {
		ui.Text(c, locales.ComponentGalleryOnThisPage()).FontSize(theme.Scale(c, 12)).TextColor(theme.Colors(c).MutedForeground)
		for _, name := range []string{"Preview", "Usage", "Examples", "Behavior and API"} {
			label := locales.ComponentGallerySection(strings.ReplaceAll(name, " ", "_"))
			mygoui.Button(c.Key(name), mygoui.ButtonOptions{Label: label, Variant: mygoui.ButtonGhost,
				Size: mygoui.ButtonSM, OnClick: func() { s.JumpSection = name }}).Label(locales.ComponentGalleryJumpSection(label)).FillWidth().Justify(ui.Start)
		}
	})
}

func adjacent(c *ui.Context, s *State) {
	ui.Row(c).FillWidth().Wrap().Justify(ui.SpaceBetween).Gap(12).Margin(24, 0).Children(func() {
		for i, entry := range Catalog {
			if entry.ID != s.Selected {
				continue
			}
			if i > 0 {
				previous := Catalog[i-1]
				mygoui.Button(c.Key("previous"), mygoui.ButtonOptions{Label: "← " + locales.ComponentGalleryName(strings.ReplaceAll(string(previous.ID), "-", "_")),
					Variant: mygoui.ButtonOutline, OnClick: func() { selectComponent(s, previous.ID) }})
			}
			if i+1 < len(Catalog) {
				next := Catalog[i+1]
				mygoui.Button(c.Key("next"), mygoui.ButtonOptions{Label: locales.ComponentGalleryName(strings.ReplaceAll(string(next.ID), "-", "_")) + " →",
					Variant: mygoui.ButtonOutline, OnClick: func() { selectComponent(s, next.ID) }})
			}
			break
		}
	})
}

func specimen(c *ui.Context, s *State, id ComponentID) {
	switch id {
	case "accordion", "breadcrumb", "collapsible", "combobox", "command", "navigation-menu", "pagination", "select", "sidebar", "tabs":
		navigationSample(c, &s.Navigation, id, true)
	case "alert", "progress", "direction", "resizable", "scroll-area", "carousel", "table", "data-table":
		dataLayoutSample(c, &s.DataLayout, id, true)
	case "calendar", "date-picker", "chart":
		PreviewDateChart(c, &s.Dates, id)
	case "code-block":
		PreviewCode(c, &s.Code)
	case "trace":
		PreviewTrace(c, &s.Trace)
	case "input", "textarea", "field", "label", "input-group", "input-otp", "checkbox", "radio-group", "switch", "slider", "native-select":
		PreviewForm(c, &s.Forms, id)
	case "dialog", "alert-dialog", "sheet", "drawer", "popover", "tooltip", "hover-card", "dropdown-menu", "context-menu", "menubar", "toast":
		PreviewOverlay(c, &s.Overlays, id)
	case "message", "bubble", "marker", "attachment", "message-scroller", "questionnaire":
		PreviewConversation(c, &s.Conversation, id)
	default:
		PreviewBasic(c, &s.Basic, id)
	}
}

func usageFor(entry Entry) string {
	switch entry.ID {
	case "calendar":
		return "Import github.com/Cyenoch/mygoui/calendar. Retain calendar.State and pass civil Date values with explicit location/constraints. Calendar.Options selects single, multiple or range mode; locale changes presentation while arithmetic remains Gregorian. The sample's fixed October 2026 date avoids clock-dependent data."
	case "date-picker":
		return "Import github.com/Cyenoch/mygoui/calendar. Retain DatePickerState; its Calendar field owns valid selection while editor drafts remain separate. DatePickerOptions supplies single/range behavior, explicit presets and clearing. Text uses strict YYYY-MM-DD civil dates; invalid input does not silently replace valid state."
	case "chart":
		return "Import github.com/Cyenoch/mygoui/chart. Supply category labels and typed series values with stable SeriesID and semantic Chart color tokens. Retain chart.State for selection/legend visibility; Options selects the native chart kind and stacking. Missing values are explicit, and zero remains a measured zero."
	case "table", "data-table":
		return "Import github.com/Cyenoch/mygoui/data. Give rows and columns stable typed IDs; Cells follow the original Columns order. Retain TableState or DataTableState. Bound the table viewport explicitly and increment Revision when existing source cells or comparator behavior change. The editable sample appends actual local rows."
	}
	switch entry.Family {
	case "Navigation":
		return "Import github.com/Cyenoch/mygoui/navigation. Keep selected values, queries and disclosure state in the caller. Use stable typed item identities and native keyboard/focus behavior. Expand Actual preview source for the component's real construction and use Examples to compare its interaction states."
	case "Layout":
		return "Import github.com/Cyenoch/mygoui/layout. Allocate a bounded native region for scrolling and resizing; keep the selected page, split size or scroll position in caller-owned state. The default specimen and editable Examples use the actual maintained layout widgets."
	case "Data":
		return "Use the component-specific native data, calendar or chart package shown in Actual preview source. Keep rows, dates and series in caller-owned state; typed callbacks update those values. Sorting, filtering and pagination operate on the local sample data displayed below."
	case "Forms":
		return "Import github.com/Cyenoch/mygoui/form. Bind the input to persistent Go state; ControlOptions supplies labels, disabled/read-only states and errors. Field associates visible labels and help with the native editor. The expandable source includes the exact imports, fixture state and calls used above."
	case "Overlays":
		return "Import github.com/Cyenoch/mygoui/overlay. Keep open state in the caller, render the trigger and overlay in the same view, and provide native body/footer callbacks. Modal dismissal restores focus; action callbacks decide when to close. The source above is the compiled gallery implementation."
	case "Conversation":
		return "Import github.com/Cyenoch/mygoui/conversation. Keep stable message identities and caller-owned state. Compose native content with callbacks; attachment actions, reactions and questionnaire submission update the local fixture. The exact compiled preview source is available above."
	case "Rich source":
		return "Import github.com/Cyenoch/mygoui/codeblock or github.com/Cyenoch/mygoui/trace. Retain presentation state across rebuilds and pass original source unchanged. CodeBlock Language selects native syntax highlighting; Trace retains disclosure state independently from the sample text. Expand the source for actual typed options and composition."
	default:
		return "Import github.com/Cyenoch/mygoui and construct the component inside a native MyGo view. Pass typed options and actual callbacks; persistent fixture state stays outside the view. Expand Actual preview source for the exact imports, constructor and variants used above."
	}
}

func examplesFor(entry Entry) string {
	switch entry.ID {
	case "button":
		return "Compare default, outline, secondary, ghost, destructive and link actions, then test loading, disabled and invalid states. The second row compares button and icon sizes."
	case "badge":
		return "Compare semantic badge variants and the action badge. Change the sample label and disable its actual action."
	case "card":
		return "Compare the standard and small card layouts. Edit the header text and activate the footer action."
	case "separator":
		return "Compare a vertical separator between two labels and a labelled horizontal separator."
	case "skeleton":
		return "An avatar placeholder and two text measures show paint-only shimmer, with stationary placeholders under reduced motion or high contrast."
	case "spinner":
		return "Compare the compact and large native busy indicators and edit the third indicator's accessible label."
	case "kbd":
		return "Compare a grouped Command-K chord with a single keycap using the editable sample text."
	case "typography":
		return "Compare heading levels, paragraph, lead, large, small, muted, inline code and blockquote using the same sample text."
	case "aspect-ratio":
		return "The sample keeps a 16:9 ratio as the available width changes; edit its label to see wrapping inside the native frame."
	case "empty":
		return "Edit the empty-state title and use Add sample to exercise the recovery action."
	case "item":
		return "Compare default, outline and muted rows with a leading file icon, description and independently actionable trailing control."
	case "avatar":
		return "Compare small, regular and large avatars, a status badge and a grouped overflow count."
	case "button-group":
		return "Switch the joined Back/Forward actions between horizontal and vertical layout; each action retains its own keyboard focus."
	case "toggle":
		return "Compare the outline toggle's persistent pressed value with its actual disabled state."
	case "toggle-group":
		return "Compare horizontal and vertical alignment choices, retained selection and an unavailable option."
	case "code-block":
		return "Edit exact Go source; append a line, change line numbers, size or destructive styling. A separate block demonstrates removed, added and highlighted lines."
	case "trace":
		return "Edit reasoning and output, change running/error status, expand both nesting levels and append output without replacing retained disclosure state."
	case "input":
		return "Edit and submit the native text field. Compare disabled, read-only, invalid and loading states, then reset its value."
	case "textarea":
		return "Edit multiple lines with native selection and undo, then compare read-only and invalid presentation."
	case "field":
		return "Compare associated display-name and notes fields inside a labelled preference set. Help and validation remain attached to their input."
	case "label":
		return "Click Focus sample input to move native focus into the associated editor; disabled input retains its disabled semantics."
	case "input-group":
		return "Edit an address beside a compact prefix, then activate the independent Use sample action in the responsive lower addon row."
	case "input-otp":
		return "Enter or paste a six-character verification code into the single native editor. The fixed slot row scrolls locally when text is enlarged."
	case "checkbox":
		return "Compare unchecked, checked and mixed states. The visible result reports both the checked and indeterminate fixture values."
	case "radio-group":
		return "Choose Fast, Balanced or Complete with the native roving radio group; the Unavailable option cannot be selected."
	case "switch":
		return "Toggle the sample option and compare its checked value with disabled and invalid presentation."
	case "slider":
		return "Compare a stepped single value and a bounded two-thumb range. Drag or use native arrow keys; the fixture reports the resulting values."
	case "native-select":
		return "Choose among grouped native menu options and confirm that the unavailable choice remains disabled."
	case "dialog", "sheet", "drawer", "alert-dialog":
		return "Disable the opener, then open the editable sample and save or dismiss it. Reset restores its independent fixture. Drawer uses a native edge modal; swipe and velocity snapping are not implemented."
	case "popover":
		return "Open the anchored sample and edit its native text field. Outside click and Escape dismiss through the maintained popup behavior."
	case "tooltip", "hover-card":
		return "Focus or hover the sample anchor to reveal its native description; activating the anchor records a separate local result."
	case "dropdown-menu", "context-menu", "menubar":
		return "Exercise the real action, checkbox, nested radio choices and disabled item. Context Menu opens with right click or Shift+F10."
	case "toast":
		return "Insert three notifications with unequal heights, expand the deck by hover or focus, update the same loading notification to success, remove a middle card, then dismiss the local samples. Actions use native callbacks and notifications stay outside document layout."
	case "message":
		return "A complete sample conversation compares opposite alignments, primary and secondary bubbles, honest avatar fallbacks, grouped replies, delivery metadata and a persistent attached reaction. Toggle Disable actions to compare reaction behavior without changing the conversation layout."
	case "bubble":
		return "Compare all implemented bubble surfaces and toggle a persistent Helpful reaction."
	case "marker":
		return "Compare plain, separator and bordered message markers plus an actionable status marker."
	case "attachment":
		return "Compare upload progress, failure/retry and removal while keeping Open and Remove actions independent."
	case "message-scroller":
		return "Append, grow and prepend stable local messages, or jump to the first sample inside the independently bounded native list."
	case "questionnaire":
		return "Answer a required single-choice/freeform question and optional multiple-choice question, then save the typed answers."
	case "accordion":
		return "Expand one or several sections, move between native disclosure triggers with the keyboard, and compare disabled content without losing retained expansion."
	case "alert":
		return "Compare a neutral informational callout and a destructive message with a real recovery action. The message remains readable independently of its color."
	case "breadcrumb":
		return "Activate earlier path segments while the current page remains nonactionable. Compare long labels and the native overflow menu."
	case "calendar":
		return "Select dates and change months with the maintained native date grid. Compare allowed and disabled dates without changing locale arithmetic."
	case "carousel":
		return "Move between sample pages with previous/next controls, observe disabled boundaries, and retain the selected page across rebuilds."
	case "chart":
		return "Inspect real sample series using the native plot, compare chart marks and selection, and read the underlying values without relying on color alone."
	case "collapsible":
		return "Expand and collapse the sample details with the native trigger; the caller-owned open state survives unrelated fixture changes."
	case "combobox":
		return "Filter sample choices, navigate matches with arrow keys, and commit a choice with Enter. Disabled options remain unavailable."
	case "command":
		return "Search local commands, move through matching results and execute an actual sample action. Empty results offer no phantom selection."
	case "data-table":
		return "Search, sort and page real local rows; selection tracks stable row identities rather than visible row positions."
	case "date-picker":
		return "Edit or choose a sample date, open the calendar popup and clear the selection. Invalid committed input remains an explicit field state."
	case "direction":
		return "Switch the native sample between left-to-right and right-to-left flow; logical content edges reverse while source code stays left-to-right."
	case "navigation-menu":
		return "Open native navigation panels with pointer or keyboard and activate a real local destination. Escape dismisses and restores its trigger."
	case "pagination":
		return "Change the sample page, activate previous/next actions and inspect the disabled first/last boundaries."
	case "progress":
		return "Change a bounded sample value and compare determinate progress with native indeterminate busy feedback. Zero remains a true zero."
	case "resizable":
		return "Drag the native pane divider or use its keyboard controls. Each pane retains its minimum size while the caller owns the split value."
	case "scroll-area":
		return "Scroll the bounded sample with wheel, keyboard and native scrollbar input; its content remains separate from surrounding documentation."
	case "select":
		return "Open the native choice popup, navigate or type ahead, and commit a value. Groups and disabled items keep their actual semantics."
	case "sidebar":
		return "Use a complete workspace shell: change the team, expand Playground's nested links, select a destination, open the account action and toggle the desktop rail or narrow sheet. Header and account footer stay fixed while groups and main content scroll independently."
	case "table":
		return "Inspect actual native rows and columns, select rows and activate the available row action without replacing the underlying data."
	case "tabs":
		return "Switch real native panels using pointer or roving keyboard focus; the unavailable tab is skipped and hidden panel state is retained."
	}
	return entry.Plan
}
