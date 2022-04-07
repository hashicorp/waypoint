package lambda

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
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
			resource.WithName("function_url"),
			// WithState is needed for `rm.Resource("function_url").State()` to succeed
			resource.WithState(&Resource_FunctionUrl{}),
			resource.WithCreate(r.resourceFunctionUrlCreate),
			resource.WithDestroy(r.resourceFunctionUrlDestroy),
		)),
	)
}

func (r *Releaser) getSession(
	_ context.Context,
	log hclog.Logger,
) (*session.Session, error) {
	return utils.GetSession(&utils.SessionConfig{
		Logger: log,
	})
}

func (r *Releaser) resourceFunctionUrlCreate(
	ctx context.Context,
	log hclog.Logger,
	sess *session.Session,
	sg terminal.StepGroup,
	ui terminal.UI,
	dep *Deployment,
) error {
	log.Info("Creating Lambda URL...", "VerArn", dep.VerArn, "FuncArn", dep.FuncArn)

	lambdasrv := lambda.New(sess)

	log.Info("Creating alias...")
	log.Info("Version: " + dep.Version)
	log.Info("FuncArn: " + dep.FuncArn)
	// create a function alias so that we can create a function url
	// https://docs.aws.amazon.com/lambda/latest/dg/API_CreateAlias.html
	qualifier := "Alias_" + dep.Version
	a, err := lambdasrv.CreateAlias(&lambda.CreateAliasInput{
		// Alias name cannot be numeric-only
		Name:            aws.String(qualifier),
		Description:     aws.String("Waypoint Lambda Alias"),
		FunctionName:    aws.String(dep.FuncArn),
		FunctionVersion: aws.String(dep.Version),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Error("Error creating alias", "error", aerr.Code(), "message", aerr.Message())
		}
		return err
	}
	log.Info("Created alias", "alias", *a.AliasArn)

	// create permissions prior to creating function url
	// https://us-east-1.console.aws.amazon.com/lambda/services/ajax?operation=addPermission&locale=en
	log.Info("Creating permission...")
	p, err := lambdasrv.AddPermission(&lambda.AddPermissionInput{
		FunctionUrlAuthType: aws.String(lambda.FunctionUrlAuthTypeNone),
		FunctionName:        a.AliasArn,
		Action:              aws.String("lambda:InvokeFunctionUrl"),
		Principal:           aws.String("*"),
		Qualifier:           aws.String(qualifier),
		StatementId:         aws.String("FunctionURLAllowPublicAccess"),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Error("Error creating permission", "error", aerr.Code(), "message", aerr.Message())
		}
	}
	log.Info("Created permission", "permission", *p.Statement)

	// https://us-east-1.console.aws.amazon.com/lambda/services/ajax?operation=createFunctionUrlConfig&locale=en
	o, err := lambdasrv.CreateFunctionUrlConfig(&lambda.CreateFunctionUrlConfigInput{
		// When you choose auth type NONE, Lambda [DASHBOARD OPERATION ONLY] automatically creates the
		// following resource-based policy and attaches it to your function.
		// This policy makes your function public to anyone with the function URL.
		// You can edit the policy later. To limit access to authenticated IAM users and roles, choose auth type AWS_IAM.
		AuthType:     aws.String(lambda.FunctionUrlAuthTypeNone),
		FunctionName: a.AliasArn,
		Cors:         &lambda.Cors{},
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Error("Error creating function url config", "error", aerr.Code(), "message", aerr.Message())
		}
		return err
	}
	log.Info("Created function url config", "url", *o.FunctionUrl)

	return nil
}

func (r *Releaser) resourceFunctionUrlDestroy(
	ctx context.Context,
	sess *session.Session,
	sg terminal.StepGroup,
	// state *Resource_LoadBalancer,
) error {
	step := sg.Add("Destroying Lambda URL...")
	step.Update("Destroyed Lambda URL...")
	step.Done()
	return nil
}

