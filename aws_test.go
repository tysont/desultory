package desultory

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

func setupAwsTest() (*session.Session, string, error) {
	stack := GetUniqueString(8)
	sess, err := GetAwsSession()
	if err != nil {
		return nil, "", err
	}
	err = CreateAwsStack(sess, stack)
	if err != nil {
		return nil, "", err
	}
	return sess, stack, nil
}

func teardownAwsTest(sess *session.Session, stack string) error {
	return DeleteAwsStack(sess, stack)
}