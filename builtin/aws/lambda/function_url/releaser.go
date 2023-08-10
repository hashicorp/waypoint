// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package function_url

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	lambdaplugin "github.com/hashicorp/waypoint/builtin/aws/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/framework/resource"
	sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/aws/utils"
)

type Releaser struct {
	config ReleaserConfig
}

// Lambda GetPolicy() returns a JSON string of its IAM policy, so
// we need to unmarshal it into this struct in order to use it.
//
// See https://github.com/aws/aws-sdk-go/issues/127
// and https://github.com/aws/aws-sdk-go-v2/issues/225
type IAMPolicy struct {
	Version   string `json:"Version"`
	ID        string `json:"Id"`
	Statement []struct {
		Sid    string `json:"Sid"`
		Effect string `json:"Effect"`
		// Principal can be either a string, like "*", or a struct
		// like { "AWS": "arn:aws:iam::123456789012:root" }
		Principal interface{} `json:"Principal"`
		Action    string      `json:"Action"`
		Resource  string      `json:"Resource"`
		Condition struct {
			StringEquals struct {
				LambdaFunctionURLAuthType string `json:"lambda:FunctionUrlAuthType"`
			} `json:"StringEquals"`
		} `json:"Condition"`
	} `json:"Statement"`
}

// Config implements Configurable
func (r *Releaser) Config() (interface{}, error) {
	return &r.config, nil
}

// ConfigSet is called after a configuration has been decoded
// we can use this to validate the config
func (r *Releaser) ConfigSet(config interface{}) error {
	rc, ok := config.(*ReleaserConfig)
	if !ok {
		// this should never happen
		return fmt.Errorf("Invalid configuration, expected *lambda.ReleaserConfig, got %s", reflect.TypeOf(config))
	}

	err := utils.Error(validation.ValidateStruct(rc,
		validation.Field(&rc.Principal,
			validation.Empty.When(rc.AuthType != "" && rc.AuthType != lambda.FunctionUrlAuthTypeAwsIam).Error("principal requires auth_type to be set to \"AWS_IAM\""),
		),
	))

	if err != nil {
		return err
	}

	return nil
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
			resource.WithName("function_permission"),
			resource.WithCreate(r.resourceFunctionPermissionCreate),
		)),

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
	dep *lambdaplugin.Deployment,
) (*session.Session, error) {
	return utils.GetSession(&utils.SessionConfig{
		Region: dep.Region,
		Logger: log,
	})
}

var (
	// If auth_type is not set we'll default to "NONE", allowing public access to the function URL
	DefaultFunctionUrlAuthType = lambda.FunctionUrlAuthTypeNone
	// If principal is not set we'll allow any authenticated AWS user to invoke the function URL
	DefaultPrincipal = "*"
)

