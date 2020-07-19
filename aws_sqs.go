package desultory

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const AwsSqsQueuePrefix = "-queue"

func GetAwsSqsQueuePath(queueName string, stack string) (string, error) {
	return GetAwsResourcePath(queueName, AwsSqsQueuePrefix, stack)
}

func GetAwsSqsQueueNameFromPath(queueName string, stack string) (string, error) {
	return GetAwsResourceNameFromPath(queueName, AwsSqsQueuePrefix, stack)
}

func CreateAwsSqsQueue(sess *session.Session, queueName string, stack string) (string, error) {
	svc := sqs.New(sess)
	qp, err := GetAwsSqsQueuePath(queueName, stack)
	if err != nil {
		return "", err
	}
	ts := make(map[string]*string, 0)
	ts[awsStackKey] = aws.String(stack)
	cqi := &sqs.CreateQueueInput{
		QueueName: aws.String(qp),
		Tags:     ts,
	}
	cqo, err := svc.CreateQueue(cqi)
	if err != nil {
		return "", err
	}
	as := []*string{aws.String(sqs.QueueAttributeNameQueueArn)}
	gqai := &sqs.GetQueueAttributesInput{
		AttributeNames: as,
		QueueUrl:       cqo.QueueUrl,
	}
	gqao, err := svc.GetQueueAttributes(gqai)
	if err != nil {
		return "", err
	}
	arn := gqao.Attributes[sqs.QueueAttributeNameQueueArn]
	return *arn, nil
}

func DeleteAwsSqsQueue(sess *session.Session, queueName string, stack string) error {
	svc := sqs.New(sess)
	qp, err := GetAwsSqsQueuePath(queueName, stack)
	if err != nil {
		return err
	}
	gqui := &sqs.GetQueueUrlInput{
		QueueName: aws.String(qp),
	}
	url, err := svc.GetQueueUrl(gqui)
	if err != nil {
		return err
	}
	dqi := &sqs.DeleteQueueInput{
		QueueUrl: url.QueueUrl,
	}
	_, err = svc.DeleteQueue(dqi)
	return err
}

func SendAwsSqsMessage(sess *session.Session, queueName string, stack string, message string) error {
	svc := sqs.New(sess)
	qp, err := GetAwsSqsQueuePath(queueName, stack)
	if err != nil {
		return err
	}
	gqui := &sqs.GetQueueUrlInput{
		QueueName: aws.String(qp),
	}
	url, err := svc.GetQueueUrl(gqui)
	if err != nil || url == nil {
		return err
	}
	smi := &sqs.SendMessageInput{
		QueueUrl: url.QueueUrl,
		MessageBody: aws.String(message),
	}
	_, err = svc.SendMessage(smi)
	return err
}

func ReceiveAwsSqsMessages(sess *session.Session, queueName string, stack string, waitSeconds int) ([]string, error) {
	svc := sqs.New(sess)
	qp, err := GetAwsSqsQueuePath(queueName, stack)
	if err != nil {
		return nil, err
	}
	gqui := &sqs.GetQueueUrlInput{
		QueueName: aws.String(qp),
	}
	url, err := svc.GetQueueUrl(gqui)
	if err != nil {
		return nil, err
	}
	rmi := &sqs.ReceiveMessageInput{
		QueueUrl: url.QueueUrl,
		WaitTimeSeconds: aws.Int64(int64(waitSeconds)),
	}
	rmo, err := svc.ReceiveMessage(rmi)
	if err != nil {
		return nil, err
	}
	ms := make([]string, 0)
	for _, m := range rmo.Messages {
		ms = append(ms, *m.Body)
	}
	return ms, nil
}
