// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package alb

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/framework/resource"
	sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/aws/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Releaser struct {
	config ReleaserConfig
}

const (
	targetGroupInitializationTimeoutSeconds         int = 120
	targetGroupInitializationPollingIntervalSeconds int = 5
)

// Config implements Configurable
func (r *Releaser) Config() (interface{}, error) {
	return &r.config, nil
}

// ReleaseFunc implements component.ReleaseManager
func (r *Releaser) ReleaseFunc() interface{} {
	return r.Release
}

// StatusFunc implements component.Status
func (r *Releaser) StatusFunc() interface{} {
	return r.Status
}

// DestroyFunc implements component.Destroyer
func (r *Releaser) DestroyFunc() interface{} {
	return r.Destroy
}

func (r *Releaser) resourceManager(log hclog.Logger) *resource.Manager {
	return resource.NewManager(
		resource.WithLogger(log.Named("resource_manager")),
		resource.WithValueProvider(r.getSession),
		resource.WithResource(resource.NewResource(
			resource.WithName("security_group"),
			resource.WithState(&Resource_SecurityGroup{}),
			resource.WithCreate(r.resourceSecurityGroupCreate),
			resource.WithDestroy(r.resourceSecurityGroupDestroy),
		)),

		resource.WithResource(resource.NewResource(
			resource.WithName("load_balancer"),
			resource.WithState(&Resource_LoadBalancer{}),
			resource.WithCreate(r.resourceLoadBalancerCreate),
			resource.WithDestroy(r.resourceLoadBalancerDestroy),
		)),

		resource.WithResource(resource.NewResource(
			resource.WithName("listener"),
			resource.WithState(&Resource_Listener{}),
			resource.WithCreate(r.resourceListenerCreate),
			resource.WithDestroy(r.resourceListenerDestroy),
		)),

		resource.WithResource(resource.NewResource(
			resource.WithName("record_set"),
			resource.WithState(&Resource_RecordSet{}),
			resource.WithCreate(r.resourceRecordSetCreate),
			resource.WithDestroy(r.resourceRecordSetDestroy),
		)),
	)
}

func (r *Releaser) resourceSecurityGroupCreate(
	ctx context.Context,
	log hclog.Logger,
	sess *session.Session,
	sg terminal.StepGroup,
	ui terminal.UI,
	lbName string,
	port int64,
	state []*Resource_SecurityGroup,
) error {
	if r.config.ListenerARN != "" {
		// a custom listner is being used, so we assume the Load Balancer is already
		// created, and alrady has a security group associated with it.
		log.Debug("custom listener arn found, skipping security group create")
		return nil
	}

	ec2srv := ec2.New(sess)

	// Figure out what subnets we're using. This will give us the VPC
	// information.
	var (
		vpc     *string
		subnets []*string
	)
	if len(r.config.Subnets) == 0 {
		_, v, err := utils.DefaultSubnets(ctx, sess)
		if err != nil {
			return err
		}

		vpc = v
	} else {
		for _, s := range r.config.Subnets {
			subnets = append(subnets, aws.String(s))
		}

		subnetInfo, err := ec2srv.DescribeSubnets(&ec2.DescribeSubnetsInput{
			SubnetIds: []*string{subnets[0]},
		})
		if err != nil {
			return err
		}

		vpc = subnetInfo.Subnets[0].VpcId
	}

	if len(r.config.SecurityGroupIDs) > 0 {
		log.Debug("user specified alb security groups found, skipping security group create")

		for _, sGroup := range r.config.SecurityGroupIDs {
			state = append(state, &Resource_SecurityGroup{Id: sGroup, Managed: false})
		}
	} else {
		sGroup, err := utils.CreateSecurityGroup(ctx, sess, fmt.Sprintf("%s-incoming", lbName), vpc, int(port))
		if err != nil {
			return err
		}
		state = append(state, &Resource_SecurityGroup{Id: *sGroup, Managed: true})
	}

	return nil
}

