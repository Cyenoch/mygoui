package chart

import (
	"fmt"
	"math"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
	native "github.com/egoist/mygo/ui/chart"
)

// HierarchyNode uses leaf Value; branch values are the sum of descendants.
// IDs must be globally unique. A node may have only one parent.
type HierarchyNode struct {
	ID, Label string
	Value     float64
	Children  []*HierarchyNode
}

// SunburstState retains focus by ID across builds and changing input trees.
type SunburstState struct {
	State
	FocusID string
}

type SunburstOptions struct {
	Options
	Root  *HierarchyNode
	Drill *SunburstState
}

type sunSegment struct {
	node  *HierarchyNode
	depth int
	slice native.Slice
	value float64
}

func hierarchy(root *HierarchyNode) (map[string]*HierarchyNode, map[string]*HierarchyNode, map[string]float64, error) {
	nodes, parents, sums := map[string]*HierarchyNode{}, map[string]*HierarchyNode{}, map[string]float64{}
	var visit func(*HierarchyNode, *HierarchyNode) (float64, error)
	visit = func(node, parent *HierarchyNode) (float64, error) {
		if node == nil || node.ID == "" || nodes[node.ID] != nil {
			return 0, fmt.Errorf("sunburst: nodes require unique IDs and an acyclic single-parent hierarchy")
		}
		if !finite(node.Value) || node.Value < 0 {
			return 0, fmt.Errorf("sunburst: node %q has an invalid value", node.ID)
		}
		nodes[node.ID], parents[node.ID] = node, parent
		total := node.Value
		if len(node.Children) > 0 {
			total = 0
			for _, child := range node.Children {
				value, err := visit(child, node)
				if err != nil {
					return 0, err
				}
				total += value
			}
		}
		if !finite(total) {
			return 0, fmt.Errorf("sunburst: node %q overflows its descendant sum", node.ID)
		}
		sums[node.ID] = total
		return total, nil
	}
	if root != nil {
		if _, err := visit(root, nil); err != nil {
			return nil, nil, nil, err
		}
	}
	return nodes, parents, sums, nil
}

func Sunburst(c *ui.Context, o SunburstOptions) (ui.Element, error) {
	if err := checkOptions(o.Options); err != nil {
		return ui.Element{}, err
	}
	nodes, parents, sums, err := hierarchy(o.Root)
	if err != nil {
		return ui.Element{}, err
	}
	focus := o.Root
	if o.Drill != nil {
		o.State = &o.Drill.State
		if o.Drill.FocusID != "" {
			focus = nodes[o.Drill.FocusID]
			if focus == nil {
				return ui.Element{}, fmt.Errorf("sunburst: focus %q is not in the hierarchy", o.Drill.FocusID)
			}
		}
	}
	segments := []sunSegment{}
	maxDepth := 1
	var partition func(*HierarchyNode, int, float64, float64)
	partition = func(node *HierarchyNode, depth int, start, end float64) {
		total := sums[node.ID]
		angle := start
		for _, child := range node.Children {
			span := 0.0
			if total > 0 {
				span = (sums[child.ID] / total) * (end - start)
			}
			segments = append(segments, sunSegment{node: child, depth: depth, slice: native.Slice{Start: angle, End: angle + span, Pad: .015}, value: sums[child.ID]})
			maxDepth = max(maxDepth, depth+1)
			partition(child, depth+1, angle, angle+span)
			angle += span
		}
	}
	if focus != nil && sums[focus.ID] > 0 {
		if len(focus.Children) == 0 {
			segments = append(segments, sunSegment{node: focus, slice: native.Slice{Start: 0, End: 2 * math.Pi}, value: sums[focus.ID]})
		} else {
			partition(focus, 0, 0, 2*math.Pi)
		}
	}
	activate := o.OnActivate
	o.OnActivate = func(i int) {
		if o.Drill != nil && len(segments[i].node.Children) > 0 {
			o.Drill.FocusID = segments[i].node.ID
			o.Drill.Selected = -1
		}
		if activate != nil {
			activate(i)
		}
	}
	t := theme.Colors(c)
	var totalLabel *ui.TextLayout
	if focus != nil {
		totalLabel = ui.ShapeText(number(sums[focus.ID]), ui.Font{Family: theme.FontFamily(c), Size: theme.Scale(c, 12), Features: "tnum"}, 0)
	}
	return ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space8).Children(func() {
		if focus != nil && o.Drill != nil {
			ancestors := []*HierarchyNode{}
			for node := focus; node != nil; node = parents[node.ID] {
				ancestors = append(ancestors, node)
			}
			ui.Row(c).FillWidth().Wrap().Gap(theme.Space4).Children(func() {
				for i := len(ancestors) - 1; i >= 0; i-- {
					node := ancestors[i]
					mygoui.Button(c.Key(node.ID), mygoui.ButtonOptions{Label: node.Label, Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall, Disabled: o.Disabled,
						OnClick: func() { o.Drill.FocusID = node.ID; o.Drill.Selected = -1 }})
				}
			})
		}
		surface(c, o.Options, len(segments), func(i int) string { return segments[i].node.Label + ": " + number(segments[i].value) },
			func(p *ui.Painter, r ui.Rect, selected int) {
				cx, cy, radius := circle(r)
				step := radius * .8 / float32(maxDepth)
				for i, segment := range segments {
					inner := radius*.2 + float32(segment.depth)*step
					fill := SeriesColor(c, i)
					if selected >= 0 && selected != i {
						fill = fill.Alpha(.5)
					}
					p.FillPath(native.Sector(cx, cy, inner, inner+step*.95, segment.slice), fill)
				}
				if totalLabel != nil {
					w, h := totalLabel.Size()
					if w <= radius*.36 && h <= radius*.36 {
						p.TextLayout(totalLabel, cx-w/2, cy-h/2, t.Foreground)
					}
				}
			}, func(x, y float32) int {
				r := o.State.Handle.Bounds(c)
				radius := min(r.W, r.H)/2 - 8
				dx, dy := (x-.5)*r.W, (y-.5)*r.H
				distance := float32(math.Hypot(float64(dx), float64(dy)))
				angle := math.Atan2(float64(dx), -float64(dy))
				if angle < 0 {
					angle += 2 * math.Pi
				}
				step := radius * .8 / float32(maxDepth)
				for i, segment := range segments {
					inner := radius*.2 + float32(segment.depth)*step
					if distance >= inner && distance <= inner+step*.95 && angle >= segment.slice.Start && angle < segment.slice.End {
						return i
					}
				}
				return -1
			})
	}), nil
}
