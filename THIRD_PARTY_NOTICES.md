# Third-party notices

Project code is [MIT](LICENSE). References and dependencies retain their own terms.

## HeroUI

[HeroUI](https://github.com/heroui-inc/heroui) and
[hero-ui-native](https://github.com/heroui-inc/heroui-native) are the visual and
component references. HeroUI web and native are Apache License 2.0. The upstream
license is retained in [licenses/HeroUI-LICENSE](licenses/HeroUI-LICENSE).
This project is an independent native Go implementation and is not an official
HeroUI or MyGo distribution.

## MyGo

`github.com/egoist/mygo v0.6.3` is a module dependency under MIT, copyright 2026 MyGo contributors. Its full [license](licenses/MyGo-LICENSE) is retained.

Spinner uses MyGo's `plugins/loading`, whose shapes and motion follow Jakub Krehel's loading.dev. The plugin's complete upstream attribution and MIT terms are retained in [licenses/loading-NOTICE.md](licenses/loading-NOTICE.md).

MyGo and its transitive dependencies retain their own licenses. Applications distributing linked binaries should inventory their actual dependencies and include the applicable notices.

## mygoicons and Lucide

`github.com/Cyenoch/mygoicons v0.2.1` supplies typed Lucide icons. Its project code is MIT; the complete [code license](licenses/mygoicons-LICENSE) is retained. The imported `lucide` package contains the `@iconify-json/lucide 1.2.147` artwork. The complete [Lucide notice](licenses/Lucide-LICENSE) includes ISC, Feather MIT and the applicable Feather-derived icon list. Retain it when distributing icons or applications containing them.

No Twemoji package is imported by this library; Twemoji artwork is not redistributed in this source tree.

## Agent skills

The optional `.agents/skills` directory contains maintenance instructions, not runtime code. MyGo maintenance instructions accompany this project. The copied `golang-*` skills identify samber's [cc-skills-golang](https://github.com/samber/cc-skills-golang) in their metadata; their declared MIT license and attribution are retained. The ponytail instructions retain their declared MIT license and source metadata. These files are development guidance and do not become an application dependency.