func (r *Releaser) resourceSecurityGroupDestroy(
	ctx context.Context,
	sess *session.Session,
	state []*Resource_SecurityGroup,
	log hclog.Logger,
	sg terminal.StepGroup,
) error {
	if len(state) >= 1 && len(r.config.SecurityGroupIDs) == 0 {
		for _, sgState := range state {
			step := sg.Add("Destroying Security Group...")
			defer step.Abort()
			log.Debug("deleting security group", "sg-id", sgState.Id)
			ec2Svc := ec2.New(sess)
			for i := 0; i < 20; i++ {
				_, err := ec2Svc.DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{
					GroupId: &sgState.Id,
				})
				if err == nil {
					step.Done()
					return nil
				}
				// if we encounter an unrecoverable error, exit now.
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case "DependencyViolation":
						time.Sleep(2 * time.Second)
						continue
					case "InvalidGroup.NotFound":
						log.Debug("security group not found", "sg-id", sgState.Id)
						return nil
					default:
						return err
					}
				}
				return err
			}
			step.Update("Destroyed Security Group")
			step.Done()
		}
	} else if len(r.config.SecurityGroupIDs) >= 0 {
		log.Debug("not deleting user-provided security groups", "sg-id", r.config.SecurityGroupIDs)
		return nil
	}
	log.Debug("no security group id found, continuing")

	return nil
}

func (r *Releaser) resourceLoadBalancerCreate(
	ctx context.Context,
	log hclog.Logger,
	sess *session.Session,
	sg terminal.StepGroup,
	ui terminal.UI,
	lbName string,
	port int64,
	target *TargetGroup,
	sgState []*Resource_SecurityGroup,
	state *Resource_LoadBalancer,
) error {

	elbsrv := elbv2.New(sess)
	// if a custom listner is being used, we assume the Load Balancer is already
	// created, and we can skip some steps
	if r.config.ListenerARN != "" {
		log.Debug("custom listener arn found, skipping load balancer creation")
		out, err := elbsrv.DescribeListeners(&elbv2.DescribeListenersInput{
			ListenerArns: []*string{aws.String(r.config.ListenerARN)},
		})
		if err != nil {
			return err
		}

		if len(out.Listeners) == 0 {
			return status.Errorf(codes.NotFound, "no listeners found for arn: %s", r.config.ListenerARN)
		}

		listener := out.Listeners[0]
		state.Arn = *listener.LoadBalancerArn
		dlb, err := elbsrv.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{
			LoadBalancerArns: []*string{listener.LoadBalancerArn},
		})
		if err != nil {
			return err
		}

		if dlb != nil && len(dlb.LoadBalancers) > 0 {
			state.DnsName = *dlb.LoadBalancers[0].DNSName
			state.ZoneId = *dlb.LoadBalancers[0].CanonicalHostedZoneId
			return nil
		}
		return status.Errorf(codes.NotFound, "no matching load balancer found")
	}

	var (
		subnets []*string
		err     error
	)
	if len(r.config.Subnets) == 0 {
		subnets, _, err = utils.DefaultSubnets(ctx, sess)
		if err != nil {
			return err
		}
	} else {
		for _, s := range r.config.Subnets {
			subnets = append(subnets, aws.String(s))
		}
	}

	var lb *elbv2.LoadBalancer
	// see if the load balancer may alread exist
	dlb, err := elbsrv.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{
		Names: []*string{&lbName},
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case elbv2.ErrCodeLoadBalancerNotFoundException:
				// fine, means we'll create it.
			default:
				return err
			}
		} else {
			return err
		}
	}

	if dlb != nil && len(dlb.LoadBalancers) > 0 {
		lb = dlb.LoadBalancers[0]
	} else {
		var sgs []*string
		for _, s := range sgState {
			sgs = append(sgs, &s.Id)
		}

		clb, err := elbsrv.CreateLoadBalancer(&elbv2.CreateLoadBalancerInput{
			Name:           aws.String(lbName),
			Subnets:        subnets,
			SecurityGroups: sgs,
		})
		if err != nil {
			return err
		}

		lb = clb.LoadBalancers[0]
	}

	if lb == nil {
		// this shouldn't happen but let's be safe from nil-pointer de-ref's below
		return status.Errorf(codes.Internal, "failed to create load balancer, this should not happen.")
	}

	state.Arn = *lb.LoadBalancerArn
	state.DnsName = *lb.DNSName
	state.ZoneId = *lb.CanonicalHostedZoneId

	return nil
}

