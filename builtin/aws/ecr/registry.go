package ecr

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/docker"
	"github.com/mattn/go-isatty"
)

// Registry represents access to an AWS registry.
type Registry struct {
	config Config
}

// Config implements Configurable
func (r *Registry) Config() (interface{}, error) {
	return &r.config, nil
}

// PushFunc implements component.Registry
func (r *Registry) PushFunc() interface{} {
	return r.Push
}

// Push pushes an image to the registry.
func (r *Registry) Push(
	ctx context.Context,
	img *docker.Image,
	ui terminal.UI,
) (*docker.Image, error) {
	stdout, _, err := ui.OutputWriters()
	if err != nil {
		return nil, err
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	cli.NegotiateAPIVersion(ctx)

	sess, err := session.NewSession(aws.NewConfig().WithRegion(r.config.Region))
	if err != nil {
		return nil, err
	}
	svc := ecr.New(sess)

	repOut, err := svc.DescribeRepositories(&ecr.DescribeRepositoriesInput{
		RepositoryNames: []*string{aws.String(r.config.Repository)},
	})
	if err != nil {
		return nil, err
	}

	if len(repOut.Repositories) == 0 {
		return nil, fmt.Errorf("Unknown repository requested: %s", r.config.Repository)
	}

	uri := repOut.Repositories[0].RepositoryUri

	target := &docker.Image{Image: *uri, Tag: r.config.Tag}

	ui.Output("Tagging Docker image: %s => %s", img.Name(), target.Name())

	err = cli.ImageTag(ctx, img.Name(), target.Name())
	if err != nil {
		return nil, err
	}

	ref, err := reference.ParseNormalizedNamed(target.Name())
	if err != nil {
		return nil, err
	}

	gat, err := svc.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return nil, err
	}

	if len(gat.AuthorizationData) == 0 {
		return nil, fmt.Errorf("No authorization tokens provided")
	}

	uptoken := *gat.AuthorizationData[0].AuthorizationToken

	data, err := base64.StdEncoding.DecodeString(uptoken)
	if err != nil {
		return nil, err
	}

	authInfo := map[string]string{
		"username": "AWS",
		"password": string(data[4:]),
	}

	authData, err := json.Marshal(authInfo)
	if err != nil {
		return nil, err
	}

	encodedAuth := base64.StdEncoding.EncodeToString(authData)

	options := types.ImagePushOptions{
		RegistryAuth: encodedAuth,
	}

	responseBody, err := cli.ImagePush(ctx, reference.FamiliarString(ref), options)
	if err != nil {
		return nil, err
	}

	defer responseBody.Close()

	var (
		termFd uintptr
		isTerm bool
	)

	if f, ok := stdout.(*os.File); ok {
		termFd = f.Fd()
		isTerm = isatty.IsTerminal(termFd)
	}

	err = jsonmessage.DisplayJSONMessagesStream(responseBody, stdout, termFd, isTerm, nil)
	if err != nil {
		return nil, err
	}

	ui.Output("Docker image pushed: %s", target.Name())

	return target, nil
}

// Config is the configuration structure for the registry.
type Config struct {
	// AWS Region to access ECR in
	Region string `hcl:"region,attr"`

	// Repository to store the image into
	Repository string `hcl:"repository,attr"`

	// Tag is the tag to apply to the image.
	Tag string `hcl:"tag,attr"`
}

func (r *Registry) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&Config{}))
	if err != nil {
		return nil, err
	}

	doc.Description("Store a docker image within an Elastic Container Registry on AWS")

	doc.Example(
		`
registry {
    use "aws-ecr" {
      region = "us-east-1"
      repository = "waypoint-example"
      tag = "latest"
    }
}
`)

	doc.Input("docker.Image")
	doc.Output("docker.Image")

	doc.SetField(
		"region",
		"the AWS region the ECR repository is in",
	)

	doc.SetField(
		"repository",
		"the ECR repository to store the image into",
		docs.Summary(
			"this ECR repository must already exist, waypoint will not create it",
		),
	)

	doc.SetField(
		"tag",
		"the docker tag to assign to the new image",
	)

	return doc, nil
}

var (
	_ component.Documented = (*Registry)(nil)
)
