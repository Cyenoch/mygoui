package form

import (
	"testing"

	"github.com/Cyenoch/mygoui/list"
)

func TestPendingFilterUsesCurrentQueryAndSkipsDisabledOptions(t *testing.T) {
	items := []list.ListBoxItem{
		{Key: "a", Label: "Alpha"},
		{Key: "disabled", Label: "Bee", Disabled: true},
	}
	sections := []list.ListBoxSection{{
		Label: "Group",
		Items: []list.ListBoxItem{{Key: "b", Label: "Beta"}, {Key: "c", Label: "Bravo"}},
	}}
	var pending pendingFilter
	pending.reset("a", nil, items, sections)
	if pending.key() != "a" {
		t.Fatalf("initial query matched %q", pending.key())
	}
	pending.reset("b", nil, items, sections)
	if pending.key() != "b" {
		t.Fatalf("updated query matched %q", pending.key())
	}
	pending.move(1)
	if pending.key() != "c" {
		t.Fatalf("next matching key = %q", pending.key())
	}
	pending.reset("missing", nil, items, sections)
	if pending.key() != "" {
		t.Fatalf("unmatched query selected %q", pending.key())
	}
}
