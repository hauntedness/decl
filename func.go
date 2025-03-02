package decl

import (
	"fmt"
	"go/types"
)

type Func struct {
	fn      *types.Func
	sig     *types.Signature
	comment *Comment
}

func NewFunc(obj types.Object) (*Func, error) {
	if obj == nil {
		return nil, fmt.Errorf("obj is nil")
	}
	if fn, ok := obj.(*types.Func); ok {
		return &Func{fn: fn, sig: fn.Type().(*types.Signature)}, nil
	} else {
		return nil, fmt.Errorf("obj is not a Func: %v", obj)
	}
}

func (f *Func) Kind() string {
	return "func"
}

func (f *Func) Params() *types.Tuple {
	return f.sig.Params()
}

func (f *Func) Param(i int) (*types.Var, error) {
	res := f.sig.Params()
	if res == nil || res.Len() < i+1 {
		return nil, fmt.Errorf("Func params length error for index=%d", i)
	}
	return res.At(i), nil
}

func (f *Func) ParamType(i int) (types.Type, error) {
	v, err := f.Param(i)
	if err != nil {
		return nil, err
	}
	return v.Type(), nil
}

func (f *Func) Recv() *types.Var {
	return f.sig.Recv()
}

func (f *Func) RecvTypeParams() *types.TypeParamList {
	return f.sig.RecvTypeParams()
}

func (f *Func) Results() *types.Tuple {
	return f.sig.Results()
}

func (f *Func) Result(i int) (*types.Var, bool) {
	res := f.sig.Results()
	if res == nil || res.Len() < i+1 {
		return nil, false
	}
	return res.At(i), true
}

func (f *Func) ResultType(i int) (types.Type, bool) {
	v, ok := f.Result(i)
	if !ok {
		return nil, false
	}
	return v.Type(), true
}

func (f *Func) ReturnError() bool {
	results := f.Results()
	if results.Len() != 2 {
		return false
	}
	if typ := results.At(1).Type(); IsError(typ) {
		return false
	}
	return false
}

func (f *Func) Signature() *types.Signature {
	return f.sig
}

func (f *Func) TypeParams() *types.TypeParamList {
	return f.sig.TypeParams()
}

func (f *Func) Underlying() types.Type {
	return f.sig.Underlying()
}

func (f *Func) Variadic() bool {
	return f.sig.Variadic()
}

func (f *Func) Unwrap() *types.Func {
	return f.fn
}

func (f *Func) Exported() bool {
	return f.fn.Exported()
}

func (f *Func) FullName() string {
	return f.fn.FullName()
}

func (f *Func) Id() string {
	return f.fn.Id()
}

func (f *Func) Name() string {
	return f.fn.Name()
}

func (f *Func) Parent() *types.Scope {
	return f.fn.Parent()
}

func (f *Func) Pkg() *types.Package {
	return f.fn.Pkg()
}

func (f *Func) String() string {
	return f.fn.String()
}

func (f *Func) Comment() *Comment {
	return f.comment
}
