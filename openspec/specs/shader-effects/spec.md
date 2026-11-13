# shader-effects Specification

## Purpose
Let an application apply the library's shader effects (an edge mask, a
loading shimmer and a reveal) to its own elements, with tokenised defaults,
a stated fallback and no change to native interaction.

## Requirements

### Requirement: An effect changes paint only

Applying an effect SHALL change only what the element paints. Its layout,
hit testing, focus order, focus rings and accessibility tree MUST stay as
they were, and an effect MUST NOT be the only signal of a state.

#### Scenario: A masked control is used from the keyboard

- **WHEN** a button inside a masked element is focused and activated with
  Enter
- **THEN** the button receives focus and its action runs once

#### Scenario: A shimmering region is read by assistive technology

- **WHEN** an element with an accessible name is given a shimmer
- **THEN** its role and name are reported as before

### Requirement: A mask fades the chosen edges to transparency

The mask effect SHALL fade the element and its native descendants to
transparency over a depth at each selected physical edge and leave the rest
of the element unchanged. No selected edge SHALL mean every edge. The depth
SHALL come from the published fade-edge token unless the call sets one, and
SHALL reach at most half of the element, so that opposite fades do not
overlap. A translucent fallback surface MUST fail explicitly.

#### Scenario: A small element is masked on every edge

- **WHEN** an element shorter than twice the fade depth is masked on every
  edge on a GPU frame
- **THEN** its middle keeps its full opacity

#### Scenario: The bottom edge is masked

- **WHEN** an element is masked on its bottom edge and drawn on a GPU frame
- **THEN** what lies behind it shows through increasingly towards the bottom
  edge, and the top of the element is unchanged

#### Scenario: A theme changes the fade depth

- **WHEN** the published fade-edge token is doubled and a mask sets no depth
- **THEN** the mask fades over the doubled depth

### Requirement: A shimmer runs only while it can and should

The shimmer effect SHALL pass one highlight across the element's subtree, or
pulse it, on the period and strength of the published tokens. It SHALL
request frames only while applied and drawn on
a GPU frame. Under reduced motion, high contrast or a non-positive period it
SHALL draw the element at rest and request no frames.

#### Scenario: Reduced motion is requested

- **WHEN** a shimmer is applied while the reduced-motion preference is on
- **THEN** the element is drawn as it is without the effect and no further
  frame is requested

#### Scenario: Loading ends

- **WHEN** the application stops applying the shimmer
- **THEN** the element is drawn as it is and no further frame is requested

### Requirement: A reveal plays once and rests

The reveal effect SHALL wipe an element in from one edge with a soft front
when the element is first built, over the published reveal duration and
curve. Afterwards the element SHALL be drawn
as it is and request no frames. Under reduced motion the element SHALL fade
in over the same duration without a moving front.

#### Scenario: Content arrives

- **WHEN** an element is first built with a downward reveal on a GPU frame
- **THEN** midway its top shows and its bottom does not, and after the
  duration all of it shows and no frame is pending

#### Scenario: Reduced motion is requested

- **WHEN** an element is first built with a reveal under reduced motion
- **THEN** midway the whole element shows at partial strength, with no edge
  further along than another

### Requirement: Every effect has a stated fallback

Where the frame does not draw on a GPU, where the element exceeds the
toolkit's effect-layer limit, and where a mask or reveal call asks for it,
each effect SHALL draw its stated fallback instead of running its shader: a
mask shades its selected edges or blends them into the opaque surface the
call names, a shimmer rests, a reveal fades in, or shows at once above the
layer limit. The content MUST stay readable in every fallback.

#### Scenario: A mask renders in software

- **WHEN** a masked element is drawn without a GPU
- **THEN** its selected edges are shaded and its content is otherwise
  unchanged

#### Scenario: A mask exceeds the layer limit

- **WHEN** a masked element wider than the effect-layer limit is drawn on a
  GPU frame
- **THEN** its selected edges are shaded and its content is otherwise
  unchanged

#### Scenario: A reveal renders in software

- **WHEN** an element is first built with a reveal and drawn without a GPU
- **THEN** it fades in over the reveal duration and then requests no frames

### Requirement: Effects cost nothing until used

Importing the effect family MUST NOT compile a shader or do other work. A
shader SHALL be compiled the first time it paints, not when a view is built
or a fallback is drawn, and shared by every element and
window afterwards.

#### Scenario: A program imports the family without applying an effect

- **WHEN** the program starts and builds views without effects
- **THEN** no shader has been compiled

### Requirement: One effect per element is stated

An element SHALL carry one effect: applying an effect replaces the one the
element had, whether the application or a component installed it. The mask
additionally owns the element's over-paint callback for its fallback.

#### Scenario: Two effects are applied to one element

- **WHEN** an application applies a shimmer and then a mask to one element
- **THEN** the element is drawn through the mask only
