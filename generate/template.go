package generate

// 생성자 함수를 만들기 위한 템플릿을 정의합니다.
var requiredArgsConstructorTmpl = `
// New{{ if not .DefaultConstructor }}{{.StructName}}WithRequiredArgs{{end}}
func New{{ if not .DefaultConstructor }}{{.StructName}}WithRequiredArgs{{end}}({{range $index, $element := .Fields}}{{if $index}}, {{end}}{{LowerCamelCase $element.Name}} {{$element.Type}}{{end}}) {{.StructName}} {
    return {{.StructName}}{
		{{range .Fields}}{{.Name}}: {{LowerCamelCase .Name}},
		{{end}}
    }
}
`

// 생성자 함수를 만들기 위한 템플릿을 정의합니다.
var allArgsConstructorTemplate = `
// New{{ if not .DefaultConstructor }}{{.StructName}}WithAllArgs{{end}}
func New{{ if not .DefaultConstructor }}{{.StructName}}WithAllArgs{{end}}({{range $index, $element := .Fields}}{{if $index}}, {{end}}{{LowerCamelCase $element.Name}} {{$element.Type}}{{end}}) {{.StructName}} {
    return {{.StructName}}{
        {{range .Fields}}{{.Name}}: {{LowerCamelCase .Name}},
		{{end}}
    }
}
`

var noArgsConstructorTemplate = `
// New{{ if not .DefaultConstructor }}{{.StructName}}WithNoArgs{{end}}
func New{{ if not .DefaultConstructor }}{{.StructName}}WithNoArgs{{end}}() {{.StructName}} {
    return {{.StructName}}{}
}
`

// Builder 패턴을 위한 템플릿을 정의합니다.
var builderTemplate = `
// {{.StructName}}Builder
// a builder for {{.StructName}}
type {{.StructName}}Builder struct {
    target *{{.StructName}}
}

{{range .Fields}}
// With{{.Name}}
// sets the {{.Name}} field of the target {{$.StructName}}
func ({{ReceiverName $.StructName}}b {{$.StructName}}Builder) With{{.Name}}({{LowerCamelCase .Name}} {{.Type}}) {{$.StructName}}Builder {
	{{ if .MustBuild }}{{ if .IsPointer }}if {{LowerCamelCase .Name}} == nil {
		panic("{{$.StructName}}Builder: {{.Name}} must not be nil")
	}{{ else }}if reflect.DeepEqual({{LowerCamelCase .Name}}, {{.Type}}{}) {
		panic("{{$.StructName}}Builder: {{.Name}} must not be empty")
	}{{ end }}{{ end }}
    {{ReceiverName $.StructName}}b.target.{{.Name}} = {{LowerCamelCase .Name}}

    return {{ReceiverName $.StructName}}b
}
{{end}}

// Build
// constructs a {{.StructName}} from the builder
func ({{ReceiverName $.StructName}}b {{.StructName}}Builder) Build() {{.StructName}} {
    return *{{ReceiverName $.StructName}}b.target
}

// New{{.StructName}}Builder
// creates a new builder instance for {{.StructName}}
func New{{.StructName}}Builder() {{.StructName}}Builder {
    return {{.StructName}}Builder{target: &{{.StructName}}{}}
}
`

var toStringTemplate = `
// String
func ({{ReceiverName $.StructName}} *{{.StructName}}) String() string {
	return fmt.Sprintf("{{.StructName}}{ {{range $index, $element := .Fields}}{{if $index}}, {{end}}{{.Name}}: %v{{end}} }", {{range $index, $element := .Fields}}{{if $index}}, {{end}}{{ReceiverName $.StructName}}.{{.Name}}{{end}})
}
`

var equalsTemplate = `
// Equals
func ({{ReceiverName $.StructName}} *{{.StructName}}) Equals({{LowerCamelCase .StructName}} {{.StructName}}) bool {
	return reflect.DeepEqual({{ReceiverName $.StructName}}, {{LowerCamelCase .StructName}})
}
`

var getterTemplate = `
{{range .Fields}}
// Get{{.Name}}
func ({{ReceiverName $.StructName}} *{{$.StructName}}) Get{{.Name}}() {{.Type}} {
	return {{ReceiverName $.StructName}}.{{.Name}}
}
{{end}}
`

var setterTemplate = `
{{range .Fields}}
// Set{{.Name}}
func ({{ReceiverName $.StructName}} *{{$.StructName}}) Set{{.Name}}({{LowerCamelCase .Name}} {{.Type}}) {
	{{ReceiverName $.StructName}}.{{.Name}} = {{LowerCamelCase .Name}}
}
{{end}}
`
