package lambda

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/subtle"
	"crypto/x509"
	"encoding/base64"
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

// ExecFunc implements component.Execer
func (p *Platform) ExecFunc() interface{} {
	return p.Exec
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

// SetupRole creates an IAM role that will be used by the Lambda function. The role
// has the basic lambda execution rules attached to it.
func (p *Platform) SetupRole(L hclog.Logger, app *component.Source, sess *session.Session) (string, error) {
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
	// Start building our deployment since we use this information
	deployment := &Deployment{}
	id, err := component.Id()
	if err != nil {
		return nil, err
	}

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: p.config.Region,
	})
	if err != nil {
		return nil, err
	}

	roleArn := p.config.RoleArn

	if roleArn == "" {
		arn, err := p.SetupRole(log, src, sess)
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

	lamSvc := lambda.New(sess)
	curFunc, err := lamSvc.GetFunction(&lambda.GetFunctionInput{
		FunctionName: aws.String(src.App),
	})

	var funcarn string

	// If the function exists, we update it's code rather than create a new one.
	// Will then publish this new code to create a stable identifier for it.
	if err == nil {

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
			FunctionName: aws.String(src.App),
			ImageUri:     aws.String(img.Name()),
		})

		if err != nil {
			return nil, err
		}

		funcarn = *funcCfg.FunctionArn

		if err != nil {
			return nil, err
		}

		ui.Output("Updated Lambda Function: %s", funcarn)
	} else {
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
			ImageConfig: &lambda.ImageConfig{},
		})

		if err != nil {
			return nil, err
		}

		funcarn = *funcOut.FunctionArn

		ui.Output("Created Lambda Function: %s", funcarn)
	}

	st := ui.Status()
	defer st.Close()

	st.Update("Waiting for Lambda function to be processed")
	// The image is never ready right away, AWS has to process it, so we wait
	// 3 seconds before trying to publish the version

	time.Sleep(3 * time.Second)

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
		return nil, fmt.Errorf("Lambda was unable to prepare the function in the aloted time")
	}

	st.Close()

	verarn := *ver.FunctionArn

	ui.Output("Published Lambda Function: %s (%s)", verarn, *ver.Version)

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

	svc := elbv2.New(sess)

	serviceName := fmt.Sprintf("%s-%s", src.App, id)

	// We have to clamp at a length of 32 because the Name field to CreateTargetGroup
	// requires that the name is 32 characters or less.
	if len(serviceName) > 32 {
		serviceName = serviceName[:32]
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

	deployment.Id = id
	deployment.FuncArn = funcarn
	deployment.VerArn = verarn
	deployment.Version = *ver.Version
	deployment.TargetGroupArn = *ctg.TargetGroups[0].TargetGroupArn

	return deployment, nil
}

// This is used by the Exec plugin to provide some helpful output
// while we proper the exec environment.
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

	// Generate 2 keys, one is the host key that the server will use, and
	// the other is the user key the client will use. That way, both parties
	// can validate they are who we think they are.
	// We armor the keys with base64 because they're going to be passed in
	// ECS environment variable fields.

	hostkey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	hoststr := base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(hostkey))

	sshHostKey, err := ssh.NewSignerFromKey(hostkey)
	if err != nil {
		return nil, err
	}

	userkey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	userstr := base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(&userkey.PublicKey))

	sshKey, err := ssh.NewSignerFromKey(userkey)
	if err != nil {
		return nil, err
	}

	var esl ecsLauncher
	esl.Region = p.config.Region
	esl.PublicKey = userstr
	esl.HostKey = hoststr
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
	cfg.User = "evan"
	cfg.Auth = []ssh.AuthMethod{
		ssh.PublicKeys(sshKey),
	}

	expectedHost := sshHostKey.PublicKey().Marshal()

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
		"the amonut of memory, in megabytes, to assign the function",
		docs.Default("265"),
	)

	doc.SetField(
		"timeout",
		"the number of seconds a function has to return a result",
		docs.Default("60"),
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
}

var (
	_ component.Platform     = (*Platform)(nil)
	_ component.Configurable = (*Platform)(nil)
)
