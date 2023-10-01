package file

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"text/template"
)

var fileTemplate = `// Code generated by gombok. DO NOT EDIT.
package {{.PackageName}}

{{- if len .ImportPackages }}
import (
{{- range .ImportPackages}}	
{{.Alias}} "{{.Path}}"
{{- end}}
)
{{- end}}

{{.Content}}
`

type ImportPackage struct {
	Alias string
	Path  string
}

type TemplateElement struct {
	PackageName    string
	ImportPackages []ImportPackage
	Content        string
}

func WriteFile(packageName string, importPackages []ImportPackage, content string, filepath string) error {
	tmpl, err := template.New("file").Parse(fileTemplate)
	if err != nil {
		return err
	}

	data := TemplateElement{
		PackageName:    packageName,
		ImportPackages: importPackages,
		Content:        content,
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, data)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath, buf.Bytes(), 0644)
	if err != nil {
		log.Printf("Error writing file %s: %v", filepath, err)
		return err
	}

	cmd := exec.Command("goimports", "-w", filepath)
	if err = cmd.Run(); err != nil {
		log.Printf("Error formatting file %s: %v", filepath, err)

		if err = os.Remove(filepath); err != nil {
			log.Printf("Error removing file %s: %v", filepath, err)
			return err
		}

		log.Printf("File %s removed", filepath)
		return err
	}

	return nil
}
