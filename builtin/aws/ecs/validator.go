package ecs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	awsarn "github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/route53"
)

func IsValidArn(arn string) bool {
	return awsarn.IsARN(arn)
}

func DoesRoute53ZoneExist(hosted_zone_id string) bool {
	privSession := CreateSession()
	svc := route53.New(privSession)
	var input route53.GetHostedZoneInput
	input.Id = aws.String(hosted_zone_id)
	_, route53Error := svc.GetHostedZone(&input)
	if route53Error != nil {
		return false
	}
	return true
}

func DoesListenerExist(arn string) bool {
	privSession := CreateSession()
	svc := elbv2.New(privSession)
	var listnerArray []string
	listnerArray = append(listnerArray, arn)

	var input elbv2.DescribeListenersInput
	input.ListenerArns = aws.StringSlice(append(listnerArray, arn))
	result, err := svc.DescribeListeners(&input)
	if err != nil {
		return false
	}
	fmt.Println(result)
	return true

}

func CreateSession() *session.Session {
	return session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
}
