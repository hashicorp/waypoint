package alb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/aws/utils"
)

const (
	managementLevelTargetGroup = "target_group"
	managementLevelListener    = "listener"
	managementLevelComplete    = "complete"
)

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

// ReleaseFunc implements component.ReleaseManager
func (r *Releaser) DestroyFunc() interface{} {
	return r.Destroy
}

// Release manages target group attachement to a configured ALB
func (r *Releaser) Release(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	ui terminal.UI,
	target *TargetGroup,
) (*Release, error) {
	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: target.Region,
	})
	if err != nil {
		return nil, err
	}

	// Start recording our state here
	result := new(Release)

	elbsrv := elbv2.New(sess)

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

	var (
		certs    []*elbv2.Certificate
		protocol string = "HTTP"
		port     int64  = 80
	)

	if r.config.CertificateId != "" {
		protocol = "HTTPS"
		port = 443
		certs = append(certs, &elbv2.Certificate{
			CertificateArn: &r.config.CertificateId,
		})
	}

	// If there is a port defined, honor it.
	if r.config.Port != 0 {
		port = int64(r.config.Port)
	}

	var (
		lb               *elbv2.LoadBalancer
		listener         *elbv2.Listener
		newListener      bool
		existingListener string
	)

	if r.config.ListenerARN != "" {
		existingListener = r.config.ListenerARN
		result.ListenerArn = existingListener
		result.ManagementLevel = managementLevelTargetGroup
	}

	tgs := []*elbv2.TargetGroupTuple{
		{
			TargetGroupArn: &target.Arn,
			Weight:         aws.Int64(100),
		},
	}

	if existingListener != "" {
		out, err := elbsrv.DescribeListeners(&elbv2.DescribeListenersInput{
			ListenerArns: []*string{aws.String(existingListener)},
		})
		if err != nil {
			return nil, err
		}

		listener = out.Listeners[0]
	} else {
		dlb, err := elbsrv.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{
			Names: []*string{&lbName},
		})
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case elbv2.ErrCodeLoadBalancerNotFoundException:
					// fine, means we'll create it.
				default:
					return nil, err
				}
			} else {
				return nil, err
			}
		}

		if dlb != nil && len(dlb.LoadBalancers) > 0 {
			lb = dlb.LoadBalancers[0]
			result.ManagementLevel = managementLevelListener
		} else {
			var (
				vpc     *string
				subnets []*string
			)

			e := ec2.New(sess)

			if len(r.config.Subnets) == 0 {
				sn, v, err := utils.DefaultSubnets(ctx, sess)
				if err != nil {
					return nil, err
				}

				var sc []string
				for _, s := range subnets {
					sc = append(sc, *s)
				}

				vpc = v
				subnets = sn
			} else {
				for _, s := range r.config.Subnets {
					subnets = append(subnets, aws.String(s))
				}

				subnetInfo, err := e.DescribeSubnets(&ec2.DescribeSubnetsInput{
					SubnetIds: []*string{subnets[0]},
				})
				if err != nil {
					return nil, err
				}

				vpc = subnetInfo.Subnets[0].VpcId
			}

			sg, err := utils.CreateSecurityGroup(ctx, sess, fmt.Sprintf("%s-incoming", lbName), vpc, int(port))
			if err != nil {
				return nil, err
			}

			result.SecurityGroupId = *sg

			clb, err := elbsrv.CreateLoadBalancer(&elbv2.CreateLoadBalancerInput{
				Name:           aws.String(lbName),
				Subnets:        subnets,
				SecurityGroups: []*string{sg},
			})
			if err != nil {
				return nil, err
			}

			lb = clb.LoadBalancers[0]
			result.ManagementLevel = managementLevelComplete
		}

		result.LoadBalancerArn = *lb.LoadBalancerArn
		listeners, err := elbsrv.DescribeListeners(&elbv2.DescribeListenersInput{
			LoadBalancerArn: lb.LoadBalancerArn,
		})
		if err != nil {
			return nil, err
		}

		if len(listeners.Listeners) > 0 {
			listener = listeners.Listeners[0]
			result.ListenerArn = *listener.ListenerArn
		} else {
			log.Info("load-balancer defined", "dns-name", *lb.DNSName)

			tgs[0].Weight = aws.Int64(100)
			lo, err := elbsrv.CreateListener(&elbv2.CreateListenerInput{
				LoadBalancerArn: lb.LoadBalancerArn,
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
				return nil, err
			}

			newListener = true
			listener = lo.Listeners[0]
			result.ListenerArn = *listener.ListenerArn
		}
	}

	log.Debug("configuring weight 100 for target group", "arn", target.Arn)

	if !newListener {
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
			return nil, err
		}
	}

	hostname := *lb.DNSName

	if r.config.ZoneId != "" {
		r53 := route53.New(sess)

		records, err := r53.ListResourceRecordSets(&route53.ListResourceRecordSetsInput{
			HostedZoneId:    aws.String(r.config.ZoneId),
			StartRecordName: aws.String(r.config.FQDN),
			StartRecordType: aws.String("A"),
			MaxItems:        aws.String("1"),
		})
		if err != nil {
			return nil, err
		}

		fqdn := r.config.FQDN

		if fqdn[len(fqdn)-1] != '.' {
			fqdn += "."
		}

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
									DNSName:              lb.DNSName,
									EvaluateTargetHealth: aws.Bool(true),
									HostedZoneId:         lb.CanonicalHostedZoneId,
								},
							},
						},
					},
					Comment: aws.String("managed by waypoint"),
				},
				HostedZoneId: aws.String(r.config.ZoneId),
			}

			out, err := r53.ChangeResourceRecordSets(input)
			if err != nil {
				return nil, err
			}
			log.Debug("record created", "change-id", *out.ChangeInfo.Id)
			result.ZoneId = r.config.ZoneId
			result.Fqdn = fqdn
		}
	}

	result.Url = "http://" + hostname
	result.TargetGroupArn = target.Arn
	return result, nil
}

