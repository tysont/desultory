package desultory

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateInvokeDeleteAwsLambdaNodeFunction(t *testing.T) {
	if !getOffline() {
		assert := assert.New(t)
		sess, stack, err := setupAwsTest()
		defer DeleteAwsStack(sess, stack)
		assert.NoError(err)
		t := "Hello, World!"
		fn := "multiply"
		fl := "node"
		fln := "index.js"
		fh := "index.handler"
		flc := `exports.handler = function(event, context) {
	context.succeed('` + t + `');
};`

		fs := map[string][]byte{
			fln: []byte(flc),
		}
		b, err := WriteZip(fs)
		assert.NoError(err)
		_, err = CreateAwsLambdaFunction(sess, fn, fl, fh, b, stack)
		assert.NoError(err)
		e, err := CheckAwsLambdaFunctionExists(sess, fn, stack)
		assert.NoError(err)
		y, err := InvokeAwsLambdaFunction(sess, fn, nil, stack)
		assert.NoError(err)
		assert.Contains(string(y), t)
		err = DeleteAwsLambdaFunction(sess, fn, stack)
		assert.NoError(err)
		e, err = CheckAwsLambdaFunctionExists(sess, fn, stack)
		assert.NoError(err)
		assert.Equal(false, e)
	}
}

