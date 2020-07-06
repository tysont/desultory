package desultory

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateWriteReadDeleteAwsS3Bucket(t *testing.T) {
	if !getOffline() {
		assert := assert.New(t)
		sess, stack, err := setupAwsTest()
		defer DeleteAwsStack(sess, stack)
		assert.NoError(err)
		bn := "test"
		_, err = CreateAwsS3Bucket(sess, bn, stack)
		assert.NoError(err)
		e, err := CheckAwsS3BucketExists(sess, bn, stack)
		assert.NoError(err)
		assert.True(e)
		bs, err := ListAwsS3Buckets(sess, stack)
		assert.NoError(err)
		assert.Equal(1, len(bs))
		ts, err := GetAwsS3BucketTags(sess, bn, stack)
		assert.NoError(err)
		assert.NotNil(ts)
		assert.NotEmpty(ts)
		assert.Equal(stack, ts[AwsStackKey])
		k := "test.txt"
		v := "multiply, world!"
		err = WriteToAwsS3Bucket(sess, bn, k, []byte(v), stack)
		assert.NoError(err)
		b, err := ReadFromAwsS3Bucket(sess, bn, k, stack)
		assert.NoError(err)
		assert.Equal(string(b), v)
		err = DeleteAwsS3Bucket(sess, bn, stack)
		assert.NoError(err)
		e, err = CheckAwsS3BucketExists(sess, bn, stack)
		assert.NoError(err)
		assert.False(e)
	}
}
