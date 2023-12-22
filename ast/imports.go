package ast

import (
	"go/ast"
	"path/filepath"
	"strings"
)

func ParseImports(imports []*ast.ImportSpec) map[string]string {
	importsMap := make(map[string]string, len(imports))

	for _, importSpec := range imports {
		var (
			path  = importSpec.Path.Value
			alias = filepath.Base(path)
		)

		if importSpec.Name != nil {
			alias = importSpec.Name.Name
		}

		alias = strings.Trim(alias, "\"")
		path = strings.Trim(path, "\"")

		importsMap[alias] = path
	}

	return importsMap
}
