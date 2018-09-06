package main

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

	"myitcv.io/gogenerate"
)

var fset = token.NewFileSet()

func inspectDir(wd string) {
	envFile, ok := os.LookupEnv(gogenerate.GOFILE)
	if !ok {
		log.Fatalf("env not correct; missing %v", gogenerate.GOFILE)
	}

	dirFiles, err := gogenerate.FilesContainingCmd(wd, optionGen)
	if err != nil {
		log.Fatalf("could not determine if we are the first file: %v", err)
	}

	if dirFiles == nil {
		log.Fatalf("cannot find any files containing the %v directive", optionGen)
	}

	if dirFiles[envFile] != 1 {
		log.Fatalf("expected a single occurrence of %v directive in %v. Got: %v", optionGen, envFile, dirFiles)
	}
}

func parseDir(dir string) {
	inspectDir(dir)

	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("unable to parse %v: %v", dir, err)
	}

	for _, pkg := range pkgs {
		log.Println("Package", pkg.Name)

		for filePath, file := range pkg.Files {
			buf := bytes.NewBufferString("")
			printer.Fprint(buf, fset, file)
			//source := buf.String()
			//log.Println(source)
			pos := fset.Position(file.Pos()).Offset

			log.Println(pos)
			if gogenerate.FileGeneratedBy(filePath, optionGen) {
				continue
			}

			var importPath []string
			for _, imp := range file.Imports {
				importPath = append(importPath, imp.Path.Value)
			}

			classList := make(map[string]bool)
			classOptionFields := make(map[string][]optionField)
			for _, d := range file.Decls {
				switch d := d.(type) {
				case *ast.FuncDecl:
					if d.Recv != nil {
						continue
					}

					// Find func which its Name match _<className>OptionDeclaration
					if strings.HasSuffix(d.Name.Name, optionDeclarationSuffix) && strings.HasPrefix(d.Name.Name, "_") {
						// Only allow return expr in class option declaration function
						if len(d.Body.List) != 1 {
							continue
						}

						stmt := d.Body.List[0]
						// Only allow return one value
						if stmt, ok := stmt.(*ast.ReturnStmt); !ok {
							continue
						} else {
							if len(stmt.Results) != 1 {
								continue
							}
							result := stmt.Results[0].(*ast.CompositeLit)
							optionFields := make([]optionField, len(result.Elts))
							for i, elt := range result.Elts {
								switch elt := elt.(type) {
								case *ast.KeyValueExpr:
									// Option Field Name
									key := elt.Key.(*ast.BasicLit)
									optionFields[i].Name = key.Value

									switch value := elt.Value.(type) {
									case *ast.FuncLit:
										optionFields[i].FieldType = FieldType_Func
										buf := bytes.NewBufferString("")
										// Option func Type
										printer.Fprint(buf, fset, value.Type)
										log.Println("Type:", buf.String())
										optionFields[i].Type = buf.String()

										// Option func Body
										buf.Reset()
										printer.Fprint(buf, fset, value.Body)
										log.Println("Body:", buf.String())
										optionFields[i].Body = buf.String()
									case *ast.CallExpr:
										optionFields[i].FieldType = FieldType_Var
										buf := bytes.NewBufferString("")

										// Option Variable Type
										printer.Fprint(buf, fset, value.Fun)
										log.Println("Type:", buf.String())
										optionFields[i].Type = buf.String()

										// Option Variable Value
										buf.Reset()
										printer.Fprint(buf, fset, value.Args[0])
										log.Println("Value:", buf.String())
										optionFields[i].Body = buf.String()
									}
								}
							}

							declarationClassName := strings.TrimPrefix(strings.TrimSuffix(d.Name.Name, optionDeclarationSuffix), "_")
							classOptionFields[declarationClassName] = optionFields
						}
					}
				case *ast.GenDecl:
					if d.Tok == token.TYPE {
						for _, spec := range d.Specs {
							if typeSpec, ok := spec.(*ast.TypeSpec); ok {
								classList[typeSpec.Name.Name] = false
							}
						}
					}
				}
			}

			for className, _ := range classOptionFields {
				if _, ok := classList[className]; !ok {
					log.Fatalf("Found %s class option declaration function, but not found class definition", className)
					delete(classOptionFields, className)
				} else {
					classList[className] = true
				}
			}

			for className, optionExist := range classList {
				if !optionExist {
					delete(classList, className)
				}
			}

			g := fileOptionGen{
				FilePath:          filePath,
				FileName:          strings.TrimSuffix(filepath.Base(filePath), ".go"),
				PkgName:           pkg.Name,
				ImportPath:        importPath,
				ClassList:         classList,
				ClassOptionFields: classOptionFields,
			}
			g.gen()
		}
	}
}
