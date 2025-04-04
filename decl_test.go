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
