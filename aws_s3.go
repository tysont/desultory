package desultory

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"net/http"
	"strings"
)

const AwsS3BucketSuffix = "-bucket"

func GetAwsS3StackTag(stack string) *s3.Tag {
	k, v := getAwsStackTagKeyValue(stack)
	t := &s3.Tag {
		Key: aws.String(k),
		Value: aws.String(v),
	}
	return t
}

func GetAwsS3BucketPath(bucketName string, stack string) (string, error) {
	return GetAwsResourcePath(bucketName, AwsS3BucketSuffix, stack)
}

func GetAwsS3BucketNameFromPath(bucketPath string, stack string) (string, error) {
	return GetAwsResourceNameFromPath(bucketPath, AwsS3BucketSuffix, stack)
}

func GetAwsS3BucketPathFromArn(sess *session.Session, bucketArn string) (string, error) {
	s := strings.Split(bucketArn, ":")
	return s[len(s) - 1], nil
}

func GetAwsS3BucketArnFromPath(bucketPath string) string {
	return "arn:aws:s3:::" + bucketPath
}

func CreateAwsS3Bucket(sess *session.Session, bucketName string, stack string) (string, error) {
	svc := s3.New(sess)
	bp, err := GetAwsS3BucketPath(bucketName, stack)
	if err != nil {
		return "", err
	}
	cbi := &s3.CreateBucketInput{
		Bucket: aws.String(bp),
	}
	_, err = svc.CreateBucket(cbi)
	if err != nil {
		return "", err
	}
	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bp),
	})
	if err != nil {
		return "", err
	}
	arn := GetAwsS3BucketArnFromPath(bp)
	if stack != "" {
		t := GetAwsS3StackTag(stack)
		ts := &s3.Tagging {
			TagSet: []*s3.Tag { t },
		}
		pbti := &s3.PutBucketTaggingInput{
			Bucket: aws.String(bp),
			Tagging: ts,
		}
		_, err = svc.PutBucketTagging(pbti)
		if err != nil {
			return arn, err
		}
	}
	return arn, nil
}

func WriteToAwsS3Bucket(sess *session.Session, bucketName string, key string, value []byte, stack string) error {
	up := s3manager.NewUploader(sess)
	bp, err := GetAwsS3BucketPath(bucketName, stack)
	if err != nil {
		return err
	}
	_, err = up.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bp),
		Key:    aws.String(key),
		Body:   bytes.NewReader(value),
	})
	return err
}

func ReadFromAwsS3Bucket(sess *session.Session, bucketName string, key string, stack string) ([]byte, error) {
	dn := s3manager.NewDownloader(sess)
	bp, err := GetAwsS3BucketPath(bucketName, stack)
	if err != nil {
		return nil, err
	}
	b := aws.NewWriteAtBuffer(make([]byte, 0))
	_, err = dn.Download(b, &s3.GetObjectInput{
		Bucket: aws.String(bp),
		Key:    aws.String(key),
	})
	return b.Bytes(), err
}

func GetAwsS3BucketTags(sess *session.Session, bucketName string, stack string) (map[string]string, error) {
	svc := s3.New(sess)
	ts := make(map[string]string, 0)
	bp, err := GetAwsS3BucketPath(bucketName, stack)
	if err != nil {
		return ts, err
	}
	gbti := &s3.GetBucketTaggingInput{
		Bucket: aws.String(bp),
	}
	gbto, err := svc.GetBucketTagging(gbti)
	if err != nil {
		return ts, err
	}
	for _, t := range gbto.TagSet {
		ts[*t.Key] = *t.Value
	}
	return ts, nil
}

func DeleteAwsS3Bucket(sess *session.Session, bucketName string, stack string) error {
	svc := s3.New(sess)
	bp, err := GetAwsS3BucketPath(bucketName, stack)
	if err != nil {
		return err
	}
	i := s3manager.NewDeleteListIterator(svc, &s3.ListObjectsInput{
		Bucket: aws.String(bp),
	})
	err = s3manager.NewBatchDeleteWithClient(svc).Delete(aws.BackgroundContext(), i)
	if err != nil {
		return err
	}
	_, err = svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bp),
	})
	if err != nil {
		return err
	}
	err = svc.WaitUntilBucketNotExists(&s3.HeadBucketInput{
		Bucket: aws.String(bp),
	})
	return err
}

func CheckAwsS3BucketExists(sess *session.Session, bucketName string, stack string) (bool, error) {
	svc := s3.New(sess)
	bp, err := GetAwsS3BucketPath(bucketName, stack)
	if err != nil {
		return false, err
	}
	hbi := &s3.HeadBucketInput{
		Bucket: aws.String(bp),
	}
	hbo, err := svc.HeadBucket(hbi)
	if err != nil {
		if e, ok := err.(s3.RequestFailure); ok {
			if e.StatusCode() == http.StatusNotFound {
				return false, nil
			}
		}
		return false, err
	}
	return hbo != nil, nil
}

func ListAwsS3Buckets(sess *session.Session, stack string) ([]string, error) {
	svc := s3.New(sess)
	bs := make([]string, 0)
	lbo, err := svc.ListBuckets(nil)
	if err != nil {
		return bs, err
	}
	for _, b := range lbo.Buckets {
		bp := *b.Name
		bn, err := GetAwsS3BucketNameFromPath(bp, stack)
		if err != nil {
			continue
		}
		ts, err := GetAwsS3BucketTags(sess, bn, stack)
		if err != nil {
			continue
			/*
			if e, ok := err.(s3.RequestFailure); ok {
				if e.StatusCode() == http.StatusNotFound {
					continue
				}
			}
			return bs, err
			*/
		}
		if ts[awsStackKey] == stack {
			bs = append(bs, bn)
		}
	}
	return bs, nil
}
