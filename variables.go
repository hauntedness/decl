package decl

import (
	"fmt"
	"go/types"
)

type Var struct {
	v       *types.Var
	comment *Comment
}

func NewVar(obj types.Object) (*Var, error) {
	if v, ok := obj.(*types.Var); ok {
		return &Var{v: v}, nil
	} else {
		return nil, fmt.Errorf("obj is not a variable: %v", obj)
	}
}

func (v *Var) Type() types.Type {
	return v.v.Type()
}

func (v *Var) Name() string {
	return v.v.Name()
}

func (v *Var) TypeName() string {
	return TypeName(v.Type())
}

func (v *Var) TypePkg() *types.Package {
	return TypePkg(v.Type())
}

// 对于func的参数和返回值来说,只应该接受部分类型, 其他一律不处理
func (v *Var) TypeKind() TypeKind {
	return TypeKindOf(v.v.Type())
}

func (v *Var) Comment() *Comment {
	return v.comment
}
