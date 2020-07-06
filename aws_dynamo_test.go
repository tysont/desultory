package desultory

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateWriteReadDeleteAwsDynamoTable(t *testing.T) {
	if !getOffline() {
		assert := assert.New(t)
		sess, stack, err := setupAwsTest()
		defer DeleteAwsStack(sess, stack)
		assert.NoError(err)
		tn := "test"
		a1 := "foo"
		a2 := "bar"
		a3 := "baz"
		v1 := "hello"
		v2 := "world"
		v3 := "from the metaverse"
		as := []string {a1, a2}

		_, err = CreateAwsDynamoTable(sess, tn, as, a1, a2, stack)
		assert.NoError(err)
		e, err := CheckAwsDynamoTableExists(sess, tn, stack)
		assert.NoError(err)
		assert.True(e)
		avs := map[string]string {
			a1: v1,
			a2: v2,
			a3: v3,
		}
		err = WriteToAwsDynamoTable(sess, tn, avs, stack)
		assert.NoError(err)
		rqavs := map[string]string {
			a1: v1,
			a2: v2,
		}
		rsavs, err := ReadFromAwsDymanoTable(sess, tn, rqavs, stack)
		assert.NoError(err)
		assert.Equal(rsavs[a3], v3)
		err = DeleteAwsDynamoTable(sess, tn, stack)
		assert.NoError(err)
		e, err = CheckAwsDynamoTableExists(sess, tn, stack)
		assert.NoError(err)
		assert.False(e)
	}
}