func (r *Releaser) resourceLoadBalancerDestroy(
	ctx context.Context,
	sess *session.Session,
	sg terminal.StepGroup,
	state *Resource_LoadBalancer,
) error {
	// if the configuration does not specify existing load balancer name and the
	// listener arn, we assume that our release created both of these and can
	// reasonably destroy the load balancer completely. Because Listeners cannot
	// exist without a Load Balancer, we assume at this point that the listener is
	// destroyed.
	if r.config.Name == "" && r.config.ListenerARN == "" {
		step := sg.Add("Destroying Load Balancer...")
		defer step.Abort()
		elbsrv := elbv2.New(sess)
		_, err := elbsrv.DeleteLoadBalancer(&elbv2.DeleteLoadBalancerInput{
			LoadBalancerArn: &state.Arn,
		})
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				if aerr.Code() == elbv2.ErrCodeLoadBalancerNotFoundException {
					step.Update("Load balancer not found: %s", state.Arn)
					step.Done()
					return nil
				}
			}
			return err
		}
		step.Update("Destroyed Load Balancer")
		step.Done()
	}
	return nil
}

func (r *Releaser) resourceListenerCreate(
	ctx context.Context,
	log hclog.Logger,
	sess *session.Session,
	sg terminal.StepGroup,
	ui terminal.UI,
	port int64,
	target *TargetGroup,
	lbState *Resource_LoadBalancer,
	state *Resource_Listener,
) error {
	state.TgArn = target.Arn

	elbsrv := elbv2.New(sess)
	var (
		certs    []*elbv2.Certificate
		protocol string = "HTTP"
		listener *elbv2.Listener
	)
	shouldCreateListener := true

	if r.config.CertificateId != "" {
		protocol = "HTTPS"
		port = 443
		certs = append(certs, &elbv2.Certificate{
			CertificateArn: &r.config.CertificateId,
		})
	}

	if r.config.ListenerARN != "" {
		listeners, err := elbsrv.DescribeListeners(&elbv2.DescribeListenersInput{
			ListenerArns: []*string{&r.config.ListenerARN},
		})
		if err != nil {
			return err
		}

		if len(listeners.Listeners) > 0 {
			listener = listeners.Listeners[0]
			state.Arn = *listener.ListenerArn
			shouldCreateListener = false
		} else {
			return status.Errorf(codes.NotFound, "no listeners found for arn: %s", r.config.ListenerARN)
		}
	}

	tgs := []*elbv2.TargetGroupTuple{
		{
			TargetGroupArn: aws.String(target.Arn),
			Weight:         aws.Int64(100),
		},
	}

	if shouldCreateListener {
		// create the listener to forward traffic to the target group
		lo, err := elbsrv.CreateListener(&elbv2.CreateListenerInput{
			LoadBalancerArn: aws.String(lbState.Arn),
			Port:            aws.Int64(port),
			Protocol:        aws.String(protocol),
			Certificates:    certs,
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
			log.Warn("error creating listener", "error", err)
			if aerr, ok := err.(awserr.Error); ok {
				// CreateListener is idempotent, but if params change, such as `DefaultActions`,
				// it will return a DuplicateListener error
				switch aerr.Code() {
				case elbv2.ErrCodeDuplicateListenerException:
					// if DuplicateListener, we actually want to ModifyListener to use
					// the newly created TargetGroup & Target pair.
					//
					// this requires us to fetch the listener ARN from the load balancer
					log.Info("duplicate listener found, attempting to modify listener")
					dlo, err := elbsrv.DescribeListeners(&elbv2.DescribeListenersInput{
						LoadBalancerArn: aws.String(lbState.Arn),
					})
					if err != nil {
						return err
					}

					// find the listener on our port
					for _, ln := range dlo.Listeners {
						if *ln.Port == port {
							listener = ln
						}
					}

					mlo, err := elbsrv.ModifyListener(&elbv2.ModifyListenerInput{
						ListenerArn:  listener.ListenerArn,
						Port:         aws.Int64(port),
						Protocol:     aws.String(protocol),
						Certificates: certs,
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
						return err
					}

					listener = mlo.Listeners[0]
					state.Arn = *listener.ListenerArn

					log.Info("modified listener", "arn", state.Arn)
					// exit
					return nil
				}
			}
			return err
		}
		listener = lo.Listeners[0]
		state.Arn = *listener.ListenerArn
	} else {
		// modify the existing listener
		def := listener.DefaultActions

		if len(def) > 0 && def[0].ForwardConfig != nil {
			for _, tg := range def[0].ForwardConfig.TargetGroups {
				if *tg.Weight > 0 && *tg.TargetGroupArn != target.Arn {
					tg.Weight = aws.Int64(0)
					tgs = append(tgs, tg)
					log.Debug("previous target group", "arn", *tg.TargetGroupArn)
				}
			}
		}

		_, err := elbsrv.ModifyListener(&elbv2.ModifyListenerInput{
			ListenerArn:  listener.ListenerArn,
			Port:         aws.Int64(port),
			Protocol:     aws.String(protocol),
			Certificates: certs,
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
			return err
		}
	}

	return nil
}

func (r *Releaser) resourceListenerDestroy(
	ctx context.Context,
	sess *session.Session,
	log hclog.Logger,
	sg terminal.StepGroup,
	state *Resource_Listener,
) error {
	elbsrv := elbv2.New(sess)
	if r.config.ListenerARN == "" {
		// if the configuration does not specify an existing listener, we can simply
		// delete the listener assuming waypoint created it
		step := sg.Add("Destroying Listener...")
		defer step.Abort()
		_, err := elbsrv.DeleteListener(&elbv2.DeleteListenerInput{
			ListenerArn: &state.Arn,
		})
		if err != nil {
			return err
		}
		step.Update("Listener destroyed")
		step.Done()
	} else {
		step := sg.Add("Removing forwarding action...")
		defer step.Abort()
		// if an existing listener is specified, we can at least remove the action
		// that forwarded to our target group
		out, err := elbsrv.DescribeListeners(&elbv2.DescribeListenersInput{
			ListenerArns: []*string{&state.Arn},
		})
		if err != nil {
			return err
		}

		if len(out.Listeners) == 0 {
			return nil
		}
		listener := out.Listeners[0]

		var actions []*elbv2.Action
		defaultActions := listener.DefaultActions
	LLOOP:
		for i, def := range defaultActions {
			if def.ForwardConfig != nil {
				for _, tg := range def.ForwardConfig.TargetGroups {
					if *tg.TargetGroupArn == state.TgArn {
						// Waypoint adds a single default action. If there are more than 1,
						// simply remove the one in question and leave the others alone
						if len(defaultActions) > 1 {
							// remove this item in-place
							copy(defaultActions[i:], defaultActions[i+1:])
						}
						actions = defaultActions[:len(defaultActions)-1]
						break LLOOP
					}
				}
			}
		}
		log.Debug("modifying actions")
		if len(actions) > 0 {
			step.Update("modifying listener: %s", state.Arn)
			_, err := elbsrv.ModifyListener(&elbv2.ModifyListenerInput{
				ListenerArn:    &state.Arn,
				DefaultActions: actions,
			})
			if err != nil {
				return err
			}
		}
		step.Update("Listener destroyed")
		step.Done()
	}
	return nil
}

func (r *Releaser) resourceRecordSetCreate(
	ctx context.Context,
	log hclog.Logger,
	sess *session.Session,
	sg terminal.StepGroup,
	ui terminal.UI,
	lbState *Resource_LoadBalancer,
	state *Resource_RecordSet,
) error {
	if r.config.ZoneId != "" {
		step := sg.Add("Updating Route53...")
		defer step.Abort()
		r53 := route53.New(sess)

		records, err := r53.ListResourceRecordSets(&route53.ListResourceRecordSetsInput{
			HostedZoneId:    aws.String(r.config.ZoneId),
			StartRecordName: aws.String(r.config.FQDN),
			StartRecordType: aws.String("A"),
			MaxItems:        aws.String("1"),
		})
		if err != nil {
			return err
		}

		fqdn := r.config.FQDN

		if fqdn[len(fqdn)-1] != '.' {
			fqdn += "."
		}

		state.FQDN = fqdn

		if len(records.ResourceRecordSets) > 0 && *(records.ResourceRecordSets[0].Name) == fqdn {
			log.Debug("found existing record, assuming it's correct")
		} else {
			log.Debug("creating new route53 record", "zone-id", r.config.ZoneId)
			input := &route53.ChangeResourceRecordSetsInput{
				ChangeBatch: &route53.ChangeBatch{
					Changes: []*route53.Change{
						{
							Action: aws.String("CREATE"),
							ResourceRecordSet: &route53.ResourceRecordSet{
								Name: aws.String(r.config.FQDN),
								Type: aws.String("A"),
								AliasTarget: &route53.AliasTarget{
									DNSName:              &lbState.DnsName,
									EvaluateTargetHealth: aws.Bool(true),
									HostedZoneId:         &lbState.ZoneId,
								},
							},
						},
					},
					Comment: aws.String("managed by waypoint"),
				},
				HostedZoneId: aws.String(r.config.ZoneId),
			}

			result, err := r53.ChangeResourceRecordSets(input)
			if err != nil {
				return err
			}
			log.Debug("record created", "change-id", *result.ChangeInfo.Id)
		}
		step.Update("Route53 record set created")
		step.Done()
	}

	return nil
}

func (r *Releaser) resourceRecordSetDestroy(
	ctx context.Context,
	sess *session.Session,
	log hclog.Logger,
	sg terminal.StepGroup,
	lbState *Resource_LoadBalancer,
	state *Resource_RecordSet,
) error {

	if state.FQDN == "" {
		// nothing to do
		log.Debug("no route 53 information")
		return nil
	}

	step := sg.Add("Removing Route53 record...")
	defer step.Abort()
	r53 := route53.New(sess)

	records, err := r53.ListResourceRecordSets(&route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(r.config.ZoneId),
		StartRecordName: aws.String(state.FQDN),
		StartRecordType: aws.String("A"),
		MaxItems:        aws.String("1"),
	})
	if err != nil {
		return err
	}

	fqdn := r.config.FQDN

	if fqdn[len(fqdn)-1] != '.' {
		fqdn += "."
	}

	if len(records.ResourceRecordSets) > 0 && *(records.ResourceRecordSets[0].Name) != fqdn {
		log.Debug("no existing record found for (%s)", fqdn)
		step.Update("No existing records found to remove")
		step.Done()
		return nil
	}

	log.Debug("removing route53 record", "zone-id", r.config.ZoneId)
	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("DELETE"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(state.FQDN),
						Type: aws.String("A"),
						AliasTarget: &route53.AliasTarget{
							DNSName:              &lbState.DnsName,
							EvaluateTargetHealth: aws.Bool(true),
							HostedZoneId:         &lbState.ZoneId,
						},
					},
				},
			},
			Comment: aws.String("managed by waypoint"),
		},
		HostedZoneId: aws.String(r.config.ZoneId),
	}

	result, err := r53.ChangeResourceRecordSets(input)
	if err != nil {
		return err
	}
	log.Debug("record destroyed", "change-id", *result.ChangeInfo.Id)

	step.Update("Route53 record set removed")
	step.Done()

	return nil
}

