package desultory

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateDeleteAwsRole(t *testing.T) {
	if !getOffline() {
		assert := assert.New(t)
		sess, stack, err := setupAwsTest()
		defer DeleteAwsStack(sess, stack)
		assert.NoError(err)
		rn := "test"
		_, err = CreateAwsIamRole(sess, rn, "lambda.amazonaws.com", stack)
		assert.NoError(err)
		err = AttachAwsIamPolicyToRole(sess, rn, "AWSLambdaExecute", stack)
		assert.NoError(err)
		e, err := CheckAwsIamRoleExists(sess, rn, stack)
		assert.NoError(err)
		assert.True(e)
		rs, err := ListAwsIamRoles(sess, stack)
		assert.NoError(err)
		assert.Equal(1, len(rs))
		err = DeleteAwsIamRole(sess, rn, stack)
		assert.NoError(err)
		e, err = CheckAwsIamRoleExists(sess, rn, stack)
		assert.NoError(err)
		assert.False(e)
	}
}
