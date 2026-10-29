# Product

<!-- impeccable:product-schema 1 -->

## Platform

adaptive

Native desktop applications on macOS, Windows and Linux, using MyGo's
platform window controls, menus, focus and services. This selector enables
Impeccable's native guidance; it does not add iOS or Android targets.
Desktop MyGo behavior takes precedence over phone-specific platform examples.
The WebAssembly build is a compatibility check, not a separate web product.

## Purpose

A Go component library for MyGo native applications. The `coss` development
branch rebuilds the coss ui documentation's component catalog, appearance,
interaction and motion in native Go, without a browser or JavaScript runtime.

## Users And Stack

Go application authors using MyGo v0.7.2. Caller-owned state, typed options and
native focus, editing, selection and accessibility are binding constraints.

## Visual Authority

https://coss.com/ui/docs and the corresponding cosscom/ui source. Previous
tokens and component designs are not design constraints for this branch.

## Capabilities

Customizable semantic theme tokens in both appearances; independently composed
components; a generic scroll fade using shadow, surface blending or a native
subtree alpha mask. Web-only
or unavailable toolkit features must be documented rather than simulated as
fully supported. GPU, OS gestures and screen-reader results require separate
evidence from headless testing.

## Native Gallery

The gallery is a native component workbench with a completely redesigned
application shell, custom titlebar and sidebar. Component discovery, live
examples, Go source and theme controls form one coherent workflow. Use real
platform window controls and draggable titlebar regions, preserve keyboard
navigation and native focus, and adapt the workspace to compact windows.
The coss component appearance remains the visual authority.

## Development Context

OpenSpec records current requirements and scoped change proposals.
Impeccable supplies design discovery, critique, audit, refinement and finish
playbooks. Existing Go tokens and native controls remain the implementation
authority; root `DESIGN.md` records that system. Browser-only live iteration
and HTML/CSS detector results do not establish native UI quality.
