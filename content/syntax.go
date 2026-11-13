package content

import (
	"strings"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// syntax describes what the tokenizer colours in one language. It is a
// lexical sketch, not a grammar.
type syntax struct {
	name string
	// keywords are space-separated, with a space at both ends.
	keywords              string
	lineComment           string
	blockOpen, blockClose string
	// quotes delimit strings that end on their line; rawQuotes delimit
	// strings that may span lines and have no escapes.
	quotes, rawQuotes string
	tripleQuotes      bool
	// markup colours tag names and attribute strings only inside tags.
	markup bool
	// plainKeys leaves a string followed by a colon uncoloured.
	plainKeys bool
}

// syntaxFor resolves a language name or common alias. Unknown names report
// false and render plain.
func syntaxFor(language string) (syntax, bool) {
	switch strings.ToLower(strings.TrimSpace(language)) {
	case "go", "golang":
		return syntax{
			name:        "Go",
			keywords:    " break case chan const continue default defer else fallthrough for func go goto if import interface map package range return select struct switch type var nil true false iota ",
			lineComment: "//", blockOpen: "/*", blockClose: "*/",
			quotes: `"'`, rawQuotes: "`",
		}, true
	case "js", "javascript", "jsx", "mjs", "cjs", "ts", "typescript", "tsx":
		name := "JavaScript"
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(language)), "t") {
			name = "TypeScript"
		}
		return syntax{
			name:        name,
			keywords:    " abstract as async await break case catch class const continue debugger declare default delete do else enum export extends false finally for from function if implements import in instanceof interface let new null of private protected public readonly return static super switch this throw true try type typeof undefined var void while yield ",
			lineComment: "//", blockOpen: "/*", blockClose: "*/",
			quotes: `"'`, rawQuotes: "`",
		}, true
	case "json", "jsonc":
		return syntax{
			name: "JSON", keywords: " true false null ",
			lineComment: "//", blockOpen: "/*", blockClose: "*/",
			quotes: `"`, plainKeys: true,
		}, true
	case "bash", "sh", "shell", "zsh":
		return syntax{
			name:        "Shell",
			keywords:    " case do done elif else esac export fi for function if in local return then until while ",
			lineComment: "#", quotes: `"'`,
		}, true
	case "py", "python":
		return syntax{
			name:        "Python",
			keywords:    " False None True and as assert async await break class continue def del elif else except finally for from global if import in is lambda nonlocal not or pass raise return try while with yield ",
			lineComment: "#", quotes: `"'`, tripleQuotes: true,
		}, true
	case "yaml", "yml":
		return syntax{
			name: "YAML", keywords: " true false null yes no on off ",
			lineComment: "#", quotes: `"'`,
		}, true
	case "html", "xml", "svg":
		return syntax{
			name: strings.ToUpper(strings.TrimSpace(language)), blockOpen: "<!--", blockClose: "-->",
			quotes: `"'`, markup: true,
		}, true
	case "css":
		return syntax{name: "CSS", blockOpen: "/*", blockClose: "*/", quotes: `"'`}, true
	}
	return syntax{}, false
}

// syntaxColors are the theme roles of the coloured token kinds. Plain code
// inherits the block's text colour.
type syntaxColors struct {
	comment, str, number, keyword ui.Color
}

func syntaxColorsOf(t theme.Palette) syntaxColors {
	return syntaxColors{
		comment: t.Muted, str: t.SuccessForeground,
		number: t.WarningForeground, keyword: t.Primary,
	}
}

func identStart(b byte) bool {
	return b == '_' || b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b >= 0x80
}

func digit(b byte) bool { return b >= '0' && b <= '9' }

