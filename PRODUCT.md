# Product

## Purpose

A Go component library for MyGo native applications. The `coss` development
branch rebuilds the coss ui documentation's component catalog, appearance,
interaction and motion in native Go, without a browser or JavaScript runtime.

## Users And Stack

Go application authors using MyGo v0.6.3. Caller-owned state, typed options and
native focus, editing, selection and accessibility are binding constraints.

## Visual Authority

https://coss.com/ui/docs and the corresponding cosscom/ui source. Previous
HeroUI tokens and component designs are not design constraints for this branch.

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
