package decl

import (
	"encoding/csv"
	"flag"
	"fmt"
	"go/types"
	"strings"
)

type ObjectKind int

const (
	KindStruct ObjectKind = iota
	KindInterface
	KindFunc
	KindVar
	KindOther
)

func Kind(obj types.Object) ObjectKind {
	switch t := obj.(type) {
	case *types.Func:
		return KindFunc
	case *types.TypeName:
		switch t.Type().Underlying().(type) {
		case *types.Struct:
			return KindStruct
		case *types.Interface:
			return KindInterface
		}
	case *types.Var:
		return KindVar
	}
	return KindOther
}

var errorType = types.Universe.Lookup("error").Type()

func IsError(obj types.Type) bool {
	return types.Identical(obj, errorType)
}

func IsErrorVar(v *types.Var) bool {
	return types.Identical(v.Type(), errorType)
}



func SearchDirectives(flagSet *flag.FlagSet, directives []string, patterns []string, name string, defaultValue string) string {
	var flag1 *flag.Flag
	for _, d := range directives {
		for _, prefix := range patterns {
			if args, ok := strings.CutPrefix(d, prefix); ok {
				reader := csv.NewReader(strings.NewReader(args))
				reader.Comma = ' '
				cmd, err := reader.Read()
				if err != nil {
					panic(fmt.Errorf("failed to read flags. flags:%s, err:%w", args, err))
				}
				err = flagSet.Parse(cmd)
				if err != nil {
					panic(err)
				}

				flag1 = flagSet.Lookup(name)
			}
		}
	}
	if flag1 == nil {
		return defaultValue
	} else if v := flag1.Value.String(); v != "" {
		return v
	} else if defaultValue != "" {
		return defaultValue
	} else {
		return flag1.DefValue
	}
}

func CamelCase(text string) string {
	if len(text) == 0 {
		return text
	} else if len(text) == 1 {
		return strings.ToLower(text)
	}
	return strings.ToLower(text[0:1]) + text[1:]
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

type TypeKind int

const (
	TypeKindStruct TypeKind = 1 << iota
	TypeKindPointer
	TypeKindInterface
	TypeKindFunc
	TypeKindBasic
	TypeKindNestedPointer
	TypeKindOther
	//
	TypeKindStructPointer = (TypeKindStruct | TypeKindPointer)
)

func TypeKindOf(typ types.Type) TypeKind {
	switch t := typ.(type) {
	case *types.Named:
		switch t.Underlying().(type) {
		case *types.Struct:
			return TypeKindStruct
		case *types.Interface:
			return TypeKindInterface
		case *types.Signature:
			return TypeKindFunc
		case *types.Basic:
			return TypeKindBasic
		case *types.Pointer:
			return TypeKindPointer
		default:
			return TypeKindOther
		}
	case *types.Pointer:
		k := TypeKindOf(t.Elem())
		if k&TypeKindPointer == TypeKindPointer {
			// drop other bits
			return TypeKindNestedPointer
		}
		return TypeKindPointer | k
	default:
		return TypeKindOther
	}
}

func IsStructPointerKind(k TypeKind) bool {
	return k&TypeKindStructPointer == TypeKindStructPointer
}

func IsInterfaceKind(k TypeKind) bool {
	return k == TypeKindInterface
}
