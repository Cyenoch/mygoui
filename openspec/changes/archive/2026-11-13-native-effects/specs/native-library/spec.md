# Spec Delta

## MODIFIED Requirements

### Requirement: Effects degrade without a GPU

Where the frame does not draw on a GPU, SkeletonGroup SHALL rest static
without requesting frames, and a mask ScrollFade SHALL show the shadow
treatment inside the viewport instead of running its shader. A mask
ScrollFade whose viewport exceeds the toolkit's effect-layer limit SHALL
show the same shadow treatment on a GPU frame.

#### Scenario: A mask fade renders in software

- **WHEN** a viewport with a mask fade has more to scroll below
- **THEN** its bottom edge is shaded and its content is otherwise unchanged

#### Scenario: A mask fade exceeds the layer limit

- **WHEN** a viewport with a mask fade is wider than the effect-layer limit
  and has more to scroll below
- **THEN** its bottom edge is shaded and its content is otherwise unchanged
