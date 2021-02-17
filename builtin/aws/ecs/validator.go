package ecs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	awsarn "github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/route53"
)

type MockRoute53Client struct {
}

func (c *MockRoute53Client) GetHostedZone(input *route53.GetHostedZoneInput) (*route53.GetHostedZoneOutput, error) {
	var id string
	id = "Z05223941XHIVUTZAMFED"
	value := *input.Id
	if id == value {
		var output route53.GetHostedZoneOutput
		return &output, nil
	}
	return nil, fmt.Errorf("Test failure:  Input is invalid")
}

type MockALBListenerClient struct {
}

func (c *MockALBListenerClient) DescribeListeners(input *elbv2.DescribeListenersInput) (*elbv2.DescribeListenersOutput, error) {
	var listener string
	listener = "arn:aws:elasticloadbalancing:us-east-1:003559363051:listener/app/EC2Co-EcsEl-Z0096VQ81O1L/a56215152ff76fb8/057269c8b4940c21"
	value := *input.ListenerArns[0]
	if value == listener {
		output := elbv2.DescribeListenersOutput{}
		return &output, nil
	}
	return nil, fmt.Errorf("Test failure: Input is an invalid listener.")
}

type ALBListenerClient interface {
	DescribeListeners(input *elbv2.DescribeListenersInput) (*elbv2.DescribeListenersOutput, error)
}

type Route53Client interface {
	GetHostedZone(input *route53.GetHostedZoneInput) (*route53.GetHostedZoneOutput, error)
}

func isValidArn(arn string) bool {
	return awsarn.IsARN(arn)
}

func doesRoute53ZoneExist(hosted_zone_id string, useMock bool) bool {
	var client Route53Client
	if !useMock {
		sess := createSession()
		client = route53.New(sess)
	}
	if useMock {
		mc := MockRoute53Client{}
		client = &mc
	}
	var input route53.GetHostedZoneInput
	input.Id = aws.String(hosted_zone_id)
	_, route53Error := client.GetHostedZone(&input)
	if route53Error != nil {
		return false
	}

	return true
}

func doesListenerExist(arn string, useMock bool) bool {
	var client ALBListenerClient
	if !useMock {
		sess := createSession()
		client = elbv2.New(sess)
	}
	if useMock {
		mc := MockALBListenerClient{}
		client = &mc
	}
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
