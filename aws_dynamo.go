package desultory

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"strings"
)

const AwsDynamoTableSuffix = "-table"

func GetAwsDynamoStackTag(stack string) *dynamodb.Tag {
	k, v := getAwsStackTagKeyValue(stack)
	t := &dynamodb.Tag {
		Key: aws.String(k),
		Value: aws.String(v),
	}
	return t
}

func GetAwsDynamoTablePath(tableName string, stack string) (string, error) {
	return GetAwsResourcePath(tableName, AwsDynamoTableSuffix, stack)
}

func GetAwsDynamoTableNameFromPath(tablePath string, stack string) (string, error) {
	return GetAwsResourceNameFromPath(tablePath, AwsDynamoTableSuffix, stack)
}

func GetAwsDynamoTablePathFromArn(sess *session.Session, bucketArn string) (string, error) {
	s := strings.Split(bucketArn, "/")
	return s[len(s) - 1], nil
}

func CreateAwsDynamoTable(sess *session.Session, tableName string, tableAttributes []string, partitionKey string, sortKey string, stack string) (string, error) {
	svc := dynamodb.New(sess)
	tp, err := GetAwsDynamoTablePath(tableName, stack)
	if err != nil {
		return "", err
	}
	tas := make([]*dynamodb.AttributeDefinition, 0)
	for _, an := range tableAttributes {
		ta := &dynamodb.AttributeDefinition{
			AttributeName: aws.String(an),
			AttributeType: aws.String("S"),
		}
		tas = append(tas, ta)
	}
	pk := &dynamodb.KeySchemaElement{
		AttributeName: aws.String(partitionKey),
		KeyType:       aws.String("HASH"),
	}
	ks := []*dynamodb.KeySchemaElement { pk }
	if sortKey != "" {
		sk := &dynamodb.KeySchemaElement{
			AttributeName: aws.String(sortKey),
			KeyType:       aws.String("RANGE"),
		}
		ks = append(ks, sk)
	}
	ts := []*dynamodb.Tag { GetAwsDynamoStackTag(stack) }
	cti := &dynamodb.CreateTableInput{
		TableName: aws.String(tp),
		BillingMode: aws.String("PAY_PER_REQUEST"),
		AttributeDefinitions: tas,
		KeySchema: ks,
		Tags: ts,
	}
	cto, err := svc.CreateTable(cti)
	if err != nil {
		return "", err
	}
	tn := *cto.TableDescription.TableArn
	dti := &dynamodb.DescribeTableInput{
		TableName: aws.String(tp),
	}
	err = svc.WaitUntilTableExists(dti)
	if err != nil {
		return "", err
	}
	return tn, nil
}

func DeleteAwsDynamoTable(sess *session.Session, tableName string, stack string) error {
	svc := dynamodb.New(sess)
	tp, err := GetAwsDynamoTablePath(tableName, stack)
	if err != nil {
		return err
	}
	dti := &dynamodb.DeleteTableInput{
		TableName: aws.String(tp),
	}
	_, err = svc.DeleteTable(dti)
	if err != nil {
		return err
	}
	dsti := &dynamodb.DescribeTableInput {
		TableName: aws.String(tp),
	}
	err = svc.WaitUntilTableNotExists(dsti)
	return err
}

func WriteToAwsDynamoTable(sess *session.Session, tableName string, itemAttributes map[string]string, stack string) error {
	svc := dynamodb.New(sess)
	tp, err := GetAwsDynamoTablePath(tableName, stack)
	if err != nil {
		return err
	}
	am := make(map[string]*dynamodb.AttributeValue, 0)
	for k, v := range itemAttributes {
		am[k] = &dynamodb.AttributeValue{
			S: aws.String(v),
		}
	}
	pii := &dynamodb.PutItemInput{
		TableName: aws.String(tp),
		Item: am,
	}
	_, err = svc.PutItem(pii)
	return err
}

func ReadFromAwsDymanoTable(sess *session.Session, tableName string, itemAttributes map[string]string, stack string) (map[string]string, error) {
	svc := dynamodb.New(sess)
	tp, err := GetAwsDynamoTablePath(tableName, stack)
	if err != nil {
		return nil, err
	}
	ram := make(map[string]*dynamodb.AttributeValue, 0)
	for k, v := range itemAttributes {
		ram[k] = &dynamodb.AttributeValue{
			S: aws.String(v),
		}
	}
	gii := &dynamodb.GetItemInput{
		TableName: aws.String(tp),
		Key:       ram,
	}
	gio, err := svc.GetItem(gii)
	if err != nil {
		return nil, err
	}
	am := make(map[string]string, 0)
	for an, av := range gio.Item {
		am[an] = *av.S
	}
	return am, nil
}

func CheckAwsDynamoTableExists(sess *session.Session, tableName string, stack string) (bool, error) {
	svc := dynamodb.New(sess)
	tp, err := GetAwsDynamoTablePath(tableName, stack)
	if err != nil {
		return false, err
	}
	dti := &dynamodb.DescribeTableInput{
		TableName: aws.String(tp),
	}
	_, err = svc.DescribeTable(dti)
	if err != nil {
		if _, ok := err.(*dynamodb.ResourceNotFoundException); ok {
			return false, nil
		}
		return false, err
	}
	return true, nil
}