func (r *Releaser) getSession(
	_ context.Context,
	log hclog.Logger,
	target *TargetGroup,
) (*session.Session, error) {
	return utils.GetSession(&utils.SessionConfig{
		Region: target.Region,
		Logger: log,
	})
}

// Release manages target group attachment to a configured ALB
func (r *Releaser) Release(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	ui terminal.UI,
	target *TargetGroup,
) (*Release, error) {
	sg := ui.StepGroup()
	defer sg.Wait()

	lbName := r.config.Name
	if lbName == "" {
		lbName = "waypoint-" + src.App
	}

	// We have to clamp at a length of 32 because the Name field to
	// CreateLoadBalancer requires that the name is 32 characters or less.
	if len(lbName) > 32 {
		lbName = lbName[:32]
		log.Debug("using a shortened value for load balancer name due to AWS's length limits", "lbName", lbName)
	}

	// If there is a port defined, honor it.
	port := int64(80)
	if r.config.Port != 0 {
		port = int64(r.config.Port)
	}

	// Create our resource manager and create
	rm := r.resourceManager(log)
	if err := rm.CreateAll(
		ctx, log, sg, ui, src,
		target, lbName, port,
	); err != nil {
		return nil, err
	}

	// Get our load balancer state to verify
	lbState := rm.Resource("load_balancer").State().(*Resource_LoadBalancer)
	if lbState == nil {
		return nil, status.Errorf(codes.Internal, "load balancer state is nil, this should never happen")
	}

	return &Release{
		Region:          target.Region,
		TargetGroupArn:  target.Arn,
		Url:             "http://" + lbState.DnsName,
		LoadBalancerArn: lbState.Arn,
		ResourceState:   rm.State(),
	}, nil
}