func (r *Releaser) Release(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	ui terminal.UI,
	// waypoint automatically injects the previous component's output
	// - https://www.waypointproject.io/docs/extending-waypoint/passing-values
	dep *Deployment,
) (*Release, error) {
	sg := ui.StepGroup()
	defer sg.Wait()

	log.Info("Creating Function URL...")
	log.Info("Deployment details", "deployment", dep)
	log.Info("Deployment details", "FuncArn", dep.FuncArn, "VerArn", dep.VerArn)

	// Create our resource manager and create
	rm := r.resourceManager(log)
	if err := rm.CreateAll(
		ctx, log, sg, ui, src,
		dep,
	); err != nil {
		log.Info("Error creating resources", "error", err)
		return nil, err
	}

	// Get our function url state to verify
	fnUrlState := rm.Resource("function_url").State().(*Resource_FunctionUrl)
	log.Info("Function URL state", "url", fnUrlState.Url)
	if fnUrlState == nil {
		return nil, status.Errorf(codes.Internal, "function url state is nil, this should never happen")
	}

	return &Release{
		Url:     fnUrlState.Url,
		FuncArn: dep.FuncArn,
		VerArn:  dep.VerArn,
	}, nil
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
		Logger: log,
	})
	if err != nil {
		return err
	}

	rm := r.resourceManager(log)

	// Destroy All
	return rm.DestroyAll(ctx, log, sg, ui, sess)
}

// ReleaserConfig is the configuration structure for the Releaser.
type ReleaserConfig struct {
	// "AWS_IAM" or "NONE"
	AuthType string `hcl:"auth_type,optional"`
	// todo(kevinwang): CORS
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
	return &report, nil
}

func (r *Releaser) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&ReleaserConfig{}))
	if err != nil {
		return nil, err
	}

	doc.Description("TODO: Add description")

	// doc.Input("alb.TargetGroup")
	// doc.Output("alb.Release")
	// doc.AddMapper(
	// 	"ec2.Deployment",
	// 	"alb.TargetGroup",
	// 	"Allow EC2 Deployments to be hooked up to an ALB",
	// )

	// doc.AddMapper(
	// 	"lambda.Deployment",
	// 	"alb.TargetGroup",
	// 	"Allow Lambda Deployments to be hooked up to an ALB",
	// )

	// doc.SetField(
	// 	"name",
	// 	"the name to assign the ALB",
	// 	docs.Summary(
	// 		"names have to be unique per region",
	// 	),
	// 	docs.Default("derived from application name"),
	// )

	// doc.SetField(
	// 	"port",
	// 	"the TCP port to configure the ALB to listen on",
	// 	docs.Default("80 for HTTP, 443 for HTTPS"),
	// )

	// doc.SetField(
	// 	"subnets",
	// 	"the subnet ids to allow the ALB to run in",
	// 	docs.Default("public subnets in the account default VPC"),
	// )

	// doc.SetField(
	// 	"certificate",
	// 	"ARN for the certificate to install on the ALB listener",
	// 	docs.Summary(
	// 		"when this is set, the port automatically changes to 443 unless",
	// 		"overriden in this configuration",
	// 	),
	// )

	// doc.SetField(
	// 	"zone_id",
	// 	"Route53 ZoneID to create a DNS record into",
	// 	docs.Summary(
	// 		"set along with domain_name to have DNS automatically setup for the ALB",
	// 	),
	// )

	// doc.SetField(
	// 	"domain_name",
	// 	"Fully qualified domain name to set for the ALB",
	// 	docs.Summary(
	// 		"set along with zone_id to have DNS automatically setup for the ALB.",
	// 		"this value should include the full hostname and domain name, for instance",
	// 		"app.example.com",
	// 	),
	// )

	// doc.SetField(
	// 	"listener_arn",
	// 	"the ARN on an existing ALB to configure",
	// 	docs.Summary(
	// 		"when this is set, no ALB or Listener is created. Instead the application is",
	// 		"configured by manipulating this existing Listener. This allows users to",
	// 		"configure their ALB outside waypoint but still have waypoint hook the application",
	// 		"to that ALB",
	// 	),
	// )

	return doc, nil
}

var (
	_ component.ReleaseManager = (*Releaser)(nil)
	_ component.Configurable   = (*Releaser)(nil)
	_ component.Destroyer      = (*Releaser)(nil)
	_ component.Documented     = (*Releaser)(nil)
	_ component.Status         = (*Releaser)(nil)
)
