## Context

coss expresses motion as CSS transitions on its registry components. MyGo
supplies element transitions, keyed animations, poses and paint-time frames.
The review compared every component that changes state with its coss source
under `references/coss/apps/ui/registry/default/ui`.

## Audit

| Component | coss | Native before | Decision |
| --- | --- | --- | --- |
| Popover, pickers | 98% fade in, fade out | same | keep |
| Tooltip | 98% fade in and out | same | keep |
| PreviewCard | 98% fade in and out | fade out only | shrink on exit |
| Menu, Select, ComboBox popups | none | none | keep |
| Select indicator | static up-down chevron | turning chevron | static icon |
| Switch | thumb slide 150 ms, track colour, 10% press stretch after 100 ms | slide and colour | add stretch |
| Slider thumb | 120% press scale | same | keep |
| Tabs indicator | 200 ms ease-in-out | 200 ms ease-out | ease-in-out |
| SegmentedControl | (tabs) 200 ms ease-in-out | same | keep |
| Accordion, Disclosure | 200 ms height, chevron turn | same | keep |
| Dialog, Sheet | 200 ms ease-in-out fade/translate | same | keep |
| Drawer | 450 ms drawer curve | same | keep |
| BottomSheet | drawer curve | 450 ms ease-out | drawer curve |
| Sidebar | 200 ms linear width | none | add |
| Progress, Meter | 500 ms value | same | keep |
| Indeterminate progress | — | full refresh rate | cap at 60 fps |
| Checkbox, Radio, Toggle | shadow/colour only | colour | keep |
| Toast, Calendar, ColorSwatch, Avatar, Skeleton | library additions | present | keep |

## Decisions

- Switch stretch is painted: the thumb is drawn in the track's paint
  callback, and MyGo poses scale uniformly. A 0–1 press value animates over
  100 ms after a 100 ms delay; the thumb widens by 10% of its width from the
  edge it rests against, which follows the slide while it moves.
- Sidebar animates a 0–1 collapse value and derives its width from it, so a
  text-scale change resizes at once; it clips only while moving, so labels
  shown on expand stay inside the growing panel.
- Reduced motion: the existing contract removes movement and fades alike.
  Keeping colour fades under reduced motion was considered; it changes a
  documented contract across many components and is left for a separate
  change if wanted.
- No public motion tokens: durations stay with the component that owns
  them, as coss writes them per component.
