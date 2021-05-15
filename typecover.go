package typecover

import (
	"fmt"
	"go/ast"
	"go/types"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

const Doc = `typecover checks that a code block is assigning to all exported 
fields of a struct or calling all exported methods of an interface.`

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
						checkMembers(pass, n, typeName)
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

func checkMembers(pass *analysis.Pass, n ast.Node, typeName string) {
	var typeNameFound bool
	ast.Inspect(n, func(n ast.Node) bool {
		compositeLit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}
		t := pass.TypesInfo.TypeOf(compositeLit.Type)
		if t == nil {
			return true
		}

		fmt.Printf("found composite literal %s \n", t.String())
		if t.String() == typeName {
			typeNameFound = true
			missing := []string{}
			str, ok := t.Underlying().(*types.Struct)
			if !ok {
				return true
			}

			for i := 0; i < str.NumFields(); i++ {
				fieldName := str.Field(i).Name()
				exists := false

				if !str.Field(i).Exported() {
					continue
				}

				for j, e := range compositeLit.Elts {
					if k, ok := e.(*ast.KeyValueExpr); ok {
						if i, ok := k.Key.(*ast.Ident); ok {
							if i.Name == fieldName {
								exists = true
								break
							}
						}
					} else {
						// Anonymous fields (e.g. Foo{1, true, "hello"})
						if j == i {
							exists = true
							break
						}
					}
				}
				if !exists {
					missing = append(missing, fieldName)
				}
			}
			if len(missing) > 0 {
				reportNodef(pass, n, "Type %s is missing %s", t.String(), strings.Join(missing, ", "))
			}
			return false // stop walking since we processed struct
		}
		return true
	})
	if !typeNameFound {
		reportNodef(pass, n, "Type %s not found in associated code block", typeName)
	}
	// allConsts := allConstsWithType(pass, typeName)
	// if len(allConsts) == 0 {
	// 	reportNodef(pass, n, "No consts found for type %v", typeName)
	// }
	// for _, want := range allConsts {
	// 	if !membersForType[want.name] && !membersForType[want.val] {
	// 		reportNodef(pass, n, "Unhandled const: %v", want)
	// 	}
	// }
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

type constVal struct {
	name string
	val  string
}

func (c constVal) String() string {
	return fmt.Sprintf("%s (%s)", c.name, c.val)
}

var allPkgs sync.Map

// TODO: do this by storing analysis.Facts about all the consts in each package?
func allConstsWithType(pass *analysis.Pass, targetType string) []constVal {
	var visit func(pkg *types.Package)
	visit = func(pkg *types.Package) {
		if _, ok := allPkgs.Load(pkg); ok {
			return
		}
		allPkgs.Store(pkg, struct{}{})
		for _, imp := range pkg.Imports() {
			visit(imp)
		}
	}
	visit(pass.Pkg)
	consts := []constVal{}
	allPkgs.Range(func(pkgKey, _ interface{}) bool {
		pkg := pkgKey.(*types.Package)
		for _, name := range pkg.Scope().Names() {
			if namedConst, ok := pkg.Scope().Lookup(name).(*types.Const); ok {
				val := unquote(namedConst.Val().ExactString())
				typeName := namedConst.Type().String()
				if typeName == targetType {
					consts = append(consts, constVal{name: namedConst.Name(), val: val})
				}
			}
		}
		return true
	})
	return consts
}
