package desultory

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"strings"
	"time"
)

const AwsIamRoleSuffix = "-role"

func GetAwsIamStackTag(stack string) *iam.Tag {
	k, v := getAwsStackTagKeyValue(stack)
	t := &iam.Tag {
		Key: aws.String(k),
		Value: aws.String(v),
	}
	return t
}

func GetAwsIamRolePath(roleName string, stack string) (string, error) {
	return GetAwsResourcePath(roleName, AwsIamRoleSuffix, stack)
}

func GetAwsIamRoleNameFromPath(rolePath string, stack string) (string, error) {
	return GetAwsResourceNameFromPath(rolePath, AwsIamRoleSuffix, stack)
}

func CreateAwsIamRole(sess *session.Session, roleName string, servicePrincipal string, stack string) (string, error) {
	svc := iam.New(sess)
	rp, err := GetAwsIamRolePath(roleName, stack)
	if err != nil {
		return "", err
	}
	pd := GetAwsAssumeRolePolicy(servicePrincipal)
	t := GetAwsIamStackTag(stack)
	ts := make([]*iam.Tag, 0)
	ts = append(ts, t)
	cri := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(pd),
		Path:                     aws.String("/"),
		RoleName:                 aws.String(rp),
		Tags: 					  ts,
	}
	res, err := svc.CreateRole(cri)
	if err != nil {
		return "", err
	}
	arn := *res.Role.Arn
	gri := &iam.GetRoleInput{
		RoleName: aws.String(rp),
	}
	err = svc.WaitUntilRoleExists(gri)
	if err != nil {
		return arn, err
	}
	return arn, nil
}

func GetAwsIamRolePolicies(sess *session.Session, roleName string, stack string) ([]string, error) {
	svc := iam.New(sess)
	ps := make([]string, 0)
	rp, err := GetAwsIamRolePath(roleName, stack)
	if err != nil {
		return ps, err
	}
	larpi := &iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(rp),
	}
	larpo, err := svc.ListAttachedRolePolicies(larpi)
	if err != nil {
		return ps, err
	}
	for _, p := range larpo.AttachedPolicies {
		ps = append(ps, *p.PolicyName)
	}
	return ps, nil
}

func DetachAwsIamRolePolicy(sess *session.Session, roleName string, policyName string, stack string) error {
	svc := iam.New(sess)
	rp, err := GetAwsIamRolePath(roleName, stack)
	if err != nil {
		return err
	}
	pa := GetAwsIamPolicyArn(policyName)
	drpi := &iam.DetachRolePolicyInput{
		RoleName:  aws.String(rp),
		PolicyArn: aws.String(pa),
	}
	_, err = svc.DetachRolePolicy(drpi)
	return err
}

func GetAwsIamPolicyArn(policyName string) string {
	return "arn:aws:iam::aws:policy/" + policyName
}

func AttachAwsIamPolicyToRole(sess *session.Session, roleName string, policyName string, stack string) error {
	svc := iam.New(sess)
	pa := GetAwsIamPolicyArn(policyName)
	rp, err := GetAwsIamRolePath(roleName, stack)
	if err != nil {
		return err
	}
	api := &iam.AttachRolePolicyInput{
		RoleName:  aws.String(rp),
		PolicyArn: aws.String(pa),
	}
	_, err = svc.AttachRolePolicy(api)
	if err != nil {
		return err
	}
	d := false
	for !d {
		ps, err := GetAwsIamRolePolicies(sess, roleName, stack)
		if err != nil {
			return err
		}
		for _, p := range ps {
			if strings.EqualFold(p, policyName) {
				d = true
				break
			}
		}
		time.Sleep(100)
	}
	return nil
}

func DeleteAwsIamRole(sess *session.Session, roleName string, stack string) error {
	svc := iam.New(sess)
	rp, err := GetAwsIamRolePath(roleName, stack)
	if err != nil {
		return err
	}
	ps, err := GetAwsIamRolePolicies(sess, roleName, stack)
	if err != nil {
		return err
	}
	for _, p := range ps {
		err = DetachAwsIamRolePolicy(sess, roleName, p, stack)
		if err != nil {
			return err
		}
	}
	dri := &iam.DeleteRoleInput{
		RoleName: aws.String(rp),
	}
	_, err = svc.DeleteRole(dri)
	return err
}

func GetAwsIamRole(sess *session.Session, roleName string, stack string) (*iam.Role, error) {
	svc := iam.New(sess)
	rp, err := GetAwsIamRolePath(roleName, stack)
	if err != nil {
		return nil, err
	}
	gri := &iam.GetRoleInput{
		RoleName: aws.String(rp),
	}
	gro, err := svc.GetRole(gri)
	if err != nil {
		if e, ok := err.(awserr.Error); ok {
			if e.Code() == "NoSuchEntity" {
				return nil, nil
			}
		}
		return nil, err
	}
	return gro.Role, nil
}

func GetAwsIamRoleTags(sess *session.Session, roleName string, stack string) (map[string]string, error) {
	ts := make(map[string]string, 0)
	r, err := GetAwsIamRole(sess, roleName, stack)
	if err != nil {
		return ts, err
	}
	for _, t := range r.Tags {
		ts[*t.Key] = *t.Value
	}
	return ts, nil
}

func CheckAwsIamRoleExists(sess *session.Session, roleName string, stack string) (bool, error) {
	r, err := GetAwsIamRole(sess, roleName, stack)
	if err != nil {
		return false, err
	}
	return r != nil, nil
}

func GetAwsAssumeRolePolicy(servicePrincipal string) string {
	return `{
  		"Version": "2012-10-17",
  		"Statement": {
			"Effect": "Allow",
    		"Principal": {"Service": "` + servicePrincipal + `"},
    		"Action": "sts:AssumeRole"
  		}
	}`
}

func ListAwsIamRoles(sess *session.Session, stack string) ([]string, error) {
	svc := iam.New(sess)
	rs := make([]string, 0)
	lro, err := svc.ListRoles(nil)
	if err != nil {
		return rs, err
	}
	for _, r := range lro.Roles {
		rp := *r.RoleName
		rn, err := GetAwsIamRoleNameFromPath(rp, stack)
		if err != nil {
			continue
		}
		ts, err := GetAwsIamRoleTags(sess, rn, stack)
		if err != nil {
			continue
			/*
			if e, ok := err.(awserr.Error); ok {
				if e.Code() == "NoSuchEntity" {
					continue
				}
			}
			return rs, err
			*/
		}
		if ts[AwsStackKey] == stack {
			rs = append(rs, rn)
		}
	}
	return rs, nil
}
