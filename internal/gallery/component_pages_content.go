package gallery

import (
	"github.com/Cyenoch/mygoicons/lucide"
	uicontent "github.com/Cyenoch/mygoui/content"
	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// contentState is the caller-owned text behind the content demos.
type contentState struct {
	ready     bool
	copyText  string
	qrValue   string
	lastLink  string
	clipError string
}

func contentStateOf(s *State) *contentState {
	st := StateOf[contentState](s, "content")
	if !st.ready {
		st.ready = true
		st.copyText = "go get github.com/Cyenoch/mygoui"
		st.qrValue = "https://coss.com/ui"
	}
	return st
}

const contentGoSample = `// Greet returns a greeting for name.
func Greet(name string) string {
	if name == "" {
		return "Hello!"
	}
	return fmt.Sprintf("Hello, %s! You are visitor %d.", name, 42)
}`

const contentMarkdownSample = "## Release notes\n\n" +
	"Version **0.9** adds *native* Markdown with [links](https://coss.com/ui) and `inline code`.\n\n" +
	"- Themed headings and lists\n- Highlighted code blocks\n- [x] Task items\n\n" +
	"> Text stays selectable.\n\n" +
	"```go\nfmt.Println(\"hello\")\n```\n\n" +
	"| Key | Action |\n| --- | --- |\n| Enter | Submit |\n| Esc | Cancel |\n"

func init() {
	RegisterComponent(ComponentPage{
		ID: "copybutton", Title: "CopyButton", Group: "Display", Order: 170,
		Description: "An icon button that copies text and confirms with a check.",
		Examples: []ComponentExample{
			{Title: "Variants", Description: "Ghost and outline faces.",
				Demo: func(c *ui.Context, s *State) { copyButtonVariantsBody(c, contentStateOf(s)) }},
			{Title: "Live text", Description: "TextFunc reads the text at activation.",
				Demo: func(c *ui.Context, s *State) { copyButtonLiveBody(c, contentStateOf(s)) }},
		},
		Code: `content.CopyButton(c, content.CopyButtonOptions{
    Text:    "go get github.com/Cyenoch/mygoui",
    Label:   "Copy install command",
    Variant: content.CopyOutline,
})`,
		Params: []ComponentParam{
			{"Text", "string", "\"\"", "Text written to the clipboard."},
			{"TextFunc", "func() string", "nil", "Supplies the text at activation instead of Text."},
			{"Label", "string", "locale Copy", "Accessible name and tooltip."},
			{"Variant", "content.CopyButtonVariant", "CopyGhost", "Ghost or outline face."},
			{"Size", "theme.ButtonSize", "ButtonMedium", "Button size."},
			{"Disabled", "bool", "false", "Stops activation."},
			{"OnError", "func(error)", "nil", "Receives a clipboard failure; no copied state shows."},
		},
		Tokens: []string{"Foreground", "Accent", "Input", "Ring"},
	})

	RegisterComponent(ComponentPage{
		ID: "snippet", Title: "Snippet", Group: "Display", Order: 180,
		Description: "One line of monospace text with a copy button.",
		Examples: []ComponentExample{
			{Title: "Command", Description: "The prefix is shown but not copied.",
				Demo: func(c *ui.Context, s *State) {
					formsFieldStack(c, func() {
						uicontent.Snippet(c, uicontent.SnippetOptions{
							Text: "go get github.com/Cyenoch/mygoui", Prefix: "$", Label: "Install command",
						}).FillWidth()
					})
				}},
			{Title: "Without copy",
				Demo: func(c *ui.Context, s *State) {
					uicontent.Snippet(c, uicontent.SnippetOptions{Text: "MYGOUI_THEME=dark", NoCopy: true})
				}},
		},
		Code: `content.Snippet(c, content.SnippetOptions{
    Text:   "go get github.com/Cyenoch/mygoui",
    Prefix: "$",
    Label:  "Install command",
})`,
		Params: []ComponentParam{
			{"Text", "string", "required", "The line shown and copied."},
			{"Prefix", "string", "\"\"", "Shown before the text and left out of the copy."},
			{"Label", "string", "\"\"", "Accessible name of the snippet."},
			{"NoCopy", "bool", "false", "Leaves out the copy button."},
		},
		Tokens: []string{"SurfaceSecondary", "Border", "Foreground", "Muted"},
	})

	RegisterComponent(ComponentPage{
		ID: "highlighter", Title: "Highlighter", Group: "Display", Order: 190,
		Description: "A read-only code block with syntax colors and a copy button.",
		Examples: []ComponentExample{
			{Title: "Titled", Description: "A header names the file; long lines scroll.",
				Demo: func(c *ui.Context, s *State) {
					uicontent.Highlighter(c, uicontent.HighlighterOptions{
						Code: contentGoSample, Language: "go", Title: "greet.go",
					})
				}},
			{Title: "Wrapped", Description: "Without a header the copy button floats.",
				Demo: func(c *ui.Context, s *State) {
					uicontent.Highlighter(c, uicontent.HighlighterOptions{
						Code: `{"name": "mygoui", "private": false, "keywords": ["native", "go", "components"]}`,
						Wrap: true,
					})
				}},
		},
		Code: `content.Highlighter(c, content.HighlighterOptions{
    Code:     "func main() {\n\tfmt.Println(\"hello\")\n}",
    Language: "go",
    Title:    "main.go",
})`,
		Params: []ComponentParam{
			{"Code", "string", "required", "Source text; blank lines around it are trimmed."},
			{"Language", "string", "\"\"", "Name or alias such as go, ts or sh; unknown renders plain."},
			{"Title", "string", "language name", "Header text, as a file name."},
			{"Wrap", "bool", "false", "Wraps long lines instead of scrolling sideways."},
			{"NoCopy", "bool", "false", "Leaves out the copy button."},
		},
		Tokens: []string{"Code", "CodeForeground", "Border", "Muted", "SuccessForeground", "WarningForeground", "Primary"},
	})

	RegisterComponent(ComponentPage{
		ID: "markdown", Title: "Markdown", Group: "Display", Order: 200,
		Description: "Native Markdown in the theme's type scale and colors.",
		Examples: []ComponentExample{
			{Title: "Document",
				Demo: func(c *ui.Context, s *State) { markdownBody(c, contentStateOf(s), false) }},
			{Title: "Compact", Description: "The denser scale for messages.",
				Demo: func(c *ui.Context, s *State) { markdownBody(c, contentStateOf(s), true) }},
		},
		Code: `content.Markdown(c, content.MarkdownOptions{
    Source:  "## Notes\n\nNative **Markdown** with ` + "`code`" + `.",
    Compact: true,
    OnLink:  func(url string) { destination = url },
})`,
		Params: []ComponentParam{
			{"Source", "string", "required", "Markdown text."},
			{"Compact", "bool", "false", "14 DIP text with tighter headings and spacing."},
			{"OnLink", "func(url string)", "nil", "Handles link clicks; nil opens web and mail links."},
		},
		Tokens: []string{"Foreground", "Muted", "Link", "Default", "Border", "MutedBackground", "WarningSoft", "Primary"},
	})

	RegisterComponent(ComponentPage{
		ID: "qrcode", Title: "QRCode", Group: "Display", Order: 210,
		Description: "A crisp QR code of a caller-owned value.",
		Examples: []ComponentExample{
			{Title: "Value", Description: "The code follows the input.",
				Demo: func(c *ui.Context, s *State) { qrCodeValueBody(c, contentStateOf(s)) }},
			{Title: "With icon", Description: "An icon raises the level to high.",
				Demo: func(c *ui.Context, s *State) { qrCodeIconBody(c) }},
		},
		Code: `if _, err := content.QRCode(c, content.QRCodeOptions{
    Value: "https://coss.com/ui",
    Size:  160,
    Level: content.QRQuartile,
}); err != nil {
    form.FieldError(c, err.Error())
}`,
		Params: []ComponentParam{
			{"Value", "string", "required", "Text to encode; too long a value returns an error."},
			{"Size", "float32", "160", "Side of the square in DIPs, quiet zone included."},
			{"Level", "content.QRLevel", "QRMedium", "Low, medium, quartile or high error correction."},
			{"Foreground", "ui.Color", "light Foreground", "Module color in both appearances."},
			{"Background", "ui.Color", "light Background", "Ground color in both appearances."},
			{"Icon", "func()", "nil", "Content over the middle; raises Level to QRHigh."},
			{"Label", "string", "Value", "Accessible name."},
		},
		Tokens: []string{"Foreground", "Background"},
	})
}

func copyButtonVariantsBody(c *ui.Context, st *contentState) {
	ui.Row(c).Gap(theme.Space12).AlignItems(ui.Center).Children(func() {
		uicontent.CopyButton(c, uicontent.CopyButtonOptions{Text: st.copyText, Label: "Copy command"})
		uicontent.CopyButton(c, uicontent.CopyButtonOptions{
			Text: st.copyText, Label: "Copy command, outline", Variant: uicontent.CopyOutline,
		})
		uicontent.CopyButton(c, uicontent.CopyButtonOptions{
			Text: st.copyText, Label: "Copy command, small",
			Variant: uicontent.CopyOutline, Size: theme.ButtonSmall,
		})
		uicontent.CopyButton(c, uicontent.CopyButtonOptions{
			Text: st.copyText, Label: "Copy command, disabled", Disabled: true,
		})
	})
}

func copyButtonLiveBody(c *ui.Context, st *contentState) {
	formsFieldStack(c, func() {
		ui.Row(c).FillWidth().MinWidth(0).Gap(theme.Space8).AlignItems(ui.Center).Children(func() {
			ui.Column(c).Grow(1).Basis(0).MinWidth(0).Children(func() {
				form.Input(c, &st.copyText, form.InputOptions{AriaLabel: "Text to copy", FullWidth: true})
			})
			uicontent.CopyButton(c, uicontent.CopyButtonOptions{
				TextFunc: func() string { return st.copyText },
				Label:    "Copy typed text", Variant: uicontent.CopyOutline,
				OnError: func(err error) { st.clipError = err.Error() },
			})
		})
		form.FieldError(c, st.clipError)
	})
}

func markdownBody(c *ui.Context, st *contentState, compact bool) {
	ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space12).Children(func() {
		uicontent.Markdown(c, uicontent.MarkdownOptions{
			Source: contentMarkdownSample, Compact: compact,
			OnLink: func(url string) { st.lastLink = url },
		})
		if st.lastLink != "" {
			formsMutedText(c, st.lastLink)
		}
	})
}

func qrCodeValueBody(c *ui.Context, st *contentState) {
	ui.Column(c).FillWidth().MinWidth(0).MaxWidth(theme.Scale(c, 320)).
		Gap(theme.Space16).AlignItems(ui.Center).Children(func() {
		_, err := uicontent.QRCode(c, uicontent.QRCodeOptions{Value: st.qrValue})
		form.Input(c, &st.qrValue, form.InputOptions{AriaLabel: "QR value", FullWidth: true})
		if err != nil {
			form.FieldError(c, err.Error())
		}
	})
}

func qrCodeIconBody(c *ui.Context) {
	// The value is a constant that always encodes.
	_, _ = uicontent.QRCode(c, uicontent.QRCodeOptions{
		Value: "https://github.com/Cyenoch/mygoui", Size: 200,
		Icon: func() {
			size := theme.Scale(c, theme.IconLG)
			ui.Icon(c, lucide.Github()).Size(size, size).TextColor(theme.ColorsForMode(false).Foreground)
		},
	})
}
