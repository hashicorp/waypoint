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
			clb, err := elbsrv.CreateLoadBalancer(&elbv2.CreateLoadBalancerInput{
				Name:           aws.String(lbName),
				Subnets:        subnets,
				SecurityGroups: []*string{sg},
			})
			if err != nil {
				return nil, err
			}

			lb = clb.LoadBalancers[0]
		}

		listeners, err := elbsrv.DescribeListeners(&elbv2.DescribeListenersInput{
			LoadBalancerArn: lb.LoadBalancerArn,
		})
		if err != nil {
			return nil, err
		}

		if len(listeners.Listeners) > 0 {
			listener = listeners.Listeners[0]
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

			result, err := r53.ChangeResourceRecordSets(input)
			if err != nil {
				return nil, err
			}
			log.Debug("record created", "change-id", *result.ChangeInfo.Id)
		}
	}

	return &Release{
		Url:             "http://" + hostname,
		LoadBalancerArn: *lb.LoadBalancerArn,
	}, nil
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