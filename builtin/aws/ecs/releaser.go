package ecs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/aws/utils"
)

// Releaser is the ReleaseManager implementation for Amazon ECS.
type Releaser struct {
	p      *Platform
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

// Release updates the load balancer for the ECS deployment
func (r *Releaser) Release(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	ui terminal.UI,
	target *Deployment,
) (*Release, error) {
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Initializing release...")
	defer s.Abort()

	if target.TargetGroupArn == "" {
		// This should only happen if someone disables the ALB in the deploy config.
		s.Update("Deployment did not define a target group - skipping release.")
		s.Done()
		return &Release{}, nil
	}

	if target.ListenerArn == "" {
		// This should only happen if someone disables the ALB in the deploy config.
		s.Update("Deployment did not define an ALB listener - skipping release.")
		s.Done()
		return &Release{}, nil
	}

	s.Update("Release initialized")
	s.Done()

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: r.p.config.Region,
		Logger: log,
	})
	if err != nil {
		return nil, err
	}
	elbsrv := elbv2.New(sess)

	var hostname string
	if r.p.config.ALB != nil && r.p.config.ALB.FQDN != "" {
		hostname = r.p.config.ALB.FQDN
	}

	tgs := []*elbv2.TargetGroupTuple{
		{
			TargetGroupArn: &target.TargetGroupArn,
			Weight:         aws.Int64(100),
		},
	}

	lo, err := elbsrv.DescribeListeners(&elbv2.DescribeListenersInput{
		ListenerArns: []*string{aws.String(target.ListenerArn)},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe listener %q", target.ListenerArn)
	}
	if len(lo.Listeners) == 0 {
		return nil, errors.Errorf("listener %q not found", target.ListenerArn)
	}
	listener := lo.Listeners[0]

	if hostname == "" {

		// We need to get the hostname from the existing alb
		if target.LoadBalancerArn == "" {
			s.Update("load balancer from deployment not specified - cannot determine hostname")
		} else {
			dlb, err := elbsrv.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{
				LoadBalancerArns: []*string{aws.String(target.LoadBalancerArn)},
			})
			if err != nil {
				return nil, errors.Wrapf(err, "failed to describe load balancer %q", target.LoadBalancerArn)
			}

			if len(dlb.LoadBalancers) == 0 {
				return nil, fmt.Errorf("no load balancers returned by DescribeLoadBalancers")
			}

			hostname = *dlb.LoadBalancers[0].DNSName
		}
	}

	log.Debug("configuring weight 100 for target group", "arn", target.TargetGroupArn)

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

	s = sg.Add("Checking that all targets are healthy...")
	targetHealth, err := elbsrv.DescribeTargetHealth(&elbv2.DescribeTargetHealthInput{
		TargetGroupArn: aws.String(target.TargetGroupArn),
	})
	if err != nil {
		log.Error("error getting target health", "err", err.Error())
		return nil, errors.Wrapf(err, "failed to describe health of target group with ARN %q", target.TargetGroupArn)
	}

	// Check each target to see if any one of them isn't healthy, before we
	// route 100% of traffic to it!
	for _, targetHealthDescription := range targetHealth.TargetHealthDescriptions {
		log.Debug("checking target health", "target", targetHealthDescription.Target.Id)
		// Possible states are: initial, healthy, unhealthy, unused, draining,
		// and unavailable.
		if *targetHealthDescription.TargetHealth.State != "healthy" {
			return nil, errors.Errorf("target (id: %s) is not healthy - will "+
				"only release when all targets in group (ARN: %q) are healthy",
				*targetHealthDescription.Target.Id, target.TargetGroupArn)
		}
	}
	s.Update("All targets are healthy!")
	s.Done()

	s = sg.Add("Modifying load balancer to introduce new target group %q", target.TargetGroupArn)
	log.Debug("modifying load balancer", "tgs", len(tgs))
	_, err = elbsrv.ModifyListener(&elbv2.ModifyListenerInput{
		ListenerArn: aws.String(target.ListenerArn),
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
		return nil, errors.Wrapf(err, "failed to modify listener %q to introduce new target group", target.ListenerArn)
	}

	s.Update("Finished ECS release")
	s.Done()

	return &Release{
		Url:             "http://" + hostname,
		LoadBalancerArn: target.LoadBalancerArn,
	}, nil
}

// ReleaserConfig is the configuration structure for the Releaser.
type ReleaserConfig struct{}

func (r *Releaser) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&ReleaserConfig{}))
	if err != nil {
		return nil, err
	}

	doc.Description("Reconfigures the ECS specific ALB to route traffic to new deployments")

	doc.Input("ecs.Deployment")
	doc.Output("ecs.Release")

	return doc, nil
}

var (
	_ component.ReleaseManager = (*Releaser)(nil)
	_ component.Configurable   = (*Releaser)(nil)
)
