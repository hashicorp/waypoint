package ecs

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint/builtin/docker"
	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/docs"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

type Platform struct {
	config Config
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
func (p *Platform) DefaultReleaserFunc() interface{} {
	return func() *Releaser { return &Releaser{p: p} }
}

type Lifecycle struct {
	Init    func(LifecycleStatus) error
	Run     func(LifecycleStatus) error
	Cleanup func(LifecycleStatus) error
}

type lStatus struct {
	ui   terminal.UI
	sg   terminal.StepGroup
	step terminal.Step
}

func (l *lStatus) Status(str string, args ...interface{}) {
	if l.sg == nil {
		l.sg = l.ui.StepGroup()
	}

	if l.step != nil {
		l.step.Done()
		l.step = nil
	}

	l.step = l.sg.Add(str, args...)
}

func (l *lStatus) Update(str string, args ...interface{}) {
	if l.sg == nil {
		l.sg = l.ui.StepGroup()
	}

	if l.step != nil {
		l.step.Update(str, args...)
	} else {
		l.step = l.sg.Add(str, args)
	}
}

func (l *lStatus) Error(str string, args ...interface{}) {
	if l.sg == nil {
		l.sg = l.ui.StepGroup()
	}

	if l.step != nil {
		l.step.Update(str, args...)
		l.step.Abort()
	} else {
		l.step = l.sg.Add(str, args)
		l.step.Abort()
	}

	l.step = nil
}

func (l *lStatus) Abort() error {
	if l.step != nil {
		l.step.Abort()
		l.step = nil
	}

	if l.sg != nil {
		l.sg.Wait()
		l.sg = nil
	}

	return nil
}

func (l *lStatus) Close() error {
	if l.step != nil {
		l.step.Done()
		l.step = nil
	}

	if l.sg != nil {
		l.sg.Wait()
		l.sg = nil
	}

	return nil
}

func (lf *Lifecycle) Execute(L hclog.Logger, ui terminal.UI) error {
	var l lStatus
	l.ui = ui

	defer l.Close()

	if lf.Init != nil {
		L.Debug("lifecycle init")

		err := lf.Init(&l)
		if err != nil {
			l.Abort()
			return err
		}

	}

	L.Debug("lifecycle run")
	err := lf.Run(&l)
	if err != nil {
		l.Abort()
		return err
	}

	if lf.Cleanup != nil {
		L.Debug("lifecycle cleanup")

		err = lf.Cleanup(&l)
		if err != nil {
			l.Abort()
			return err
		}
	}

	return nil
}

type LifecycleStatus interface {
	Status(str string, args ...interface{})
	Update(str string, args ...interface{})
	Error(str string, args ...interface{})
}

func (p *Platform) Deploy(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	img *docker.Image,
	deployConfig *component.DeploymentConfig,
	ui terminal.UI,
) (*Deployment, error) {
	var (
		sess *session.Session
		dep  *Deployment

		role, cluster, logGroup string

		err error
	)

	if p.config.ALB != nil {
		if p.config.ALB.ListenerARN != "" {
			if p.config.ALB.ZoneId != "" || p.config.ALB.FQDN != "" {
				return nil, fmt.Errorf("When using an existing listener, Route53 setup is not available")
			}

			if p.config.ALB.CertificateId != "" {
				return nil, fmt.Errorf("When using an existing listener, certification configuration is not available")
			}
		}
	}

	lf := &Lifecycle{
		Init: func(s LifecycleStatus) error {
			sess = session.New(aws.NewConfig().WithRegion(p.config.Region))

			cluster, err = p.SetupCluster(ctx, s, sess)
			if err != nil {
				return err
			}

			role, err = p.SetupRole(ctx, s, log, sess, src)
			if err != nil {
				return err
			}

			logGroup, err = p.SetupLogs(ctx, s, log, sess)
			if err != nil {
				return err
			}

			return nil
		},

		Run: func(s LifecycleStatus) error {
			dep, err = p.Launch(ctx, s, log, ui, sess, src, img, deployConfig, role, cluster, logGroup)
			return err
		},

		Cleanup: func(s LifecycleStatus) error { return nil },
	}

	if err := lf.Execute(log, ui); err != nil {
		return nil, err
	}

	return dep, nil
}

func defaultSubnets(ctx context.Context, sess *session.Session) ([]*string, error) {
	svc := ec2.New(sess)

	desc, err := svc.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("default-for-az"),
				Values: []*string{aws.String("true")},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	var subnets []*string

	for _, subnet := range desc.Subnets {
		subnets = append(subnets, subnet.SubnetId)
	}

	return subnets, nil
}

func (p *Platform) SetupCluster(ctx context.Context, s LifecycleStatus, sess *session.Session) (string, error) {
	ecsSvc := ecs.New(sess)

	cluster := p.config.Cluster
	if cluster == "" {
		cluster = "waypoint"
	}

	desc, err := ecsSvc.DescribeClusters(&ecs.DescribeClustersInput{
		Clusters: []*string{aws.String(cluster)},
	})

	if err != nil {
		return "", err
	}

	if len(desc.Clusters) > 1 {
		s.Status("Found existing ECS cluster: %s", cluster)
		return cluster, nil
	}

	if p.config.EC2Cluster {
		return "", fmt.Errorf("EC2 clusters can not be automatically created")
	}

	s.Status("Creating new ECS cluster: %s", cluster)

	_, err = ecsSvc.CreateCluster(&ecs.CreateClusterInput{
		ClusterName: aws.String(cluster),
	})

	if err != nil {
		return "", err
	}

	s.Update("Created new ECS cluster: %s", cluster)
	return cluster, nil
}

const rolePolicy = `{
  "Version": "2012-10-17",
  "Statement": [
    {
		  "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "ecs-tasks.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}`

var fargateResources = map[int][]int{
	512:  {256},
	1024: {256, 512},
	2048: {256, 512, 1024},
	3072: {512, 1024},
	4096: {512, 1024},
	5120: {1024},
	6144: {1024},
	7168: {1024},
	8192: {1024},
}

func init() {
	for i := 4096; i < 16384; i += 1024 {
		fargateResources[i] = append(fargateResources[i], 2048)
	}

	for i := 8192; i <= 30720; i += 1024 {
		fargateResources[i] = append(fargateResources[i], 4096)
	}
}

func (p *Platform) SetupRole(ctx context.Context, s LifecycleStatus, L hclog.Logger, sess *session.Session, app *component.Source) (string, error) {
	svc := iam.New(sess)

	roleName := p.config.RoleName

	if roleName == "" {
		roleName = "ecr-" + app.App
	}

	// p.updateStatus("setting up IAM role")
	L.Debug("attempting to retrieve existing role", "role-name", roleName)

	queryInput := &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	}

	getOut, err := svc.GetRole(queryInput)
	if err == nil {
		s.Status("Found existing IAM role to use: %s", roleName)
		return *getOut.Role.Arn, nil
	}

	L.Debug("creating new role")
	s.Status("Creating IAM role: %s (%s)", roleName)

	input := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(rolePolicy),
		Path:                     aws.String("/"),
		RoleName:                 aws.String(roleName),
	}

	result, err := svc.CreateRole(input)
	if err != nil {
		return "", err
	}

	roleArn := *result.Role.Arn

	L.Debug("created new role", "arn", roleArn)

	aInput := &iam.AttachRolePolicyInput{
		RoleName:  aws.String(roleName),
		PolicyArn: aws.String("arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"),
	}

	_, err = svc.AttachRolePolicy(aInput)
	if err != nil {
		return "", err
	}

	L.Debug("attached execution role policy")

	s.Update("Created IAM role: %s", roleName)
	return roleArn, nil
}

