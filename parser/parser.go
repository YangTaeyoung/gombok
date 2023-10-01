package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

	filepkg "github.com/YangTaeyoung/gombok/file"
	"github.com/YangTaeyoung/gombok/generate"
)

func Run() {
	root, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		fmt.Println(info.Name())

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "", content, parser.ParseComments)
		if err != nil {
			return err
		}

		importPkgs := make([]filepkg.ImportPackage, 0)
		ast.Inspect(file, func(n ast.Node) bool {

			// 주석을 찾는다.
			switch x := n.(type) {
			case *ast.ImportSpec:
				if x.Path != nil {
					importPath := strings.Trim(x.Path.Value, "\"")

					// 중복을 방지하기 위해 이미 importPkgs에 포함되어 있는지 확인
					alreadyIncluded := false
					for _, pkg := range importPkgs {
						if pkg.Path == importPath {
							alreadyIncluded = true
							break
						}
					}
					if !alreadyIncluded {
						var alias string
						if x.Name != nil {
							alias = x.Name.Name
						}

						importPkgs = append(importPkgs, filepkg.ImportPackage{
							Alias: alias,
							Path:  importPath,
						})
					}
				}
			case *ast.GenDecl:
				if x.Tok != token.TYPE {
					return true
				}

				for _, spec := range x.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}

					structType, ok := typeSpec.Type.(*ast.StructType)
					if !ok {
						continue
					}

					if x.Doc == nil {
						continue
					}

					var (
						fileContent       string
						requireReflectPkg bool
						result            string
					)

					for _, comment := range x.Doc.List {
						if strings.Contains(comment.Text, "@AllArgsConstructor") {
							log.Printf("Found @AllArgsConstructor in %s\n", typeSpec.Name.Name)
							result, err = generate.AllArgsConstructor(typeSpec.Name.Name, structType.Fields.List)
							if err != nil {
								log.Println("Error generating AllArgsConstructor:", err)
								continue
							}

							fileContent = fileContent + result
						}
						if strings.Contains(comment.Text, "@RequiredArgsConstructor") {
							log.Printf("Found @RequiredArgsConstructor in %s\n", typeSpec.Name.Name)
							result, err = generate.RequiredArgsConstructor(typeSpec.Name.Name, structType.Fields.List)
							if err != nil {
								log.Println("Error generating RequiredArgsConstructor:", err)
								continue
							}

							fileContent = fileContent + result
						}
						if strings.Contains(comment.Text, "@NoArgsConstructor") {
							log.Printf("Found @NoArgsConstructor in %s\n", typeSpec.Name.Name)
							result, err = generate.NoArgsConstructor(typeSpec.Name.Name)
							if err != nil {
								log.Println("Error generating NoArgsConstructor:", err)
								continue
							}

							fileContent = fileContent + result
						}
						if strings.Contains(comment.Text, "@Builder") {
							log.Printf("Found @Builder in %s\n", typeSpec.Name.Name)
							result, err = generate.Builder(typeSpec.Name.Name, structType.Fields.List)
							if err != nil {
								log.Println("Error generating Builder:", err)
								continue
							}

							fileContent = fileContent + result
						}
						if strings.Contains(comment.Text, "@ToString") {
							log.Printf("Found @ToString in %s", typeSpec.Name.Name)
							result, err = generate.ToString(typeSpec.Name.Name, structType.Fields.List)
							if err != nil {
								log.Println("Error generating ToString:", err)
								continue
							}

							fileContent = fileContent + result
						}
						if strings.Contains(comment.Text, "@Equals") {
							log.Printf("Found @Equals in %s", typeSpec.Name.Name)
							result, err = generate.Equals(typeSpec.Name.Name)
							if err != nil {
								log.Println("Error generating Equals:", err)
								continue
							}

							requireReflectPkg = true
							fileContent = fileContent + result
						}
						if strings.Contains(comment.Text, "@Getter") {
							log.Printf("Found @Getter in %s", typeSpec.Name.Name)
							result, err = generate.Getter(typeSpec.Name.Name, structType.Fields.List)
							if err != nil {
								log.Println("Error generating Getter:", err)
								continue
							}

							fileContent = fileContent + result
						}
						if strings.Contains(comment.Text, "@Setter") {
							log.Printf("Found @Setter in %s", typeSpec.Name.Name)
							result, err = generate.Setter(typeSpec.Name.Name, structType.Fields.List)
							if err != nil {
								log.Println("Error generating Setter:", err)
								continue
							}

							fileContent = fileContent + result
						}
					}

					if fileContent != "" {
						newFileName := strings.TrimSuffix(info.Name(), ".go") + "_gombok.go"
						newFilePath := filepath.Join(filepath.Dir(path), newFileName)
						if requireReflectPkg {
							importPkgs = append(importPkgs, filepkg.ImportPackage{
								Path: "reflect",
							})
						}

						if err = filepkg.WriteFile(file.Name.Name, importPkgs, fileContent, newFilePath); err != nil {
							log.Println("Error writing file:", err)
							continue
						}
					}
				}
			}

			return true
		})

		return nil
	})

	if err != nil {
		fmt.Println("Error processing files:", err)
	}

}
