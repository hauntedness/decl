package decl

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"

	"golang.org/x/tools/go/packages"
)

func TestComment(t *testing.T) {
	pkgs, err := packages.Load(&packages.Config{Mode: math.MaxInt}, reflect.TypeFor[json.Decoder]().PkgPath())
	if err != nil {
		t.Fatal(err)
	}
	pkg, err := Load(pkgs[0])
	if err != nil {
		t.Fatal(err)
	}
	if len(pkg.funcs) == 0 {
		t.Fatal("no function found.")
	}
}
