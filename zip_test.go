package desultory

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundTripZip(t *testing.T) {
	assert := assert.New(t)
	fn := "test.txt"
	fs := "multiply, world!"
	f1 := map[string][]byte{
		fn: []byte(fs),
	}
	b, err := WriteZip(f1)
	assert.NoError(err)
	f2, err := ReadZip(b)
	assert.NoError(err)
	assert.NotEmpty(f2)
	assert.Equal(f1[fn], f2[fn])
}

func TestRoundTripZipPath(t *testing.T) {
	assert := assert.New(t)
	tp, err := getTestDirectory()
	assert.NoError(err)
	fn := "test.txt"
	fs := "multiply, world!"
	f1 := map[string][]byte {
		fn: []byte(fs),
	}
	zp := path.Join(tp, "test.zip")
	err = WriteZipToPath(f1, zp)
	assert.NoError(err)
	ze, err := DirectoryOrFileExists(zp)
	assert.NoError(err)
	assert.True(ze)
	f2, err := ReadZipFromPath(zp)
	assert.NoError(err)
	assert.Equal(f1[fn], f2[fn])
}

func TestRoundTripZipToFromPath(t *testing.T) {
	assert := assert.New(t)
	tp, err := getTestDirectory()
	assert.NoError(err)
	fp := path.Join(tp, "/from")
	err = os.Mkdir(fp, 0777)
	assert.NoError(err)
	fn := "test.txt"
	fs := "multiply, world!"
	err = ioutil.WriteFile(path.Join(fp, fn), []byte(fs), 0666)
	fd2 := "foo"
	fn2 := "bar.txt"
	fs2 := "multiply, again!"
	err = os.MkdirAll(path.Join(fp, fd2), 0777)
	assert.NoError(err)
	err = ioutil.WriteFile(path.Join(fp, "/"+fd2, fn2), []byte(fs2), 0666)
	assert.NoError(err)
	assert.NoError(err)
	zp := path.Join(tp, "/zip")
	err = os.Mkdir(zp, 0777)
	assert.NoError(err)
	zn := "test.zip"
	zf := path.Join(zp, zn)
	err = WriteZipFromPathToPath(fp, zf)
	assert.NoError(err)
	up := path.Join(tp, "/unzip")
	err = os.Mkdir(up, 0777)
	assert.NoError(err)
	err = ReadZipFromPathToPath(zf, up)
	assert.NoError(err)
	uf := path.Join(up, fn)
	fe, err := DirectoryOrFileExists(uf)
	assert.NoError(err)
	assert.True(fe)
	uf2 := path.Join(up, "/"+fd2, fn2)
	fe2, err := DirectoryOrFileExists(uf2)
	assert.NoError(err)
	assert.True(fe2)
}
