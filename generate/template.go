package generate

// 생성자 함수를 만들기 위한 템플릿을 정의합니다.
var requiredArgsConstructorTmpl = `func New{{.StructName}}WithRequiredArgs({{range $index, $element := .Fields}}{{if $index}}, {{end}}{{$element.Name}} {{$element.Type}}{{end}}) {{.StructName}} {
    return {{.StructName}}{
		{{range .Fields}}{{.Name}}: {{.Name}},
		{{end}}
    }
}
`

// 생성자 함수를 만들기 위한 템플릿을 정의합니다.
var allArgsConstructorTemplate = `func New{{.StructName}}WithAllArgs({{range $index, $element := .Fields}}{{if $index}}, {{end}}{{$element.Name}} {{$element.Type}}{{end}}) {{.StructName}} {
    return {{.StructName}}{
        {{range .Fields}}{{.Name}}: {{.Name}},
		{{end}}
    }
}
`

var noArgsConstructorTemplate = `func New{{.StructName}}WithNoArgs() {{.StructName}} {
    return {{.StructName}}{}
}
`

// Builder 패턴을 위한 템플릿을 정의합니다.
var builderTemplate = `
// {{.StructName}}Builder is a builder for {{.StructName}}
type {{.StructName}}Builder struct {
    target *{{.StructName}}
}

{{range .Fields}}
// Set{{.Name}} sets the {{.Name}} field of the target {{$.StructName}}
func (b *{{$.StructName}}Builder) With{{.Name}}(value {{.Type}}) *{{$.StructName}}Builder {
    b.target.{{.Name}} = value

    return b
}
{{end}}

// Build constructs a {{.StructName}} from the builder
func (b *{{.StructName}}Builder) Build() *{{.StructName}} {
    return b.target
}

// New{{.StructName}}Builder creates a new builder instance for {{.StructName}}
func New{{.StructName}}Builder() *{{.StructName}}Builder {
    return &{{.StructName}}Builder{target: &{{.StructName}}{}}
}
`

var toStringTemplate = `
func (s {{.StructName}}) String() string {
	return fmt.Sprintf("{{.StructName}}{ {{range $index, $element := .Fields}}{{if $index}}, {{end}}{{.Name}}: %v{{end}} }", {{range $index, $element := .Fields}}{{if $index}}, {{end}}s.{{.Name}}{{end}})
}
`

var equalsTemplate = `
func (s {{.StructName}}) Equals(other {{.StructName}}) bool {
	return reflect.DeepEqual(s, other)
}
`

var getterTemplate = `
{{range .Fields}}
func (s {{$.StructName}}) Get{{.Name}}() {{.Type}} {
	return s.{{.Name}}
}
{{end}}
`

var setterTemplate = `
{{range .Fields}}
func (s *{{$.StructName}}) Set{{.Name}}(value {{.Type}}) {
	s.{{.Name}} = value
}
{{end}}
`
