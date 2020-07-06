package desultory

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateDeleteAwsResourceGroup(t *testing.T) {
	if !getOffline() {
		assert := assert.New(t)
		sess, stack, err := setupAwsTest()
		defer DeleteAwsStack(sess, stack)
		assert.NoError(err)
		gp, err := GetAwsResourceGroupPath("master", stack)
		assert.NoError(err)
		gn, err := GetAwsResourceGroupNameFromPath(gp, stack)
		assert.NoError(err)
		e, err := CheckAwsResourceGroupExists(sess, gn, stack)
		assert.NoError(err)
		assert.True(e)
	}
}
