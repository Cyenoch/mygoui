package chart

import (
	"fmt"
	"math"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

type FlowNode struct{ ID, Label string }

type FlowLink struct {
	Source, Target string
	Value          float64
}

type SankeyOptions struct {
	Options
	Nodes []FlowNode
	Links []FlowLink
}

type flowBox struct{ x, y, width, height float32 }
type flowRibbon struct{ x1, y1, x2, y2, height float32 }

func layoutFlow(o SankeyOptions) ([]flowBox, []flowRibbon, error) {
	n := len(o.Nodes)
	ids := map[string]int{}
	for i, node := range o.Nodes {
		if _, exists := ids[node.ID]; node.ID == "" || exists {
			return nil, nil, fmt.Errorf("sankey: node %d needs a unique nonempty ID", i)
		}
		ids[node.ID] = i
	}
	indegree, depth := make([]int, n), make([]int, n)
	in, out := make([]float64, n), make([]float64, n)
	edges := make([][]int, n)
	for i, link := range o.Links {
		source, sourceOK := ids[link.Source]
		target, targetOK := ids[link.Target]
		if !sourceOK || !targetOK || source == target || !finite(link.Value) || link.Value <= 0 {
			return nil, nil, fmt.Errorf("sankey: link %d requires distinct existing nodes and positive finite flow", i)
		}
		indegree[target]++
		edges[source] = append(edges[source], i)
		out[source] += link.Value
		in[target] += link.Value
		if !finite(out[source]) || !finite(in[target]) {
			return nil, nil, fmt.Errorf("sankey: aggregate flow overflows at link %d", i)
		}
	}
	queue := []int{}
	for i, d := range indegree {
		if d == 0 {
			queue = append(queue, i)
		}
	}
	for pos := 0; pos < len(queue); pos++ {
		for _, edge := range edges[queue[pos]] {
			target := ids[o.Links[edge].Target]
			depth[target] = max(depth[target], depth[queue[pos]]+1)
			indegree[target]--
			if indegree[target] == 0 {
				queue = append(queue, target)
			}
		}
	}
	if len(queue) != n {
		return nil, nil, fmt.Errorf("sankey: flow graph contains a cycle")
	}
	maxDepth := 0
	for _, d := range depth {
		maxDepth = max(maxDepth, d)
	}
	columns := make([][]int, maxDepth+1)
	for i, d := range depth {
		columns[d] = append(columns[d], i)
	}
	scale := math.Inf(1)
	weights := make([]float64, n)
	for _, column := range columns {
		total := 0.0
		for _, i := range column {
			weights[i] = max(in[i], out[i])
			total += weights[i]
			if !finite(total) {
				return nil, nil, fmt.Errorf("sankey: column flow overflows")
			}
		}
		if total > 0 {
			scale = min(scale, .9/total)
		}
	}
	if math.IsInf(scale, 1) {
		scale = 0
	}
	boxes := make([]flowBox, n)
	nodeWidth := min(float32(.06), .45/float32(max(1, maxDepth)))
	for d, column := range columns {
		total := float32(0)
		for _, i := range column {
			total += float32(weights[i] * scale)
		}
		gap := (1 - total) / float32(len(column)+1)
		y := gap
		for _, i := range column {
			h := float32(weights[i] * scale)
			boxes[i] = flowBox{x: float32(d) / float32(max(1, maxDepth)) * (1 - nodeWidth), y: y, width: nodeWidth, height: h}
			y += h + gap
		}
	}
	ribbons := make([]flowRibbon, len(o.Links))
	sourceOffsets, targetOffsets := make([]float32, n), make([]float32, n)
	for i, link := range o.Links {
		s, target := ids[link.Source], ids[link.Target]
		a, b := boxes[s], boxes[target]
		h := float32(link.Value * scale)
		ribbons[i] = flowRibbon{x1: a.x + a.width, y1: a.y + sourceOffsets[s], x2: b.x, y2: b.y + targetOffsets[target], height: h}
		sourceOffsets[s] += h
		targetOffsets[target] += h
	}
	return boxes, ribbons, nil
}

func Sankey(c *ui.Context, o SankeyOptions) (ui.Element, error) {
	if err := checkOptions(o.Options); err != nil {
		return ui.Element{}, err
	}
	boxes, ribbons, err := layoutFlow(o)
	if err != nil {
		return ui.Element{}, err
	}
	t := theme.Colors(c)
	fontSize := theme.Scale(c, 12)
	describe := func(i int) string {
		if i < len(o.Nodes) {
			return o.Nodes[i].Label
		}
		v := o.Links[i-len(o.Nodes)]
		return v.Source + " → " + v.Target + ": " + number(v.Value)
	}
	return surface(c, o.Options, len(o.Nodes)+len(o.Links), describe, func(p *ui.Painter, r ui.Rect, selected int) {
		for i, ribbon := range ribbons {
			x1, x2 := r.X+ribbon.x1*r.W, r.X+ribbon.x2*r.W
			y1, y2 := r.Y+ribbon.y1*r.H, r.Y+ribbon.y2*r.H
			h, middle := ribbon.height*r.H, (x1+x2)/2
			path := (&ui.Path{}).MoveTo(x1, y1).CubeTo(middle, y1, middle, y2, x2, y2).
				LineTo(x2, y2+h).CubeTo(middle, y2+h, middle, y1+h, x1, y1+h).Close()
			fill := SeriesColor(c, i).Alpha(.3)
			if selected == len(o.Nodes)+i || (selected >= 0 && selected < len(o.Nodes) && (o.Nodes[selected].ID == o.Links[i].Source || o.Nodes[selected].ID == o.Links[i].Target)) {
				fill = SeriesColor(c, i).Alpha(.65)
			}
			p.FillPath(path, fill)
		}
		for i, box := range boxes {
			rb := ui.Rect{X: r.X + box.x*r.W, Y: r.Y + box.y*r.H, W: box.width * r.W, H: max(1, box.height*r.H)}
			p.Fill(rb, SeriesColor(c, i), 2)
			p.Text(max(r.X, min(r.X+r.W-65, rb.X+rb.W+4)), rb.Y, o.Nodes[i].Label, fontSize, t.Foreground)
		}
	}, func(x, y float32) int {
		for i, b := range boxes {
			if x >= b.x && x <= b.x+b.width && y >= b.y && y <= b.y+b.height {
				return i
			}
		}
		for i, ribbon := range ribbons {
			if x < ribbon.x1 || x > ribbon.x2 {
				continue
			}
			// Solve the monotone Bezier x to use the same ribbon as the paint.
			lo, hi := float32(0), float32(1)
			midX := (ribbon.x1 + ribbon.x2) / 2
			for range 16 {
				t := (lo + hi) / 2
				bx := cubic(ribbon.x1, midX, midX, ribbon.x2, t)
				if bx < x {
					lo = t
				} else {
					hi = t
				}
			}
			top := cubic(ribbon.y1, ribbon.y1, ribbon.y2, ribbon.y2, (lo+hi)/2)
			if y >= top && y <= top+ribbon.height {
				return len(o.Nodes) + i
			}
		}
		return -1
	}), nil
}

func cubic(a, b, c, d, t float32) float32 {
	u := 1 - t
	return u*u*u*a + 3*u*u*t*b + 3*u*t*t*c + t*t*t*d
}