func (r *Releaser) Status(
	ctx context.Context,
	log hclog.Logger,
	release *Release,
	ui terminal.UI,
) (*sdk.StatusReport, error) {

	var report sdk.StatusReport
	report.External = true
	defer func() {
		report.GeneratedTime = timestamppb.Now()
	}()

	if release.Region == "" {
		log.Debug("Region is not available for this release. Unable to determine status.")
		return &report, nil
	}

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: release.Region,
		Logger: log,
	})
	if err != nil {
		return nil, err
	}

	elbsrv := elbv2.New(sess)

	sg := ui.StepGroup()
	defer sg.Wait()

	step := sg.Add("Gathering health report for AWS/ALB platform...")
	defer step.Abort()

	step.Update("Waiting for at least one target to pass initialization...")

	var targetHealthDescriptions []*elbv2.TargetHealthDescription

	startTime := time.Now().Unix()
	for {
		if targetHealthDescriptions != nil {
			sleepDuration := time.Second * time.Duration(targetGroupInitializationPollingIntervalSeconds)
			log.Debug("Sleeping %0.f seconds to give the following target group time to initialize:\n%s", sleepDuration.Seconds(), release.TargetGroupArn)
			time.Sleep(sleepDuration)
		}

		if startTime+int64(targetGroupInitializationTimeoutSeconds) <= time.Now().Unix() {
			report.HealthMessage = fmt.Sprintf("timed out after %d seconds waiting for the following target group to initialize:\n%s", time.Now().Unix()-startTime, release.TargetGroupArn)
			report.Health = sdk.StatusReport_UNKNOWN
			step.Status(terminal.StatusWarn)
			step.Update(report.HealthMessage)
			return &report, nil
		}

		tgHealthResp, err := elbsrv.DescribeTargetHealthWithContext(ctx, &elbv2.DescribeTargetHealthInput{
			TargetGroupArn: &release.TargetGroupArn,
		})
		if err != nil {
			return nil, status.Errorf(codes.Unknown, "failed to describe target group %s health: %s", release.TargetGroupArn, err)
		}

		targetHealthDescriptions = tgHealthResp.TargetHealthDescriptions

		// We may not have any targets if the target group was created very recently.
		if len(targetHealthDescriptions) == 0 {
			step.Update("Waiting for registered targets with health...")
			continue
		}

		initializingCount := 0
		for _, tgHealth := range targetHealthDescriptions {
			// NOTE(izaaklauer) potentially unsafe dereference
			if *tgHealth.TargetHealth.State == elbv2.TargetHealthStateEnumInitial {
				initializingCount++
			}
		}
		if initializingCount == len(targetHealthDescriptions) {
			step.Update("Waiting for at least one target to finish initializing...")
			continue
		}

		step.Update("Target group has been initialized.")
		break
	}

	report.Resources = []*sdk.StatusReport_Resource{}

	healthyCount := 0
	for _, tgHealth := range targetHealthDescriptions {

		targetId := *tgHealth.Target.Id

		var health sdk.StatusReport_Health

		switch *tgHealth.TargetHealth.State {
		case elbv2.TargetHealthStateEnumHealthy:
			healthyCount++
			health = sdk.StatusReport_READY
		case elbv2.TargetHealthStateEnumUnavailable:
			// Lambda functions present this way. Defaulting to UNKNOWN seems reasonable here too.
			healthyCount++
			health = sdk.StatusReport_READY
		case elbv2.TargetHealthStateEnumUnhealthy:
			health = sdk.StatusReport_DOWN
		default:
			// There are more TargetHealthStateEnums, but they do not cleanly map to our states.
			health = sdk.StatusReport_UNKNOWN
		}

		var healthMessage string
		if tgHealth.TargetHealth.Description != nil {
			healthMessage = *tgHealth.TargetHealth.Description
		}

		report.Resources = append(report.Resources, &sdk.StatusReport_Resource{
			Health:        health,
			HealthMessage: healthMessage,
			Name:          targetId,
		})
	}

	step.Update("Finished building report for AWS/ALB platform")
	step.Done()

	// NOTE(briancain): Replace ui.Status with StepGroups once this bug
	// has been fixed: https://github.com/hashicorp/waypoint/issues/1536
	st := ui.Status()
	defer st.Close()

	// If AWS registers targets slowly or incrementally, we may report an artificially low number of total targets.
	totalTargets := len(targetHealthDescriptions)

	if healthyCount == totalTargets {
		report.Health = sdk.StatusReport_READY
		report.HealthMessage = fmt.Sprintf("All %d targets are healthy.", totalTargets)
	} else if healthyCount > 0 {
		report.Health = sdk.StatusReport_PARTIAL
		report.HealthMessage = fmt.Sprintf("Only %d/%d targets are healthy.", healthyCount, totalTargets)
		st.Step(terminal.StatusWarn, report.HealthMessage)
	} else {
		report.Health = sdk.StatusReport_DOWN
		report.HealthMessage = fmt.Sprintf("All targets are unhealthy, however your application might be available or still starting up.")
		st.Step(terminal.StatusWarn, report.HealthMessage)
	}

	return &report, nil
}

