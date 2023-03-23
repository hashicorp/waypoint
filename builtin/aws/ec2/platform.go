// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ec2

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/aws/ami"
	"github.com/hashicorp/waypoint/builtin/aws/utils"
)

const (
	labelId    = "waypoint.hashicorp.com/id"
	labelNonce = "waypoint.hashicorp.com/nonce"
)

// Platform is the Platform implementation for Amazon EC2.
type Platform struct {
	config PlatformConfig
}

// Config implements Configurable
func (p *Platform) Config() (interface{}, error) {
	return &p.config, nil
}

// DeployFunc implements component.Platform
func (p *Platform) DeployFunc() interface{} {
	return p.Deploy
}

// DestroyFunc implements component.Destroyer
func (p *Platform) DestroyFunc() interface{} {
	return p.Destroy
}

// ValidateAuthFunc implements component.Authenticator
func (p *Platform) ValidateAuthFunc() interface{} {
	return p.ValidateAuth
}

// AuthFunc implements component.Authenticator
func (p *Platform) AuthFunc() interface{} {
	return p.Auth
}

func (p *Platform) Auth() error {
	return nil
}

func (p *Platform) ValidateAuth() error {
	return nil
}

// DefaultReleaserFunc implements component.PlatformReleaser
// func (p *Platform) DefaultReleaserFunc() interface{} {
// return func() *Releaser { return &Releaser{p: p} }
// }