// highlight splits src into spans: comments muted italic, strings, numbers
// and keywords in their roles, the rest in one uncoloured run between them.
// It only cuts at ASCII bytes, so multi-byte runes stay whole.
func highlight(src string, s syntax, colors syntaxColors) []ui.Span {
	spans := make([]ui.Span, 0, 16)
	plain := 0
	emit := func(from, to int, span ui.Span) {
		if from > plain {
			spans = append(spans, ui.Span{Text: src[plain:from]})
		}
		span.Text = src[from:to]
		spans = append(spans, span)
		plain = to
	}
	n := len(src)
	inTag := false
	for i := 0; i < n; {
		ch := src[i]
		rest := src[i:]
		if s.blockOpen != "" && strings.HasPrefix(rest, s.blockOpen) {
			end := n
			if k := strings.Index(rest[len(s.blockOpen):], s.blockClose); k >= 0 {
				end = i + len(s.blockOpen) + k + len(s.blockClose)
			}
			emit(i, end, ui.Span{Color: colors.comment, Italic: true})
			i = end
			continue
		}
		if s.markup {
			switch {
			case !inTag && ch == '<':
				inTag = true
				i++
				if i < n && src[i] == '/' {
					i++
				}
				j := i
				for j < n && (identStart(src[j]) || digit(src[j]) || src[j] == '-') {
					j++
				}
				if j > i {
					emit(i, j, ui.Span{Color: colors.keyword, Weight: theme.WeightMedium})
				}
				i = j
				continue
			case !inTag:
				i++
				continue
			case ch == '>':
				inTag = false
				i++
				continue
			}
		}
		// A hash opens a comment only at the start of a word, so "$#" and
		// URL fragments stay code.
		if s.lineComment != "" && strings.HasPrefix(rest, s.lineComment) &&
			(s.lineComment != "#" || i == 0 || src[i-1] == ' ' || src[i-1] == '\t' || src[i-1] == '\n') {
			end := n
			if k := strings.IndexByte(rest, '\n'); k >= 0 {
				end = i + k
			}
			emit(i, end, ui.Span{Color: colors.comment, Italic: true})
			i = end
			continue
		}
		switch {
		case strings.IndexByte(s.rawQuotes, ch) >= 0:
			end := n
			if k := strings.IndexByte(rest[1:], ch); k >= 0 {
				end = i + k + 2
			}
			emit(i, end, ui.Span{Color: colors.str})
			i = end
		case strings.IndexByte(s.quotes, ch) >= 0:
			end := stringEnd(src, i, s.tripleQuotes)
			if end < 0 {
				// No closing quote on the line: an apostrophe, not a string.
				i++
				continue
			}
			if !s.plainKeys || !beforeColon(src, end) {
				emit(i, end, ui.Span{Color: colors.str})
			}
			i = end
		case digit(ch):
			j := i + 1
			for j < n && (identStart(src[j]) || digit(src[j]) || src[j] == '.') {
				j++
			}
			emit(i, j, ui.Span{Color: colors.number})
			i = j
		case identStart(ch):
			j := i + 1
			for j < n && (identStart(src[j]) || digit(src[j])) {
				j++
			}
			// A member, a flag or a hyphenated name is not a keyword.
			joined := i > 0 && strings.IndexByte(".-$", src[i-1]) >= 0 || j < n && src[j] == '-'
			if !joined && s.keywords != "" && isKeyword(s.keywords, src[i:j]) {
				emit(i, j, ui.Span{Color: colors.keyword, Weight: theme.WeightMedium})
			}
			i = j
		default:
			i++
		}
	}
	if plain < n {
		spans = append(spans, ui.Span{Text: src[plain:]})
	}
	return spans
}

// stringEnd returns the index after the string opening at i, or -1 when it
// does not close on its line. Triple-quoted strings may span lines and run
// to the end of the source when left open.
func stringEnd(src string, i int, triple bool) int {
	q := src[i]
	if triple && i+2 < len(src) && src[i+1] == q && src[i+2] == q {
		if k := strings.Index(src[i+3:], src[i:i+3]); k >= 0 {
			return i + 3 + k + 3
		}
		return len(src)
	}
	for j := i + 1; j < len(src); j++ {
		switch src[j] {
		case '\\':
			j++
		case '\n':
			return -1
		case q:
			return j + 1
		}
	}
	return -1
}

// beforeColon reports whether the next byte after spaces is a colon, as
// after an object key.
func beforeColon(src string, i int) bool {
	for i < len(src) && (src[i] == ' ' || src[i] == '\t') {
		i++
	}
	return i < len(src) && src[i] == ':'
}

// isKeyword finds word in a space-delimited list without allocating.
func isKeyword(list, word string) bool {
	for from := 0; ; {
		k := strings.Index(list[from:], word)
		if k < 0 {
			return false
		}
		at := from + k
		if list[at-1] == ' ' && list[at+len(word)] == ' ' {
			return true
		}
		from = at + 1
	}
}
