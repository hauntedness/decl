package decl

import (
	"go/ast"

	"golang.org/x/tools/go/packages"
)

type Package struct {
	structs    []*Struct
	funcs      []*Func
	interfaces []*Interface
	vars       []*Var
	implements []*ImplementStmt
}

func Load(pkg *packages.Package) (*Package, error) {
	c := &Package{}
	commentMap := map[*ast.Ident]*Comment{}
	fieldsCommentMap := map[*ast.Ident][]*Comment{}
	for _, syntax := range pkg.Syntax {
		comments := ast.NewCommentMap(pkg.Fset, syntax, syntax.Comments)
		for _, decl := range syntax.Decls {
			switch declImpl := decl.(type) {
			case *ast.GenDecl:
				for i := range declImpl.Specs {
					switch spec := declImpl.Specs[i].(type) {
					case *ast.TypeSpec:
						commentMap[spec.Name] = NewComment(comments.Filter(decl))
						// here we need more comments for fields
						switch specType := spec.Type.(type) {
						case *ast.StructType:
							for _, f := range specType.Fields.List {
								for range f.Names {
									fieldsCommentMap[spec.Name] = append(fieldsCommentMap[spec.Name], NewComment(comments.Filter(f)))
								}
							}
						case *ast.InterfaceType:
							for _, f := range specType.Methods.List {
								for range f.Names {
									fieldsCommentMap[spec.Name] = append(fieldsCommentMap[spec.Name], NewComment(comments.Filter(f)))
								}
							}
						}
					case *ast.ValueSpec:
						if len(spec.Names) == 1 {
							commentMap[spec.Names[0]] = NewComment(comments.Filter(decl))
						}
					}
				}
			case *ast.FuncDecl:
				commentMap[declImpl.Name] = NewComment(comments.Filter(decl))
			default:
				// omit other case
			}
		}
	}
	for id, def := range pkg.TypesInfo.Defs {
		if id == nil {
			continue
		}
		switch Kind(def) {
		case KindStruct:
			st, err := NewStruct(def)
			if err != nil {
				return nil, err
			}
			st.comment = commentMap[id]
			fieldsComment := fieldsCommentMap[id]
			if len(fieldsComment) == len(st.fields) {
				for i := range st.fields {
					st.fields[i].comment = fieldsComment[i]
				}
			}
			c.structs = append(c.structs, st)
		case KindInterface:
			it, err := NewInterface(def)
			if err != nil {
				continue
			}
			it.comment = commentMap[id]
			fieldsComment := fieldsCommentMap[id]
			if len(fieldsComment) == it.iface.NumMethods() {
				for i := range it.iface.NumMethods() {
					it.methodComments[i] = fieldsComment[i]
				}
			}
			c.interfaces = append(c.interfaces, it)
		case KindFunc:
			// TODO exclude interface method.
			fn, err := NewFunc(def)
			if err != nil {
				return nil, err
			}
			fn.comment = (commentMap[id])
			c.funcs = append(c.funcs, fn)
		case KindVar:
			var1, err := NewVar(def)
			if err != nil {
				return nil, err
			}
			var1.comment = commentMap[id]
			c.vars = append(c.vars, var1)
		}
	}
	return c, nil
}
