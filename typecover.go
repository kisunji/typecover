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
					var exclude []string
					matches := commentRegex.FindAllStringSubmatch(comment.Text, 1)
					if len(matches) == 1 && len(matches[0]) == 2 {
						// Check if there are any member to exclude
						ss := strings.Split(comment.Text, "-exclude")
						if len(ss) > 1 {
							exclude = strings.Split(strings.TrimLeft(ss[1], "= "), ",")
						}
						if len(ss) > 2 {
							reportNodef(pass, n, "Detected more than one '~'. Separate arguments with commas")
						}

						typeName := fullTypeName(pass, file, n, strings.TrimSpace(matches[0][1]))
						t := findType(pass, typeName)
						if t == nil {
							reportNodef(pass, n, "Type %s not found in project scope", typeName)
							return false
						}
						missing := checkMembers(pass, n, t, exclude)
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

func checkMembers(pass *analysis.Pass, n ast.Node, target types.Type, exclude []string) []string {
	membersFound := map[string]bool{}

	for _, t := range []types.Type{target, types.NewPointer(target)} {
		mset := types.NewMethodSet(t)
		for i := 0; i < mset.Len(); i++ {
			switch u := mset.At(i).Obj().(type) {
			case *types.Func:
				if u.Exported() {
					membersFound[u.Name()] = false
				}
			}
		}
	}

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
			switch nt := n.(type) {
			case *ast.CompositeLit: // nt = MyType{Field: 1}
				t := pass.TypesInfo.TypeOf(nt.Type)
				if t == nil || strings.TrimPrefix(t.String(), "*") != target.String() {
					return true
				}

				for _, e := range nt.Elts {
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
			case *ast.SelectorExpr: // nt.Field = val
				t := pass.TypesInfo.TypeOf(nt.X)
				if t == nil || strings.TrimPrefix(t.String(), "*") != target.String() {
					return true
				}
				if nt.Sel != nil {
					if _, ok := membersFound[nt.Sel.Name]; ok {
						membersFound[nt.Sel.Name] = true
					}
				}
			}
			return true
		})
	}

	// Mark all excluded members as found.
	for _, e := range exclude {
		if _, ok := membersFound[strings.TrimSpace(e)]; ok {
			membersFound[e] = true
		}
	}

	var missing []string
	for member, found := range membersFound {
		if !found {
			missing = append(missing, member)
		}
	}
	return missing
}

func findType(pass *analysis.Pass, targetType string) types.Type {
	lastDot := strings.LastIndex(targetType, ".")
	if lastDot == -1 {
		panic(fmt.Sprintf("ill-formed targetType %q", targetType))
	}
	pkgName := targetType[:lastDot]
	typeName := targetType[lastDot+1:]
	if pkgMatch(pass.Pkg.Path(), pkgName) {
		o := pass.Pkg.Scope().Lookup(typeName)
		if o != nil {
			return o.Type()
		}
	}

	for _, imp := range pass.Pkg.Imports() {
		if pkgMatch(imp.Path(), pkgName) {
			o := imp.Scope().Lookup(typeName)
			if o != nil {
				return o.Type()
			}
		}
	}
	return nil
}

func pkgMatch(path, pkgName string) bool {
	for strings.Contains(path, "/vendor/") {
		path = strings.Split(path, "/vendor/")[1]
	}
	return path == pkgName
}

func fullTypeName(pass *analysis.Pass, file *ast.File, n ast.Node, typeName string) string {
	selectorParts := strings.Split(typeName, ".")
	if len(selectorParts) == 2 {
		for _, fimport := range file.Imports {
			var pkgName string
			if fimport.Name != nil {
				if fimport.Name.Name == "." {
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
