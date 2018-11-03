package codengine_test

import (
	"os"
	"testing"

	"github.com/wusphinx/codengine"
)

func TestCommon(t *testing.T) {
	var ce codengine.CodeEngine
	cwd, _ := os.Getwd()
	cwd = cwd + "/"
	ce.Exec(cwd, "util.go", cwd)
}
