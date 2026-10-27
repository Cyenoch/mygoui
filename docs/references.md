# Reference provenance

## Current visual reference

[coss ui documentation](https://coss.com/ui/docs) and the MIT `apps/ui`
registry at commit **`3d8ca28e7d90123e9e0f44665944dc240b026ccb`** of
[`cosscom/ui`](https://github.com/cosscom/ui/tree/3d8ca28e7d90123e9e0f44665944dc240b026ccb)
are the reconstruction reference. Previous HeroUI design and tokens are not
constraints for this branch.

The research checkout is `/tmp/mygoui-coss-upstream`. It is temporary and
not part of the public module or a prerequisite for consuming it. The live
site was also downloaded to `/tmp/mygoui-coss-docs.html` and its LLM catalog
to `/tmp/mygoui-coss-llms.txt`. Those local snapshots establish the catalog
observed during reconstruction, not a promise about future website changes.

| Purpose | Path in the pinned upstream checkout |
| --- | --- |
| Directory license declaration | `LICENSING.md`, corroborated by root `README.md` |
| UI license declaration | `apps/ui/README.md` and `apps/ui/package.json` |
| Component catalog order | `apps/ui/content/docs/components/meta.json` |
| Component behavior, anatomy and examples | `apps/ui/content/docs/components/*.mdx` |
| Registry component source | `apps/ui/registry/default/ui/*.tsx` |
| Registry inventory and dependencies | `apps/ui/registry.json` |
| Live example source | `apps/ui/registry/default/` entries referenced by the MDX previews |
| Published color/type guidance | `apps/ui/content/docs/(root)/styling.mdx` |
| Hook documentation | `apps/ui/content/docs/hooks/use-media-query.mdx`, `use-copy-to-clipboard.mdx` |
| Chat composition/application boundaries | `apps/ui/content/docs/(root)/chat.mdx` |

The [component matrix](components.md) uses the union of the downloaded site
links, `llms.txt` and component metadata: 57 component subjects plus two
hooks. The site/metadata have Composer and Segmented Control but omit Sidebar;
`llms.txt` includes Sidebar and omits the other two. Sidebar also exists at
`apps/ui/registry/default/ui/sidebar.tsx`. Date Picker is a documented
composition without a standalone registry UI file. Do not use one stale
catalog alone as a completeness check.

Chat UI, particles, installation, skills, roadmap and migration pages add
examples or application concerns. They do not turn the native library into
a chat backend, browser hook runtime or website installer.

## Mixed licensing boundary

At this commit, upstream's
[LICENSING.md](https://github.com/cosscom/ui/blob/3d8ca28e7d90123e9e0f44665944dc240b026ccb/LICENSING.md)
explicitly exempts `apps/ui/` and `apps/origin/` as MIT. Its root README says
the remaining directories use AGPLv3. The UI README and package metadata also
declare MIT. The root `LICENSE` is AGPL, not a license for the UI exception.

Consequently, the native reference material is the MIT **registry and docs**.
Do not copy `packages/*` source or font files as if the whole monorepo were
MIT. Transitive imports in a registry file do not extend the directory license
exception to their targets. Base UI, Tailwind and upstream React packages
are web implementation dependencies, not Go module dependencies here.

There is no dedicated `apps/ui/LICENSE` in this pinned tree. The explicit
upstream MIT license text retained at `apps/origin/LICENSE.md` credits coss.com
and Origin UI. [licenses/coss-ui-NOTICE.md](../licenses/coss-ui-NOTICE.md)
records both that text's exact source and the separate UI directory declaration;
it does not invent a UI-specific copyright notice or relicense AGPL code.
The project's own [LICENSE](../LICENSE) remains unchanged.

`apps/ui/app/globals.css` imports a stylesheet resolved under
`packages/ui/src/styles/globals.css`. The runtime dark background/popover mixes
differ from the published MIT styling page. [theme.md](theme.md) records the
numeric observation and selected native values separately. That comparison
is not permission to copy the imported package implementation.

## Fonts

The site links Cal Sans variable and Paper Mono regular web assets. Do not
copy those from coss packages. [fonts.md](fonts.md) records independently
licensed sources, pinned commits, license attribution, file paths and known
version ambiguity:

- Cal Sans v2: official `calcom/sans` at
  `9fe03d44f00806aadeb6d3653e37dce358554de0`, OFL 1.1,
  `fonts/calsans-cossui/CalSansVF.ttf` and its italic companion.
- Paper Mono: official `paper-design/paper-mono` at
  `e6eaeceaef02e77e3db997711e07a16378de2bd7`, OFL 1.1,
  independent regular/variable TTFs. No Paper Mono addition is implied.

Application startup owns registration. Naming a family in Preferences does
not register it or establish browser-equivalent shaping.

## Reviewing an update

1. Identify the exact new upstream commit and read its directory licensing
   declaration before reusing any source or assets.
2. Compare catalogs, relevant MIT registry files, MDX examples and published
   styling. Record new/removed subjects and real native limits in the matrix.
3. Record numeric/runtime discrepancies explicitly instead of silently mixing
   token sources. Keep independent font provenance separate from coss source.
4. Update capability/design/customization docs and the unreleased changelog.
   Record actual verification evidence in `VALIDATION.md`.

The `references/heroui` and `references/heroui-native` submodules and retained
Apache license are historical reference artifacts. Their presence does not
make them current visual authority. `docs/mygo-0.6.md` preserves the earlier
integration record; its historical validation claims do not validate the coss
reconstruction.
