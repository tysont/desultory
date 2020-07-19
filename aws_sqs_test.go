package desultory

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateSendReceiveDeleteAwsSqsQueue(t *testing.T) {
	if !getOffline() {
		assert := assert.New(t)
		sess, stack, err := setupAwsTest()
		defer DeleteAwsStack(sess, stack)
		assert.NoError(err)
		qn := "test"
		_, err = CreateAwsSqsQueue(sess, qn, stack)
		assert.NoError(err)
		m := "peter piper picked a pack of pickled peppers"
		err = SendAwsSqsMessage(sess, qn, stack, m)
		assert.NoError(err)
		ms, err := ReceiveAwsSqsMessages(sess, qn, stack, 5)
		assert.NoError(err)
		assert.Equal(len(ms), 1)
		assert.Equal(m, ms[0])
		//err = DeleteAwsSqsQueue(sess, qn, stack)
		//assert.NoError(err)
	}
}