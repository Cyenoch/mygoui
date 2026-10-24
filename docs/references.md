# Upstream references

Visual and API references are vendored as Git submodules so upstream changes are
reviewable rather than silently inferred from memory.

| Path | Upstream | Pinned commit | Branch |
| --- | --- | --- | --- |
| `references/heroui` | `https://github.com/heroui-inc/heroui` | `2f03bfbdccbad04167ee3f22df088aef56d2bc96` | `v3` |
| `references/heroui-native` | `https://github.com/heroui-inc/heroui-native` | `dc6833e906f36cb0f11283df0c7fdda2170a6413` | `main` |

The web source owns the design tokens (`packages/styles/themes/default/variables.css`),
component CSS (`packages/styles/components/*.css`), variant defaults
(`packages/styles/src/components/*.styles.ts`) and compound APIs
(`packages/react/src/components`). The native source adds components that exist
only on mobile (bottom sheet, control field, list group, pressable feedback,
link button, menu/sub-menu, skeleton group, dialog).

Clone references with:

```sh
git submodule update --init --depth 1
```

## Updating

```sh
git -C references/heroui fetch --depth 1 origin v3
git -C references/heroui checkout <upstream-commit>
git -C references/heroui-native fetch --depth 1 origin main
git -C references/heroui-native checkout <upstream-commit>
git add references/heroui references/heroui-native
```

Review `git -C references/heroui log --oneline <old>..<new>` and the changelogs
before translating changes. Record the translation impact in `docs/upstream.md`.
Read the relevant reference files for every component you touch; do not infer
current behavior from an older copy.
