package decl

import (
	"fmt"
	"go/types"
	"reflect"
)

type Struct struct {
	structType *types.Struct
	typeName   *types.TypeName
	named      *types.Named
	fields     []*Var
	comment    *Comment
}

func NewStruct(obj types.Object) (*Struct, error) {
	t, ok := obj.(*types.TypeName)
	if !ok {
		return nil, fmt.Errorf("obj is not a struct: %v", obj)
	}
	n, ok := t.Type().(*types.Named)
	if !ok {
		return nil, fmt.Errorf("obj is not a named type: %v", obj)
	}
	s, ok := n.Underlying().(*types.Struct)
	if !ok {
		return nil, fmt.Errorf("obj is not a struct: %v", obj)
	}
	var fields []*Var
	for f := range s.Fields() {
		fields = append(fields, &Var{v: f})
	}
	return &Struct{typeName: t, named: n, structType: s, fields: fields}, nil
}

func (s *Struct) Unwrap() *types.Struct {
	return s.structType
}

func (s *Struct) UnwrapNamed() *types.Named {
	return s.named
}

func (s *Struct) Fields() []*Var {
	return s.fields
}

func (s *Struct) Tag(i int) reflect.StructTag {
	return reflect.StructTag(s.structType.Tag(i))
}

func (s *Struct) Pkg() *types.Package {
	return s.typeName.Pkg()
}

func (s *Struct) String() string {
	return s.typeName.String()
}

func (s *Struct) Name() string {
	return s.typeName.Name()
}

func (s *Struct) ParentScope() *types.Scope {
	return s.typeName.Parent()
}

func (s *Struct) Type() types.Type {
	return s.typeName.Type()
}

func (s *Struct) Exported() bool {
	return s.typeName.Exported()
}

func (s *Struct) Method(i int) *types.Func {
	return s.named.Method(i)
}

func (s *Struct) NumMethods() int {
	return s.named.NumMethods()
}

func (s *Struct) Origin() *types.Named {
	return s.named.Origin()
}

func (s *Struct) TypeArgs() *types.TypeList {
	return s.named.TypeArgs()
}

func (s *Struct) TypeParams() *types.TypeParamList {
	return s.named.TypeParams()
}