func (r *Releaser) resourceFunctionPermissionCreate(
	ctx context.Context,
	log hclog.Logger,
	sess *session.Session,
	sg terminal.StepGroup,
	ui terminal.UI,
	dep *lambdaplugin.Deployment,
) error {
	step := sg.Add("Checking function permission...")
	defer step.Abort()

	lambdasrv := lambda.New(sess)

	var reset bool
	var revisionId string

	// use a single StatementId for now
	statementId := "waypoint-function-url-access"

	authtype := strings.ToUpper(r.config.AuthType)
	if authtype == "" {
		authtype = DefaultFunctionUrlAuthType
	}
	principal := r.config.Principal
	if principal == "" {
		principal = DefaultPrincipal
	}

	// check if principal or auth type have changed
	if gpo, err := lambdasrv.GetPolicy(&lambda.GetPolicyInput{
		FunctionName: aws.String(dep.FuncArn),
	}); err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			// noop if ResourceNotFoundException
			// This will likely happen if the function has 0 permissions
			case lambda.ErrCodeResourceNotFoundException:
				step.Update("No permissions found. Creating one...")
				reset = true
			default:
				step.Update("Failed to get policy: %s", err)
				return err
			}
		} else {
			return err
		}
	} else {
		policy := IAMPolicy{}
		revisionId = *gpo.RevisionId
		if err := json.Unmarshal([]byte(*gpo.Policy), &policy); err != nil {
			// failed to unmarshal policy, this should never happen
			log.Info("Failed to unmarshal policy: %s", err)
			return err
		}

		statementFound := false
		// determine if we should update permissions
		// find the statement
		for _, st := range policy.Statement {
			if st.Sid == statementId {
				// found the previous statement, check if the principal has changed
				statementFound = true
				pType := reflect.TypeOf(st.Principal)

				switch pType {
				case reflect.TypeOf(""):
					if st.Principal != principal {
						// principal has changed, we should update permissions
						reset = true
					}
				case reflect.TypeOf(map[string]interface{}{}):
					if st.Principal.(map[string]interface{})["AWS"].(string) != principal {
						// principal has changed, we should update permissions
						reset = true
					}
				default:
					// this should never happen
					return status.Errorf(codes.Unknown, "Unknown principal type from AWS policy: %s", pType)
				}

				if st.Condition.StringEquals.LambdaFunctionURLAuthType != authtype {
					// auth type has changed, we should update permissions
					reset = true
				}
			}
		}
		if !statementFound {
			// statement not found, we should create permissions
			reset = true
		}
	}

	if reset {
		step.Update("Updating permissions to invoke lambda URL...")

		// attempt to remove the old permission
		if _, err := lambdasrv.RemovePermission(&lambda.RemovePermissionInput{
			FunctionName: aws.String(dep.FuncArn),
			RevisionId:   aws.String(revisionId),
			StatementId:  aws.String(statementId),
		}); err != nil {
			// no-op if there is no permission
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case "ResourceNotFoundException":
					// Non-fatal
					log.Warn("Failed to remove permission", "error", err)
				default:
					log.Error("Failed to remove permission", "error", err)
					return err
				}
			}
		}

		// add new permission
		if _, err := lambdasrv.AddPermission(&lambda.AddPermissionInput{
			Action:              aws.String("lambda:InvokeFunctionUrl"),
			FunctionUrlAuthType: aws.String(authtype),
			FunctionName:        aws.String(dep.FuncArn),
			Principal:           aws.String(principal),
			StatementId:         aws.String(statementId),
		}); err != nil {
			log.Error("Error creating permission", "error", err)
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case "ResourceConflictException":
					// permissions already exist. This should not happen since we remove the permission first
					step.Update("Permissions to invoke lambda URL access already exist")
				default:
					step.Update("Error creating permissions: %s", err)
					return err
				}
			} else {
				return err
			}
		} else {
			step.Update("Updated permissions to invoke lambda URL")
		}
	} else {
		step.Update("No permissions need to be updated")
	}

	step.Done()
	return nil
}

func (r *Releaser) resourceFunctionUrlCreate(
	ctx context.Context,
	log hclog.Logger,
	sess *session.Session,
	sg terminal.StepGroup,
	ui terminal.UI,
	dep *lambdaplugin.Deployment,
	state *Resource_FunctionUrl,
) error {
	lambdasrv := lambda.New(sess)

	functionUrlAuthType := DefaultFunctionUrlAuthType
	if r.config.AuthType != "" {
		functionUrlAuthType = strings.ToUpper(r.config.AuthType)
	}

	corsCfg := r.config.Cors
	if corsCfg == nil {
		corsCfg = &ReleaserConfigCors{}
	}

	cors := lambda.Cors{
		AllowCredentials: corsCfg.AllowCredentials,
		AllowHeaders:     corsCfg.AllowHeaders,
		AllowMethods:     corsCfg.AllowMethods,
		AllowOrigins:     corsCfg.AllowOrigins,
		ExposeHeaders:    corsCfg.ExposeHeaders,
		MaxAge:           corsCfg.MaxAge,
	}

	step := sg.Add("Creating Lambda URL...")
	defer step.Abort()

	createFunctionUrlConfigInput := lambda.CreateFunctionUrlConfigInput{
		AuthType:     aws.String(functionUrlAuthType),
		FunctionName: aws.String(dep.FuncArn),
		Cors:         &cors,
	}

	// Create or Update the lambda URL
	shouldUpdate := false
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
					// compare remote config to incoming config
					if functionUrlAuthType != *gfc.AuthType || !reflect.DeepEqual(&cors, gfc.Cors) {
						shouldUpdate = true
					} else {
						step.Update("Reusing existing Lambda URL: %q", *gfc.FunctionUrl)
						state.Url = *gfc.FunctionUrl
					}
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

	if shouldUpdate {
		step.Update("Updating Lambda URL Config...")
		if ufc, err := lambdasrv.UpdateFunctionUrlConfig(&lambda.UpdateFunctionUrlConfigInput{
			AuthType:     aws.String(functionUrlAuthType),
			FunctionName: aws.String(dep.FuncArn),
			Cors:         &cors,
		}); err != nil {
		} else {
			state.Url = *ufc.FunctionUrl
			step.Update("Updated Lambda URL Config: %q", state.Url)
		}
	}

	step.Done()

	return nil
}

