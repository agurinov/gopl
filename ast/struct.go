package ast

import "go/ast"

func GetStructByName(structName string, node ast.Node) *ast.StructType {
	typeSpec, ok := node.(*ast.TypeSpec)
	if !ok {
		return nil
	}

	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return nil
	}

	if structName != typeSpec.Name.Name {
		return nil
	}

	return structType
}