// Deploy deploys an image to Amazon EC2.
func (p *Platform) Deploy(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	job *component.JobInfo,
	img *ami.Image,
	deployConfig *component.DeploymentConfig,
	ui terminal.UI,
) (*Deployment, error) {
	cid, err := component.Id()
	if err != nil {
		return nil, err
	}

	// We'll update the user in real time
	st := ui.Status()
	defer st.Close()

	st.Update("Creating EC2 instances in ASG...")

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: p.config.Region,
		Logger: log,
	})
	if err != nil {
		return nil, err
	}
	e := ec2.New(sess)

	var (
		subnetsComma string
		vpc          *string
	)

	if p.config.Subnet == "" {
		subnets, v, err := utils.DefaultSubnets(ctx, sess)
		if err != nil {
			return nil, err
		}

		var sc []string
		for _, s := range subnets {
			sc = append(sc, *s)
		}

		subnetsComma = strings.Join(sc, ",")

		vpc = v

	} else {
		subnetInfo, err := e.DescribeSubnets(&ec2.DescribeSubnetsInput{
			SubnetIds: []*string{aws.String(p.config.Subnet)},
		})
		if err != nil {
			return nil, err
		}

		subnetsComma = p.config.Subnet

		vpc = subnetInfo.Subnets[0].VpcId
	}

	ports := append([]int{p.config.ServicePort}, p.config.ExtraPorts...)

	sec, err := utils.CreateSecurityGroup(ctx, sess, fmt.Sprintf("waypoint-%s", src.App), vpc, ports...)
	if err != nil {
		return nil, err
	}

	groups := []*string{sec}

	for _, g := range p.config.SecurityGroups {
		groups = append(groups, aws.String(g))
	}

	st.Update("Launching instance...")

	ud, err := UserData(deployConfig.Env())
	if err != nil {
		return nil, err
	}

	var key *string

	if p.config.Key != "" {
		key = aws.String(p.config.Key)
	}

	rand := cid[len(cid)-(31-len(src.App)):]

	serviceName := fmt.Sprintf("%s-%s", src.App, rand)

	elbsrv := elbv2.New(sess)

	ctg, err := elbsrv.CreateTargetGroup(&elbv2.CreateTargetGroupInput{
		Name:       aws.String(serviceName),
		Port:       aws.Int64(int64(p.config.ServicePort)),
		Protocol:   aws.String("HTTP"),
		TargetType: aws.String("instance"),
		VpcId:      vpc,
	})
	if err != nil {
		return nil, err
	}

	tgArn := ctg.TargetGroups[0].TargetGroupArn

	as := autoscaling.New(sess)

	_, err = as.CreateLaunchConfiguration(&autoscaling.CreateLaunchConfigurationInput{
		AssociatePublicIpAddress: aws.Bool(true),
		LaunchConfigurationName:  aws.String(serviceName),
		ImageId:                  &img.Image,
		InstanceType:             aws.String(p.config.InstanceType),
		KeyName:                  key,
		SecurityGroups:           groups,
		UserData:                 aws.String(ud),
	})

	if err != nil {
		return nil, err
	}

	var min, max, desired int64

	if p.config.Count == nil {
		min, max = 1, 1
	} else {
		min = p.config.Count.Min
		max = p.config.Count.Max
		desired = p.config.Count.Desired

		// Calculate sensible defaults. If no max is configured but a min is,
		// then it's a staticly sized group. If neither is configured, set
		// everything to 1.
		// And if desired isn't specified, max it the max.
		if max == 0 {
			if min == 0 {
				min = 1
				max = 1
			} else {
				max = min
			}
		}

	}

	if desired == 0 {
		desired = max
	}

	_, err = as.CreateAutoScalingGroup(&autoscaling.CreateAutoScalingGroupInput{
		AutoScalingGroupName:    aws.String(serviceName),
		DesiredCapacity:         aws.Int64(desired),
		MaxSize:                 aws.Int64(max),
		MinSize:                 aws.Int64(min),
		LaunchConfigurationName: aws.String(serviceName),
		VPCZoneIdentifier:       aws.String(subnetsComma),
		TargetGroupARNs:         []*string{tgArn},
	})
	if err != nil {
		return nil, err
	}

	st.Update("Waiting for the first instance to start...")

	err = as.WaitUntilGroupExists(&autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{aws.String(serviceName)},
	})
	if err != nil {
		return nil, err
	}

	var instances []*string

	// activityTicker periodically checks the ASG activity API to look for failure
	// or cancellation
	activityTicker := time.NewTicker(10 * time.Second)
	defer activityTicker.Stop()

	for {
		select {
		case <-activityTicker.C:
			dsa, err := as.DescribeScalingActivities(&autoscaling.DescribeScalingActivitiesInput{
				AutoScalingGroupName: aws.String(serviceName),
			})
			if err != nil {
				return nil, err
			}

			// check scaling activities, capture any non-nil status code and messages
			// to check for failures
			var lastStatusCode, lastStatusMessage string
			if len(dsa.Activities) > 0 {
				lastStatusCode = aws.StringValue(dsa.Activities[0].StatusCode)
				lastStatusMessage = aws.StringValue(dsa.Activities[0].StatusMessage)
			}

			if lastStatusCode == autoscaling.ScalingActivityStatusCodeFailed ||
				lastStatusCode == autoscaling.ScalingActivityStatusCodeCancelled {
				return nil, status.Errorf(codes.FailedPrecondition, "error setting up autoscaling group. Last status code (%s): \n Message: %s", lastStatusCode, lastStatusMessage)
			}
		default:
		}

		asg, err := as.DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{
			AutoScalingGroupNames: []*string{aws.String(serviceName)},
		})
		if err != nil {
			return nil, err
		}

		if len(asg.AutoScalingGroups[0].Instances) == 0 {
			time.Sleep(5 * time.Second)
			continue
		}

		for _, inst := range asg.AutoScalingGroups[0].Instances {
			instances = append(instances, inst.InstanceId)
		}

		break
	}

	err = e.WaitUntilInstanceRunning(&ec2.DescribeInstancesInput{
		InstanceIds: instances[:1],
	})
	if err != nil {
		return nil, err
	}

	out, err := e.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: instances[:1],
	})
	if err != nil {
		return nil, err
	}

	inst := out.Reservations[0].Instances[0]

	publicIp := *inst.NetworkInterfaces[0].Association.PublicIp
	publicDns := *inst.NetworkInterfaces[0].Association.PublicDnsName

	st.Close()

	ui.Output("EC2 ASG Instances launched: %s", publicIp, terminal.WithSuccessStyle())

	result := &Deployment{
		ServiceName:    serviceName,
		Region:         p.config.Region,
		PublicIp:       publicIp,
		PublicDns:      publicDns,
		TargetGroupArn: *tgArn,
	}

	return result, nil
}

