package desultory

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateUpdateGetDeleteFaunaDatabase(t *testing.T) {
	assert := assert.New(t)
	sn := GetUniqueString(6)
	db := "entity-" + sn
	err := CreateFaunaDatabase(db)
	defer DeleteFaunaDatabase(db)
	assert.NoError(err)
	cn := "collection-" + sn
	err = CreateFaunaCollection(db, cn)
	assert.NoError(err)
	txt := "foo"
	nm := 7
	ts1 := &TestStruct{
		Text: txt,
		Number: nm,
	}
	in := "index-" + sn
	pk := "Text"
	err = CreateFaunaIndex(db, cn, in, pk)
	assert.NoError(err)
	err = CreateFaunaInstance(db, cn, ts1)
	ts2 := &TestStruct{
		Text: "bar",
		Number: 9,
	}
	err = CreateFaunaInstance(db, cn, ts2)
	assert.NoError(err)
	ts1.Number = nm * 2
	err = UpdateFaunaInstance(db, in, txt, ts1)
	assert.NoError(err)
	ts := &TestStruct{}
	err = GetFaunaInstance(db, in, txt, ts)
	assert.NoError(err)
	assert.Equal(txt, ts.Text)
	assert.Equal(nm, ts.Number * 2)
	err = DeleteFaunaIndex(db, in)
	assert.NoError(err)
	err = DeleteFaunaCollection(db, cn)
	assert.NoError(err)
}
