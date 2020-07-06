package desultory

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

/*
func TestGoSimpleBuild(t *testing.T) {
	assert := assert.New(t)
	id, err := ioutil.TempDir("/tmp", "test-")
	defer os.RemoveAll(id)
	assert.NoError(err)
	m := "Hello, world!"
	c := `
package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

func main() {
    spew.Dump("` + m + `")
}`

	fs := map[string][]byte {
		"main.go": []byte(c),
	}
	err = WriteFilesToDirectory(fs, id)
	cp := NewGoCompiler()
	cc, err := cp.Compile(id)
	f := cc.Functions["main.main"]
	_, od, err := buildGoAwsLambdaFunction(f, "darwin", "amd64")
	assert.NoError(err)
	err = CaptureOutput()
	defer RestoreOutput()
	assert.NoError(err)
	err = RunCommand("main", nil, nil, od)
	assert.NoError(err)
	o, err := RestoreOutput()
	assert.NoError(err)
	assert.True(strings.Contains(o, m))
}
*/

func TestParseTraverseGoPackages(t *testing.T) {
	assert := assert.New(t)
	s := "hello world"
	ns := "main"
	fn := ns + ".go"
	fc := `package ` + ns + `
import "fmt"
func main() {
	s := "` + s + `"
	fmt.Println(s)
}`

	fs := map[string][]byte {fn: []byte(fc)}
	id, err := ioutil.TempDir("/tmp", "test-")
	defer os.RemoveAll(id)
	assert.NoError(err)
	pkgs, err := ParseGoPackages(fs, id)
	assert.NoError(err)
	ctx := struct{
		b bool
	} {
		b: false,
	}
	h := func(n ast.Node, _ interface{}) error {
		if l, ok := n.(*ast.BasicLit); ok {
			if strings.Contains(l.Value, s) {
				ctx.b = true
			}
		}
		return nil
	}
	err = TraverseGoPackage(pkgs[ns], h, ctx)
	assert.True(ctx.b)
}
