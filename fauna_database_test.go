package desultory

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatePutGetDeleteFaunaDatabase(t *testing.T) {
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
		Number: 14,
	}
	err = CreateFaunaInstance(db, cn, ts2)
	assert.NoError(err)
	ts := &TestStruct{}
	err = GetFaunaInstance(db, in, txt, ts)
	assert.NoError(err)
	assert.Equal(txt, ts.Text)
	assert.Equal(nm, ts.Number)
	err = DeleteFaunaIndex(db, in)
	assert.NoError(err)
	err = DeleteFaunaCollection(db, cn)
	assert.NoError(err)

	/*
	bn := GetUniqueString(6)
	_, err = CreateAwsS3Bucket(sess, bn, sn)
	assert.NoError(err)
	f := noosly.CreateFeed("www.stackoverflow.com")
	k := f.GetKey()
	b, err := json.Marshal(f)
	assert.NoError(err)
	err = WriteToAwsS3Bucket(sess, bn, k, b, sn)
	assert.NoError(err)
	b2, err := ReadFromAwsS3Bucket(sess, bn, k, sn)
	f2 := new(noosly.Feed)
	err = json.Unmarshal(b2, f2)
	assert.NoError(err)
	assert.NotNil(f2)
	assert.Equal(f.Url, f2.Url)
	*/
}
