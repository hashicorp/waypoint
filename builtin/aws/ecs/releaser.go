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

	if target.TargetGroupArn == "" {
		log.Info("No target group configured, skipping release")
		return &Release{}, nil
	}

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

	// existingListener, if discovered, will be modified to introduce the new target group.
	var existingListener *elbv2.Listener
	var lbArn string

	switch lbRef := target.LbReference.(type) {
	case *Deployment_LoadBalancerArn:
		// We have a load balancer. Either discover the existing listener, or create a new one.
		lbArn = lbRef.LoadBalancerArn

		dlb, err := elbsrv.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{
			LoadBalancerArns: []*string{&lbRef.LoadBalancerArn},
		})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to describe load balancer %q", lbRef.LoadBalancerArn)
		}

		var lb *elbv2.LoadBalancer

		if len(dlb.LoadBalancers) == 0 {
			return nil, fmt.Errorf("no load balancers returned by DescribeLoadBalancers")
		}

		lb = dlb.LoadBalancers[0]

		// Now that we have the LB, set the hostname if necessary
		if hostname == "" {
			hostname = *lb.DNSName
		}

		listeners, err := elbsrv.DescribeListeners(&elbv2.DescribeListenersInput{
			LoadBalancerArn: lb.LoadBalancerArn,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to describe listener for lb %q", lb.LoadBalancerArn)
		}

		if len(listeners.Listeners) > 0 {
			if len(listeners.Listeners) > 1 {
				log.Warn("ALB has multiple listeners - only modifying the first listaner and ignoring all others")
			}
			existingListener = listeners.Listeners[0]
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
				return nil, errors.Wrapf(err, "failed to create listener")
			}

			// Do not set existingListener
		}

	case *Deployment_ListenerArn:
		lo, err := elbsrv.DescribeListeners(&elbv2.DescribeListenersInput{
			ListenerArns: []*string{&lbRef.ListenerArn},
		})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to describe listener %q", lbRef.ListenerArn)
		}
		if len(lo.Listeners) == 0 {
			return nil, errors.Errorf("listener %q not found", lbRef.ListenerArn)
		}
		existingListener = lo.Listeners[0]
		lbArn = *existingListener.LoadBalancerArn

		if hostname == "" {
			// We need to get the hostname from the existing alb

			dlb, err := elbsrv.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{
				LoadBalancerArns: []*string{existingListener.LoadBalancerArn},
			})
			if err != nil {
				return nil, errors.Wrapf(err, "failed to describe load balancer for listener %q", *existingListener.LoadBalancerArn)
			}

			if len(dlb.LoadBalancers) == 0 {
				return nil, fmt.Errorf("no load balancers returned by DescribeLoadBalancers")
			}

			hostname = *dlb.LoadBalancers[0].DNSName
		}

	}

	if existingListener != nil {
		log.Debug("configuring weight 100 for target group", "arn", target.TargetGroupArn)

		def := existingListener.DefaultActions

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

		log.Debug("modifying load balancer", "tgs", len(tgs))
		_, err = elbsrv.ModifyListener(&elbv2.ModifyListenerInput{
			ListenerArn: existingListener.ListenerArn,
			Port:        existingListener.Port,
			Protocol:    existingListener.Protocol,
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
			return nil, errors.Wrapf(err, "failed to modify listener %q to introduce new target group", existingListener.ListenerArn)
		}
	}

	return &Release{
		Url:             "http://" + hostname,
		LoadBalancerArn: lbArn,
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