// Destroy deletes the EC2 deployment.
func (p *Platform) Destroy(
	ctx context.Context,
	log hclog.Logger,
	deployment *Deployment,
	ui terminal.UI,
) error {
	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: p.config.Region,
		Logger: log,
	})
	if err != nil {
		return err
	}
	as := autoscaling.New(sess)

	_, err = as.DeleteAutoScalingGroup(&autoscaling.DeleteAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(deployment.ServiceName),
		ForceDelete:          aws.Bool(true),
	})
	if err != nil {
		log.Error("error deleting ASG", "error", err, "name", deployment.ServiceName)
		return err
	}

	_, err = as.DeleteLaunchConfiguration(&autoscaling.DeleteLaunchConfigurationInput{
		LaunchConfigurationName: aws.String(deployment.ServiceName),
	})
	if err != nil {
		ui.Output("error deleting lc: %s", err)
		return err
	}

	if deployment.TargetGroupArn != "" {
		_, err := elbv2.New(sess).DeleteTargetGroup(&elbv2.DeleteTargetGroupInput{
			TargetGroupArn: &deployment.TargetGroupArn,
		})
		if err != nil {
			ui.Output("error deleting tg: %s", err)
			return err
		}
	}

	return err
}

type countConfig struct {
	Desired int64 `hcl:"desired,optional"`
	Min     int64 `hcl:"min,optional"`
	Max     int64 `hcl:"max,optional"`
}

// Config is the configuration structure for the Platform.
type PlatformConfig struct {
	// AWS region to operate in
	Region string `hcl:"region"`

	Count *countConfig `hcl:"count,block"`

	// The type of instance to create
	InstanceType string `hcl:"instance_type"`

	// The key to associate with the instance
	Key string `hcl:"key,optional"`

	// The port that the service runs on within the instance.
	ServicePort int `hcl:"service_port"`

	// Additional ports to allow into the instance
	ExtraPorts []int `hcl:"extra_ports,optional"`

	// Additional security groups to add to the EC2 instance.
	SecurityGroups []string `hcl:"security_groups,optional"`

	// Subnet to put the instance into. Defaults to a public subnet in the default VPC.
	Subnet string `hcl:"subnet,optional"`
}

func (p *Platform) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&PlatformConfig{}), docs.FromFunc(p.DeployFunc()))
	if err != nil {
		return nil, err
	}

	doc.Description("Deploy the application into an AutoScaling Group on EC2")

	doc.Input("ami.Image")
	doc.Output("ec2.Deployment")

	doc.SetField(
		"region",
		"the AWS region to deploy into",
	)

	doc.SetField(
		"count",
		"how many EC2 instances to configure the ASG with",
		docs.Summary(
			"the fields here (desired, min, max) map directly to the typical ASG configuration",
		),
	)

	doc.SetField(
		"instance_type",
		"the EC2 instance type to deploy",
	)

	doc.SetField(
		"key",
		"the name of an SSH Key to associate with the instances, as preconfigured in EC2",
	)

	doc.SetField(
		"service_port",
		"the TCP port on the instances that the app will be running on",
	)

	doc.SetField(
		"extra_ports",
		"additional TCP ports to allow into the EC2 instances",
		docs.Summary(
			"these additional ports are usually used to allow secondary services, such as ssh",
		),
	)

	doc.SetField(
		"security_groups",
		"additional security groups to attached to the EC2 instances",
		docs.Summary(
			"this plugin creates security groups that match the above ports by default.",
			"this field allows additional security groups to be specified for the instances",
		),
	)

	doc.SetField(
		"subnet",
		"the subnet to place the instances into",
		docs.Default("a public subnet in the dafault VPC"),
	)

	return doc, nil
}

var (
	_ component.Platform     = (*Platform)(nil)
	_ component.Configurable = (*Platform)(nil)
	_ component.Destroyer    = (*Platform)(nil)
	_ component.Documented   = (*Platform)(nil)
)
