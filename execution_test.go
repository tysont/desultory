package desultory

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

var capturing bool
var standardOut *os.File
var standardError *os.File
var captureRead *os.File
var captureWrite *os.File
var captureChannel chan string

func CaptureOutput() error {
	if capturing {
		return errors.New("attempted to capture standard out while already capturing")
	}
	standardOut = os.Stdout
	standardError = os.Stderr
	var err error
	captureRead, captureWrite, err = os.Pipe()
	if err != nil {
		RestoreOutput()
		return err
	}
	capturing = true
	captureChannel = make(chan string)
	os.Stdout = captureWrite
	os.Stderr = captureWrite
	print()
	go func() {
		var b bytes.Buffer
		io.Copy(&b, captureRead)
		captureChannel <- b.String()
	}()
	return nil
}

func RestoreOutput() (string, error) {
	if !capturing {
		return "", errors.New("attempted to stop capturing standard out while not capturing")
	}
	if standardOut == nil {
		return "", errors.New("attempted to restore standard out to invalid state")
	} else if standardError == nil {
		return "", errors.New("attempted to restore standard out to invalid state")
	}
	capturing = false
	os.Stdout = standardOut
	os.Stderr = standardError
	captureWrite.Close()
	out := <- captureChannel
	standardOut = nil
	captureRead = nil
	captureWrite = nil
	captureChannel = nil
	return out, nil
}

func TestRunCommandLs(t *testing.T) {
	assert := assert.New(t)
	ifs := make(map[string][]byte, 0)
	fn := "test.txt"
	ifs[fn] = []byte("testing.")
	err := CaptureOutput()
	defer RestoreOutput()
	assert.NoError(err)
	ofs, err := RunCommandInTempDirectory("ls", nil, ifs)
	assert.NoError(err)
	o, err := RestoreOutput()
	assert.NoError(err)
	assert.NotEmpty(ofs)
	assert.True(strings.Contains(o, fn))
}

func TestRunCommandGoBuild(t *testing.T) {
	assert := assert.New(t)
	d, err := ioutil.TempDir("/tmp", "test-")
	defer os.RemoveAll(d)
	assert.NoError(err)
	s := "Hello, world!"
	c := `
package main

import "fmt"

func main() {
    fmt.Println("` + s + `")
}`

	ifs := map[string][]byte {
		"main.go": []byte(c),
	}
	err = WriteFilesToDirectory(ifs, d)
	assert.NoError(err)
	err = RunCommand("go", []string{"build", "-o", "hello"}, map[string]string {"GOPATH": d, "GOOS": "darwin", "GOARCH": "amd64",}, d)
	assert.NoError(err)
	err = CaptureOutput()
	defer RestoreOutput()
	assert.NoError(err)
	err = RunCommand("hello", nil, map[string]string {"PATH": d}, d)
	assert.NoError(err)
	o, err := RestoreOutput()
	assert.NoError(err)
	assert.True(strings.Contains(o, s))
}
