package utils

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func DefaultSubnets(ctx context.Context, sess *session.Session) ([]*string, *string, error) {
	svc := ec2.New(sess)

	desc, err := svc.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("default-for-az"),
				Values: []*string{aws.String("true")},
			},
		},
	})

	if err != nil {
		return nil, nil, err
	}

	var (
		subnets []*string
		vpc     *string
	)

	for _, subnet := range desc.Subnets {
		if vpc == nil {
			vpc = subnet.VpcId
		}

		subnets = append(subnets, subnet.SubnetId)
	}

	return subnets, vpc, nil
}

func DefaultPublicSubnets(ctx context.Context, sess *session.Session) ([]*string, *string, error) {
	svc := ec2.New(sess)

	desc, err := svc.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("default-for-az"),
				Values: []*string{aws.String("true")},
			},
		},
	})

	if err != nil {
		return nil, nil, err
	}

	var (
		subnets []*string
		vpc     *string
	)

	for _, subnet := range desc.Subnets {
		if vpc == nil {
			vpc = subnet.VpcId
		}

		if *subnet.MapPublicIpOnLaunch {
			subnets = append(subnets, subnet.SubnetId)
		}
	}

	return subnets, vpc, nil
}
