package decl

import (
	"fmt"
	"go/types"
	"iter"
)

type Interface struct {
	typeName       *types.TypeName
	named          *types.Named
	iface          *types.Interface
	methodComments []*Comment
	comment        *Comment
}

func NewInterface(obj types.Object) (*Interface, error) {
	t, ok := obj.(*types.TypeName)
	if !ok {
		return nil, fmt.Errorf("obj is not a struct: %v", obj)
	}
	n, ok := t.Type().(*types.Named)
	if !ok {
		return nil, fmt.Errorf("obj is not a named type: %v", obj)
	}
	i, ok := n.Underlying().(*types.Interface)
	if !ok {
		return nil, fmt.Errorf("obj is not a struct: %v", obj)
	}
	return &Interface{typeName: t, named: n, iface: i, methodComments: make([]*Comment, i.NumMethods())}, nil
}

func (t *Interface) Exported() bool {
	return t.typeName.Exported()
}

func (t *Interface) Name() string {
	return t.typeName.Name()
}

func (t *Interface) Parent() *types.Scope {
	return t.typeName.Parent()
}

func (t *Interface) Pkg() *types.Package {
	return t.typeName.Pkg()
}

func (t *Interface) String() string {
	return t.typeName.String()
}

func (t *Interface) Type() types.Type {
	return t.typeName.Type()
}

func (t *Interface) Complete() *types.Interface {
	return t.iface.Complete()
}

func (t *Interface) EmbeddedType(i int) types.Type {
	return t.iface.EmbeddedType(i)
}

func (t *Interface) Empty() bool {
	return t.iface.Empty()
}

func (t *Interface) ExplicitMethod(i int) *types.Func {
	return t.iface.ExplicitMethod(i)
}

func (t *Interface) IsComparable() bool {
	return t.iface.IsComparable()
}

func (t *Interface) IsImplicit() bool {
	return t.iface.IsImplicit()
}

func (t *Interface) IsMethodSet() bool {
	return t.iface.IsMethodSet()
}

func (t *Interface) Method(i int) *types.Func {
	return t.iface.Method(i)
}

func (t *Interface) MethodComments(i int) []*Comment {
	return t.methodComments
}

func (t *Interface) MethodsWithComment() iter.Seq2[*types.Func, *Comment] {
	return func(yield func(*types.Func, *Comment) bool) {
		for i := range t.iface.NumMethods() {
			if !yield(t.Method(i), t.methodComments[i]) {
				return
			}
		}
	}
}

func (t *Interface) NumEmbeddeds() int {
	return t.iface.NumEmbeddeds()
}

func (t *Interface) NumExplicitMethods() int {
	return t.iface.NumExplicitMethods()
}

func (t *Interface) NumMethods() int {
	return t.iface.NumMethods()
}

func (t *Interface) Underlying() types.Type {
	return t.iface.Underlying()
}
