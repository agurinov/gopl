package py

import (
	"fmt"
	"path/filepath"

	"github.com/agurinov/gopl/py/cpy"
)

const MinorVersion = cpy.MinorVersion

func VenvPath(dir string) string {
	return filepath.Join(
		dir,
		fmt.Sprintf("lib/python%s/site-packages", MinorVersion),
	)
}
