package decl

import (
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {
	pkg, err := Load("github.com/hauntedness/decl/internal/test")
	if err != nil {
		t.Fatal(err)
	}

	for obj, comments := range pkg.Definitions() {
		fmt.Println(obj, comments)
	}
}

func TestPackage_Structs(t *testing.T) {
	pkg, err := Load("github.com/hauntedness/decl/internal/test")
	if err != nil {
		t.Fatal(err)
	}

	for info, comments := range pkg.Structs() {
		fmt.Println("struct comments:", comments)
		for field := range info.Underlying.Fields() {
			fieldComments := pkg.CommentsAt(field.Pos())
			fmt.Println("field comments:", fieldComments)
		}
	}
}

func TestPackage_DefinedTypes(t *testing.T) {
	pkg, err := Load("github.com/hauntedness/decl/internal/test")
	if err != nil {
		t.Fatal(err)
	}

	for info, comments := range pkg.DefinedTypes() {
		fmt.Println("type:", info.Named.Obj().Name(), ";", "comments:", comments)
	}
}

func TestPackage_Interfaces(t *testing.T) {
	pkg, err := Load("github.com/hauntedness/decl/internal/test")
	if err != nil {
		t.Fatal(err)
	}

	for info, comments := range pkg.Interfaces() {
		fmt.Println("type:", info.Named.Obj().Name(), ";", "comments:", comments)
		for method := range info.Underlying.Methods() {
			fmt.Println("method comments:", pkg.CommentsAt(method.Pos()))
		}
	}
}

func TestPackage_Funcs(t *testing.T) {
	pkg, err := Load("github.com/hauntedness/decl/internal/test")
	if err != nil {
		t.Fatal(err)
	}

	for info, comments := range pkg.Funcs() {
		fmt.Println("func:", info.Func.Name(), ";", "comments:", comments)
		for method := range info.Func.Signature().Params().Variables() {
			fmt.Println("func param comments:", pkg.CommentsAt(method.Pos()))
		}
	}
}
