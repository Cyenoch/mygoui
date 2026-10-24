package data

import (
	"slices"

	"github.com/egoist/mygo/ui"
)

// keyRow returns caller row index's identity.
func (s *TableState) keyRow(rowKey func(int) any, row int) any {
	if rowKey != nil {
		return rowKey(row)
	}
	return row
}

// syncSort reconciles TableState.Sort with the native sort the user changed,
// and reports whether the user changed it.
func (s *TableState) syncSort() bool {
	native := TableSort{Column: s.nativeSort.Column, Descending: s.nativeSort.Descending}
	switch {
	case s.Sort != s.lastSort:
		// The app changed the sort, or set it for the first time.
		s.nativeSort = ui.SortOrder{Column: s.Sort.Column, Descending: s.Sort.Descending}
		s.lastSort = s.Sort
		return false
	case native != s.lastSort:
		// A header click changed the native sort.
		s.Sort = native
		s.lastSort = native
		return true
	}
	return false
}

// syncDisabled rebuilds the disabled rows when the row count, the revision or
// the presence of RowDisabled changed.
func (s *TableState) syncDisabled(o *TableOptions) {
	has := o.RowDisabled != nil
	if s.disabledDone && s.disabledHas == has && s.disabledRows == o.Rows && s.disabledRev == s.Revision {
		return
	}
	s.disabled = s.disabled[:0]
	if has {
		if s.disabledSet == nil {
			s.disabledSet = map[int]struct{}{}
		} else {
			clear(s.disabledSet)
		}
		for i := range o.Rows {
			if o.RowDisabled(i) {
				s.disabled = append(s.disabled, i)
				s.disabledSet[i] = struct{}{}
			}
		}
	} else if s.disabledSet != nil {
		clear(s.disabledSet)
	}
	s.disabledRows = o.Rows
	s.disabledRev = s.Revision
	s.disabledHas = has
	s.disabledDone = true
}

func (s *TableState) isDisabled(row int) bool {
	if s.disabledSet == nil {
		return false
	}
	_, ok := s.disabledSet[row]
	return ok
}

// syncOrder rebuilds the displayed row order when rows, the revision or the
// sort changed. A single selection follows its row through the reordering.
func (s *TableState) syncOrder(o *TableOptions) {
	dirty := !s.ordered || s.orderRows != o.Rows || s.orderRev != s.Revision || s.orderSort != s.Sort
	if !dirty {
		return
	}
	under := -1
	if s.Selection.Row >= 0 && s.Selection.Row < len(s.order) {
		under = s.order[s.Selection.Row]
	}
	if cap(s.order) < o.Rows {
		s.order = make([]int, o.Rows)
	} else {
		s.order = s.order[:o.Rows]
	}
	for i := range s.order {
		s.order[i] = i
	}
	if o.Compare != nil && s.Sort.Column != "" {
		descending := s.Sort.Descending
		slices.SortStableFunc(s.order, func(a, b int) int {
			ra := TableRow{Index: a, Key: s.keyRow(o.RowKey, a), Disabled: s.isDisabled(a)}
			rb := TableRow{Index: b, Key: s.keyRow(o.RowKey, b), Disabled: s.isDisabled(b)}
			cmp := o.Compare(ra, rb)
			if descending {
				return -cmp
			}
			return cmp
		})
	}
	s.orderRows = o.Rows
	s.orderRev = s.Revision
	s.orderSort = s.Sort
	s.ordered = true
	if under >= 0 {
		s.Selection.Row = s.position(under)
	}
}

// position returns the displayed position of caller row under, or -1.
func (s *TableState) position(under int) int {
	for i, row := range s.order {
		if row == under {
			return i
		}
	}
	return -1
}

// underlying returns the caller row index at displayed position pos, or -1.
func (s *TableState) underlying(pos int) int {
	if pos < 0 || pos >= len(s.order) {
		return -1
	}
	return s.order[pos]
}

// rowAt returns the metadata of the row at displayed position pos.
func (s *TableState) rowAt(o *TableOptions, pos int) TableRow {
	under := s.underlying(pos)
	if under < 0 {
		return TableRow{Index: -1}
	}
	key := s.keyRow(o.RowKey, under)
	selected := false
	switch s.Selection.Mode {
	case TableSelectionSingle:
		selected = s.Selection.Row == pos
	case TableSelectionMultiple:
		selected = s.Selection.Rows.Has(key)
	}
	return TableRow{Index: under, Key: key, Selected: selected, Disabled: s.isDisabled(under)}
}

// wire points the native table state at this table's selection, sort and keys.
func (s *TableState) wire(o *TableOptions) {
	s.List.Selected = nil
	s.List.Selection = nil
	switch s.Selection.Mode {
	case TableSelectionSingle:
		s.List.Selected = &s.Selection.Row
	case TableSelectionMultiple:
		s.List.Selected = &s.Selection.Row
		s.List.Selection = &s.Selection.Rows
	}
	s.List.Sort = &s.nativeSort
	s.List.Key = func(pos int) any {
		return s.keyRow(o.RowKey, s.underlying(pos))
	}
	if o.RowMenu != nil {
		s.List.CellMenu = func(m *ui.Menu, row, _ int) {
			o.RowMenu(m, s.rowAt(o, row))
		}
	} else {
		s.List.CellMenu = nil
	}
}