func (p *Platform) SetupLogs(ctx context.Context, s LifecycleStatus, L hclog.Logger, sess *session.Session) (string, error) {
	// e.updateStatus("setting up CloudWatchLogs")

	logGroup := p.config.LogGroup
	if logGroup == "" {
		logGroup = "waypoint-logs"
	}

	cwl := cloudwatchlogs.New(sess)
	groups, err := cwl.DescribeLogGroups(&cloudwatchlogs.DescribeLogGroupsInput{
		Limit:              aws.Int64(1),
		LogGroupNamePrefix: aws.String(logGroup),
	})

	if err != nil {
		return "", err
	}

	if len(groups.LogGroups) == 0 {
		s.Status("Creating CloudWatchLogs group to store logs in: %s", logGroup)

		L.Debug("creating log group", "group", logGroup)
		_, err = cwl.CreateLogGroup(&cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(logGroup),
		})
		if err != nil {
			return "", err
		}

		s.Update("Created CloudWatchLogs group to store logs in: %s", logGroup)
	}

	return logGroup, nil

}

func createSG(
	ctx context.Context,
	s LifecycleStatus,
	sess *session.Session,
	name string,
	vpcId *string,
	ports ...int,
) (*string, error) {
	ec2srv := ec2.New(sess)

	dsg, err := ec2srv.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("group-name"),
				Values: []*string{aws.String(name)},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	var groupId *string

	if len(dsg.SecurityGroups) != 0 {
		groupId = dsg.SecurityGroups[0].GroupId
		s.Status("Using existing security group: %s", name)
	} else {
		s.Status("Creating security group: %s", name)
		out, err := ec2srv.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
			Description: aws.String("created by waypoint"),
			GroupName:   aws.String(name),
			VpcId:       vpcId,
		})
		if err != nil {
			return nil, err
		}

		groupId = out.GroupId
		s.Update("Created security group: %s", name)
	}

	s.Update("Authorizing ports to security group")
	for _, port := range ports {
		_, err = ec2srv.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
			CidrIp:     aws.String("0.0.0.0/0"),
			FromPort:   aws.Int64(int64(port)),
			ToPort:     aws.Int64(int64(port)),
			GroupId:    groupId,
			IpProtocol: aws.String("tcp"),
		})
	}

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "InvalidPermission.Duplicate":
				// fine, means we already added it.
			default:
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	s.Update("Configured security group: %s", name)

	return groupId, nil
}

