package lambda

import (
	"context"
	"crypto/subtle"
	"fmt"
	"io"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint/builtin/aws/ecr"
	"github.com/hashicorp/waypoint/builtin/aws/utils"
	wpssh "github.com/hashicorp/waypoint/internal/ssh"
	"github.com/pkg/errors"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

// Platform is the Platform implementation for AWS Lambda
type Platform struct {
	config Config
}

// ConfigSet is called after a configuration has been decoded
// we can use this to validate the config
func (p *Platform) ConfigSet(config interface{}) error {
	_, ok := config.(*Config)
	if !ok {
		// this should never happen
		return fmt.Errorf("Invalid configuration, expected *lambda.Config, got %s", reflect.TypeOf(config))
	}

	return nil
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

// DestroyWorkspaceFunc implements component.WorkspaceDestroyer
func (p *Platform) DestroyWorkspaceFunc() interface{} {
	return p.DestroyWorkspace
}

// ExecFunc implements component.Execer
func (p *Platform) ExecFunc() interface{} {
	return p.Exec
}

// LogsFunc implements component.LogsPlatform
func (p *Platform) LogsFunc() interface{} {
	return p.Logs
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

// TODO sort out the right way to validate AWS config and use it in all the AWS types
func (p *Platform) ValidateAuth(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	ui terminal.UI,
) error {
	return nil
}

const (
	// The default amount of memory to give to the function invocation, 256MB
	DefaultMemory = 256

	// How long the function should run before terminating it.
	DefaultTimeout = 60

	// The instruction set architecture that the function supports.
	DefaultArchitecture = lambda.ArchitectureX8664
)

const lambdaRolePolicy = `{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"Principal": {
				"Service": "lambda.amazonaws.com"
			},
			"Action": "sts:AssumeRole"
		}
	]
}`

// setupRole creates an IAM role that will be used by the Lambda function. The role
// has the basic lambda execution rules attached to it.
func (p *Platform) setupRole(L hclog.Logger, app *component.Source, sess *session.Session) (string, error) {
	svc := iam.New(sess)

	roleName := "lambda-" + app.App

	L.Info("attempting to retrieve existing role", "role-name", roleName)

	queryInput := &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	}

	var roleArn string

	getOut, err := svc.GetRole(queryInput)
	if err == nil {
		roleArn = *getOut.Role.Arn
		L.Info("found existing role", "arn", roleArn)
		return roleArn, nil
	}

	L.Info("creating new role")

	input := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(lambdaRolePolicy),
		Path:                     aws.String("/"),
		RoleName:                 aws.String(roleName),
	}

	result, err := svc.CreateRole(input)
	if err != nil {
		return "", err
	}

	roleArn = *result.Role.Arn

	L.Info("created new role", "arn", roleArn)

	aInput := &iam.AttachRolePolicyInput{
		RoleName:  aws.String(roleName),
		PolicyArn: aws.String("arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"),
	}

	_, err = svc.AttachRolePolicy(aInput)
	if err != nil {
		return "", err
	}

	L.Info("attached execution role policy")

	return roleArn, nil
}

// Deploy deploys an image to AWS Lambda.
func (p *Platform) Deploy(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	img *ecr.Image,
	deployConfig *component.DeploymentConfig,
	ui terminal.UI,
) (*Deployment, error) {

	sg := ui.StepGroup()
	defer sg.Wait()

	step := sg.Add("Connecting to AWS")

	// We put this in a function because if/when step is reassigned, we want to
	// abort the new value.
	defer func() {
		step.Abort()
	}()

	// Start building our deployment since we use this information
	deployment := &Deployment{}
	id, err := component.Id()
	if err != nil {
		return nil, err
	}

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: p.config.Region,
		Logger: log,
	})
	if err != nil {
		return nil, err
	}

	roleArn := p.config.RoleArn

	if roleArn == "" {
		arn, err := p.setupRole(log, src, sess)
		if err != nil {
			return nil, err
		}

		roleArn = arn
	}

	mem := int64(p.config.Memory)
	if mem == 0 {
		mem = DefaultMemory
	}

	timeout := int64(p.config.Timeout)
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	architecture := p.config.Architecture
	if architecture == "" {
		architecture = DefaultArchitecture
	}

	step.Done()

	step = sg.Add("Reading Lambda function: %s", src.App)

	lamSvc := lambda.New(sess)
	curFunc, err := lamSvc.GetFunction(&lambda.GetFunctionInput{
		FunctionName: aws.String(src.App),
	})

	var funcarn string

	// If the function exists (ie we read it), we update it's code rather than create a new one.
	if err == nil {
		step.Update("Updating Lambda function with new code")

		var reset bool
		var update lambda.UpdateFunctionConfigurationInput

		if *curFunc.Configuration.MemorySize != mem {
			update.MemorySize = aws.Int64(mem)
			reset = true
		}

		if *curFunc.Configuration.Timeout != timeout {
			update.Timeout = aws.Int64(timeout)
			reset = true
		}

		if reset {
			update.FunctionName = curFunc.Configuration.FunctionArn

			_, err = lamSvc.UpdateFunctionConfiguration(&update)
			if err != nil {
				return nil, errors.Wrapf(err, "unable to update function configuration")
			}
		}

		funcCfg, err := lamSvc.UpdateFunctionCode(&lambda.UpdateFunctionCodeInput{
			FunctionName:  aws.String(src.App),
			ImageUri:      aws.String(img.Name()),
			Architectures: aws.StringSlice([]string{architecture}),
		})

		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case "ValidationException":
					// likely here if Architectures was invalid
					if architecture != lambda.ArchitectureX8664 && architecture != lambda.ArchitectureArm64 {
						return nil, fmt.Errorf("architecture must be either x86_64 or arm64")
					}
					return nil, err
				}
			}
			return nil, err
		}

		funcarn = *funcCfg.FunctionArn

		if err != nil {
			return nil, err
		}

		// We couldn't read the function before, so we'll go ahead and create one.
	} else {
		step.Update("Creating new Lambda function")

		// Run this in a loop to guard against eventual consistency errors with the specified
		// role not showing up within lambda right away.
		for i := 0; i < 30; i++ {
			funcOut, err := lamSvc.CreateFunction(&lambda.CreateFunctionInput{
				Description:  aws.String(fmt.Sprintf("waypoint %s", src.App)),
				FunctionName: aws.String(src.App),
				Role:         aws.String(roleArn),
				Timeout:      aws.Int64(timeout),
				MemorySize:   aws.Int64(mem),
				Tags: map[string]*string{
					"waypoint.app": aws.String(src.App),
				},
				PackageType: aws.String("Image"),
				Code: &lambda.FunctionCode{
					ImageUri: aws.String(img.Name()),
				},
				ImageConfig:   &lambda.ImageConfig{},
				Architectures: aws.StringSlice([]string{architecture}),
			})

			if err != nil {
				// if we encounter an unrecoverable error, exit now.
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case "ResourceConflictException":
						return nil, err
					case "ValidationException":
						// likely here if Architectures was an invalid
						if architecture != lambda.ArchitectureX8664 && architecture != lambda.ArchitectureArm64 {
							return nil, fmt.Errorf("architecture must be either x86_64 or arm64")
						}
						return nil, err
					}
				}

				// otherwise sleep and try again
				time.Sleep(2 * time.Second)
			} else {
				funcarn = *funcOut.FunctionArn
				break
			}
		}
	}

	if funcarn == "" {
		return nil, fmt.Errorf("Unable to create function, timed out trying")
	}

	step.Done()

	step = sg.Add("Waiting for Lambda function to be processed")
	// The image is never ready right away, AWS has to process it, so we wait
	// 3 seconds before trying to publish the version

	time.Sleep(3 * time.Second)

	// no publish this new code to create a stable identifier for it. Otherwise
	// if a manually pushes to the function and we use $LATEST, we'll accidentally
	// start running their manual code rather then the fixed one we have here.
	var ver *lambda.FunctionConfiguration

	// Only try 30 times.
	for i := 0; i < 30; i++ {
		ver, err = lamSvc.PublishVersion(&lambda.PublishVersionInput{
			FunctionName: aws.String(src.App),
		})

		if err == nil {
			break
		}

		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "ResourceConflictException":
				// It's updating, wait a sec and try again
				time.Sleep(time.Second)
			default:
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if ver == nil {
		return nil, fmt.Errorf("Lambda was unable to prepare the function in the allotted time")
	}

	verarn := *ver.FunctionArn

	step.Update("Published Lambda function: %s (%s)", verarn, *ver.Version)
	step.Done()

	_, err = lamSvc.AddPermission(&lambda.AddPermissionInput{
		Action:       aws.String("lambda:InvokeFunction"),
		FunctionName: aws.String(verarn),
		Principal:    aws.String("elasticloadbalancing.amazonaws.com"),
		StatementId:  aws.String("load-balancer"),
	})

	if err != nil {
		return nil, err
	}

	// Now generate a new TargetGroup so the Lambda can be attached to an ALB easily.

	step = sg.Add("Creating TargetGroup for Lambda version")
	svc := elbv2.New(sess)

	serviceName := fmt.Sprintf("%s-%s", src.App, id)

	// We have to clamp at a length of 32 because the Name field to CreateTargetGroup
	// requires that the name is 32 characters or less.
	if len(serviceName) > 32 {
		serviceName = serviceName[:32]
		log.Debug("using a shortened value for service name due to AWS's length limits", "serviceName", serviceName)
	}

	ctg, err := svc.CreateTargetGroup(&elbv2.CreateTargetGroupInput{
		Name:       aws.String(serviceName),
		TargetType: aws.String(elbv2.TargetTypeEnumLambda),
	})

	if err != nil {
		return nil, err
	}

	_, err = svc.RegisterTargets(&elbv2.RegisterTargetsInput{
		TargetGroupArn: ctg.TargetGroups[0].TargetGroupArn,
		Targets: []*elbv2.TargetDescription{
			{
				Id: aws.String(verarn),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	step.Done()

	deployment.Region = p.config.Region
	deployment.Id = id
	deployment.FuncArn = funcarn
	deployment.VerArn = verarn
	deployment.Version = *ver.Version
	deployment.TargetGroupArn = *ctg.TargetGroups[0].TargetGroupArn

	return deployment, nil
}

// This is used by the Exec plugin to provide some helpful output
// while we prepare the exec environment.
func rewriteLine(w io.Writer, str string, args ...interface{}) {
	fmt.Fprintf(w, "\r\033[K"+str, args...)
}

// Exec creates an ECS task using the given deployments ECR image and then
// ssh's to it to provide the shell.
func (p *Platform) Exec(
	ctx context.Context,
	log hclog.Logger,
	deployment *Deployment,
	image *ecr.Image,
	ui terminal.UI,
	es *component.ExecSessionInfo,
	app *component.Source,
) (*component.ExecResult, error) {
	rewriteLine(es.Output, "Launching ECS task to provide shell...")

	sshMaterial, err := wpssh.GenerateKeys()
	if err != nil {
		return nil, err
	}

	var esl ecsLauncher
	esl.Region = p.config.Region
	esl.PublicKey = sshMaterial.UserPublic
	esl.HostKey = sshMaterial.HostPrivate
	esl.Image = image.Name()
	esl.LogOutput = es.Output

	ti, err := esl.Launch(ctx, log, ui, app, deployment)
	if err != nil {
		log.Error("error launching ECS task", "error", err)
		fmt.Fprintf(es.Output, "\r\nError launching ECS task: %s", err)
		return nil, err
	}

	log.Info("starting exec session for aws lambda", "args", es.Arguments, "ip", ti.IP)

	var cfg ssh.ClientConfig
	cfg.User = "waypoint"
	cfg.Auth = []ssh.AuthMethod{
		ssh.PublicKeys(sshMaterial.UserPrivate),
	}

	expectedHost := sshMaterial.HostPublic.Marshal()

	cfg.HostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		// Weirdly this is how you make sure the host key is what you think it should be.
		// Think of this as where normal ssh client would do the "Do you want to trust this
		// host?" popup.
		if subtle.ConstantTimeCompare(expectedHost, key.Marshal()) == 1 {
			return nil
		}

		return fmt.Errorf("wrong host key detected")
	}

	rewriteLine(es.Output, "Connecting to ECS task...")

	// This is how long to setup the deadline for a TCP connect is.
	// If the remote side refuses the connection (which is what we're
	// expecting is the most common error here) Dial will return before
	// the timeout, because it has a result. It won't try again until
	// the timeout is filled.
	cfg.Timeout = 5 * time.Second

	var client *ssh.Client

	// We retry 30 times, sleeping for a second each time, thusly about
	// 30 seconds.
	for i := 0; i < 30; i++ {
		client, err = ssh.Dial("tcp", ti.IP, &cfg)
		if err == nil {
			break
		}

		if _, ok := err.(interface{ Temporary() bool }); ok {
			time.Sleep(time.Second)
			continue
		}

		log.Error("error dialing ssh", "error", err)
		return nil, err
	}

	sess, err := client.NewSession()
	if err != nil {
		log.Error("error starting ssh session", "error", err)
		return nil, err
	}

	// Wire the SSH session I/O up to the I/O for the Exec session directly.

	sess.Stderr = es.Error
	sess.Stdout = es.Output
	sess.Stdin = es.Input

	if es.IsTTY {
		// TODO(evanphx) should we be setting other modes here?
		modes := make(ssh.TerminalModes)

		err = sess.RequestPty(
			es.Term,
			es.InitialWindowSize.Height, es.InitialWindowSize.Width,
			modes,
		)
		if err != nil {
			log.Error("error requesting ssh pty", "error", err)
			return nil, err
		}
	}

	// Pull the environment pairs back part and send them via the traditional
	// ssh environment variable path.
	for _, pair := range es.Environment {
		idx := strings.IndexByte(pair, '=')
		if idx != -1 {
			sess.Setenv(pair[:idx], pair[idx+1:])
		}
	}

	// Just an FYI for the user of the session to know what deployment they're
	// using.
	sess.Setenv("WAYPOINT_DEPLOYMENT", deployment.Id)

	log.Info("starting shell")

	// SSH takes the command to run as just a string that is shell escaped, etc.
	// So we need to turn the arguments into something the remote side can
	// pull apart to use.
	var parts []string

	for _, arg := range es.Arguments {
		if shouldQuote(arg) {
			arg = strconv.Quote(arg)
		}

		parts = append(parts, arg)
	}

	rewriteLine(es.Output, "")

	err = sess.Run(strings.Join(parts, " "))

	log.Info("shell finished")

	var ec component.ExecResult
	if err != nil {
		if ee, ok := err.(*ssh.ExitError); ok {
			log.Info("exec command exited non-zero", "code", ee.ExitStatus())
			ec.ExitCode = ee.ExitStatus()
		} else {
			log.Error("error running ssh shell", "error", err)
			return nil, err
		}

	}

	return &ec, nil
}

// Used for our simple shell quoting loop above.
func shouldQuote(s string) bool {
	for _, c := range s {
		if (c >= '0' && c <= '9') || (c >= 'A' && c <= 'z') || c == '-' || c == '_' {
			continue
		}

		return true
	}

	return false
}

// Destroy deletes the AWS Lambda revision
func (p *Platform) Destroy(
	ctx context.Context,
	log hclog.Logger,
	deployment *Deployment,
	ui terminal.UI,
) error {
	// We'll update the user in real time
	st := ui.Status()
	defer st.Close()

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: p.config.Region,
		Logger: log,
	})
	if err != nil {
		return err
	}

	elbsrv := elbv2.New(sess)

	log.Debug("deleting target group", "arn", deployment.TargetGroupArn)
	st.Update("Deleting target group...")

	_, err = elbsrv.DeleteTargetGroup(&elbv2.DeleteTargetGroupInput{
		TargetGroupArn: &deployment.TargetGroupArn,
	})
	if err != nil {
		return err
	}

	st.Step(terminal.StatusOK, "Deleted target group")
	st.Update("Deleting Lambda function version " + deployment.Version)

	lamSvc := lambda.New(sess)

	if deployment.Version != "" {
		_, err = lamSvc.DeleteFunction(&lambda.DeleteFunctionInput{
			FunctionName: aws.String(deployment.FuncArn),
			Qualifier:    aws.String(deployment.Version),
		})
	}
	st.Step(terminal.StatusOK, "Deleted Lambda function version")

	return err
}

// DestroyWorkspace deletes other bits we created
func (p *Platform) DestroyWorkspace(
	ctx context.Context,
	log hclog.Logger,
	deployment *Deployment,
	app *component.Source,
	ui terminal.UI,
) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	step := sg.Add("Connecting to AWS")
	defer func() {
		step.Abort()
	}()

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: p.config.Region,
		Logger: log,
	})
	if err != nil {
		return err
	}

	lamSvc := lambda.New(sess)

	step.Done()
	step = sg.Add("Deleting Lambda function")

	_, err = lamSvc.DeleteFunction(&lambda.DeleteFunctionInput{
		FunctionName: aws.String(app.App),
	})

	if err != nil {
		return err
	}

	// If the user specified a role, we don't delete it
	if p.config.RoleArn == "" {
		svc := iam.New(sess)

		roleName := "lambda-" + app.App

		step.Done()
		step = sg.Add("Deleting automatically created IAM role...")

		log.Info("attempting to delete role", "role-name", roleName)

		_, err = svc.DetachRolePolicy(&iam.DetachRolePolicyInput{
			RoleName:  aws.String(roleName),
			PolicyArn: aws.String("arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"),
		})
		if err != nil {
			return err
		}

		_, err = svc.DeleteRole(&iam.DeleteRoleInput{
			RoleName: aws.String(roleName),
		})
		if err != nil {
			return err
		}

		step.Update("IAM role deleted")
		step.Done()
	}

	return err
}

