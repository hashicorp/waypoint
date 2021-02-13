package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func ECSTaskPublicIPs(sess *session.Session, tasks []*ecs.Task) ([]string, error) {
	mapping, err := getPublicIPsFromENIs(sess, tasks)
	if err != nil {
		return nil, err
	}

	var ips []string

	for _, v := range mapping {
		ips = append(ips, v)
	}

	return ips, nil
}

const (
	eniIDKey          = "networkInterfaceId"
	ENIStatusAttached = "ATTACHED"
	ENIAttachmentType = "ElasticNetworkInterface"
)

// processAttachment takes the attachment and associates the ID of an attached ENI with the TaskArn
// Mutates: eniIDs, taskENIs
func processAttachment(taskENIs map[string]string, eniIDs *[]*string, ecsTask *ecs.Task, attachment *ecs.Attachment) {
	if aws.StringValue(attachment.Status) == ENIStatusAttached && aws.StringValue(attachment.Type) == ENIAttachmentType {
		for _, detail := range attachment.Details {
			if aws.StringValue(detail.Name) == eniIDKey {
				eniID := detail.Value
				*eniIDs = append(*eniIDs, eniID)
				taskENIs[aws.StringValue(eniID)] = aws.StringValue(ecsTask.TaskArn)
			}
		}
	}
}

func getPublicIPsFromENIs(sess *session.Session, ecsTasks []*ecs.Task) (map[string]string, error) {
	taskPublicIPs := make(map[string]string)
	var eniIDs []*string
	taskENIs := make(map[string]string)
	for _, ecsTask := range ecsTasks {
		if aws.StringValue(ecsTask.LaunchType) == "FARGATE" && aws.StringValue(ecsTask.LastStatus) == ecs.DesiredStatusRunning {
			for _, attachment := range ecsTask.Attachments {
				processAttachment(taskENIs, &eniIDs, ecsTask, attachment)
			}
		}
	}

	if len(eniIDs) == 0 {
		return taskPublicIPs, nil
	}

	ecc := ec2.New(sess)
	netInterfaces, err := ecc.DescribeNetworkInterfaces(&ec2.DescribeNetworkInterfacesInput{
		NetworkInterfaceIds: eniIDs,
	})
	if err != nil {
		return taskPublicIPs, nil
	}

	for _, eni := range netInterfaces.NetworkInterfaces {
		if eni.Association != nil {
			taskArn := taskENIs[aws.StringValue(eni.NetworkInterfaceId)]
			taskPublicIPs[taskArn] = aws.StringValue(eni.Association.PublicIp)
		}
	}

	return taskPublicIPs, nil
}
