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
	// package provider type information we need.
	*packages.Package
	// position map is used to find the associated ident.
	cmap map[token.Pos][]*ast.CommentGroup
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

func LoadMap(cfg *packages.Config, patterns ...string) (map[string]*Package, error) {
	if cfg == nil {
		cfg = &packages.Config{Mode: DefaultLoadMode}
	}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, err
	}

	var mp = make(map[string]*Package)
	for i := range pkgs {
		p, err := LoadPackage(pkgs[i])
		if err != nil {
			return mp, err
		}
		mp[p.PkgPath] = p
	}

	return mp, nil
}

// LoadPackage make sure pkg was loaded with correct LoadMode.
func LoadPackage(pkg *packages.Package) (*Package, error) {
	c := &Package{
		Package: pkg,
		cmap:    map[token.Pos][]*ast.CommentGroup{},
	}
	for _, syntax := range pkg.Syntax {
		ast.Inspect(syntax, func(n ast.Node) bool {
			switch nodeType := n.(type) {
			case *ast.GenDecl:
				for _, spec := range nodeType.Specs {
					switch specType := spec.(type) {
					case *ast.TypeSpec:
						c.cmap[specType.Name.Pos()] = []*ast.CommentGroup{nodeType.Doc}
					case *ast.ValueSpec:
						for _, name := range specType.Names {
							c.cmap[name.Pos()] = []*ast.CommentGroup{nodeType.Doc}
						}
					}
				}
			case *ast.Field:
				for _, ident := range nodeType.Names {
					c.cmap[ident.Pos()] = []*ast.CommentGroup{nodeType.Comment, nodeType.Doc}
				}
			}
			return true
		})
	}

	return c, nil
}

// CommentsRaw return raw comment group, including nil.
func (pkg *Package) CommentsRaw(pos token.Pos) []*ast.CommentGroup {
	return pkg.cmap[pos]
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

func (pkg *Package) DefinedTypes() iter.Seq2[*types.Named, Comments] {
	return pkg.NamedTypes()
}

// NamedTypes return iterator over each defined type and its comments.
func (pkg *Package) NamedTypes() iter.Seq2[*types.Named, Comments] {
	return func(yield func(*types.Named, Comments) bool) {
		for ident, def := range pkg.TypesInfo.Defs {
			if ident == nil || def == nil {
				continue
			}
			if tn, ok := def.(*types.TypeName); ok {
				if info, ok := tn.Type().(*types.Named); ok {
					if !yield(info, pkg.CommentsAt(ident.Pos())) {
						return
					}
				}
			}
		}
	}
}

type InterfaceInfo struct {
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

// LookupValue remove the prefix from 1st matched comment and return the remaining.
func (c Comments) LookupValue(prefix string) (string, bool) {
	return c.Lookup(CutPrefix(prefix))
}

// Lookup execute fn on comment and return the 1st matched result.
func (c Comments) Lookup(fn func(line string) (string, bool)) (string, bool) {
	for i := range c {
		remain, ok := fn(c[i])
		if ok {
			return remain, true
		}
	}
	return "", false
}

// Collect call fn to convert comment into literal 'key1=value1 key2=value2',
// then collect the key and value pairs into map.
// Collect return true if any line converted by fn.
func (c Comments) Collect(fn func(line string) (string, bool)) (map[string]string, bool) {
	mp, found := map[string]string{}, false
	for i := range c {
		if remain, ok := fn(c[i]); ok {
			for field := range strings.FieldsSeq(remain) {
				key, value, _ := strings.Cut(field, "=")
				mp[key] = value
			}
			found = true
		}
	}
	return mp, found
}

// At return comment at index, if index = -1, return the last one.
func (c Comments) At(index int) string {
	if index < 0 {
		return c[len(c)+index]
	}
	return c[index]
}

func HasPrefix(prefix string) func(string) bool {
	return func(line string) bool {
		return strings.HasPrefix(line, prefix)
	}
}

func CutPrefix(prefix string) func(string) (string, bool) {
	return func(line string) (string, bool) {
		return strings.CutPrefix(line, prefix)
	}
}
