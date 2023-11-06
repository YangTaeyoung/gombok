package generate

import (
	"bytes"
	stringpkg "github.com/YangTaeyoung/gombok/strings"
	"go/ast"
	"go/printer"
	"go/token"
	"reflect"
	"strings"
	"text/template"
)

type Field struct {
	Name      string
	Type      string
	MustBuild bool
	IsPointer bool
}

type StructFields struct {
	StructName         string
	Fields             []Field
	DefaultConstructor bool
}

func exprToString(expr ast.Expr) string {
	var buf bytes.Buffer

	_ = printer.Fprint(&buf, token.NewFileSet(), expr)

	return buf.String()
}

func AllArgsConstructor(name string, fields []*ast.Field, isDefault bool) (string, error) {
	// 모든 필드를 리스트에 추가합니다.
	allFields := make([]Field, 0)
	for _, field := range fields {
		if field.Tag != nil {
			tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))

			// 필드에 constructor 태그가 있고 ignore로 정의되어 있다면 필드를 추가하지 않습니다.
			if value, exists := tag.Lookup("constructor"); exists && strings.Contains(value, "ignore") {
				continue
			}
		}

		// embedded 필드
		if field.Names == nil {
			allFields = append(allFields, Field{Name: exprToString(field.Type), Type: exprToString(field.Type)})
			continue
		}

		// 일반 필드
		for _, fieldName := range field.Names {
			allFields = append(allFields, Field{Name: fieldName.Name, Type: exprToString(field.Type)})
		}
	}

	// 템플릿 파싱.
	tmpl, err := template.New("allArgsConstructorTemplate").Funcs(template.FuncMap{
		"LowerCamelCase": stringpkg.LowerCamel,
		"ReceiverName":   stringpkg.ReceiverName,
	}).Parse(allArgsConstructorTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, StructFields{
		StructName:         name,
		Fields:             allFields,
		DefaultConstructor: isDefault,
	})

	if err != nil {
		return "", err
	}

	// 생성된 문자열을 반환
	return buf.String(), nil
}

func RequiredArgsConstructor(name string, fields []*ast.Field, isDefault bool) (string, error) {
	requiredFields := make([]Field, 0)

	for _, field := range fields {
		if field.Tag != nil {
			tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))

			// 필드에 constructor 태그가 있고 ignore로 정의되어 있다면 필드를 추가하지 않습니다.
			if value, exists := tag.Lookup("constructor"); exists && strings.Contains(value, "ignore") {
				continue
			}

			// 필드에 validate 태그가 있고, required로 정의되어 있다면 필드를 추가.
			if value, exists := tag.Lookup("validate"); exists && strings.Contains(value, "required") {
				// embedded 필드
				if field.Names == nil {
					requiredFields = append(requiredFields, Field{Name: exprToString(field.Type), Type: exprToString(field.Type)})
					continue
				}

				// 일반 필드
				requiredFields = append(requiredFields, Field{Name: field.Names[0].Name, Type: exprToString(field.Type)})
			}
		}
	}

	// 템플릿을 파싱합니다.
	tmpl, err := template.New("requiredArgsConstructor").Funcs(template.FuncMap{
		"LowerCamelCase": stringpkg.LowerCamel,
		"ReceiverName":   stringpkg.ReceiverName,
	}).Parse(requiredArgsConstructorTmpl)
	if err != nil {
		return "", err
	}

	// 템플릿에 데이터를 적용하여 문자열을 생성합니다.
	var buf bytes.Buffer

	err = tmpl.Execute(&buf, StructFields{
		StructName:         name,
		Fields:             requiredFields,
		DefaultConstructor: isDefault,
	})

	if err != nil {
		return "", err
	}

	// 생성된 문자열을 반환합니다.
	return buf.String(), nil
}

func NoArgsConstructor(name string, isDefault bool) (string, error) {
	tmpl, err := template.New("noArgsConstructorTemplate").Parse(noArgsConstructorTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, StructFields{
		StructName:         name,
		DefaultConstructor: isDefault,
	})

	if err != nil {
		return "", err
	}

	// 생성된 문자열을 반환합니다.
	return buf.String(), nil
}