// Destroy will modify or delete Listeners, so that the platform can destroy the
// target groups
func (r *Releaser) Destroy(
	ctx context.Context,
	log hclog.Logger,
	release *Release,
	ui terminal.UI,
) error {
	sg := ui.StepGroup()
	defer sg.Wait()
	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: release.Region,
		Logger: log,
	})
	if err != nil {
		return err
	}

	rm := r.resourceManager(log)

	// If we don't have resource state, this state is from an older version
	// and we need to manually recreate it. Older versions did not save anything
	// for Security Group Ids
	if release.ResourceState == nil {
		rm.Resource("security_group").SetState(&Resource_SecurityGroup{})
		rm.Resource("load_balancer").SetState(&Resource_LoadBalancer{
			Arn: release.LoadBalancerArn,
			// resourceRecordSetDestroy relies on the Load Balancer state having the
			// Zone Id present.
			ZoneId: r.config.ZoneId,
		})
		l, err := findListenerArnFromTargetGroup(sess, log, release.TargetGroupArn, release.LoadBalancerArn)
		if err != nil {
			return err
		}
		rm.Resource("listener").SetState(&Resource_Listener{
			TgArn: release.TargetGroupArn,
			Arn:   l,
		})
		// older releases don't have any security group information
		rm.Resource("record_set").SetState(&Resource_RecordSet{
			FQDN: r.config.FQDN,
		})
	} else {
		// Load our set state
		if err := rm.LoadState(release.ResourceState); err != nil {
			return err
		}
	}

	// Destroy All
	return rm.DestroyAll(ctx, log, sg, ui, sess)
}

