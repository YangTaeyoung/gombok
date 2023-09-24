package generate

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"reflect"
	"strings"
	"text/template"
)

type Field struct {
	Name string
	Type string // 필드의 타입을 문자열로 저장합니다.
}

type StructFields struct {
	StructName string
	Fields     []Field
}

func exprToString(expr ast.Expr) string {
	var buf bytes.Buffer

	_ = printer.Fprint(&buf, token.NewFileSet(), expr)

	return buf.String()
}

func AllArgsConstructor(name string, fields []*ast.Field) (string, error) {
	// 모든 필드를 리스트에 추가합니다.
	allFields := make([]Field, 0)
	for _, field := range fields {
		for _, fieldName := range field.Names {
			allFields = append(allFields, Field{Name: fieldName.Name, Type: exprToString(field.Type)})
		}
	}

	// 템플릿 파싱.
	tmpl, err := template.New("allArgsConstructorTemplate").Parse(allArgsConstructorTemplate)
	if err != nil {
		return "", err
	}

	// 템플릿 데이터를 적용 후 문자열 생성
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, StructFields{
		StructName: name,
		Fields:     allFields,
	})

	if err != nil {
		return "", err
	}

	// 생성된 문자열을 반환
	return buf.String(), nil
}

func RequiredArgsConstructor(name string, fields []*ast.Field) (string, error) {
	requiredFields := make([]Field, 0)

	for _, field := range fields {
		if field.Tag != nil {
			tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))

			// 필드에 validate 태그가 있고, required로 정의되어 있다면 필드를 추가.
			if value, exists := tag.Lookup("validate"); exists && strings.Contains(value, "required") {
				requiredFields = append(requiredFields, Field{Name: field.Names[0].Name, Type: exprToString(field.Type)})
			}
		}
	}

	// 템플릿을 파싱합니다.
	tmpl, err := template.New("requiredArgsConstructor").Parse(requiredArgsConstructorTmpl)
	if err != nil {
		return "", err
	}

	// 템플릿에 데이터를 적용하여 문자열을 생성합니다.
	var buf bytes.Buffer

	err = tmpl.Execute(&buf, StructFields{
		StructName: name,
		Fields:     requiredFields,
	})

	if err != nil {
		return "", err
	}

	// 생성된 문자열을 반환합니다.
	return buf.String(), nil
}

func NoArgsConstructor(name string) (string, error) {
	// 템플릿을 파싱합니다.
	tmpl, err := template.New("noArgsConstructorTemplate").Parse(noArgsConstructorTemplate)
	if err != nil {
		return "", err
	}

	// 템플릿에 데이터를 적용하여 문자열을 생성합니다.
	var buf bytes.Buffer

	err = tmpl.Execute(&buf, StructFields{
		StructName: name,
	})

	if err != nil {
		return "", err
	}

	// 생성된 문자열을 반환합니다.
	return buf.String(), nil
}

func Builder(name string, fields []*ast.Field) (string, error) {
	allFields := make([]Field, len(fields))
	for i, field := range fields {
		for _, fieldName := range field.Names {
			allFields[i] = Field{Name: fieldName.Name, Type: exprToString(field.Type)}
		}
	}

	tmpl, err := template.New("builderTemplate").Parse(builderTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, StructFields{
		StructName: name,
		Fields:     allFields,
	})

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func ToString(name string, fields []*ast.Field) (string, error) {
	allFields := make([]Field, len(fields))
	for i, field := range fields {
		for _, fieldName := range field.Names {
			allFields[i] = Field{Name: fieldName.Name, Type: exprToString(field.Type)}
		}
	}

	tmpl, err := template.New("toStringTemplate").Parse(toStringTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, StructFields{
		StructName: name,
		Fields:     allFields,
	})

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func Equals(name string, fields []*ast.Field) (string, error) {
	allFields := make([]Field, len(fields))
	for i, field := range fields {
		for _, fieldName := range field.Names {
			allFields[i] = Field{Name: fieldName.Name, Type: exprToString(field.Type)}
		}
	}

	tmpl, err := template.New("equalsTemplate").Parse(equalsTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, StructFields{
		StructName: name,
		Fields:     allFields,
	})

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func Getter(name string, fields []*ast.Field) (string, error) {
	allFields := make([]Field, len(fields))
	for i, field := range fields {
		for _, fieldName := range field.Names {
			allFields[i] = Field{Name: fieldName.Name, Type: exprToString(field.Type)}
		}
	}

	tmpl, err := template.New("getterTemplate").Parse(getterTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, StructFields{
		StructName: name,
		Fields:     allFields,
	})

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func Setter(name string, fields []*ast.Field) (string, error) {
	allFields := make([]Field, len(fields))
	for i, field := range fields {
		for _, fieldName := range field.Names {
			allFields[i] = Field{Name: fieldName.Name, Type: exprToString(field.Type)}
		}
	}

	tmpl, err := template.New("setterTemplate").Parse(setterTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, StructFields{
		StructName: name,
		Fields:     allFields,
	})

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
