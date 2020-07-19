package desultory

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileOrDirectoryExists(t *testing.T) {
	assert := assert.New(t)
	tp, err := getTestDirectory()
	assert.NoError(err)
	dp := path.Join(tp, "foo")
	err = os.Mkdir(dp, 0777)
	de, err := DirectoryOrFileExists(dp)
	assert.NoError(err)
	assert.True(de)
	fp := path.Join(tp, "bar.txt")
	f, err := os.Create(fp)
	assert.NoError(err)
	f.WriteString("multiply world!")
	f.Close()
	fe, err := DirectoryOrFileExists(fp)
	assert.NoError(err)
	assert.True(fe)
	np := path.Join(tp, "baz.txt")
	ne, err := DirectoryOrFileExists(np)
	assert.NoError(err)
	assert.False(ne)
}

func TestDeleteDirectoryContents(t *testing.T) {
	assert := assert.New(t)
	tp, err := getTestDirectory()
	assert.NoError(err)
	dp := path.Join(tp, "foo")
	err = os.Mkdir(dp, 0777)
	de, err := DirectoryOrFileExists(dp)
	assert.NoError(err)
	assert.True(de)
	fp := path.Join(dp, "bar.txt")
	f, err := os.Create(fp)
	assert.NoError(err)
	f.WriteString("multiply world!")
	f.Close()
	fe, err := DirectoryOrFileExists(fp)
	assert.NoError(err)
	assert.True(fe)
	err = DeleteDirectoryContents(dp)
	assert.NoError(err)
	ne, err := DirectoryOrFileExists(fp)
	assert.NoError(err)
	assert.False(ne)
}

func TestSerializeDeserializeObject(t *testing.T) {
	assert := assert.New(t)
	tp, err := getTestDirectory()
	assert.NoError(err)
	f := "foo.yml"
	o := &TestStruct{
		Text: "ok computer",
		Number: 42,
	}
	err = SerializeObject(o, f, tp)
	assert.NoError(err)
	p := &TestStruct{}
	err = DeserializeObject(p, f, tp)
	assert.NoError(err)
	assert.Equal(o.Text, p.Text)
	assert.Equal(o.Number, p.Number)
}