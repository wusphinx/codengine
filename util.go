package codengine

const (
	fileSuffix = "_generated.go"
)

var tmpl = `/* Created by CodeEngine - DO NOT EDIT. */

package {{.PackageName}} 

func (*{{.DstName}}) GetName() string {
	return "{{.DstName}}"
}
`

type HelloWorld struct{}
