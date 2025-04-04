package decl

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"iter"
	"strings"

	"golang.org/x/tools/go/packages"
)

// Package with doc associated.
type Package struct {
	*packages.Package
	// position map is used to find the associated ident.
	pmap map[token.Pos]*ast.Ident
	// (TODO) try remove this as currently we don't modify ident position.
	cmap ast.CommentMap
}

var DefaultLoadMode = packages.LoadFiles | packages.LoadSyntax | packages.LoadImports

func Load(pkg string) (*Package, error) {
	pkgs, err := packages.Load(&packages.Config{Mode: DefaultLoadMode}, pkg)
	if err != nil {
		return nil, err
	}
	if len(pkgs) != 1 {
		return nil, fmt.Errorf("expect <%v> is 1 package but not", pkg)
	}
	return LoadPackage(pkgs[0])
}

// LoadPackage make sure pkg was loaded with correct LoadMode.
func LoadPackage(pkg *packages.Package) (*Package, error) {
	c := &Package{
		Package: pkg,
		cmap:    ast.CommentMap{},
		pmap:    map[token.Pos]*ast.Ident{},
	}
	for _, syntax := range pkg.Syntax {
		ast.Inspect(syntax, func(n ast.Node) bool {
			switch nodeType := n.(type) {
			case *ast.GenDecl:
				for _, spec := range nodeType.Specs {
					switch specType := spec.(type) {
					case *ast.TypeSpec:
						c.cmap[specType.Name] = []*ast.CommentGroup{nodeType.Doc}
					case *ast.ValueSpec:
						for _, name := range specType.Names {
							c.cmap[name] = []*ast.CommentGroup{nodeType.Doc}
						}
					}
				}
			case *ast.Field:
				for _, ident := range nodeType.Names {
					c.cmap[ident] = []*ast.CommentGroup{nodeType.Comment, nodeType.Doc}
				}
			}
			return true
		})
	}
	for ident, def := range pkg.TypesInfo.Defs {
		if ident == nil || def == nil {
			continue
		}
		c.pmap[ident.Pos()] = ident
	}
	return c, nil
}

// CommentsRaw return raw comment group, including nil.
func (pkg *Package) CommentsRaw(pos token.Pos) []*ast.CommentGroup {
	return pkg.cmap[pkg.pmap[pos]]
}

// Comments return comments in text slices, remove then '\n' in the end and nil values.
func (pkg *Package) Comments(obj types.Object) Comments {
	return pkg.commentLines(pkg.CommentsRaw(obj.Pos()))
}

// CommentsAt return comments at pos, remove line feed and nil comment group.
func (pkg *Package) CommentsAt(pos token.Pos) Comments {
	return pkg.commentLines(pkg.CommentsRaw(pos))
}

func (pkg *Package) commentLines(comments []*ast.CommentGroup) Comments {
	ss := make([]string, 0, 1)
	for _, v := range comments {
		if v != nil {
			for i := range v.List {
				ss = append(ss, strings.TrimSuffix(v.List[i].Text, "\n"))
			}
		}
	}
	return ss
}

// Definitions return iterator over each object and its comments.
func (pkg *Package) Definitions() iter.Seq2[types.Object, Comments] {
	return func(yield func(types.Object, Comments) bool) {
		for ident, def := range pkg.TypesInfo.Defs {
			if ident == nil || def == nil {
				continue
			}
			if !yield(def, pkg.CommentsAt(ident.Pos())) {
				return
			}
		}
	}
}

type NamedInfo = DefinedInfo

type DefinedInfo struct {
	Definition types.Object
	TypeName   *types.TypeName
	// Defined Type also called Named Type in go.
	// represent for a declaration such as:
	// type S struct { ... }
	Named *types.Named
}

