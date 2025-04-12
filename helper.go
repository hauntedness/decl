package decl

import "go/types"

func TypePkg(typ types.Type) *types.Package {
	switch t := typ.(type) {
	case *types.Named:
		return t.Origin().Obj().Pkg()
	case *types.Pointer:
		return TypePkg(t.Elem())
	case *types.Basic:
		return nil
	case *types.Alias:
		return t.Obj().Pkg()
	default:
		return nil
	}
}

func TypeName(typ types.Type) string {
	switch t := typ.(type) {
	case *types.Named:
		return t.Obj().Name()
	case *types.Pointer:
		return TypeName(t.Elem())
	case *types.Basic:
		return t.Name()
	case *types.Alias:
		return t.Obj().Name()
	default:
		return ""
	}
}
