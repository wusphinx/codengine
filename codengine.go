package codengine

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type CodeEngine struct {
	FileSet     *token.FileSet
	AstFile     *ast.File
	PackageName string
	DstArray    []string
	Tmpl        string
}

func (ce *CodeEngine) Exec(srcFileDir, srcFileName, dstFileDir string) {
	ce.FileSet = token.NewFileSet()
	ce.Tmpl = tmpl
	ce.ReadSrcFile(srcFileDir, srcFileName)
	ce.ExtractDst()

	for _, dst := range ce.DstArray {
		ce.GenerateCodeFile(dst, dstFileDir)
	}
}

func (ce *CodeEngine) ReadSrcFile(filePathDir, fileName string) {
	data, err := ioutil.ReadFile(filePathDir + fileName)
	if err != nil {
		println("ReadFile-err:%s", err.Error())
		os.Exit(1)
	}

	if AstFile, err := parser.ParseFile(ce.FileSet, fileName, data, parser.ParseComments); err != nil {
		println("ParseFile-err:%s", err.Error())
		os.Exit(1)
	} else {
		ce.AstFile = AstFile
	}

	ce.PackageName = ce.AstFile.Name.Name
}

func (ce *CodeEngine) ExtractDst() {
	vistor := func(node ast.Node) bool {
		v, ok := node.(*ast.GenDecl)
		if !ok {
			return true
		}

		if v.Tok == token.IMPORT || v.Tok == token.VAR || v.Tok == token.CONST {
			return true
		}

		var dstNames []string
		for _, spec := range v.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				return true
			}
			if _, ok = typeSpec.Type.(*ast.StructType); !ok {
				return true
			}
			dstNames = append(dstNames, typeSpec.Name.Name)
		}
		ce.DstArray = dstNames
		return true
	}

	ce.walk(vistor)
}

func (ce *CodeEngine) walk(fn func(ast.Node) bool) {
	ast.Walk(walker(fn), ce.AstFile)
}

func (ce *CodeEngine) GenerateCodeFile(DstArray, dstFileDir string) {
	filePath := dstFileDir + strings.ToLower(DstArray) + fileSuffix
	fp, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer fp.Close()

	code, err := ce.RenderTemple(DstArray)
	if err != nil {
		println("RenderTemple-err:%s", err.Error())
		os.Exit(1)
	}

	if _, err := fp.Write(code); err != nil {
		println("Write-err:%s", err.Error())
		os.Exit(1)
	}

	fp.Sync()
	exec.Command("gofmt", "-w", filePath).CombinedOutput()
}

func (ce *CodeEngine) RenderTemple(dstName string) ([]byte, error) {
	Tmpl := template.Must(template.New(dstName).Parse(ce.Tmpl))
	var tmpBuf []byte
	buf := bytes.NewBuffer(tmpBuf)

	type m struct {
		PackageName string
		DstName     string
	}

	if err := Tmpl.Execute(buf, m{
		PackageName: ce.PackageName,
		DstName:     dstName,
	}); err != nil {
		println("Execute-err:%s", err.Error())
		return nil, err
	}

	return buf.Bytes(), nil
}

type walker func(ast.Node) bool

func (w walker) Visit(node ast.Node) ast.Visitor {
	if w(node) {
		return w
	}

	return nil
}