// DefinedTypes return iterator over each defined type and its comments.
func (pkg *Package) DefinedTypes() iter.Seq2[DefinedInfo, Comments] {
	return func(yield func(DefinedInfo, Comments) bool) {
		for ident, def := range pkg.TypesInfo.Defs {
			if ident == nil || def == nil {
				continue
			}
			if tn, ok := def.(*types.TypeName); ok {
				if nm, ok := tn.Type().(*types.Named); ok {
					info := DefinedInfo{
						Definition: def,
						TypeName:   tn,
						Named:      nm,
					}
					if !yield(info, pkg.CommentsAt(ident.Pos())) {
						return
					}
				}
			}
		}
	}
}

type InterfaceInfo struct {
	Definition types.Object
	TypeName   *types.TypeName
	// Defined Type also called Named Type in go.
	// represent for a declaration such as:
	// type S struct { ... }
	Named      *types.Named
	Underlying *types.Interface
}

// Interfaces return iterator over each struct and its comments.
func (pkg *Package) Interfaces() iter.Seq2[InterfaceInfo, Comments] {
	return func(yield func(InterfaceInfo, Comments) bool) {
		for ident, def := range pkg.TypesInfo.Defs {
			if ident == nil || def == nil {
				continue
			}
			if tn, ok := def.(*types.TypeName); ok {
				if nm, ok := tn.Type().(*types.Named); ok {
					if st, ok := nm.Underlying().(*types.Interface); ok {
						info := InterfaceInfo{
							Definition: def,
							TypeName:   tn,
							Named:      nm,
							Underlying: st,
						}
						if !yield(info, pkg.CommentsAt(ident.Pos())) {
							return
						}
					}
				}
			}
		}
	}
}

type StructInfo struct {
	Definition types.Object
	TypeName   *types.TypeName
	// Defined Type also called Named Type in go.
	// represent for a declaration such as:
	// type S struct { ... }
	Named      *types.Named
	Underlying *types.Struct
}

// Structs return iterator over each struct and its comments.
func (pkg *Package) Structs() iter.Seq2[StructInfo, Comments] {
	return func(yield func(StructInfo, Comments) bool) {
		for ident, def := range pkg.TypesInfo.Defs {
			if ident == nil || def == nil {
				continue
			}
			if tn, ok := def.(*types.TypeName); ok {
				if nm, ok := tn.Type().(*types.Named); ok {
					if st, ok := nm.Underlying().(*types.Struct); ok {
						info := StructInfo{
							Definition: def,
							TypeName:   tn,
							Named:      nm,
							Underlying: st,
						}
						if !yield(info, pkg.CommentsAt(ident.Pos())) {
							return
						}
					}
				}
			}
		}
	}
}

type FuncInfo struct {
	Definition types.Object
	Func       *types.Func
}

// Funcs return iterator over each funcs (including abstract ones) and its comments.
func (pkg *Package) Funcs() iter.Seq2[FuncInfo, Comments] {
	return func(yield func(FuncInfo, Comments) bool) {
		for ident, def := range pkg.TypesInfo.Defs {
			if ident == nil || def == nil {
				continue
			}
			if fn, ok := def.(*types.Func); ok {
				info := FuncInfo{
					Definition: def,
					Func:       fn,
				}
				if !yield(info, pkg.CommentsAt(ident.Pos())) {
					return
				}
			}
		}
	}
}

type Comments []string

func (c Comments) String() string {
	return strings.Join(c, "\n")
}

func (c Comments) Filter(f func(line string) bool) Comments {
	var nc Comments
	for i := range c {
		if f(c[i]) {
			nc = append(nc, c[i])
		}
	}
	return nc
}

// FilterPrefix can used to find deirectives like //go:linkname.
func (c Comments) FilterPrefix(prefix string) Comments {
	return c.Filter(func(s string) bool {
		return strings.HasPrefix(s, prefix)
	})
}

// At return comment at index, if index = -1, return the last one.
func (c Comments) At(index int) string {
	if index < 0 {
		return c[len(c)+index]
	}
	return c[index]
}
