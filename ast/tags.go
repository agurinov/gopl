package ast

import (
	"go/ast"
	"strings"
)

func ParseStructTags(tags *ast.BasicLit) map[string]string {
	if tags == nil {
		return nil
	}

	tagsString := tags.Value
	tagsString = strings.Trim(tagsString, "`")

	tagsBlocks := strings.Split(tagsString, " ")
	tagsMap := make(map[string]string, len(tagsBlocks))

	for _, tag := range tagsBlocks {
		key, value, found := strings.Cut(tag, ":")
		if !found {
			continue
		}

		tagsMap[key] = strings.Trim(value, "\"")
	}

	return tagsMap
}
