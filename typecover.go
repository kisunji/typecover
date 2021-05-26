package typecover

import (
	"fmt"
	"go/ast"
	"go/types"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

const Doc = `typecover checks that a code block is assigning to all exported fields of a struct or calling all exported methods of an interface.`

var Analyzer = &analysis.Analyzer{
	Doc:      Doc,
	Name:     "typecover",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var commentRegex = regexp.MustCompile(`typecover:([\w.]+)`)

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		commentMap := ast.NewCommentMap(pass.Fset, file, file.Comments)
		ast.Inspect(file, func(n ast.Node) bool {
			for _, comments := range commentMap[n] {
				for _, comment := range comments.List {
					matches := commentRegex.FindAllStringSubmatch(comment.Text, 1)
					if len(matches) == 1 && len(matches[0]) == 2 {
						typeName := fullTypeName(pass, file, n, strings.TrimSpace(matches[0][1]))
						t := findType(pass, typeName)
						if t == nil {
							reportNodef(pass, n, "Type %s not found in project scope", typeName)
							return false
						}
						missing := checkMembers(pass, n, t)
						if len(missing) > 0 {
							reportNodef(pass, n, "Type %s is missing %s", typeName, strings.Join(missing, ", "))
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

func findType(pass *analysis.Pass, targetType string) types.Type {
	ss := strings.Split(targetType, ".")
	pkgName := ss[0]
	typeName := ss[1]
	if pass.Pkg.Path() == pkgName {
		o := pass.Pkg.Scope().Lookup(typeName)
		if o != nil {
			return o.Type()
		}
	}

	for _, imp := range pass.Pkg.Imports() {
		if imp.Path() == pkgName {
			o := imp.Scope().Lookup(typeName)
			if o != nil {
				return o.Type()
			}
		}
	}
	return nil
}

func checkMembers(pass *analysis.Pass, n ast.Node, target types.Type) []string {
	var missing []string
	membersFound := map[string]bool{}

	switch u := target.Underlying().(type) {
	case *types.Interface:
		for i := 0; i < u.NumMethods(); i++ {
			if u.Method(i).Exported() {
				membersFound[u.Method(i).Name()] = false
			}
		}

		ast.Inspect(n, func(n ast.Node) bool {
			if se, ok := n.(*ast.SelectorExpr); ok {
				t1 := pass.TypesInfo.TypeOf(se.X)
				if t1 == nil {
					return true
				}

				// either the type itself or the pointer of type should implement interface u
				if !types.Implements(t1, u) && !types.Implements(types.NewPointer(t1), u) {
					return true
				}

				if se.Sel != nil {
					if _, ok := membersFound[se.Sel.Name]; ok {
						membersFound[se.Sel.Name] = true
					}
				}
			}
			return true
		})

	case *types.Struct:
		for i := 0; i < u.NumFields(); i++ {
			if u.Field(i).Exported() {
				membersFound[u.Field(i).Name()] = false
			}
		}

		ast.Inspect(n, func(n ast.Node) bool {
			switch nodeType := n.(type) {
			case *ast.CompositeLit: // nodeType = MyType{Field: 1}
				t := pass.TypesInfo.TypeOf(nodeType.Type)
				if t == nil || strings.TrimPrefix(t.String(), "*")  != target.String() {
					return true
				}

				for _, e := range nodeType.Elts {
					if k, ok := e.(*ast.KeyValueExpr); ok {
						if i, ok2 := k.Key.(*ast.Ident); ok2 {
							if _, ok3 := membersFound[i.Name]; ok3 {
								membersFound[i.Name] = true
							}
						}
					} else {
						// todo: support CompositeLit with anonymous fields
					}
				}
			case *ast.AssignStmt: // nodeType.Field = val
				for _, s := range nodeType.Lhs {
					if se, ok := s.(*ast.SelectorExpr); ok {
						t := pass.TypesInfo.TypeOf(se.X)
						if t == nil || strings.TrimPrefix(t.String(), "*") != target.String() {
							return true
						}
						if se.Sel != nil {
							if _, ok := membersFound[se.Sel.Name]; ok {
								membersFound[se.Sel.Name] = true
							}
						}
					}
				}
			}
			return true
		})
	}

	for member, found := range membersFound {
		if !found {
			missing = append(missing, member)
		}
	}

	return missing
}

func fullTypeName(pass *analysis.Pass, file *ast.File, n ast.Node, typeName string) string {
	selectorParts := strings.Split(typeName, ".")
	if len(selectorParts) == 2 {
		for _, fimport := range file.Imports {
			var pkgName string
			if fimport.Name != nil {
				if fimport.Name.Name == "." {
					// TODO: handle dot imports
					reportNodef(pass, n, "Dot imports are unhandled!")
				}
				pkgName = fimport.Name.Name
			} else {
				components := strings.Split(unquote(fimport.Path.Value), "/")
				pkgName = components[len(components)-1]
			}
			if selectorParts[0] == pkgName {
				typeName = unquote(fimport.Path.Value) + "." + selectorParts[1]
			}
		}
	} else {
		typeName = pass.Pkg.Path() + "." + typeName
	}
	return typeName
}

func reportNodef(pass *analysis.Pass, node ast.Node, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	pass.Report(analysis.Diagnostic{Pos: node.Pos(), End: node.End(), Message: msg})
}

func unquote(str string) string {
	if unquoted, err := strconv.Unquote(str); err == nil {
		return unquoted
	}
	return str
}
