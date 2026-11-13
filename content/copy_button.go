package content

import (
	"time"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// CopyButtonVariant selects the copy button face.
type CopyButtonVariant uint8

const (
	CopyGhost CopyButtonVariant = iota
	CopyOutline
)

// CopyButtonOptions configures a copy button. TextFunc, when set, supplies
// the text at activation instead of Text. Label names what is copied for
// assistive technology and the tooltip; it defaults to the locale's "Copy".
// OnError receives a clipboard failure, after which no copied state shows.
type CopyButtonOptions struct {
	Text     string
	TextFunc func() string
	Label    string
	Variant  CopyButtonVariant
	Size     theme.ButtonSize
	Disabled bool
	OnError  func(error)
}

// copiedFeedback is how long a copy button shows its copied state.
const copiedFeedback = 1500 * time.Millisecond

// copyState is the transient copied feedback: when the text reached the
// clipboard, and whether that is still to be announced.
type copyState struct {
	at       time.Time
	announce bool
}

type copyStateKey struct{}

// CopyButton is an icon-only button that writes the caller's text to the
// clipboard and shows a check for a moment once the text got there.
func CopyButton(c *ui.Context, o CopyButtonOptions) ui.Element {
	m := locale.M()
	root := ui.Row(c).Shrink(0)
	st := ui.Local(root, copyStateKey{}, func() copyState { return copyState{} })
	label, icon := o.Label, lucide.Copy()
	if label == "" {
		label = m.Copy
	}
	if !st.at.IsZero() {
		if rest := copiedFeedback - c.Now().Sub(st.at); rest > 0 {
			label, icon = m.Copied, lucide.Check()
			c.After(rest)
			if st.announce {
				c.Announce(m.Copied)
			}
		} else {
			st.at = time.Time{}
		}
		st.announce = false
	}
	variant := mygoui.ButtonGhost
	if o.Variant == CopyOutline {
		variant = mygoui.ButtonOutline
	}
	return root.Children(func() {
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: label, Icon: icon, IconOnly: true, Variant: variant,
			Size: o.Size, Disabled: o.Disabled,
			OnClick: func() {
				text := o.Text
				if o.TextFunc != nil {
					text = o.TextFunc()
				}
				// The result arrives outside a build, where the context
				// has no clock.
				now := c.Now()
				c.WriteClipboardThen(text, func(err error) {
					if err != nil {
						if o.OnError != nil {
							o.OnError(err)
						}
						return
					}
					*st = copyState{at: now, announce: true}
				})
			},
		})
	})
}
