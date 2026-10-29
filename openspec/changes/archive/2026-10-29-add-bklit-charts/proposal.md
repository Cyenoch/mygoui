# Proposal

## Why

Applications using mygoui cannot yet compose native charts. Bklit's open chart
catalog supplies a concrete capability reference while coss remains the visual
authority for the library and workbench.

## What Changes

- Pin the Bklit source reference and document its MIT/proprietary boundary.
- Add a public native `chart` family covering the 19 documented chart subjects,
  their composable axes, legend, tooltip, brush, indicators and annotations.
- Add discoverable Gallery examples, public usage and native interaction tests.
- Record native adaptations explicitly rather than claiming React API parity.

## Capabilities

### New Capabilities

- `native-charts`: Caller-owned chart data, interaction, composition and themed
  rendering across the Bklit chart and utility catalog.

### Modified Capabilities

None. Existing gallery, ownership and theme contracts continue to apply.

## Impact

New `chart/` package, Gallery consumers, reference submodule and documentation.
Reuse MyGo 0.7.2 chart primitives, standard library and existing coss tokens;
no runtime dependency or application initialization is added.