// afterBuild corrects a selection that landed on a disabled row and reports
// whether the user changed the selection. nativeChanged is the toolkit's
// change flag; it is also set by a sort click, which sortChanged identifies.
// It updates the restore snapshot.
func (s *TableState) afterBuild(o *TableOptions, nativeChanged, sortChanged bool) bool {
	if !s.snapshotOK {
		s.takeSnapshot()
		s.leadUnder = s.underlying(s.Selection.Row)
		s.rememberLead(o, s.leadUnder)
	}
	restored := false
	if pos := s.Selection.Row; pos >= 0 && pos < len(s.order) && s.isDisabled(s.order[pos]) {
		switch s.Selection.Mode {
		case TableSelectionMultiple:
			s.restoreSnapshot()
		default:
			s.Selection.Row = s.savedPosition()
		}
		restored = true
	} else if s.Selection.Mode == TableSelectionMultiple {
		for _, under := range s.disabled {
			key := s.keyRow(o.RowKey, under)
			s.Selection.Rows.Remove(key)
		}
	}

	under := s.underlying(s.Selection.Row)
	changed := false
	if !restored {
		// With a stable identity, the chosen row is its key: inserting or
		// sorting rows changes its position, not the selection, so no
		// callback reports it. Without one, the position is the identity.
		if o.RowKey != nil {
			changed = s.leadKeyDiffers(o, under)
		} else {
			changed = under != s.leadUnder
		}
		if s.Selection.Mode == TableSelectionMultiple {
			changed = changed || s.Selection.Rows.Len() != s.snapshotLen
		}
		// A multiple selection can change to another of the same size, as a
		// Shift range moving, which neither the lead row nor the count
		// shows. Compare the chosen keys then. A sort click sets the
		// toolkit's change flag too, but sortChanged identifies it, and its
		// intermediate pass changes no keys.
		if !changed && nativeChanged && !sortChanged && s.Selection.Mode == TableSelectionMultiple {
			changed = s.selectionSetChanged()
		}
	}
	s.leadUnder = under
	s.rememberLead(o, under)
	if changed {
		s.takeSnapshot()
	}
	return changed
}

// rememberLead records the identity of the row the selection is on.
func (s *TableState) rememberLead(o *TableOptions, under int) {
	s.leadHas = under >= 0
	s.leadKey = nil
	if s.leadHas {
		s.leadKey = s.keyRow(o.RowKey, under)
	}
}

// leadKeyDiffers reports whether the row of the last frame's selection is not
// the row of under. RowKey must not be nil.
func (s *TableState) leadKeyDiffers(o *TableOptions, under int) bool {
	if under < 0 {
		return s.leadHas
	}
	if !s.leadHas {
		return true
	}
	return s.keyRow(o.RowKey, under) != s.leadKey
}

// savedPosition returns the displayed position of the row the selection is
// restored to, or -1 when it is gone or disabled.
func (s *TableState) savedPosition() int {
	if s.leadUnder < 0 || s.isDisabled(s.leadUnder) {
		return -1
	}
	return s.position(s.leadUnder)
}

func (s *TableState) restoreSnapshot() {
	s.Selection.Rows.Clear()
	for _, key := range s.snapshot {
		s.Selection.Rows.Add(key)
	}
	s.Selection.Row = s.savedPosition()
}

func (s *TableState) takeSnapshot() {
	s.snapshot = s.snapshot[:0]
	s.snapshotLen = 0
	s.snapshotOK = true
	if s.Selection.Mode != TableSelectionMultiple {
		if s.snapshotSet != nil {
			clear(s.snapshotSet)
		}
		return
	}
	if s.snapshotSet == nil {
		s.snapshotSet = make(map[any]struct{}, s.Selection.Rows.Len())
	} else {
		clear(s.snapshotSet)
	}
	for key := range s.Selection.Rows.All() {
		s.snapshot = append(s.snapshot, key)
		s.snapshotSet[key] = struct{}{}
	}
	s.snapshotLen = s.Selection.Rows.Len()
}

// selectionSetChanged reports whether a multiple selection no longer holds the
// keys of the last snapshot, at the same size.
func (s *TableState) selectionSetChanged() bool {
	if s.Selection.Rows.Len() != len(s.snapshotSet) {
		return true
	}
	for key := range s.Selection.Rows.All() {
		if _, ok := s.snapshotSet[key]; !ok {
			return true
		}
	}
	return false
}

// activeRow returns the row Enter or a double click activated.
func (s *TableState) activeRow(o *TableOptions) (TableRow, bool) {
	pos := s.Selection.Row
	if pos < 0 || pos >= len(s.order) {
		return TableRow{}, false
	}
	under := s.order[pos]
	if s.isDisabled(under) {
		return TableRow{}, false
	}
	return s.rowAt(o, pos), true
}