func Builder(name string, fields []*ast.Field) (string, error) {
	allFields := make([]Field, 0)
	for _, field := range fields {
		if field.Tag != nil {
			tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))

			// 필드에 builder 태그가 있고 ignore로 정의되어 있다면 필드를 추가하지 않습니다.
			if value, exists := tag.Lookup("builder"); exists && strings.Contains(value, "ignore") {
				continue
			}

			// 필드에 builder 태그가 있고 must로 정의되어 있다면 필드를 추가합니다.
			if value, exists := tag.Lookup("builder"); exists && strings.Contains(value, "must") {
				var isPointer bool
				if strings.Contains(exprToString(field.Type), "*") {
					isPointer = true
				}
				// embedded 필드
				if field.Names == nil {
					allFields = append(allFields, Field{Name: exprToString(field.Type), Type: exprToString(field.Type), MustBuild: true, IsPointer: isPointer})
					continue
				}

				// 일반 필드
				for _, fieldName := range field.Names {
					allFields = append(allFields, Field{Name: fieldName.Name, Type: exprToString(field.Type), MustBuild: true, IsPointer: isPointer})
				}
				continue
			}
		}

		// embedded 필드
		if field.Names == nil {
			allFields = append(allFields, Field{Name: exprToString(field.Type), Type: exprToString(field.Type)})
			continue
		}

		// 일반 필드
		for _, fieldName := range field.Names {
			allFields = append(allFields, Field{Name: fieldName.Name, Type: exprToString(field.Type)})
		}
	}

	tmpl, err := template.New("z").Funcs(template.FuncMap{
		"LowerCamelCase": stringpkg.LowerCamel,
		"ReceiverName":   stringpkg.ReceiverName,
	}).Parse(builderTemplate)
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
	allFields := make([]Field, 0)
	for _, field := range fields {
		if field.Tag != nil {
			tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))

			// 필드에 to_string 태그가 있고 ignore로 정의되어 있다면 필드를 추가하지 않습니다.
			if value, exists := tag.Lookup("to_string"); exists && strings.Contains(value, "ignore") {
				continue
			}
		}

		// embedded 필드
		if field.Names == nil {
			allFields = append(allFields, Field{Name: exprToString(field.Type), Type: exprToString(field.Type)})
			continue
		}
		// 일반 필드
		for _, fieldName := range field.Names {
			allFields = append(allFields, Field{Name: fieldName.Name, Type: exprToString(field.Type)})
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

func Equals(name string) (string, error) {
	tmpl, err := template.New("equalsTemplate").Funcs(template.FuncMap{
		"LowerCamelCase": stringpkg.LowerCamel,
		"ReceiverName":   stringpkg.ReceiverName,
	}).Parse(equalsTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, StructFields{
		StructName: name,
	})

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func Getter(name string, fields []*ast.Field) (string, error) {
	allFields := make([]Field, 0)
	for _, field := range fields {
		if field.Tag != nil {
			tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))

			// 필드에 getter 태그가 있고 ignore로 정의되어 있다면 필드를 추가하지 않습니다.
			if value, exists := tag.Lookup("getter"); exists && strings.Contains(value, "ignore") {
				continue
			}
		}

		// embedded 필드
		if field.Names == nil {
			allFields = append(allFields, Field{Name: exprToString(field.Type), Type: exprToString(field.Type)})
			continue
		}

		// 일반 필드
		for _, fieldName := range field.Names {
			allFields = append(allFields, Field{Name: fieldName.Name, Type: exprToString(field.Type)})
		}
	}

	tmpl, err := template.New("getterTemplate").Funcs(template.FuncMap{
		"LowerCamelCase": stringpkg.LowerCamel,
		"ReceiverName":   stringpkg.ReceiverName,
	}).Parse(getterTemplate)
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
	allFields := make([]Field, 0)
	for _, field := range fields {
		if field.Tag != nil {
			tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))

			// 필드에 setter 태그가 있고 ignore로 정의되어 있다면 필드를 추가하지 않습니다.
			if value, exists := tag.Lookup("setter"); exists && strings.Contains(value, "ignore") {
				continue
			}
		}

		// embedded 필드
		if field.Names == nil {
			allFields = append(allFields, Field{Name: exprToString(field.Type), Type: exprToString(field.Type)})
			continue
		}

		// 일반 필드
		for _, fieldName := range field.Names {
			allFields = append(allFields, Field{Name: fieldName.Name, Type: exprToString(field.Type)})
		}
	}

	tmpl, err := template.New("setterTemplate").Funcs(template.FuncMap{
		"LowerCamelCase": stringpkg.LowerCamel,
		"ReceiverName":   stringpkg.ReceiverName,
	}).Parse(setterTemplate)
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
