package desultory

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
	"strings"
)

func ParseGoPackagesInDirectory(directory string) (map[string]*ast.Package, error) {
	pkgs, err := parser.ParseDir(token.NewFileSet(), directory, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	return pkgs, nil
}

func ParseGoPackages(files map[string][]byte, directory string) (map[string]*ast.Package, error) {
	err := WriteFilesToDirectory(files, directory)
	if err != nil {
		return nil, err
	}
	pkgs, err := ParseGoPackagesInDirectory(directory)
	if err != nil {
		return nil, err
	}
	return pkgs, nil
}

func TraverseGoPackage(pkg *ast.Package, handler func(ast.Node, interface{}) error, context interface{}) error {
	var err error
	ast.Inspect(pkg, func(node ast.Node) bool {
		err = handler(node, context)
		if err != nil {
			return false
		}
		return true
	})
	if err != nil {
		return err
	}
	return nil
}

func GetGoPathDirectory(name string, domain string) (string, error) {
	gp := os.Getenv("GOPATH")
	if gp == "" {
		return "", fmt.Errorf("couldn't create go project directory because 'GOPATH' environment variable wasn't set")
	}
	fn := strings.ToLower(name)
	pd := path.Join(gp, "/src/", domain,  "/", fn)
	if _, err := os.Stat(pd); os.IsNotExist(err) {
		err = os.Mkdir(pd, os.ModePerm)
		if err != nil {
			return "", err
		}
	} else {
		err = DeleteDirectoryContents(pd)
		if err != nil {
			return "", err
		}
	}
	return pd, nil
}

func InstallGoFunctionDependencies(dependencies []string, directory string) error {
	err := RunCommand("govendor", []string{"init"}, nil, directory)
	if err != nil {
		return err
	}
	for _, d := range dependencies {
		err := RunCommand("govendor", []string{"fetch", d}, nil, directory)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetGoNodeCode(node ast.Node) string {
	b := &strings.Builder{}
	printer.Fprint(b, token.NewFileSet(), node)
	return b.String()
}