func (r *Releaser) Release(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	ui terminal.UI,
	dep *lambdaplugin.Deployment,
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
	// Only permitted if AuthType is "AWS_IAM" otherwise defaults to "*"
	Principal string `hcl:"principal,optional"`
	// Configuration options for function url CORS
	Cors *ReleaserConfigCors `hcl:"cors,block"`
}

// Based on the Cors type from the AWS SDK, but with our HCL mappings.
// https://pkg.go.dev/github.com/aws/aws-sdk-go/service/lambda#Cors
type ReleaserConfigCors struct {
	// Whether to allow cookies or other credentials in requests to your function
	// URL. The default is false.
	AllowCredentials *bool `hcl:"allow_credentials,optional"`

	// The HTTP headers that origins can include in requests to your function URL.
	// For example: Date, Keep-Alive, X-Custom-Header.
	AllowHeaders []*string `hcl:"allow_headers,optional"`

	// The HTTP methods that are allowed when calling your function URL. For example:
	// GET, POST, DELETE, or the wildcard character (*).
	AllowMethods []*string `hcl:"allow_methods,optional"`

	// The origins that can access your function URL. You can list any number of
	// specific origins, separated by a comma. For example: https://www.example.com,
	// http://localhost:60905.
	//
	// Alternatively, you can grant access to all origins using the wildcard character
	// (*).
	AllowOrigins []*string `hcl:"allow_origins,optional"`

	// The HTTP headers in your function response that you want to expose to origins
	// that call your function URL. For example: Date, Keep-Alive, X-Custom-Header.
	ExposeHeaders []*string `hcl:"expose_headers,optional"`

	// The maximum amount of time, in seconds, that web browsers can cache results
	// of a preflight request. By default, this is set to 0, which means that the
	// browser doesn't cache results.
	MaxAge *int64 `hcl:"max_age,optional"`
}

func (r *Releaser) Status(
	ctx context.Context,
	log hclog.Logger,
	release *Release,
	ui terminal.UI,
) (*sdk.StatusReport, error) {
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Gathering health report for lambda: %q", release.Url)
	defer s.Done()

	report := sdk.StatusReport{
		External: true,
	}

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: release.Region,
		Logger: log,
	})
	if err != nil {
		return nil, err
	}

	log.Info("Checking function status...", "url", release.Url)
	lambdasrv := lambda.New(sess)
	if err := lambdasrv.WaitUntilFunctionActiveV2(&lambda.GetFunctionInput{
		FunctionName: aws.String(release.FuncArn),
	}); err != nil {
		log.Error("Error waiting for function to become active", "error", err)
		report.Health = sdk.StatusReport_DOWN
		report.HealthMessage = "Failed to wait for function to become active"
	} else {
		log.Info("Function is active")
		report.Health = sdk.StatusReport_READY
		report.HealthMessage = "Function is active"
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
	use "lambda-function-url" {
		auth_type = "NONE"
		cors {
			allow_methods = ["*"]
		}
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

	doc.SetField(
		"principal",
		"the principal to use when auth_type is `AWS_IAM`",
		docs.Summary(
			"The Principal parameter specifies the principal that is allowed to invoke the function.",
		),
		docs.Default("*"),
	)

	doc.SetField(
		"cors",
		"CORS configuration for the function URL",
		docs.Default("NONE"),
		docs.SubFields(func(d *docs.SubFieldDoc) {
			d.SetField(
				"allow_credentials",
				"Whether to allow cookies or other credentials in requests to your function URL.",
				docs.Default("false"),
			)
			d.SetField(
				"allow_headers",
				"The HTTP headers that origins can include in requests to your function URL. For example: Date, Keep-Alive, X-Custom-Header.",
				docs.Default("[]"),
			)
			d.SetField(
				"allow_methods",
				"The HTTP methods that are allowed when calling your function URL. For example: GET, POST, DELETE, or the wildcard character (*).",
				docs.Default("[]"),
			)
			d.SetField(
				"allow_origins",
				"The origins that can access your function URL. You can list any number of specific origins, separated by a comma. You can grant access to all origins using the wildcard character (*).",
				docs.Default("[]"),
			)
			d.SetField(
				"expose_headers",
				"The HTTP headers in your function response that you want to expose to origins that call your function URL. For example: Date, Keep-Alive, X-Custom-Header.",
				docs.Default("[]"),
			)
			d.SetField(
				"max_age",
				"The maximum amount of time, in seconds, that web browsers can cache results of a preflight request.",
				docs.Default("0"),
			)
		}),
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
