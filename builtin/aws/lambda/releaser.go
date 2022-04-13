package lambda

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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
			resource.WithState(&Resource_FunctionUrl{}),
			resource.WithCreate(r.resourceFunctionUrlCreate),
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

var (
	DefaultFunctionUrlAuthType = lambda.FunctionUrlAuthTypeNone
)

func (r *Releaser) resourceFunctionUrlCreate(
	ctx context.Context,
	log hclog.Logger,
	sess *session.Session,
	sg terminal.StepGroup,
	ui terminal.UI,
	dep *Deployment,
	state *Resource_FunctionUrl,
) error {
	lambdasrv := lambda.New(sess)

	functionUrlAuthType := DefaultFunctionUrlAuthType
	if r.config.AuthType != "" {
		functionUrlAuthType = strings.ToUpper(r.config.AuthType)
	}

	// TODO(thiskevinwang): source cors from HCL config
	cors := lambda.Cors{}

	addPermissionInput := lambda.AddPermissionInput{
		Action:              aws.String("lambda:InvokeFunctionUrl"),
		FunctionUrlAuthType: aws.String(functionUrlAuthType),
		FunctionName:        aws.String(dep.FuncArn),
		Principal:           aws.String("*"),
		StatementId:         aws.String("FunctionURLAllowPublicAccess"),
	}

	createFunctionUrlConfigInput := lambda.CreateFunctionUrlConfigInput{
		AuthType:     aws.String(functionUrlAuthType),
		FunctionName: &dep.FuncArn,
		// TODO(thiskevinwang): make Cors configurable via HCL
		Cors: &cors,
	}

	step := sg.Add("Creating permissions for public access to the lambda URL...")
	defer step.Abort()

	// Grant public/anonymous access to the lambda URL
	if _, err := lambdasrv.AddPermission(&addPermissionInput); err != nil {
		log.Error("Error creating permission", "error", err)
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "ResourceConflictException":
				// permissions already exist. likely safe to continue
				step.Update("Permissions for public access access already exist")
			default:
				step.Update("Error creating permissions: %q, %q", aerr.Code(), aerr.Message())
				return err
			}
		} else {
			return err
		}
	} else {
		step.Update("Created permissions for public access access to the lambda URL")
	}
	step.Done()

	step = sg.Add("Creating Lambda URL...")
	defer step.Abort()

	cfo, err := lambdasrv.CreateFunctionUrlConfig(&createFunctionUrlConfigInput)
	if err != nil {
		log.Error("Error creating function url config", "error", err)
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "ResourceConflictException":
				// function url config already exists. get it. maybe safe to continue
				step.Update("Function url config already exists")

				// retrieve existing function url to update state
				if gfc, err := lambdasrv.GetFunctionUrlConfig(&lambda.GetFunctionUrlConfigInput{
					FunctionName: aws.String(dep.FuncArn),
				}); err != nil {
					// this should not realistically occur
					log.Error("Error getting function url config", "error", err)
					return err
				} else {
					state.Url = *gfc.FunctionUrl
					step.Update("Reusing existing Lambda URL: %q", state.Url)
				}
			default:
				step.Update("Error creating function url config: %q, %q", aerr.Code(), aerr.Message())
				return err
			}
		} else {
			return err
		}
	} else {
		// update state
		state.Url = *cfo.FunctionUrl
		step.Update("Created Lambda URL: %q", state.Url)
	}

	step.Done()

	return nil
}

func (r *Releaser) Release(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	ui terminal.UI,
	dep *Deployment,
) (*Release, error) {
	sg := ui.StepGroup()
	defer sg.Wait()

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
		Url:           fnUrlState.Url,
		FuncArn:       dep.FuncArn,
		VerArn:        dep.VerArn,
		ResourceState: rm.State(),
	}, nil
}

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
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Gathering health report for lambda url deployment: %q", release.Url)
	defer s.Done()

	report := sdk.StatusReport{
		External: true,
	}

	// check the function url status
	log.Info("Checking function url status...", "url", release.Url)
	if resp, err := http.Get(release.Url); err != nil {
		log.Error("Failed to get a response from the Lambda URL: %s", err)
		report.Health = sdk.StatusReport_UNKNOWN
		report.HealthMessage = "Failed to get a response from the Lambda URL"
	} else {
		switch resp.StatusCode {
		case http.StatusOK:
			log.Info("Lambda URL returned a 200")
			report.Health = sdk.StatusReport_READY
			report.HealthMessage = "Lambda URL appears to be healthy"
			break
		default:
			log.Error("Lambda URL returned a non-200 response: %d", resp.StatusCode)
			report.Health = sdk.StatusReport_DOWN
			report.HealthMessage = fmt.Sprintf("Lambda URL returned a non-200 response: %d", resp.StatusCode)
			break
		}
	}

	return &report, nil
}

func (r *Releaser) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&ReleaserConfig{}))
	if err != nil {
		return nil, err
	}

	doc.Description("Create an AWS Lambda function URL")

	doc.Example(
		`
release {
	use "aws-lambda" {
		auth_type = "NONE"
	}
}
`)

	doc.Input("lambda.Deployment")
	doc.Output("lambda.Release")

	doc.SetField(
		"auth_type",
		"the Lambda function URL auth type",
		docs.Summary(
			"The AuthType parameter determines how Lambda authenticates or authorizes requests to your function URL. Must be either `AWS_IAM` or `NONE`.",
		),
		docs.Default("NONE"),
	)

	// todo(thiskevinwang): CORS
	// doc.SetField(
	// 	"cors",
	// 	"the CORS configuration for the Lambda function URL",
	// 	docs.Summary(
	// 		"Use CORS to allow access to your function URL from any domain. You can also use CORS to control access for specific HTTP headers and methods in requests to your function URL",
	// 	),
	// 	docs.Default("{}"),
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
