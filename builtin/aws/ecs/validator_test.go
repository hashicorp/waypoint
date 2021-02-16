package ecs

import (
	"fmt"
	"testing"

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

func TestRoute53HostedZone(t *testing.T) {
	mc := MockRoute53Client{}
	result := DoesRoute53ZoneExist("Z05223941XHIVUTZAMFED", &mc)
	if result == false {
		t.Fatal("Error: Exiiting Route53 zone was not found.")
	}

}

func TestLoadBalancerListener(t *testing.T) {
	mc := MockALBListenerClient{}
	result := DoesListenerExist("arn:aws:elasticloadbalancing:us-east-1:003559363051:listener/app/EC2Co-EcsEl-Z0096VQ81O1L/a56215152ff76fb8/057269c8b4940c21", &mc)
	if result == false {
		t.Fatal("Error: Loadblancer Listener does not exist.")
	}
}

func TestLoadBalancerArn(t *testing.T) {
	result := IsValidArn("arn:aws:elasticloadbalancing:us-east-1:003559363051:listener/app/EC2Co-EcsEl-Z0096VQ81O1L/a56215152ff76fb8/057269c8b4940c21")
	if result == false {
		t.Fatal("Error:  ARN supplied is not a valid ARN.")
	}
}