func (p *Platform) Launch(
	ctx context.Context,
	s LifecycleStatus,
	L hclog.Logger,
	ui terminal.UI,
	sess *session.Session,
	app *component.Source,
	img *docker.Image,
	deployConfig *component.DeploymentConfig,
	roleArn, clusterName, logGroup string,
) (*Deployment, error) {
	id, err := component.Id()
	if err != nil {
		return nil, err
	}

	ecsSvc := ecs.New(sess)

	streamPrefix := fmt.Sprintf("waypoint-%d", time.Now().Nanosecond())

	env := []*ecs.KeyValuePair{
		{
			Name:  aws.String("PORT"),
			Value: aws.String("3000"),
		},
	}

	for k, v := range deployConfig.Env() {
		env = append(env, &ecs.KeyValuePair{
			Name:  aws.String(k),
			Value: aws.String(v),
		})
	}

	def := ecs.ContainerDefinition{
		Essential: aws.Bool(true),
		Name:      aws.String(app.App),
		Image:     aws.String(img.Name()),
		PortMappings: []*ecs.PortMapping{
			{
				ContainerPort: aws.Int64(3000),
			},
		},
		Environment: env,
		LogConfiguration: &ecs.LogConfiguration{
			LogDriver: aws.String("awslogs"),
			Options: map[string]*string{
				"awslogs-group":         aws.String(logGroup),
				"awslogs-region":        aws.String(p.config.Region),
				"awslogs-stream-prefix": aws.String(streamPrefix),
			},
		},
	}

	L.Debug("registring task definition", "id", id)

	var cpuShares int

	runtime := aws.String("FARGATE")
	if p.config.EC2Cluster {
		runtime = aws.String("EC2")
	} else {
		if p.config.Memory == 0 {
			return nil, fmt.Errorf("Memory value required for fargate")
		}
		cpuValues, ok := fargateResources[p.config.Memory]
		if !ok {
			var (
				allValues  []int
				goodValues []string
			)

			for k := range fargateResources {
				allValues = append(allValues, k)
			}

			sort.Ints(allValues)

			for _, k := range allValues {
				goodValues = append(goodValues, strconv.Itoa(k))
			}

			return nil, fmt.Errorf("Invalid memory value: %d (valid values: %s)",
				p.config.Memory, strings.Join(goodValues, ", "))
		}

		if p.config.CPU == 0 {
			cpuShares = cpuValues[0]
		} else {
			var (
				valid      bool
				goodValues []string
			)

			for _, c := range cpuValues {
				goodValues = append(goodValues, strconv.Itoa(c))
				if c == p.config.CPU {
					valid = true
					break
				}
			}

			if !valid {
				return nil, fmt.Errorf("Invalid cpu value: %d (valid values: %s)",
					p.config.Memory, strings.Join(goodValues, ", "))
			}

			cpuShares = p.config.CPU
		}
	}

	cpus := strconv.Itoa(cpuShares)
	mems := strconv.Itoa(p.config.Memory)

	family := "waypoint-" + app.App

	s.Status("Registering Task definition: %s", family)

	taskOut, err := ecsSvc.RegisterTaskDefinition(&ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: []*ecs.ContainerDefinition{&def},

		ExecutionRoleArn: aws.String(roleArn),
		Cpu:              aws.String(cpus),
		Memory:           aws.String(mems),
		Family:           aws.String(family),

		NetworkMode:             aws.String("awsvpc"),
		RequiresCompatibilities: []*string{runtime},

		Tags: []*ecs.Tag{
			{
				Key:   aws.String("waypoint-app"),
				Value: aws.String(app.App),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	s.Update("Registered Task definition: %s", family)
	rand := id[len(id)-(31-len(app.App)):]

	serviceName := fmt.Sprintf("%s-%s", app.App, rand)

	taskArn := *taskOut.TaskDefinition.TaskDefinitionArn

	var subnets []*string

	if len(p.config.Subnets) == 0 {
		s.Update("Using default subnets for Service networking")
		subnets, err = defaultSubnets(ctx, sess)
		if err != nil {
			return nil, err
		}
	} else {
		for i := range p.config.Subnets {
			subnets[i] = &p.config.Subnets[i]
		}
	}

	ec2srv := ec2.New(sess)

	subnetInfo, err := ec2srv.DescribeSubnets(&ec2.DescribeSubnetsInput{
		SubnetIds: subnets,
	})
	if err != nil {
		return nil, err
	}

	vpcId := subnetInfo.Subnets[0].VpcId

	s.Update("Creating ALB target group")
	L.Debug("creating target group", "name", serviceName)

	elbsrv := elbv2.New(sess)
	ctg, err := elbsrv.CreateTargetGroup(&elbv2.CreateTargetGroupInput{
		HealthCheckEnabled: aws.Bool(true),
		Name:               aws.String(serviceName),
		Port:               aws.Int64(3000),
		Protocol:           aws.String("HTTP"),
		TargetType:         aws.String("ip"),
		VpcId:              vpcId,
	})
	if err != nil {
		return nil, err
	}

	tgArn := ctg.TargetGroups[0].TargetGroupArn

	s.Update("Created ALB target group")

	// Create the load balancer OR modify the existing one to have this new target
	// group but with a weight of 0

	L.Debug("creating security group for ports 80 and 443")
	sgweb, err := createSG(ctx, s, sess, fmt.Sprintf("%s-inbound", app.App), vpcId, 80, 443)
	if err != nil {
		return nil, err
	}

	tgs := []*elbv2.TargetGroupTuple{
		{
			TargetGroupArn: tgArn,
			Weight:         aws.Int64(0),
		},
	}

	var (
		certs    []*elbv2.Certificate
		protocol string = "HTTP"
		port     int64  = 80
	)

	if p.config.ALB != nil && p.config.ALB.CertificateId != "" {
		protocol = "HTTPS"
		port = 443
		certs = append(certs, &elbv2.Certificate{
			CertificateArn: &p.config.ALB.CertificateId,
		})
	}

	var existingListener string

	if p.config.ALB != nil && p.config.ALB.ListenerARN != "" {
		existingListener = p.config.ALB.ListenerARN
	}

	var (
		lb          *elbv2.LoadBalancer
		listener    *elbv2.Listener
		newListener bool
	)

	if existingListener != "" {
		out, err := elbsrv.DescribeListeners(&elbv2.DescribeListenersInput{
			ListenerArns: []*string{aws.String(existingListener)},
		})
		if err != nil {
			return nil, err
		}

		listener = out.Listeners[0]
		s.Update("Using configured ALB Listener: %s (load-balancer: %s)",
			*listener.ListenerArn, *listener.LoadBalancerArn)
	} else {
		lbName := "waypoint-ecs-" + app.App
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
			s.Update("Using existing ALB %s (%s, dns-name: %s)",
				lbName, *lb.LoadBalancerArn, *lb.DNSName)
		} else {
			s.Update("Creating new ALB: %s", lbName)

			clb, err := elbsrv.CreateLoadBalancer(&elbv2.CreateLoadBalancerInput{
				Name:           aws.String(lbName),
				Subnets:        subnets,
				SecurityGroups: []*string{sgweb},
			})
			if err != nil {
				return nil, err
			}

			lb = clb.LoadBalancers[0]

			s.Update("Created new ALB: %s (dns-name: %s)", lbName, *lb.DNSName)
		}

		listeners, err := elbsrv.DescribeListeners(&elbv2.DescribeListenersInput{
			LoadBalancerArn: lb.LoadBalancerArn,
		})
		if err != nil {
			return nil, err
		}

		if len(listeners.Listeners) > 0 {
			listener = listeners.Listeners[0]
			s.Update("Using existing ALB Listener")
		} else {
			s.Update("Creating new ALB Listener")

			L.Info("load-balancer defined", "dns-name", *lb.DNSName)

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

			s.Update("Created new ALB Listener")
		}
	}

	if !newListener {
		def := listener.DefaultActions

		if len(def) > 0 && def[0].ForwardConfig != nil {
			for _, tg := range def[0].ForwardConfig.TargetGroups {
				if *tg.Weight > 0 {
					tgs = append(tgs, tg)
					L.Debug("previous target group", "arn", *tg.TargetGroupArn)
				}
			}
		}

		if len(tgs) == 0 {
			tgs[0].Weight = aws.Int64(100)
		}

		s.Update("Modifing ALB Listener to introduce target group")

		_, err = elbsrv.ModifyListener(&elbv2.ModifyListenerInput{
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

		s.Update("Modified ALB Listener to introduce target group")
	}

	if p.config.ALB != nil && p.config.ALB.ZoneId != "" {
		r53 := route53.New(sess)

		records, err := r53.ListResourceRecordSets(&route53.ListResourceRecordSetsInput{
			HostedZoneId:    aws.String(p.config.ALB.ZoneId),
			StartRecordName: aws.String(p.config.ALB.FQDN),
			StartRecordType: aws.String("A"),
			MaxItems:        aws.String("1"),
		})
		if err != nil {
			return nil, err
		}

		fqdn := p.config.ALB.FQDN

		if fqdn[len(fqdn)-1] != '.' {
			fqdn += "."
		}

		if len(records.ResourceRecordSets) > 0 && *(records.ResourceRecordSets[0].Name) == fqdn {
			s.Status("Found existing Route53 record: %s", *records.ResourceRecordSets[0].Name)
			L.Debug("found existing record, assuming it's correct")
		} else {
			s.Status("Creating new Route53 record: %s (zone-id: %s)",
				p.config.ALB.FQDN, p.config.ALB.ZoneId)

			L.Debug("creating new route53 record", "zone-id", p.config.ALB.ZoneId)
			input := &route53.ChangeResourceRecordSetsInput{
				ChangeBatch: &route53.ChangeBatch{
					Changes: []*route53.Change{
						{
							Action: aws.String("CREATE"),
							ResourceRecordSet: &route53.ResourceRecordSet{
								Name: aws.String(p.config.ALB.FQDN),
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
				HostedZoneId: aws.String(p.config.ALB.ZoneId),
			}

			result, err := r53.ChangeResourceRecordSets(input)
			if err != nil {
				return nil, err
			}
			L.Debug("record created", "change-id", *result.ChangeInfo.Id)

			s.Update("Created new Route53 record: %s (zone-id: %s)",
				p.config.ALB.FQDN, p.config.ALB.ZoneId)
		}
	}

	// Create the service

	L.Debug("creating service", "arn", *taskOut.TaskDefinition.TaskDefinitionArn)
	sg3000, err := createSG(ctx, s, sess, fmt.Sprintf("%s-inbound-internal", app.App), vpcId, 3000)
	if err != nil {
		return nil, err
	}

	count := int64(p.config.Count)
	if count == 0 {
		count = 1
	}

	netCfg := &ecs.AwsVpcConfiguration{
		Subnets:        subnets,
		SecurityGroups: []*string{sg3000},
	}

	netCfg.AssignPublicIp = aws.String("ENABLED")

	s.Status("Creating ECS Service (%s, cluster-name: %s)", serviceName, clusterName)
	servOut, err := ecsSvc.CreateService(&ecs.CreateServiceInput{
		Cluster:        &clusterName,
		DesiredCount:   aws.Int64(count),
		LaunchType:     runtime,
		ServiceName:    aws.String(serviceName),
		TaskDefinition: aws.String(taskArn),
		NetworkConfiguration: &ecs.NetworkConfiguration{
			AwsvpcConfiguration: netCfg,
		},
		LoadBalancers: []*ecs.LoadBalancer{
			{
				ContainerName:  aws.String(app.App),
				ContainerPort:  aws.Int64(3000),
				TargetGroupArn: tgArn,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	s.Update("Created ECS Service (%s, cluster-name: %s)", serviceName, clusterName)
	L.Debug("service started", "arn", servOut.Service.ServiceArn)

	dep := &Deployment{
		Cluster:         clusterName,
		TaskArn:         taskArn,
		ServiceArn:      *servOut.Service.ServiceArn,
		TargetGroupArn:  *tgArn,
		LoadBalancerArn: *lb.LoadBalancerArn,
	}

	return dep, nil
}

func (p *Platform) Destroy(
	ctx context.Context,
	log hclog.Logger,
	deployment *Deployment,
	ui terminal.UI,
) error {
	log.Debug("removing deployment target group from load balancer")

	sess := session.New(aws.NewConfig().WithRegion(p.config.Region))
	elbsrv := elbv2.New(sess)

	listeners, err := elbsrv.DescribeListeners(&elbv2.DescribeListenersInput{
		LoadBalancerArn: &deployment.LoadBalancerArn,
	})
	if err != nil {
		return err
	}

	var listener *elbv2.Listener

	if len(listeners.Listeners) > 0 {
		listener = listeners.Listeners[0]

		def := listener.DefaultActions

		var tgs []*elbv2.TargetGroupTuple

		if len(def) > 0 && def[0].ForwardConfig != nil {
			var active bool

			for _, tg := range def[0].ForwardConfig.TargetGroups {
				if *tg.TargetGroupArn != deployment.TargetGroupArn {
					tgs = append(tgs, tg)
					if *tg.Weight > 0 {
						active = true
					}
				}
			}

			// If there are no target groups active, then we just activate the first
			// one, otherwise we can't modify the listener.
			if !active && len(tgs) > 0 {
				tgs[0].Weight = aws.Int64(100)
			}

			log.Debug("modifying listener to remove target group", "target-groups", len(tgs))

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
				return err
			}
		}
	}

	log.Debug("deleting target group", "arn", deployment.TargetGroupArn)

	_, err = elbsrv.DeleteTargetGroup(&elbv2.DeleteTargetGroupInput{
		TargetGroupArn: &deployment.TargetGroupArn,
	})
	if err != nil {
		return err
	}

	log.Debug("deleting ecs service", "arn", deployment.ServiceArn)

	_, err = ecs.New(sess).DeleteService(&ecs.DeleteServiceInput{
		Cluster: &deployment.Cluster,
		Force:   aws.Bool(true),
		Service: &deployment.ServiceArn,
	})
	return nil
}

type ALBConfig struct {
	// Certificate ARN to attach to the load balancer
	CertificateId string `hcl:"certificate"`

	// Route53 Zone to setup record in
	ZoneId string `hcl:"zone_id"`

	// Fully qualified domain name of the record to create in the target zone id
	FQDN string `hcl:"domain_name"`

	// When set, waypoint will configure the target group into the specified
	// ALB Listener ARN. This allows for usage of existing ALBs.
	ListenerARN string `hcl:"listener_arn,optional"`
}

type Config struct {
	// AWS Region to deploy into
	Region string `hcl:"region"`

	// Name of the Log Group to store logs into
	LogGroup string `hcl:"log_group,optional"`

	// Name of the ECS cluster to install the service into
	Cluster string `hcl:"cluster,optional"`

	// Name of the IAM Role to associate with the ECS Service
	RoleName string `hcl:"role_name,optional"`

	// Subnets to place the service into. Defaults to the subnets in the default VPC.
	Subnets []string `hcl:"subnets,optional"`

	// How many tasks of the service to run. Default 1.
	Count int `hcl:"count,optional"`

	// How much memory to assign to the containers
	Memory int `hcl:"memory"`

	// How much CPU to assign to the containers
	CPU int `hcl:"cpu,optional"`

	// Assign each task a public IP. Default false.
	// TODO to access ECR you need a nat gateway or a public address and so if you
	// set this to false in the default subnets, ECS can't pull the image. Leaving
	// it disabled until we figure out how to handle that onramp case.
	// AssignPublicIp bool `hcl:"assign_public_ip,optional"`

	// Indicate that service should be deployed on an EC2 cluster.
	EC2Cluster bool `hcl:"ec2_cluster,optional"`

	// Configuration options for how the ALB will be configured.
	ALB *ALBConfig `hcl:"alb,block"`
}

func (p *Platform) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&Config{}))
	if err != nil {
		return nil, err
	}

	doc.Description("Deploy the application into an ECS cluster on AWS")

	doc.Input("docker.Image")
	doc.Output("ecs.Deployment")

	doc.SetField(
		"region",
		"the AWS region for the ECS cluster",
	)

	doc.SetField(
		"log_group",
		"the CloudWatchLogs log group to store container logs into",
		docs.Default("derived from the application name"),
	)

	doc.SetField(
		"cluster",
		"the name of the ECS cluster to deploy into",
		docs.Summary(
			"the ECS cluster that will run the application as a Service.",
			"if there is no ECS cluster with this name, the ECS cluster will be",
			"created and configured to use Fargate to run containers.",
		),
	)

	doc.SetField(
		"role_name",
		"the name of the IAM role to use for ECS execution",
		docs.Default("create a new IAM role based on the application name"),
	)

	doc.SetField(
		"subnets",
		"the VPC subnets to use for the application",
		docs.Default("public subnets in the default VPC"),
	)

	doc.SetField(
		"count",
		"how many instances of the application should run",
	)

	doc.SetField(
		"memory",
		"how much memory to assign to the container running the application",
		docs.Summary(
			"when running in Fargate, this must be one of a few values, specified in MB:",
			"512, 1024, 2048, 3072, 4096, 5120, and up to 16384 in increments of 1024.",
			"The memory value also controls the possible values for cpu",
		),
	)

	doc.SetField(
		"ec2_cluster",
		"indicate if the ECS cluster should be EC2 type rather than Fargate",
		docs.Summary(
			"this controls if we should verify the ECS cluster in EC2 type. The cluster",
			"will not be created if it doesn't exist, only that there as existing cluster",
			"this is using EC2 and not Fargate",
		),
	)

	doc.SetField(
		"alb.certificate",
		"the ARN of an AWS Certificate Manager cert to associate with the ALB",
	)

	doc.SetField(
		"alb.zone_id",
		"Route53 ZoneID to create a DNS record into",
		docs.Summary(
			"set along with alb.domain_name to have DNS automatically setup for the ALB",
		),
	)

	doc.SetField(
		"alb.domain_name",
		"Fully qualified domain name to set for the ALB",
		docs.Summary(
			"set along with zone_id to have DNS automatically setup for the ALB.",
			"this value should include the full hostname and domain name, for instance",
			"app.example.com",
		),
	)

	doc.SetField(
		"alb.listener_arn",
		"the ARN on an existing ALB to configure",
		docs.Summary(
			"when this is set, no ALB or Listener is created. Instead the application is",
			"configured by manipulating this existing Listener. This allows users to",
			"configure their ALB outside waypoint but still have waypoint hook the application",
			"to that ALB",
		),
	)

	var memvals []int

	for k := range fargateResources {
		memvals = append(memvals, k)
	}

	sort.Ints(memvals)

	var sb strings.Builder

	for _, mem := range memvals {
		cpu := fargateResources[mem]

		var cpuVals []string

		for _, c := range cpu {
			cpuVals = append(cpuVals, strconv.Itoa(c))
		}

		fmt.Fprintf(&sb, "%dMB: %s\n", mem, strings.Join(cpuVals, ", "))
	}

	doc.SetField(
		"cpu",
		"how many cpu shares the container running the application is allowed",
		docs.Summary(
			"on Fargate, possible values for this are configured by the amount of memory",
			"the container is using. Here is a complete listing of possible values:\n",
			sb.String(),
		),
	)

	return doc, nil
}
