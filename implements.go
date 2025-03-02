package decl

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"
)

// ImplementStmt denote that user explicitly specify type x implements interface y
//
// it is derived from ast with blow convention.
//  1. var _ TheInterface = (*TheImplement)(nil)
//  2. func (impl *TheImplement) Implements() TheInterface { return impl }
type ImplementStmt struct {
	namedInterface *types.Named
	namedImpl      *types.Named
	comment        *Comment
	isPointerImpl  bool
}

func (s *ImplementStmt) String() string {
	return fmt.Sprintf("var _ %s = (*%s)(nil)", s.namedInterface.String(), s.namedImpl.String())
}

func (s *ImplementStmt) IfacePkg() *types.Package {
	return s.namedInterface.Obj().Pkg()
}

func (s *ImplementStmt) ImplPkg() *types.Package {
	return s.namedImpl.Obj().Pkg()
}

func (s *ImplementStmt) ImplType() types.Type {
	return s.namedImpl.Obj().Type()
}

func (s *ImplementStmt) ImplName() string {
	return s.namedImpl.Obj().Name()
}

func (s *ImplementStmt) IfaceName() string {
	return s.namedInterface.Obj().Name()
}

func (s *ImplementStmt) IsPointerImpl() bool {
	return s.isPointerImpl
}

func AsInterface(_type ast.Expr, pkg *packages.Package) (*types.Named, bool) {
	var t types.Type
	uses := pkg.TypesInfo.Uses
	switch expr := _type.(type) {
	case *ast.SelectorExpr:
		obj := uses[expr.Sel]
		if obj == nil {
			return nil, false
		}
		t = obj.Type()
	case *ast.Ident:
		obj := uses[expr]
		if obj == nil {
			return nil, false
		}
		t = obj.Type()
	default:
		return nil, false
	}
	if t == nil {
		return nil, false
	}
	named, ok := t.(*types.Named)
	if !ok {
		return nil, false
	}
	_, ok = t.Underlying().(*types.Interface)
	return named, ok
}

func AsImpl(value ast.Expr, pkg *packages.Package) (*types.Named, bool) {
	var t types.Type
	uses := pkg.TypesInfo.Uses
	// var _ liu.Namer = (*guan.Guan)(nil)
	callExpr, ok := value.(*ast.CallExpr)
	if !ok {
		return nil, false
	}
	parenExpr, ok := callExpr.Fun.(*ast.ParenExpr)
	if !ok {
		return nil, false
	}
	switch expr := parenExpr.X.(type) {
	case *ast.StarExpr:
		// either selector or ident
		switch x := expr.X.(type) {
		case *ast.SelectorExpr:
			obj := uses[x.Sel]
			if obj == nil {
				return nil, false
			}
			t = obj.Type()
		case *ast.Ident:
			obj := uses[x]
			if obj == nil {
				return nil, false
			}
			t = obj.Type()
		default:
			return nil, false
		}
	default:
		return nil, false
	}
	named, ok := t.(*types.Named)
	if !ok {
		return nil, false
	}
	if len(callExpr.Args) != 1 {
		return nil, false
	}
	ident, ok := callExpr.Args[0].(*ast.Ident)
	if !ok {
		return nil, false
	}
	_, ok = uses[ident].(*types.Nil)
	return named, ok
}

func NewImplementStmtList(fn *Func, pkg *packages.Package) ([]*ImplementStmt, bool) {
	//
	if fn.fn.Name() != "Implements" {
		return nil, false
	}
	recv := fn.Recv()
	if recv == nil {
		return nil, false
	}
	cache := &ImplementStmt{}
	switch typ := recv.Type().(type) {
	case *types.Named:
		cache.namedImpl = typ
		cache.isPointerImpl = false
	case *types.Pointer:
		switch named := typ.Elem().(type) {
		case *types.Named:
			cache.namedImpl = named
			cache.isPointerImpl = true
		default:
			return nil, false
		}
	default:
		return nil, false
	}
	stmtList := make([]*ImplementStmt, 0, 1)
	tuple := fn.Results()
	for i := range tuple.Len() {
		var1 := tuple.At(i)
		ifaceNamed, ok := var1.Type().(*types.Named)
		if !ok {
			return nil, false
		}
		_, ok = ifaceNamed.Underlying().(*types.Interface)
		if !ok {
			return nil, false
		}
		stmt := &ImplementStmt{}
		stmt.namedImpl = cache.namedImpl
		stmt.isPointerImpl = cache.isPointerImpl
		stmt.namedInterface = ifaceNamed
		typeInterface := stmt.namedInterface.Underlying().(*types.Interface)
		var typeImpl types.Type = stmt.namedImpl
		if stmt.isPointerImpl {
			typeImpl = types.NewPointer(typeImpl)
		}
		if !types.Implements(typeImpl, typeInterface) {
			return nil, false
		}
		stmtList = append(stmtList, stmt)
	}
	return stmtList, true
}

func NewImplementStmt(spec *ast.ValueSpec, pkg *packages.Package) (*ImplementStmt, bool) {
	stmt := &ImplementStmt{}
	if len(spec.Names) != 1 {
		return nil, false
	}
	// var _ TheInterface = (*TheImplement)(nil)
	name := spec.Names[0]
	if name.Name != "_" {
		return nil, false
	}
	_type, ok := AsInterface(spec.Type, pkg)
	if !ok {
		return nil, false
	}
	stmt.namedInterface = _type
	value, ok := AsImpl(spec.Values[0], pkg)
	if !ok {
		return nil, false
	}
	stmt.namedImpl = value
	iface := _type.Underlying().(*types.Interface)
	if types.Implements(value, iface) {
		stmt.isPointerImpl = false
	} else if types.Implements(types.NewPointer(value), iface) {
		stmt.isPointerImpl = true
	} else {
		return nil, false
	}
	return stmt, true
}
