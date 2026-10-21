package conversation

import (
	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	locales "github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

type MessageID string
type ScopeID string

// ScrollerState belongs to one native engine. Stable caller keys let MyGo
// preserve the reading anchor when history is prepended or messages grow.
type ScrollerState struct {
	List        ui.ListState
	scope       ScopeID
	initialized bool
	jumpKey     MessageID
	jumpAlign   ui.Align
	jumpEnd     bool
	anchor      MessageID
}

type ScrollerOptions struct {
	Scope                       ScopeID
	Label                       string
	Count                       int
	Key                         func(int) MessageID
	RowLabel                    func(int) string
	Row                         func(*ui.Context, int)
	Empty                       func(*ui.Context)
	DisableFollow               bool
	Smooth                      bool
	MoreHistory, LoadingHistory bool
	OnLoadOlder                 func()
	// Anchor names a new turn that should start near the top instead of at
	// the document bottom. It is consumed once per changed nonempty identity.
	Anchor MessageID
}

func (s *ScrollerState) JumpTo(key MessageID, align ui.Align) {
	s.jumpKey, s.jumpAlign, s.jumpEnd = key, align, false
}
func (s *ScrollerState) JumpToLatest() { s.jumpEnd = true; s.jumpKey = "" }

func Scroller(c *ui.Context, s *ScrollerState, o ScrollerOptions) ui.Element {
	if !s.initialized || s.scope != o.Scope {
		s.List = ui.ListState{}
		s.initialized, s.scope = true, o.Scope
		s.jumpKey = ""
		s.jumpEnd = !o.DisableFollow
		s.anchor = ""
	}
	if o.Anchor != "" && o.Anchor != s.anchor {
		s.anchor = o.Anchor
		s.jumpKey = o.Anchor
		s.jumpAlign = ui.Start
		s.jumpEnd = false
	}
	s.List.FollowEnd = !o.DisableFollow
	s.List.FollowAnimation = nil
	if o.Smooth && !theme.ReduceMotion(c) {
		s.List.FollowAnimation = &ui.Smooth
	}
	s.List.Key = func(i int) any {
		if o.Key != nil {
			return o.Key(i)
		}
		return i
	}
	s.List.Label = func(i int) string {
		if o.RowLabel != nil {
			return o.RowLabel(i)
		}
		return ""
	}
	container := ui.Column(c).Fill().MinWidth(0).MinHeight(0)
	container.Children(func() {
		if o.MoreHistory {
			label := locales.ComponentLoadEarlierMessages()
			if o.LoadingHistory {
				label = locales.ComponentLoadingEarlierMessages()
			}
			mygoui.Button(c.Key("older-messages"), mygoui.ButtonOptions{Label: label, Variant: mygoui.ButtonGhost, Loading: o.LoadingHistory, OnClick: o.OnLoadOlder})
		}
		list := ui.List(c.Key("message-list"), &s.List, o.Count).Grow(1).FillWidth().MinWidth(0).MinHeight(0).Gap(16).Padding(16)
		label := o.Label
		if label == "" {
			label = locales.ComponentMessages()
		}
		list.Label(label)
		list.Rows(func(row ui.ListRow) {
			if o.Row != nil {
				o.Row(row.Context, row.Index)
			}
		})
		if o.Count == 0 && o.Empty != nil {
			list.Children(func() { o.Empty(c) })
		}
		if s.jumpEnd {
			s.List.ScrollToEnd()
			s.jumpEnd = false
		}
		target := s.jumpKey
		if target != "" && o.Key != nil {
			for i := range o.Count {
				if o.Key(i) == target {
					s.List.ScrollTo(i, s.jumpAlign)
					s.jumpKey = ""
					break
				}
			}
		}
		if !s.List.AtEnd() && o.Count > 0 {
			mygoui.Button(c.Key("jump-latest"), mygoui.ButtonOptions{Label: locales.ComponentJumpLatestMessage(), Icon: lucide.ChevronDown(), Variant: mygoui.ButtonSecondary, Size: mygoui.ButtonSM, OnClick: func() { s.List.ScrollToEnd() }}).AlignSelf(ui.End).Margin(4, 12)
		}
	})
	return container
}
