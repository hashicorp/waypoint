package ecs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/aws/utils"
)

// Releaser is the ReleaseManager implementation for Amazon ECS.
type Releaser struct {
	config ReleaserConfig
}

// Config implements Configurable
func (r *Releaser) Config() (interface{}, error) {
	return &r.config, nil
}

// ReleaseFunc implements component.ReleaseManager
func (r *Releaser) ReleaseFunc() interface{} {
	return r.Release
}

// Release creates a Kubernetes service configured for the deployment
func (r *Releaser) Release(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	ui terminal.UI,
	target *Deployment,
) (*Release, error) {
	log.Debug("releasing deployment target group in to load balancer")

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: r.config.Region,
	})
	if err != nil {
		return nil, err
	}
	elbsrv := elbv2.New(sess)

	dlb, err := elbsrv.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{
		LoadBalancerArns: []*string{&target.LoadBalancerArn},
	})
	if err != nil {
		return nil, err
	}

	var lb *elbv2.LoadBalancer

	if len(dlb.LoadBalancers) == 0 {
		return nil, fmt.Errorf("No load balancers returned by DescribeLoadBalancers")
	}

	lb = dlb.LoadBalancers[0]

	listeners, err := elbsrv.DescribeListeners(&elbv2.DescribeListenersInput{
		LoadBalancerArn: lb.LoadBalancerArn,
	})
	if err != nil {
		return nil, err
	}

	tgs := []*elbv2.TargetGroupTuple{
		{
			TargetGroupArn: &target.TargetGroupArn,
			Weight:         aws.Int64(100),
		},
	}

	log.Debug("configuring weight 100 for target group", "arn", target.TargetGroupArn)

	scheme := "http"

	if len(listeners.Listeners) > 0 {
		for _, listener := range listeners.Listeners {
			def := listener.DefaultActions

			if len(def) > 0 && def[0].ForwardConfig != nil {
				for _, tg := range def[0].ForwardConfig.TargetGroups {
					// Drain any target groups to 0 but leave them registered.
					// This loop also inherently removes any target groups already
					// set to 0 that ARE NOT the one we're releasing.
					if *tg.Weight > 0 && *tg.TargetGroupArn != target.TargetGroupArn {
						tg.Weight = aws.Int64(0)
						tgs = append(tgs, tg)
						log.Debug("previous target group", "arn", *tg.TargetGroupArn)
					}
				}
			}

			log.Debug("modifying listener", "tgs", len(tgs))
			_, err = elbsrv.ModifyListener(&elbv2.ModifyListenerInput{
				ListenerArn: listener.ListenerArn,
				Port:        listener.Port,
				Protocol:    listener.Protocol,
				DefaultActions: []*elbv2.Action{
					{
						ForwardConfig: &elbv2.ForwardActionConfig{
							TargetGroups: tgs,
						},
						Type: aws.String("forward"),
					},
				},
			})
			if err != nil {
				return nil, err
			}

			if listener.Port == aws.Int64(443) {
				scheme = "https"
			}
		}
	} else {
		log.Info("load-balancer defined", "dns-name", *lb.DNSName)

		_, err := elbsrv.CreateListener(&elbv2.CreateListenerInput{
			LoadBalancerArn: lb.LoadBalancerArn,
			Port:            aws.Int64(80),
			Protocol:        aws.String("HTTP"),
			DefaultActions: []*elbv2.Action{
				{
					ForwardConfig: &elbv2.ForwardActionConfig{
						TargetGroups: tgs,
					},
					Type: aws.String("forward"),
				},
			},
		})

		if err != nil {
			return nil, err
		}
	}

	hostname := *lb.DNSName
	log.Debug("ALB hostname", hostname)

	return &Release{
		Url:             scheme + "://" + hostname,
		LoadBalancerArn: *lb.LoadBalancerArn,
	}, nil
}

// ReleaserConfig is the configuration structure for the Releaser.
type ReleaserConfig struct {
	// AWS Region to deploy into
	Region string `hcl:"region"`
}

func (r *Releaser) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&ReleaserConfig{}))
	if err != nil {
		return nil, err
	}

	doc.Description("Reconfigures the ECS specific ALB to route traffic to new deployments")

	doc.Input("ecs.Deployment")
	doc.Output("ecs.Release")

	doc.SetField(
		"region",
		"the AWS region for the ECS cluster",
	)

	return doc, nil
}

var (
	_ component.ReleaseManager = (*Releaser)(nil)
	_ component.Configurable   = (*Releaser)(nil)
)
