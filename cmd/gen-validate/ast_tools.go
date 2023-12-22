package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
)

type (
	rawTarget struct {
		typeSpec   *ast.TypeSpec
		structType *ast.StructType
		imports    []*ast.ImportSpec
	}
	rawTargets map[string]rawTarget
)

func (c Config) getStructTypes() (rawTargets, error) {
	var srcDir string

	switch c.sourcePkg {
	case c.pkg:
		srcDir = "."
	case "main":
		srcDir = c.pkg
	default:
		srcDir = filepath.Join(c.sourcePkg, c.pkg)
	}

	fset := token.NewFileSet()

	pkgs, _ := parser.ParseDir(fset, srcDir, func(info fs.FileInfo) bool {
		return !info.IsDir()
	}, parser.AllErrors)

	for pkgName, pkg := range pkgs {
		fmt.Println("GGGGHJ", pkgName, pkg.Imports, pkg.Files)
	}

	typesMap := make(rawTargets, len(c.types))
	for _, typ := range c.types {
		typesMap[typ] = rawTarget{}
	}

	filepath.Walk(srcDir, func(filePath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fset := token.NewFileSet()
		parsedAST, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
		if err != nil {
			return err
		}

		ast.Inspect(parsedAST, func(n ast.Node) bool {
			typeSpec, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				return true
			}

			structName := typeSpec.Name.Name

			if _, desired := typesMap[structName]; !desired {
				delete(typesMap, structName)

				return true
			}

			typesMap[structName] = rawTarget{
				structType: structType,
				typeSpec:   typeSpec,
				imports:    parsedAST.Imports,
			}

			return true
		})

		return nil
	})

	if len(typesMap) != len(c.types) {
		return nil, fmt.Errorf("some structs was not found")
	}

	return typesMap, nil
}