func findListenerArnFromTargetGroup(
	sess *session.Session,
	log hclog.Logger,
	tgArn string,
	lbArn string,
) (string, error) {

	var listener *elbv2.Listener

	elbsrv := elbv2.New(sess)
	out, err := elbsrv.DescribeListeners(&elbv2.DescribeListenersInput{
		LoadBalancerArn: &lbArn,
	})
	if err != nil {
		return "", err
	}

LLOOP:
	// We don't know the ARN of the listener, only the Target Group and the
	// LoadBalancer. We need to loop through each Listener's actions to find the
	// target group that matches.
	for j, ol := range out.Listeners {
		for _, def := range ol.DefaultActions {
			if def.ForwardConfig != nil {
				for _, tg := range def.ForwardConfig.TargetGroups {
					if *tg.TargetGroupArn == tgArn {
						listener = out.Listeners[j]
						break LLOOP
					}
				}
			}
		}
	}

	if listener != nil {
		return *listener.ListenerArn, nil
	}
	return "", status.Errorf(codes.NotFound,
		"unable to find listener")
}

// ReleaserConfig is the configuration structure for the Releaser.
type ReleaserConfig struct {
	Name string `hcl:"name,optional"`

	// Port configures the port that is used to access the service.
	// The default is 80.
	Port int `hcl:"port,optional"`

	// Subnets to place the service into. Defaults to the subnets in the default VPC.
	Subnets []string `hcl:"subnets,optional"`

	// Certificate ARN to attach to the load balancer
	CertificateId string `hcl:"certificate,optional"`

	// Route53 Zone to setup record in
	ZoneId string `hcl:"zone_id,optional"`

	// Fully qualified domain name of the record to create in the target zone id
	FQDN string `hcl:"domain_name,optional"`

	// When set, waypoint will configure the target group into the specified
	// ALB Listener ARN. This allows for usage of existing ALBs.
	ListenerARN string `hcl:"listener_arn,optional"`

	// Existing Security Group ID to use for ALB.
	SecurityGroupIDs []string `hcl:"security_group_ids,optional"`
}

