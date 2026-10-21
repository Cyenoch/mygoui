package gallery

// ComponentID distinguishes catalog identity from editable labels and source text.
type ComponentID string

// Entry describes a native specimen and its documented capabilities.
type Entry struct {
	ID                 ComponentID
	Name, Family, Plan string
}

// Catalog lists the 66 implemented native component specimens.
var Catalog = []Entry{
	{ID: "accordion", Name: "Accordion", Family: "Navigation", Plan: "Native disclosure group with single/multiple selection and arrow-key focus."},
	{ID: "alert", Name: "Alert", Family: "Presentation", Plan: "Native callout with semantic status, description and optional real action."},
	{ID: "alert-dialog", Name: "Alert Dialog", Family: "Overlays", Plan: "Modal focus, explicit action/cancel, no outside-click acceptance."},
	{ID: "aspect-ratio", Name: "Aspect Ratio", Family: "Presentation", Plan: "Native aspect-ratio layout with caller-owned media."},
	{ID: "attachment", Name: "Attachment", Family: "Conversation", Plan: "File/image metadata, lifecycle status and independent open/remove/retry actions."},
	{ID: "avatar", Name: "Avatar", Family: "Presentation", Plan: "Caller-loaded native image, initials fallback, badge and overflow count."},
	{ID: "badge", Name: "Badge", Family: "Presentation", Plan: "Semantic variants; passive text or actual keyboard-accessible action."},
	{ID: "breadcrumb", Name: "Breadcrumb", Family: "Navigation", Plan: "Native path actions, current-page semantics and overflow menu."},
	{ID: "bubble", Name: "Bubble", Family: "Conversation", Plan: "Selectable content, grouped alignment, reactions and retained expansion."},
	{ID: "button", Name: "Button", Family: "Presentation", Plan: "Native variants, sizes, loading/disabled states and real actions."},
	{ID: "button-group", Name: "Button Group", Family: "Presentation", Plan: "Joined native actions preserving independent keyboard focus."},
	{ID: "calendar", Name: "Calendar", Family: "Data", Plan: "Localized native date grid, selection and month keyboard navigation."},
	{ID: "card", Name: "Card", Family: "Presentation", Plan: "Composable header, content and footer with shared native spacing."},
	{ID: "carousel", Name: "Carousel", Family: "Layout", Plan: "Scrollable items with actual next/previous selection and bounded navigation."},
	{ID: "chart", Name: "Chart", Family: "Data", Plan: "Maintained native charts, semantic series colors and keyboard selection."},
	{ID: "checkbox", Name: "Checkbox", Family: "Forms", Plan: "Native checked/mixed/disabled states and accessible label."},
	{ID: "code-block", Name: "Code Block", Family: "Rich source", Plan: "Exact source copy, wrap, numbers, line variants and native selection."},
	{ID: "collapsible", Name: "Collapsible", Family: "Navigation", Plan: "Native disclosure with persistent caller-owned expansion."},
	{ID: "combobox", Name: "Combobox", Family: "Navigation", Plan: "Native filtering, keyboard highlight and single/multiple choices."},
	{ID: "command", Name: "Command", Family: "Navigation", Plan: "Native search input and list with real command callbacks."},
	{ID: "context-menu", Name: "Context Menu", Family: "Overlays", Plan: "Native right-click and Shift+F10 menu, checked/disabled/submenu actions."},
	{ID: "data-table", Name: "Data Table", Family: "Data", Plan: "Native table with real Go sorting, filtering, selection and pagination."},
	{ID: "date-picker", Name: "Date Picker", Family: "Data", Plan: "Native date editor and calendar popup with caller-owned date/range."},
	{ID: "dialog", Name: "Dialog", Family: "Overlays", Plan: "Native modal with editable content, focus trap and focus restoration."},
	{ID: "direction", Name: "Direction", Family: "Layout", Plan: "Native logical layout direction; code remains left to right."},
	{ID: "drawer", Name: "Drawer", Family: "Overlays", Plan: "Native edge modal; gestures require separate native input evidence."},
	{ID: "dropdown-menu", Name: "Dropdown Menu", Family: "Overlays", Plan: "Native menu with actual actions, checks and submenus."},
	{ID: "empty", Name: "Empty", Family: "Presentation", Plan: "Native empty-state media, explanation and real recovery action."},
	{ID: "field", Name: "Field", Family: "Forms", Plan: "Native associated label, description, errors and field groups."},
	{ID: "hover-card", Name: "Hover Card", Family: "Overlays", Plan: "Native anchored preview with keyboard route and delayed hover."},
	{ID: "input", Name: "Input", Family: "Forms", Plan: "Native editor with IME, selection, undo and actual bound value."},
	{ID: "input-group", Name: "Input Group", Family: "Forms", Plan: "Native editor with independent inline/block addon actions."},
	{ID: "input-otp", Name: "Input OTP", Family: "Forms", Plan: "One native editor with slot guides, paste filtering and completion."},
	{ID: "item", Name: "Item", Family: "Presentation", Plan: "Native content/media/actions; static and actionable variants."},
	{ID: "kbd", Name: "Kbd", Family: "Presentation", Plan: "Native keycap labels and chord grouping."},
	{ID: "label", Name: "Label", Family: "Forms", Plan: "Native focus handle association and accessible input naming."},
	{ID: "marker", Name: "Marker", Family: "Conversation", Plan: "Textual status, bordered row and labelled separator."},
	{ID: "menubar", Name: "Menubar", Family: "Overlays", Plan: "Native menu actions; macOS application placement is platform-owned."},
	{ID: "message", Name: "Message", Family: "Conversation", Plan: "Caller-composed avatar/header/content/footer and alignment."},
	{ID: "message-scroller", Name: "Message Scroller", Family: "Conversation", Plan: "Native stable list keys, intentional follow, anchors and history preservation."},
	{ID: "native-select", Name: "Native Select", Family: "Forms", Plan: "Native platform choice menu with grouped/disabled options."},
	{ID: "navigation-menu", Name: "Navigation Menu", Family: "Navigation", Plan: "Native actions and popovers with keyboard traversal."},
	{ID: "pagination", Name: "Pagination", Family: "Navigation", Plan: "Actual caller page changes, boundary disabling and current-page semantics."},
	{ID: "popover", Name: "Popover", Family: "Overlays", Plan: "Native anchored editable popup, edge flipping and dismissal."},
	{ID: "progress", Name: "Progress", Family: "Presentation", Plan: "Native bounded value and accessible progress text."},
	{ID: "questionnaire", Name: "Questionnaire", Family: "Conversation", Plan: "Typed Go questions/answers, validation, previous/skip/submit and scoped shortcuts."},
	{ID: "radio-group", Name: "Radio Group", Family: "Forms", Plan: "Native single-choice selection and roving keyboard focus."},
	{ID: "resizable", Name: "Resizable", Family: "Layout", Plan: "Maintained native Split with pointer and keyboard divider resizing."},
	{ID: "scroll-area", Name: "Scroll Area", Family: "Layout", Plan: "Native scroll container with wheel, keyboard and scrollbar input."},
	{ID: "select", Name: "Select", Family: "Forms", Plan: "Native popup choice, typeahead, groups and selected-value semantics."},
	{ID: "separator", Name: "Separator", Family: "Presentation", Plan: "Native horizontal/vertical semantic rule."},
	{ID: "sheet", Name: "Sheet", Family: "Overlays", Plan: "Native edge modal preserving focus and editable content."},
	{ID: "sidebar", Name: "Sidebar", Family: "Navigation", Plan: "Native source list with collapse, item actions and narrow overlay."},
	{ID: "skeleton", Name: "Skeleton", Family: "Presentation", Plan: "Native static placeholders with text/avatar/card shapes."},
	{ID: "slider", Name: "Slider", Family: "Forms", Plan: "Native drag/keyboard stepped single value and two-thumb range."},
	{ID: "spinner", Name: "Spinner", Family: "Presentation", Plan: "Native busy indicator with reduced-motion-safe text."},
	{ID: "switch", Name: "Switch", Family: "Forms", Plan: "Native boolean control with actual checked state."},
	{ID: "table", Name: "Table", Family: "Data", Plan: "Native virtual rows and columns with selection and row actions."},
	{ID: "tabs", Name: "Tabs", Family: "Navigation", Plan: "Native tab selection, roving keyboard focus and actual panels."},
	{ID: "textarea", Name: "Textarea", Family: "Forms", Plan: "Native multiline editor with caller-owned source and read-only selection."},
	{ID: "toast", Name: "Toast", Family: "Overlays", Plan: "Native toast manager with actual action, update and dismissal."},
	{ID: "toggle", Name: "Toggle", Family: "Presentation", Plan: "Native persistent pressed-state button."},
	{ID: "toggle-group", Name: "Toggle Group", Family: "Presentation", Plan: "Native single/multiple selection and disabled-skipping navigation."},
	{ID: "tooltip", Name: "Tooltip", Family: "Overlays", Plan: "Native hover/focus description with real anchor and Escape dismissal."},
	{ID: "trace", Name: "Trace", Family: "Rich source", Plan: "Native status/disclosure composition with persistent source-keyed expansion."},
	{ID: "typography", Name: "Typography", Family: "Presentation", Plan: "Native semantic type roles and append-stable prose rhythm."},
}

func entryFor(id ComponentID) Entry {
	for _, entry := range Catalog {
		if entry.ID == id {
			return entry
		}
	}
	return Catalog[0]
}
