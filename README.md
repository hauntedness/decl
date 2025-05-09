# decl
A package to load comments in source code. Can be used as directive parser for generate tools.


```golang
package main

import (
	"flag"
	"fmt"
	"go/format"
	"os"
	"strings"
	"time"

	"github.com/hauntedness/decl"
	"github.com/hauntedness/std/hv"
	"github.com/valyala/fasttemplate"
)

//gen:str
type Book struct {
	name *string
	//gen:str --format="2016"
	year  *time.Time
	words uint
	sale  hv.Option[bool]
	//gen:str --ignore
	dirty map[any]any
}

//go:generate go run github.com/hauntedness/ebg/cmd/stringer

//nolint:funlen //+nolint
func main() {
	//
	pkg := flag.String("pkg", ".", "package to scan")
	writeTo := flag.String("w", "generated_struct_stringer.go", "package to scan")
	flag.Parse()
	//
	thepkg, err := decl.Load(*pkg)
	if err != nil {
		panic(err)
	}
	wb := &strings.Builder{}
	prefix := "//gen:str"
	for info, comment := range thepkg.Structs() {
		if _, ok := comment.LookupValue(prefix); !ok {
			continue
		}
		def := []string{}
		val := []string{}
		for f := range info.Underlying.Fields() {
			cmd, _ := thepkg.CommentsAt(f.Pos()).Collect(decl.CutPrefix(prefix))
			if _, ignore := cmd["--ignore"]; !ignore {
				fname := f.Name()
				def = append(def, fname+" string")
				typ := f.Type().String()
				switch typ {
				case "time.Time":
					format := cmd["--format"]
					tmpl := `
						v.{{fname}} = s.{{fname}}.Format({{fformat}})
					`
					val = append(val, Render(tmpl, map[string]any{"fname": fname, "fformat": format}))
				case "int", "int8", "int16", "int32", "int64":
					tmpl := `
						v.{{fname}} = strconv.FormatInt(int64(s.{{fname}}), 10)
					`
					val = append(val, Render(tmpl, map[string]any{"fname": fname}))
				case "uint", "uint8", "uint16", "uint32", "uint64":
					tmpl := `
						v.{{fname}} = strconv.FormatUint(uint64(s.{{fname}}), 10)
					`
					val = append(val, Render(tmpl, map[string]any{"fname": fname}))
				case "*string":
					tmpl := `
						if s.{{fname}} != nil {
							v.{{fname}} = *s.{{fname}}
						}
					`
					val = append(val, Render(tmpl, map[string]any{"fname": fname}))
				case "*time.Time":
					format := cmd["--format"]
					tmpl := `
						if s.{{fname}} != nil {
							v.{{fname}} = s.{{fname}}.Format({{fformat}})
						}
					`
					val = append(val, Render(tmpl, map[string]any{"fname": fname, "fformat": format}))
				default:
					if strings.HasPrefix(typ, "*") {
						tmpl := `
							if s.{{fname}} != nil {
								v.{{fname}} = fmt.Sprint(*s.{{fname}})
							}
						`
						val = append(val, Render(tmpl, map[string]any{"fname": fname}))
					} else if strings.HasPrefix(typ, "github.com/hauntedness/std/hv.Option") {
						tmpl := `
							if s.{{fname}}.IsPresent() {
								v.{{fname}} = fmt.Sprint(s.{{fname}}.MustGet())
							}
						`
						val = append(val, Render(tmpl, map[string]any{"fname": fname}))
					} else {
						tmpl := `
								v.{{fname}} = fmt.Sprint(s.{{fname}})
						`
						val = append(val, Render(tmpl, map[string]any{"fname": fname}))
					}
				}
			}
		}
		result := Render(body, map[string]any{
			"typeName": info.TypeName.Name(),
			"fieldDef": strings.Join(def, "; "),
			"fieldVal": strings.Join(val, "\n"),
		})
		data, err := format.Source([]byte(result))
		if err != nil {
			fmt.Println(result)
			panic(err)
		}
		wb.Write(data)
		wb.WriteString("\n")
	}
	if wb.Len() > 0 {
		file, err := os.Create(*writeTo)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		_, err = file.WriteString(Render(header, map[string]any{"package": thepkg.Package.Name}))
		if err != nil {
			panic(err)
		}
		_, err = file.WriteString(wb.String())
		if err != nil {
			panic(err)
		}
	}
}

func Render(template string, m map[string]any) string {
	return fasttemplate.ExecuteString(template, "{{", "}}", m)
}

var header = `// Code generated by cmd/stringer. DO NOT EDIT.
package {{package}}

import (
	"fmt"	
	"strconv"
	"time"

	"github.com/hauntedness/std/hv"
)

// force import 
var _ hv.Option[struct{}]
var _ time.Time
var _ = strconv.Atoi
`

var body = `
func (s *{{typeName}}) String() string {
	var v = struct {
		{{fieldDef}}
	} {}
	{{fieldVal}}
	return fmt.Sprintf("%+v", v)
}
`
```