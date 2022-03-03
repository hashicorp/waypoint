package ecr

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/aws/utils"
	"github.com/hashicorp/waypoint/builtin/docker"
)

// Registry represents access to an AWS registry.
type Registry struct {
	config Config
}

// Config implements Configurable
func (r *Registry) Config() (interface{}, error) {
	return &r.config, nil
}

func (r *Registry) AccessInfoFunc() interface{} {
	return r.AccessInfo
}

// PushFunc implements component.Registry
func (r *Registry) PushFunc() interface{} {
	return r.Push
}

type ecrInfo struct {
	repo, token string
}

func (r *Registry) setupRepo(
	ctx context.Context,
	log hclog.Logger,
	ui terminal.UI,
	sg terminal.StepGroup,
	src *component.Source,
) (*ecrInfo, error) {
	// If there is no region setup.  Try and load it from environment variables.
	if r.config.Region == "" {
		r.config.Region = os.Getenv("AWS_REGION")

		if r.config.Region == "" {
			r.config.Region = os.Getenv("AWS_REGION_DEFAULT")
		}
	}

	if r.config.Region == "" {
		return nil, status.Error(
			codes.FailedPrecondition,
			"Please set your aws region in the deployment config, or set the environment variable 'AWS_REGION' or 'AWS_DEFAULT_REGION'")
	}

	step := sg.Add("Connecting to AWS")
	defer func() { step.Abort() }()

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: r.config.Region,
		Logger: log,
	})
	if err != nil {
		return nil, err
	}
	svc := ecr.New(sess)

	step.Update("Getting ECR Authentication token...")
	gat, err := svc.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return nil, err
	}
	if len(gat.AuthorizationData) == 0 {
		return nil, fmt.Errorf("No authorization tokens provided")
	}

	step.Update("All services available.")
	step.Done()

	step = sg.Add("Calculating repository name")
	repoName := r.config.Repository
	if repoName == "" {
		log.Info("infering ECR repo name from app name")
		repoName = "waypoint-" + src.App
	}

	step.Update("Set ECR Repository name to '%s'", repoName)
	step.Done()

	repOut, err := svc.DescribeRepositories(&ecr.DescribeRepositoriesInput{
		RepositoryNames: []*string{aws.String(repoName)},
	})
	if err != nil {
		_, ok := err.(*ecr.RepositoryNotFoundException)
		if !ok {
			return nil, err
		}
	}

	var repo *ecr.Repository
	if repOut == nil || len(repOut.Repositories) == 0 {
		log.Info("no ECR repository detected, creating", "name", repoName)

		step = sg.Add("Creating new repository: %s", repoName)
		out, err := svc.CreateRepository(&ecr.CreateRepositoryInput{
			RepositoryName: aws.String(repoName),
			Tags: []*ecr.Tag{
				{
					Key:   aws.String("waypoint-app"),
					Value: aws.String(src.App),
				},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("Unable to create repository: %w", err)
		}

		step.Done()

		repo = out.Repository
	} else {
		repo = repOut.Repositories[0]
	}

	uptoken := *gat.AuthorizationData[0].AuthorizationToken
	data, err := base64.StdEncoding.DecodeString(uptoken)
	if err != nil {
		return nil, err
	}

	info := &ecrInfo{
		repo:  *repo.RepositoryUri,
		token: string(data[4:]),
	}

	return info, nil
}

// AccessInfo returns information on how to access the registry
func (r *Registry) AccessInfo(
	ctx context.Context,
	log hclog.Logger,
	ui terminal.UI,
	src *component.Source,
) (*docker.AccessInfo, error) {
	sg := ui.StepGroup()
	defer sg.Wait()

	info, err := r.setupRepo(ctx, log, ui, sg, src)
	if err != nil {
		return nil, err
	}

	ai := &docker.AccessInfo{
		Image: info.repo,
		Tag:   r.config.Tag,

		Auth: &docker.AccessInfo_UserPass_{
			UserPass: &docker.AccessInfo_UserPass{
				Username: "AWS",
				Password: info.token,
			},
		},
	}

	return ai, nil
}

// Push pushes an image to the registry.
func (r *Registry) Push(
	ctx context.Context,
	log hclog.Logger,
	img *docker.Image,
	ui terminal.UI,
	src *component.Source,
) (*Image, error) {
	sg := ui.StepGroup()
	defer sg.Wait()

	info, err := r.setupRepo(ctx, log, ui, sg, src)
	if err != nil {
		return nil, err
	}

	// Use the authorization token to create the base64 package that
	// Docker requires to perform authentication.
	authInfo := map[string]string{
		"username": "AWS",
		"password": info.token,
	}

	authData, err := json.Marshal(authInfo)
	if err != nil {
		return nil, err
	}
	encodedAuth := base64.StdEncoding.EncodeToString(authData)

	// Create the Docker registry plugin. We use the Docker registry plugin
	// directly so we can inherit all of its logic around Dockerless if necessary
	// and other complexities since at this point we're doing a standard
	// tag and push.
	dockerReg := &docker.Registry{}
	raw, err := dockerReg.Config()
	if err != nil {
		return nil, err
	}
	dockerConfig := raw.(*docker.Config)
	dockerConfig.EncodedAuth = encodedAuth
	dockerConfig.Image = info.repo
	dockerConfig.Tag = r.config.Tag

	// Build
	dockerImg, err := dockerReg.Push(ctx, img, ui, log)
	if err != nil {
		return nil, err
	}

	return &Image{
		Image:        dockerImg.Image,
		Tag:          dockerImg.Tag,
		Architecture: dockerImg.Architecture,
	}, nil
}

// Config is the configuration structure for the registry.
type Config struct {
	// AWS Region to access ECR in
	Region string `hcl:"region,optional"`

	// Repository to store the image into
	Repository string `hcl:"repository,optional"`

	// Tag is the tag to apply to the image.
	Tag string `hcl:"tag,attr"`
}

func (r *Registry) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&Config{}), docs.FromFunc(r.PushFunc()))
	if err != nil {
		return nil, err
	}

	doc.Description("Store a docker image within an Elastic Container Registry on AWS")

	doc.Example(
		`
registry {
    use "aws-ecr" {
      region = "us-east-1"
      tag = "latest"
    }
}
`)

	doc.Input("docker.Image")
	doc.Output("ecr.Image")

	doc.AddMapper(
		"ecr.Image",
		"docker.Image",
		"Allow an ECR Image to be used as a standard docker.Image",
	)

	doc.SetField(
		"region",
		"the AWS region the ECR repository is in",
		docs.Summary("if not set uses the environment variable AWS_REGION or AWS_REGION_DEFAULT"),
		docs.EnvVar("AWS_REGION"),
		docs.EnvVar("AWS_REGION_DEFAULT"),
	)

	doc.SetField(
		"repository",
		"the ECR repository to store the image into",
		docs.Summary(
			"This defaults to waypoint- then the application name. The repository will be automatically created if needed",
		),
	)

	doc.SetField(
		"tag",
		"the docker tag to assign to the new image",
	)

	return doc, nil
}

var _ component.Documented = (*Registry)(nil)
var _ component.Registry = (*Registry)(nil)
var _ component.RegistryAccess = (*Registry)(nil)
