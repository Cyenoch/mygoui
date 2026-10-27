package mygoui

// Side is a placement edge. SideAuto selects the component's default edge.
type Side uint8

const (
	SideAuto Side = iota
	SideTop
	SideBottom
	SideLeft
	SideRight
)

// Align is placement alignment along the chosen edge.
type Align uint8

const (
	AlignCenter Align = iota
	AlignStart
	AlignEnd
)

// Placement anchors a floating panel. Offset is the gap from the anchor; zero
// uses the component's default gap, typically 4 DIPs.
type Placement struct {
	Side   Side
	Align  Align
	Offset float32
}
