package gallery

import (
	"embed"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
	"sync"
)

// The page files carry every example, so a demo's code is read from the
// source that runs it and cannot drift from it.
//
//go:embed component_pages_*.go page_*.go
var pageSources embed.FS

// exampleSources maps a page ID to its examples' Go source by title, parsed
// once on first use.
var exampleSources = sync.OnceValue(func() map[string]map[string]string {
	fset := token.NewFileSet()
	type parsed struct {
		file *ast.File
		text []byte
	}
	var files []parsed
	funcs := map[string]parsed{} // top-level function name to its file
	bodies := map[string]*ast.BlockStmt{}
	// The pattern also matches the pages' test files; they are skipped here.
	entries, err := pageSources.ReadDir(".")
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		text, err := pageSources.ReadFile(entry.Name())
		if err != nil {
			panic(err)
		}
		file, err := parser.ParseFile(fset, entry.Name(), text, parser.SkipObjectResolution)
		if err != nil {
			panic(err)
		}
		p := parsed{file, text}
		files = append(files, p)
		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Recv == nil && fn.Body != nil {
				funcs[fn.Name.Name], bodies[fn.Name.Name] = p, fn.Body
			}
		}
	}
	source := func(text []byte, body *ast.BlockStmt) string {
		from, to := fset.Position(body.Lbrace).Offset+1, fset.Position(body.Rbrace).Offset
		return dedent(string(text[from:to]))
	}
	out := map[string]map[string]string{}
	for _, p := range files {
		ast.Inspect(p.file, func(n ast.Node) bool {
			page, ok := n.(*ast.CompositeLit)
			if !ok || !isIdent(page.Type, "ComponentPage") {
				return true
			}
			id, examples := "", (*ast.CompositeLit)(nil)
			for _, element := range page.Elts {
				switch key, value := keyValue(element); key {
				case "ID":
					id = stringLit(value)
				case "Examples":
					examples, _ = value.(*ast.CompositeLit)
				}
			}
			if id == "" || examples == nil {
				return true
			}
			for _, element := range examples.Elts {
				example, ok := element.(*ast.CompositeLit)
				if !ok {
					continue
				}
				title, code := "", ""
				for _, field := range example.Elts {
					switch key, value := keyValue(field); key {
					case "Title":
						title = stringLit(value)
					case "Demo":
						// A demo that is or wraps one page function, or only
						// passes its own arguments on to one, shows that function.
						name := ""
						switch demo := value.(type) {
						case *ast.FuncLit:
							code = source(p.text, demo.Body)
							name = soleCall(demo.Body)
						case *ast.CallExpr:
							if len(demo.Args) != 1 {
								break
							}
							if lit, ok := demo.Args[0].(*ast.FuncLit); ok {
								code = source(p.text, lit.Body)
								name = soleCall(lit.Body)
							} else {
								name = identName(demo.Args[0])
							}
						case *ast.Ident:
							name = demo.Name
						}
						if bodies[name] != nil {
							code = source(funcs[name].text, bodies[name])
						}
					}
				}
				if title != "" && code != "" {
					if out[id] == nil {
						out[id] = map[string]string{}
					}
					out[id][title] = code
				}
			}
			return true
		})
	}
	return out
})

// exampleSource returns the Go source of one example, or "" when the page
// builds its examples in a way the source does not spell out.
func exampleSource(pageID, title string) string { return exampleSources()[pageID][title] }

func isIdent(expr ast.Expr, name string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == name
}

func keyValue(expr ast.Expr) (string, ast.Expr) {
	kv, ok := expr.(*ast.KeyValueExpr)
	if !ok {
		return "", nil
	}
	key, _ := kv.Key.(*ast.Ident)
	if key == nil {
		return "", nil
	}
	return key.Name, kv.Value
}

func stringLit(expr ast.Expr) string {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return ""
	}
	text, _ := strconv.Unquote(lit.Value)
	return text
}

// soleCall names the function body consists of calling with nothing but
// the demo's own arguments, which says less than that function's body.
func soleCall(body *ast.BlockStmt) string {
	if len(body.List) != 1 {
		return ""
	}
	stmt, ok := body.List[0].(*ast.ExprStmt)
	if !ok {
		return ""
	}
	call, ok := stmt.X.(*ast.CallExpr)
	if !ok {
		return ""
	}
	for _, arg := range call.Args {
		if !passedOn(arg) {
			return ""
		}
	}
	return identName(call.Fun)
}

// passedOn reports whether arg is one of the demo's own arguments or the
// state looked up from them, such as cossStateOf(s, "button").
func passedOn(arg ast.Expr) bool {
	if identName(arg) != "" {
		return true
	}
	lookup, ok := arg.(*ast.CallExpr)
	if !ok {
		return false
	}
	for _, key := range lookup.Args {
		if identName(key) == "" && stringLit(key) == "" {
			return false
		}
	}
	return true
}

func identName(expr ast.Expr) string {
	if id, ok := expr.(*ast.Ident); ok {
		return id.Name
	}
	return ""
}

// dedent trims blank edges and the indentation the lines share.
func dedent(text string) string {
	lines := strings.Split(strings.Trim(text, "\n"), "\n")
	indent := -1
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		tabs := len(line) - len(strings.TrimLeft(line, "\t"))
		if indent < 0 || tabs < indent {
			indent = tabs
		}
	}
	for i, line := range lines {
		if len(line) >= indent && indent > 0 {
			lines[i] = line[indent:]
		}
		lines[i] = strings.ReplaceAll(strings.TrimRight(lines[i], " \t"), "\t", "    ")
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}