func (r *Releaser) Destroy(
	ctx context.Context,
	log hclog.Logger,
	release *Release,
	ui terminal.UI,
) error {
	// We'll update the user in real time
	st := ui.Status()
	defer st.Close()

	sess, err := utils.GetSession(&utils.SessionConfig{})
	if err != nil {
		return err
	}

	elbsrv := elbv2.New(sess)

	switch release.ManagementLevel {
	case managementLevelTargetGroup:
		// Existing listener was supplied, just detach. Since we just
		// inserted our target group with a weight of 100 before in the
		// list, just locate that target group in the list and remove it
		// there.
		existingListener := release.ListenerArn
		log.Debug("detaching target group from supplied listener", "arn", existingListener)
		st.Update("Detaching target group...")
		out, err := elbsrv.DescribeListeners(&elbv2.DescribeListenersInput{
			ListenerArns: []*string{aws.String(existingListener)},
		})
		if err != nil {
			st.Step(terminal.StatusError, "Error detaching target group")
			return err
		}

		actions := out.Listeners[0].DefaultActions
		tgtMissing := true
		for i := 0; i < len(actions[0].ForwardConfig.TargetGroups[1:]); i++ {
			tgt := actions[0].ForwardConfig.TargetGroups[i]
			if *tgt.TargetGroupArn == release.TargetGroupArn {
				actions[0].ForwardConfig.TargetGroups = append(
					actions[0].ForwardConfig.TargetGroups[:i],
					actions[0].ForwardConfig.TargetGroups[i+1:]...,
				)
				tgtMissing = false
			}
		}

		if tgtMissing {
			log.Debug("target group not found in listener, not modifying", "arn", existingListener)
			break
		}

		_, err = elbsrv.ModifyListener(&elbv2.ModifyListenerInput{
			ListenerArn:    aws.String(existingListener),
			DefaultActions: actions,
		})
		if err != nil {
			st.Step(terminal.StatusError, "Error detaching target group")
			return err
		}

		st.Step(terminal.StatusOK, "Detached target group")

	case managementLevelListener:
		// Existing load balancer was supplied (via the "name" field),
		// but no listeners existed. Simply remove the listener we
		// created.
		log.Debug("removing managed listener from existing load balancer", "arn", release.LoadBalancerArn)
		st.Update("Removing listener from existing load balancer...")
		out, err := elbsrv.DescribeListeners(&elbv2.DescribeListenersInput{
			LoadBalancerArn: aws.String(release.LoadBalancerArn),
		})
		if err != nil {
			st.Step(terminal.StatusError, "Error removing listener from existing load balancer")
			return err
		}

		var listenerArns []string
	nextListener:
		for _, l := range out.Listeners {
			for _, a := range l.DefaultActions {
				if *a.TargetGroupArn == release.TargetGroupArn {
					listenerArns = append(listenerArns, *l.ListenerArn)
					continue nextListener
				}

				for _, tgt := range a.ForwardConfig.TargetGroups {
					if *tgt.TargetGroupArn == release.TargetGroupArn {
						listenerArns = append(listenerArns, *l.ListenerArn)
						continue nextListener
					}
				}
			}
		}

		if len(listenerArns) == 0 {
			log.Debug("no listeners found with assigned target group, not modifying", "arn", release.LoadBalancerArn)
			break
		}

		for _, listenerArn := range listenerArns {
			log.Debug("deleting listener", "arn", listenerArn)
			_, err = elbsrv.DeleteListener(&elbv2.DeleteListenerInput{
				ListenerArn: aws.String(listenerArn),
			})
			if err != nil {
				st.Step(terminal.StatusError, "Error removing listener from existing load balancer")
				return err
			}
		}

		st.Step(terminal.StatusOK, "Removed listeners from existing load balancer")

	default:
		// Assume managementLevelComplete here, as it's the only state
		// otherwise. Just remove the load balancer wholesale.
		log.Debug("deleting load balancer", "arn", release.LoadBalancerArn)
		st.Update("Deleting load balancer...")
		if _, err := elbsrv.DeleteLoadBalancer(&elbv2.DeleteLoadBalancerInput{
			LoadBalancerArn: aws.String(release.LoadBalancerArn),
		}); err != nil {
			st.Step(terminal.StatusError, "Error deleting load balancer")
			return err
		}

		st.Step(terminal.StatusOK, "Deleted load balancer")

		// Delete the security group as well
		log.Debug("deleting security group", "id", release.SecurityGroupId)
		st.Update("Deleting security group...")
		e := ec2.New(sess)
		if _, err := e.DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{
			GroupId: aws.String(release.SecurityGroupId),
		}); err != nil {
			st.Step(terminal.StatusError, "Error deleting security group")
			return err
		}

		st.Step(terminal.StatusOK, "Deleted security group")
	}

	if release.ZoneId != "" && release.Fqdn != "" {
		// Delete the Route53 record we created
		log.Debug("deleting route53 record", "zone_id", release.ZoneId, "fqdn", release.Fqdn)
		st.Update("Deleting DNS record...")

		r53 := route53.New(sess)

		records, err := r53.ListResourceRecordSets(&route53.ListResourceRecordSetsInput{
			HostedZoneId:    aws.String(release.ZoneId),
			StartRecordName: aws.String(release.Fqdn),
			StartRecordType: aws.String("A"),
			MaxItems:        aws.String("1"),
		})
		if err != nil {
			st.Step(terminal.StatusError, "Error deleting DNS record")
			return err
		}

		if len(records.ResourceRecordSets) > 0 && *(records.ResourceRecordSets[0].Name) == release.Fqdn {
			if _, err := r53.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
				HostedZoneId: aws.String(release.ZoneId),
				ChangeBatch: &route53.ChangeBatch{
					Changes: []*route53.Change{
						{
							Action:            aws.String(route53.ChangeActionDelete),
							ResourceRecordSet: records.ResourceRecordSets[0],
						},
					},
				},
			}); err != nil {
				st.Step(terminal.StatusError, "Error deleting DNS record")
				return err
			}

			st.Step(terminal.StatusOK, "Deleted DNS record")
		} else {
			log.Debug("route53 record not found, ignoring", "zone_id", release.ZoneId, "fqdn", release.Fqdn)
		}
	}

	return nil
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

	return doc, nil
}

var (
	_ component.ReleaseManager = (*Releaser)(nil)
	_ component.Configurable   = (*Releaser)(nil)
	_ component.Documented     = (*Releaser)(nil)
)
