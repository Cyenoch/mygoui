# Third-party notices

Project code retains its [MIT license](LICENSE). References, dependencies and
optional font assets retain their own terms. This is an independent native Go
implementation; no upstream affiliation is implied.

## coss ui reference

The visual/component reference is
[`cosscom/ui`](https://github.com/cosscom/ui/tree/3d8ca28e7d90123e9e0f44665944dc240b026ccb)
commit `3d8ca28e7d90123e9e0f44665944dc240b026ccb`, specifically its
`apps/ui` registry and documentation.

That repository uses **mixed licensing**. `LICENSING.md` and the root README
identify `apps/ui/` and `apps/origin/` as MIT; the remaining directories use
AGPLv3. `apps/ui/README.md` and `apps/ui/package.json` independently declare
MIT. The monorepo root AGPL license must not be mistaken for the UI exception,
and the UI exception must not be extended to `packages/*` source or fonts.

[licenses/coss-ui-NOTICE.md](licenses/coss-ui-NOTICE.md) records these exact
sources and retains the upstream MIT text available at
`apps/origin/LICENSE.md`, including its coss.com and Origin UI attribution.
There is no separate `apps/ui/LICENSE` in the pinned tree; the retained text
is identified by its actual path rather than represented as a UI-specific
copyright declaration. [docs/references.md](docs/references.md) records
registry paths and the published/runtime token discrepancy.

## Fonts

The theme's native system defaults require no bundled coss font assets.
Applications can explicitly register their chosen faces. Independent sources
are documented in [docs/fonts.md](docs/fonts.md):

- **Cal Sans v2:** official `calcom/sans` commit
  `9fe03d44f00806aadeb6d3653e37dce358554de0`,
  `fonts/calsans-cossui/CalSansVF.ttf` and italic companion. Upstream
  `OFL.txt` declares SIL OFL 1.1, copyright 2021 The Cal Sans Project Authors.
  The optional `fonts` package embeds those faces and retains the complete
  [OFL notice](fonts/OFL.txt). This source is independent of coss `packages/*`.
  Registration remains an explicit host startup action.
- **Paper Mono:** official `paper-design/paper-mono` commit
  `e6eaeceaef02e77e3db997711e07a16378de2bd7`, SIL OFL 1.1.
  Its `LICENSE.txt` credits Lost Coast Labs, Inc. (Paper Design), 2025, and
  Geist Mono's Vercel/basement.studio attribution, 2023; its `OFL.txt` also
  credits the Paper-Mono.Git Project Authors, 2025. It is an optional host
  source, not a Paper Mono asset added by this documentation work.

OFL font files remain OFL when bundled with software; the project's MIT
license does not replace their terms. Family names in theme preferences do
not download, install or register the corresponding fonts.

## MyGo

`github.com/egoist/mygo v0.7.1` is a module dependency under MIT, copyright
2026 MyGo contributors. Its full [license](licenses/MyGo-LICENSE) is retained.

Spinner uses MyGo's `plugins/loading`, whose shapes and motion follow Jakub
Krehel's loading.dev. The complete upstream attribution and MIT terms are
retained in [licenses/loading-NOTICE.md](licenses/loading-NOTICE.md).

MyGo and its transitive dependencies retain their own licenses. Applications
distributing linked binaries should inventory their actual dependencies and
include applicable notices.

## mygoicons and Lucide

`github.com/Cyenoch/mygoicons v0.2.1` supplies typed Lucide icons. Its project
code is MIT; the complete [code license](licenses/mygoicons-LICENSE) is
retained. The imported `lucide` package contains `@iconify-json/lucide 1.2.147`
artwork. The complete [Lucide notice](licenses/Lucide-LICENSE) includes ISC,
Feather MIT and the applicable Feather-derived icon list. Retain it when
distributing icons or applications containing them.

No Twemoji package is imported by this library.

## Earlier reference attribution

The [retained Apache 2.0 notice](licenses/legacy-Apache-2.0.txt) preserves the
original NextUI Inc. copyright for earlier reference work. It is not a current
design source. The only reference submodule is `references/coss`; the mixed
upstream licensing boundary is described in [docs/references.md](docs/references.md).

## Agent skills

The optional `.agents/skills` directory contains maintenance instructions,
not runtime code. MyGo maintenance instructions accompany the project.
Copied `golang-*` skills identify samber's cc-skills-golang in their metadata
and retain their declared MIT license/attribution. Ponytail instructions
retain their declared MIT license and source metadata. These development
files do not become application dependencies.
