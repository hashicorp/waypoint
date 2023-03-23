// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package alb

import (
	"github.com/hashicorp/waypoint/builtin/aws/ec2"
	"github.com/hashicorp/waypoint/builtin/aws/lambda"
)

func EC2TGMapper(src *ec2.Deployment) *TargetGroup {
	return &TargetGroup{
		Region: src.Region,
		Arn:    src.TargetGroupArn,
	}
}

func LambdaTGMapper(src *lambda.Deployment) *TargetGroup {
	return &TargetGroup{
		Region: src.Region,
		Arn:    src.TargetGroupArn,
	}
}