func (p *Platform) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&Config{}), docs.FromFunc(p.DeployFunc()))
	if err != nil {
		return nil, err
	}

	doc.Description("Deploy functions as OCI Images to AWS Lambda")

	doc.Example(
		`
deploy {
	use "aws-lambda" {
		region = "us-east-1"
		memory = 512
	}
}
`)

	doc.Input("ecr.Image")
	doc.Output("lambda.Deployment")

	doc.SetField(
		"region",
		"the AWS region for the ECS cluster",
	)

	doc.SetField(
		"iam_role",
		"an IAM Role specified by ARN that will be used by the Lambda at execution time",
		docs.Default("created automatically"),
	)

	doc.SetField(
		"memory",
		"the amount of memory, in megabytes, to assign the function",
		docs.Default("265"),
	)

	doc.SetField(
		"timeout",
		"the number of seconds a function has to return a result",
		docs.Default("60"),
	)

	doc.SetField(
		"architecture",
		"The instruction set architecture that the function supports. Valid values are: \"x86_64\", \"arm64\"",
		docs.Default("x86_64"),
	)

	return doc, nil
}

// Config is the configuration structure for the Platform.
// Validation tags are provided by Go Pkg Validator
// https://pkg.go.dev/gopkg.in/go-playground/validator.v10?tab=doc
type Config struct {
	// The AWS region to create the Lambda function in
	Region string `hcl:"region"`

	// The IAM role to associate with the Lambda, specified by ARN. If no value is provided,
	// a role will be generated automatically.
	RoleArn string `hcl:"iam_role,optional"`

	// The amount of memory, measured in megabytes, to allocate the function.
	// Defaults to 256
	Memory int `hcl:"memory,optional"`

	// The number of seconds to wait for a function to complete it's work.
	// Defaults to 256
	Timeout int `hcl:"timeout,optional"`

	// The instruction set architecture that the function supports.
	// Valid values are: "x86_64", "arm64"
	// Defaults to "x86_64".
	Architecture string `hcl:"architecture,optional"`
}

var (
	_ component.Platform     = (*Platform)(nil)
	_ component.Configurable = (*Platform)(nil)
)
