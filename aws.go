package desultory

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/resourcegroups"
	"strings"
)

var awsResourcePrefix = "desultory-"
var awsStackKey = awsResourcePrefix + "stack"
var awsRegion = "us-west-2"

func SetAwsResourcePrefix(prefix string) {
	awsResourcePrefix = prefix
}

func SetAwsStackKey(key string) {
	awsStackKey = key
}

func SetAwsRegion(region string) {
	awsRegion = region
}

func GetAwsSession() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)
	return sess, err
}

func getAwsStackTagKeyValue(stack string) (string, string) {
	return awsStackKey, stack
}

func GetAwsResourcePath(resourceName string, resourceSuffix string, stack string) (string, error) {
	if strings.Contains(resourceName, "-") {
		return "", fmt.Errorf("resource name '%v' contains a dash, which is not supported", resourceName)
	}
	s := awsResourcePrefix + resourceName + "-" + stack + resourceSuffix
	if len(s) > 64 {
		return "", fmt.Errorf("implied resource path '%v' exceeds the 64 character limit", s)
	}
	return s, nil
}

func GetAwsResourceNameFromPath(resourcePath string, resourceSuffix string, stack string) (string, error) {
	s := strings.Replace(resourcePath, awsResourcePrefix, "", 1)
	s = strings.Replace(s, "-" + stack + resourceSuffix, "", 1)
	if strings.Contains(s, "-") {
		return "", fmt.Errorf("implied resource name '%v' contains a dash, which is not supported", s)
	} else if s == resourcePath {
		return "", fmt.Errorf("implied resource name '%v' matches the path, which is invalid", s)
	}
	return s, nil
}

func CreateAwsStack(sess *session.Session, stack string) error {
	k, v := getAwsStackTagKeyValue(stack)
	gn := "master"
	_, err := CreateAwsResourceGroup(sess, gn, k, v, stack)
	return err
}

func DeleteAwsStack(sess *session.Session, stack string) error {
	svc := resourcegroups.New(sess)
	gn := "master"
	gp, err := GetAwsResourceGroupPath(gn, stack)
	if err != nil {
		return err
	}
	lgri := &resourcegroups.ListGroupResourcesInput{
		GroupName:  aws.String(gp),
	}
	lgro, err := svc.ListGroupResources(lgri)
	if err != nil {
		return err
	}
	for _, ri := range lgro.ResourceIdentifiers {
		rt := *ri.ResourceType
		arn := *ri.ResourceArn
		switch rt {
		case "AWS::S3::Bucket":
			bp, err := GetAwsS3BucketPathFromArn(sess, arn)
			if err != nil {
				return err
			}
			bn, err := GetAwsS3BucketNameFromPath(bp, stack)
			if err != nil {
				return err
			}
			err = DeleteAwsS3Bucket(sess, bn, stack)
			if err != nil {
				return err
			}
		case "AWS::Lambda::Function":
			fp, err := GetAwsLambdaFunctionPathFromArn(sess, arn)
			if err != nil {
				return err
			}
			fn, err := GetAwsLambdaFunctionNameFromPath(fp, stack)
			if err != nil {
				return err
			}
			err = DeleteAwsLambdaFunction(sess, fn, stack)
			if err != nil {
				return err
			}
		case "AWS::DynamoDB::Table":
			tp, err := GetAwsDynamoTablePathFromArn(sess, arn)
			if err != nil {
				return nil
			}
			tn, err := GetAwsDynamoTableNameFromPath(tp, stack)
			if err != nil {
				return nil
			}
			err = DeleteAwsDynamoTable(sess, tn, stack)
			if err != nil {
				return nil
			}
		case "AWS::SQS::Queue":
			qp, err := GetAwsSqsQueuePathFromArn(sess, arn)
			if err != nil {
				return nil
			}
			qn, err := GetAwsSqsQueueNameFromPath(qp, stack)
			if err != nil {
				return nil
			}
			err = DeleteAwsSqsQueue(sess, qn, stack)
			if err != nil {
				return nil
			}
		}
	}
	rs, err := ListAwsIamRoles(sess, stack)
	if err != nil {
		return err
	}
	for _, rn := range rs {
		err = DeleteAwsIamRole(sess, rn, stack)
		if err != nil {
			return err
		}
	}
	err = DeleteAwsResourceGroup(sess, gn, stack)
	if err != nil {
		return err
	}
	return nil
}