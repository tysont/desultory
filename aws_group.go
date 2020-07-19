package desultory

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/resourcegroups"
)

const AwsResourceGroupSuffix = "-group"

func GetAwsResourceGroupPath(groupName string, stack string) (string, error) {
	return GetAwsResourcePath(groupName, AwsResourceGroupSuffix, stack)
}

func GetAwsResourceGroupNameFromPath(groupPath string, stack string) (string, error) {
	return GetAwsResourceNameFromPath(groupPath, AwsResourceGroupSuffix, stack)
}

func CreateAwsResourceGroup(sess *session.Session, groupName string, tagName string, tagValue string, stack string) (string, error) {
	svc := resourcegroups.New(sess)
	gp, err := GetAwsResourceGroupPath(groupName, stack)
	if err != nil {
		return "", err
	}
	rs := []string {"AWS::ResourceGroups::Group", "AWS::Lambda::Function", "AWS::S3::Bucket", "AWS::DynamoDB::Table", "AWS::SQS::Queue"}
	rts := ""
	for _, r := range rs {
		if rts != "" {
			rts = rts + ","
		}
		rts = rts + "\"" + r + "\""
	}
	cgi := &resourcegroups.CreateGroupInput{
		Name: aws.String(gp),
		ResourceQuery: &resourcegroups.ResourceQuery{
			Query: aws.String("{\"ResourceTypeFilters\":[" + rts + "],\"TagFilters\":[{\"Key\":\"" + tagName + "\",\"Values\":[\"" + tagValue + "\"]}]}"),
			Type:  aws.String("TAG_FILTERS_1_0"),
		},
		Tags : make(map[string]*string, 0),
	}
	cgi.Tags[tagName] = aws.String(tagValue)
	cgo, err := svc.CreateGroup(cgi)
	if err != nil {
		return "", err
	}
	arn := cgo.Group.GroupArn
	return *arn, nil
}

func DeleteAwsResourceGroup(sess *session.Session, groupName string, stack string) error {
	svc := resourcegroups.New(sess)
	gp, err := GetAwsResourceGroupPath(groupName, stack)
	if err != nil {
		return err
	}
	dgi := &resourcegroups.DeleteGroupInput{
		GroupName: aws.String(gp),
	}
	_, err = svc.DeleteGroup(dgi)
	return err
}

func GetAwsResourceGroup(sess *session.Session, groupName string, stack string) (*resourcegroups.Group, error) {
	svc := resourcegroups.New(sess)
	gp, err := GetAwsResourceGroupPath(groupName, stack)
	if err != nil {
		return nil, err
	}
	ggi := &resourcegroups.GetGroupInput{
		GroupName: aws.String(gp),
	}
	ggo, err := svc.GetGroup(ggi)
	if err != nil {
		return nil, err
	}
	return ggo.Group, nil
}

func CheckAwsResourceGroupExists(sess *session.Session, groupName string, stack string) (bool, error) {
	g, err := GetAwsResourceGroup(sess, groupName, stack)
	if err != nil {
		return false, err
	}
	return g != nil, nil
}
