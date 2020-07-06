package desultory

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"net/http"
	"strings"
	"time"
)

const AwsLambdaFunctionSuffix = "-function"

func GetAwsLambdaFunctionPath(functionName string, stack string) (string, error) {
	return GetAwsResourcePath(functionName, AwsLambdaFunctionSuffix, stack)
}

func GetAwsLambdaFunctionNameFromPath(functionPath string, stack string) (string, error) {
	return GetAwsResourceNameFromPath(functionPath, AwsLambdaFunctionSuffix, stack)
}

func GetAwsLambdaFunctionArnFromPath(sess *session.Session, functionPath string) (string, error) {
	svc := lambda.New(sess)
	gfi := &lambda.GetFunctionInput {
		FunctionName: aws.String(functionPath),
	}
	gfo, err := svc.GetFunction(gfi)
	if err != nil {
		return "", err
	}
	arn := *gfo.Configuration.FunctionArn
	return arn, nil
}

func GetAwsLambdaFunctionPathFromArn(sess *session.Session, functionArn string) (string, error) {
	s := strings.Split(functionArn, ":")
	fp := s[len(s) - 1]
	return fp, nil
}

func CreateAwsLambdaFunction(sess *session.Session, functionName string, functionLanguage string, functionHandler string, functionCodeZip *bytes.Buffer, stack string) (string, error) {
	zn := functionName + ".zip"
	bn := "code"
	bp, err := GetAwsS3BucketPath(bn, stack)
	if err != nil {
		return "", err
	}
	be, err := CheckAwsS3BucketExists(sess, bn, stack)
	if err != nil {
		return "", err
	}
	if !be {
		_, err = CreateAwsS3Bucket(sess, bn, stack)
		if err != nil {
			return "", err
		}
	}
	b := functionCodeZip.Bytes()
	err = WriteToAwsS3Bucket(sess, bn, zn, b, stack)
	if err != nil {
		return "", err
	}
	rn := functionName
	rarn, err := CreateAwsIamRole(sess, rn, "lambda.amazonaws.com", stack)
	if err != nil {
		return "", err
	}
	err = AttachAwsIamPolicyToRole(sess, rn, "AWSLambdaExecute", stack)
	if err != nil {
		return "", err
	}

	fc := &lambda.FunctionCode{
		S3Bucket: aws.String(bp),
		S3Key:    aws.String(zn),
	}
	fp, err := GetAwsLambdaFunctionPath(functionName, stack)
	if err != nil {
		return "", err
	}
	l := ""
	switch strings.ToLower(functionLanguage) {
	case "go":
		l = "go1.x"
	case "java":
		l = "java11"
	case "node":
		l = "nodejs12.x"
	default:
		return "", fmt.Errorf("function language '%v' is not supported", functionLanguage)
	}
	cfi := &lambda.CreateFunctionInput{
		FunctionName: aws.String(fp),
		Code:         fc,
		Handler:      aws.String(functionHandler),
		Role:         aws.String(rarn),
		Runtime:      aws.String(l),
	}
	svc := lambda.New(sess)
	//https://stackoverflow.com/questions/36419442/the-role-defined-for-the-function-cannot-be-assumed-by-lambda
	s := time.Now()
	w := 15 * time.Second
	d := false
	var res *lambda.FunctionConfiguration
	for !d && s.Add(w).After(time.Now()) {
		time.Sleep(1 * time.Second)
		res, err = svc.CreateFunction(cfi)
		if err != nil {
			if _, ok := err.(*lambda.InvalidParameterValueException); ok {
				continue
			}
			return "", err
		}
		d = true
	}
	arn := *res.FunctionArn
	ts := make(map[string]*string, 0)
	ts[AwsStackKey] = aws.String(stack)
	tri := &lambda.TagResourceInput{
		Resource: aws.String(arn),
		Tags:     ts,
	}
	_, err = svc.TagResource(tri)
	if err != nil {
		return arn, err
	}
	return arn, nil
}

func DeleteAwsLambdaFunction(sess *session.Session, functionName string, stack string) error {
	svc := lambda.New(sess)
	fp, err := GetAwsLambdaFunctionPath(functionName, stack)
	if err != nil {
		return err
	}
	dfi := &lambda.DeleteFunctionInput {
		FunctionName: aws.String(fp),
	}
	_, err = svc.DeleteFunction(dfi)
	if err != nil {
		return err
	}
	return nil
}

func CheckAwsLambdaFunctionExists(sess *session.Session, functionName string, stack string) (bool, error) {
	svc := lambda.New(sess)
	fp, err := GetAwsLambdaFunctionPath(functionName, stack)
	if err != nil {
		return false, err
	}
	gfi := &lambda.GetFunctionInput {
		FunctionName: aws.String(fp),
	}
	gfo, err := svc.GetFunction(gfi)
	if err != nil {
		if e, ok := err.(*lambda.ResourceNotFoundException); ok {
			if e.StatusCode() == http.StatusNotFound {
				return false, nil
			}
		}
		return false, err
	}
	return gfo.Configuration != nil, nil
}

func InvokeAwsLambdaFunction(sess *session.Session, functionName string, functionInput []byte, stack string) ([]byte, error) {
	svc := lambda.New(sess)
	fp, err := GetAwsLambdaFunctionPath(functionName, stack)
	if err != nil {
		return nil, err
	}
	ii := &lambda.InvokeInput {
		FunctionName: aws.String(fp),
		Payload:      functionInput,
	}
	io, err := svc.Invoke(ii)
	return io.Payload, nil
}
