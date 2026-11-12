# Semantic Theme

## ADDED Requirements

### Requirement: Elevation follows the coss scale in both appearances

Popup and dialog elevation SHALL use the coss large shadow at 5% alpha and
tooltips the medium shadow at 5%, each over a 1px edge: a dark lower edge in
light appearance and a light upper edge in dark appearance. Each call SHALL
return an independent copy.

#### Scenario: A popup opens in dark appearance

- **WHEN** a component requests overlay elevation for dark appearance
- **THEN** it receives drop-shadow layers followed by the light inset edge

#### Scenario: A caller edits returned layers

- **WHEN** a caller mutates the returned shadow slice
- **THEN** later requests return the unchanged published elevation
