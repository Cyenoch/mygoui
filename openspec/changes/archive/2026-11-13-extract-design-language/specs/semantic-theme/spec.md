## RENAMED Requirements

- FROM: `### Requirement: Elevation follows the coss scale in both appearances`
- TO: `### Requirement: Elevation follows the recorded scale in both appearances`

## MODIFIED Requirements

### Requirement: Elevation follows the recorded scale in both appearances

Popup and dialog elevation SHALL use the large shadow (10/15/-3 and 4/6/-4
DIP, each black at 5% alpha) and tooltips the medium shadow (4/6/-1 and
2/4/-2 DIP at 5%), each followed by a 1px inset edge: black 4% along the
top inside edge in light appearance and white 6% along the bottom inside
edge in dark appearance. Each call SHALL
return an independent copy.

#### Scenario: A popup opens in dark appearance

- **WHEN** a component requests overlay elevation for dark appearance
- **THEN** it receives drop-shadow layers followed by the light inset edge

#### Scenario: A caller edits returned layers

- **WHEN** a caller mutates the returned shadow slice
- **THEN** later requests return the unchanged published elevation
