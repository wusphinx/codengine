package codengine

const (
	fileSuffix = "_generated.go"
)

var tmpl = `package {{.PackageName}} 

func (*{{.DstName}}) GetName() string {
	return "{{.DstName}}"
}
`

type HelloWorld struct{}
