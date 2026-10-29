package gallery

const cossSourceRoot = "https://github.com/cosscom/ui/blob/3d8ca28e7d90123e9e0f44665944dc240b026ccb/apps/ui/"

// cossSources records the upstream docs/registry counterpart. Everything else
// is a native addition, rather than an assertion of upstream feature parity.
var cossSources = map[string]string{
	"accordion": "accordion", "alert": "alert", "alertdialog": "alert-dialog",
	"autocomplete": "autocomplete", "avatar": "avatar", "badge": "badge",
	"breadcrumbs": "breadcrumb", "button": "button", "calendar": "calendar",
	"card": "card", "checkbox": "checkbox", "checkboxgroup": "checkbox-group",
	"disclosure": "collapsible", "combobox": "combobox", "composer": "composer",
	"command": "command", "contextmenu": "context-menu", "datepicker": "date-picker",
	"dialog": "dialog", "drawer": "drawer", "emptystate": "empty", "field": "field",
	"fieldset": "fieldset", "form": "form", "frame": "frame", "group": "group",
	"input": "input", "inputgroup": "input-group", "kbd": "kbd", "label": "label",
	"menu": "menu", "meter": "meter", "numberfield": "number-field",
	"inputotp": "otp-field", "pagination": "pagination", "popover": "popover",
	"previewcard": "preview-card", "progressbar": "progress", "radiogroup": "radio-group",
	"scrollarea": "scroll-area", "segmentedcontrol": "segmented-control",
	"select": "select", "separator": "separator", "sheet": "sheet", "sidebar": "sidebar",
	"skeleton": "skeleton", "slider": "slider", "spinner": "spinner", "switch": "switch",
	"table": "table", "tabs": "tabs", "textarea": "textarea", "toast": "toast",
	"togglebutton": "toggle", "togglebuttongroup": "toggle-group",
	"toolbar": "toolbar", "tooltip": "tooltip",
}

func sourceURL(id string) string {
	if url := chartSource(id); url != "" {
		return url
	}
	if id == "theming" {
		return cossSourceRoot + "content/docs/(root)/styling.mdx"
	}
	slug := cossSources[id]
	if slug == "" {
		return ""
	}
	if slug == "sidebar" {
		return cossSourceRoot + "registry/default/ui/sidebar.tsx"
	}
	return cossSourceRoot + "content/docs/components/" + slug + ".mdx"
}

func referenceLabel(id string) string {
	if chartSource(id) != "" {
		return "Bklit source"
	}
	return "Coss source"
}

// Coss names are the public semantic roles. Native derived colors remain
// documented on addition pages, while source pages refer to their base roles.
func cossTokens(p ComponentPage) []string {
	if cossSources[p.ID] == "" {
		return p.Tokens
	}
	switch p.ID {
	case "checkbox", "checkboxgroup", "radiogroup", "switch", "slider":
		return []string{"Background", "Foreground", "Input", "Primary", "PrimaryForeground", "Destructive", "Ring"}
	case "calendar":
		return []string{"Background", "Foreground", "Primary", "PrimaryForeground", "Accent", "MutedForeground", "Ring"}
	case "datepicker":
		return []string{"Background", "Input", "Ring", "Popover", "Primary", "PrimaryForeground", "Accent"}
	case "progressbar", "meter":
		return []string{"Primary", "Input", "Foreground", "MutedForeground"}
	case "spinner":
		return []string{"Primary", "Info", "Success", "Warning", "Destructive"}
	case "segmentedcontrol", "tabs":
		return []string{"MutedBackground", "MutedForeground", "Background", "Foreground", "Input", "Ring"}
	}
	aliases := map[string]string{
		"Surface": "Card", "SurfaceForeground": "CardForeground",
		"SurfaceSecondary": "MutedBackground", "SurfaceSecondaryForeground": "MutedForeground",
		"SurfaceTertiary": "MutedBackground", "SurfaceTertiaryForeground": "MutedForeground",
		"SurfaceHover": "Accent", "Overlay": "Popover", "OverlayForeground": "PopoverForeground",
		"Muted": "MutedForeground", "Default": "Secondary", "DefaultForeground": "SecondaryForeground",
		"DefaultHover": "Accent", "DefaultSoft": "Secondary", "DefaultSoftForeground": "SecondaryForeground",
		"DefaultSoftHover": "Accent", "AccentSoft": "Accent", "AccentSoftHover": "Accent",
		"AccentSoftForeground": "AccentForeground", "AccentHover": "Accent",
		"Danger": "Destructive", "DangerHover": "Destructive", "DangerForeground": "DestructiveForeground",
		"DangerSoft": "Destructive", "DangerSoftForeground": "DestructiveForeground",
		"FieldBackground": "Background", "FieldForeground": "Foreground", "FieldPlaceholder": "MutedForeground",
		"FieldBorder": "Input", "FieldBorderHover": "Input", "FieldBorderFocus": "Ring",
		"FieldFocus": "Ring", "FieldHover": "Accent", "Focus": "Ring",
		"Separator": "Border", "SeparatorSecondary": "Border", "SeparatorTertiary": "Border",
		"Segment": "Background", "SegmentForeground": "Foreground", "Selection": "Accent",
	}
	tokens := make([]string, 0, len(p.Tokens))
	seen := map[string]bool{}
	for _, token := range p.Tokens {
		if canonical, ok := aliases[token]; ok {
			token = canonical
		}
		if !seen[token] {
			seen[token] = true
			tokens = append(tokens, token)
		}
	}
	return tokens
}
