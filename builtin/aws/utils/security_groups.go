// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package utils

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func CreateSecurityGroup(
	ctx context.Context,
	sess *session.Session,
	name string,
	vpcId *string,
	ports ...int,
) (*string, error) {
	ec2srv := ec2.New(sess)

	dsg, err := ec2srv.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("group-name"),
				Values: []*string{aws.String(name)},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	var groupId *string

	if len(dsg.SecurityGroups) != 0 {
		groupId = dsg.SecurityGroups[0].GroupId
	} else {
		out, err := ec2srv.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
			Description: aws.String("created by waypoint"),
			GroupName:   aws.String(name),
			VpcId:       vpcId,
		})
		if err != nil {
			return nil, err
		}

		groupId = out.GroupId
	}

	for _, port := range ports {
		_, err = ec2srv.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
			CidrIp:     aws.String("0.0.0.0/0"),
			FromPort:   aws.Int64(int64(port)),
			ToPort:     aws.Int64(int64(port)),
			GroupId:    groupId,
			IpProtocol: aws.String("tcp"),
		})
	}

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "InvalidPermission.Duplicate":
				// fine, means we already added it.
			default:
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return groupId, nil
}