func (r *Releaser) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&ReleaserConfig{}))
	if err != nil {
		return nil, err
	}

	doc.Description("Release target groups by attaching them to an ALB")

	doc.Input("alb.TargetGroup")
	doc.Output("alb.Release")
	doc.AddMapper(
		"ec2.Deployment",
		"alb.TargetGroup",
		"Allow EC2 Deployments to be hooked up to an ALB",
	)

	doc.AddMapper(
		"lambda.Deployment",
		"alb.TargetGroup",
		"Allow Lambda Deployments to be hooked up to an ALB",
	)

	doc.SetField(
		"name",
		"the name to assign the ALB",
		docs.Summary(
			"names have to be unique per region",
		),
		docs.Default("derived from application name"),
	)

	doc.SetField(
		"port",
		"the TCP port to configure the ALB to listen on",
		docs.Default("80 for HTTP, 443 for HTTPS"),
	)

	doc.SetField(
		"subnets",
		"the subnet ids to allow the ALB to run in",
		docs.Default("public subnets in the account default VPC"),
	)

	doc.SetField(
		"certificate",
		"ARN for the certificate to install on the ALB listener",
		docs.Summary(
			"when this is set, the port automatically changes to 443 unless",
			"overriden in this configuration",
		),
	)

	doc.SetField(
		"zone_id",
		"Route53 ZoneID to create a DNS record into",
		docs.Summary(
			"set along with domain_name to have DNS automatically setup for the ALB",
		),
	)

	doc.SetField(
		"domain_name",
		"Fully qualified domain name to set for the ALB",
		docs.Summary(
			"set along with zone_id to have DNS automatically setup for the ALB.",
			"this value should include the full hostname and domain name, for instance",
			"app.example.com",
		),
	)

	doc.SetField(
		"listener_arn",
		"the ARN on an existing ALB to configure",
		docs.Summary(
			"when this is set, no ALB or Listener is created. Instead the application is",
			"configured by manipulating this existing Listener. This allows users to",
			"configure their ALB outside waypoint but still have waypoint hook the application",
			"to that ALB",
		),
	)

	doc.SetField(
		"security_group_ids",
		"the existing security groups to add to the ALB",
		docs.Summary(
			"a set of existing security groups to add to the ALB",
		),
	)

	return doc, nil
}

var (
	_ component.ReleaseManager = (*Releaser)(nil)
	_ component.Configurable   = (*Releaser)(nil)
	_ component.Destroyer      = (*Releaser)(nil)
	_ component.Documented     = (*Releaser)(nil)
	_ component.Status         = (*Releaser)(nil)
)
