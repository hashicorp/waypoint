package ecs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	awsarn "github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/route53"
)

type ALBListenerClient interface {
	DescribeListeners(input *elbv2.DescribeListenersInput) (*elbv2.DescribeListenersOutput, error)
}

type Route53Client interface {
	GetHostedZone(input *route53.GetHostedZoneInput) (*route53.GetHostedZoneOutput, error)
}

func isValidArn(arn string) bool {
	return awsarn.IsARN(arn)
}

func doesRoute53ZoneExist(hosted_zone_id string, client Route53Client) bool {

	var input route53.GetHostedZoneInput
	input.Id = aws.String(hosted_zone_id)
	_, route53Error := client.GetHostedZone(&input)
	if route53Error != nil {
		return false
	}
	return true
}

func doesListenerExist(arn string, client ALBListenerClient) bool {

	var listnerArray []string
	listnerArray = append(listnerArray, arn)

	var input elbv2.DescribeListenersInput
	input.ListenerArns = aws.StringSlice(append(listnerArray, arn))
	result, err := client.DescribeListeners(&input)
	if err != nil {
		return false
	}
	fmt.Println(result)
	return true

}

func createSession() *session.Session {
	return session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
}
