# Cal Sans

This optional package embeds the upright and italic Cal Sans variable fonts
and their SIL Open Font License 1.1. The files are unmodified from
[calcom/sans at revision 9fe03d44f00806aadeb6d3653e37dce358554de0](https://github.com/calcom/sans/tree/9fe03d44f00806aadeb6d3653e37dce358554de0):

- `CalSansVF.ttf` and `CalSansVF-Italic.ttf` are from `fonts/calsans-cossui/`.
- `OFL.txt` is from the repository root.

The font metadata credits Designer Mark Davis DBA Wordmark (2026); the
license credits the Cal Sans Project Authors (2021).

Call `fonts.Register()` explicitly and handle its error, then use `fonts.Family`
with MyGo's `Font` or an application's font preferences. Importing the package
does not register the font. Registration is process-wide, and repeated calls
return the first result. If the italic face fails to register, the regular
face remains registered.

Browser builds return an error: MyGo 0.7.1 registers both faces as normal
weight and style and does not report asynchronous font-loading errors.
