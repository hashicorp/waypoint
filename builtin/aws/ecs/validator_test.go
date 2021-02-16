package ecs

import (
	"testing"
)

func TestRoute53HostedZone(t *testing.T) {

	result := DoesRoute53ZoneExist("Z05223941XHIVUTZAMFED")
	if result == false {
		t.Fatal("Error: Exiiting Route53 zone was not found.")
	}

}

func TestLoadBalancerListener(t *testing.T) {
	result := DoesListenerExist("arn:aws:elasticloadbalancing:us-east-1:003559363051:listener/app/EC2Co-EcsEl-Z0096VQ81O1L/a56215152ff76fb8/057269c8b4940c21")
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
