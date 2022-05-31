package ecrpull

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aws/aws-sdk-go/aws"
	awsecr "github.com/aws/aws-sdk-go/service/ecr"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/aws/ecr"
	"github.com/hashicorp/waypoint/builtin/aws/utils"
)

// Builder uses `docker build` to build a Docker iamge.
type Builder struct {
	config BuilderConfig
}

// BuildFunc implements component.Builder
func (b *Builder) BuildFunc() interface{} {
	return b.Build
}

// Config is the configuration structure for the registry.
type BuilderConfig struct {
	Region     string `hcl:"region,attr"`
	Repository string `hcl:"repository,attr"`
	Tag        string `hcl:"tag,attr"`
}

func (b *Builder) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&BuilderConfig{}), docs.FromFunc(b.BuildFunc()))
	if err != nil {
		return nil, err
	}

	doc.Description(`
Use an existing, pre-built AWS ECR image.

This builder attempts to find an image by repository and tag in the
specified region. If found, it will pass along the image information
to the next step.

This builder will not modify the image. 

If you wish to rename or retag an image, please use the "docker-pull" component
in conjunction with the "aws-ecr" registry option.
`)

	doc.Example(`
build {
  use "aws-ecr-pull" {
	region     = "us-east-1"
    repository = "deno-http"
    tag        = "latest"
  }
}
`)

	doc.Input("component.Source")
	doc.Output("ecr.Image")

	doc.SetField(
		"region",
		"the AWS region the ECR repository is in",
		docs.Summary(
			"if not set uses the environment variable AWS_REGION or AWS_REGION_DEFAULT.",
		),
	)

	doc.SetField(
		"repository",
		"the AWS ECR repository name",
	)

	doc.SetField(
		"tag",
		"the tag of the image to pull",
	)
	return doc, nil
}

// Config implements Configurable
func (b *Builder) Config() (interface{}, error) {
	return &b.config, nil
}

// Build
func (b *Builder) Build(ui terminal.UI, log hclog.Logger) (*ecr.Image, error) {

	// If there is no region setup. Try and load it from environment variables.
	if b.config.Region == "" {
		b.config.Region = os.Getenv("AWS_REGION")

		if b.config.Region == "" {
			b.config.Region = os.Getenv("AWS_REGION_DEFAULT")
		}
	}

	if b.config.Region == "" {
		return nil, status.Error(
			codes.FailedPrecondition,
			"Please set your aws region in the deployment config, or set the environment variable 'AWS_REGION' or 'AWS_DEFAULT_REGION'")
	}

	sg := ui.StepGroup()
	step := sg.Add("")
	defer func() {
		if step != nil {
			step.Abort()
		}
	}()

	step.Update("Connecting to AWS")
	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: b.config.Region,
		Logger: log,
	})

	if err != nil {
		log.Error("error connecting to AWS", "error", err)
		return nil, err
	}

	step.Done()

	step = sg.Add("Verifying image exists")
	svc := awsecr.New(sess)

	cfgTag := b.config.Tag
	cfgRepository := b.config.Repository

	imgs, err := svc.DescribeImages(&awsecr.DescribeImagesInput{
		RepositoryName: aws.String(cfgRepository),
		Filter: &awsecr.DescribeImagesFilter{
			TagStatus: aws.String("TAGGED"),
		},
	})
	if err != nil {
		log.Error("error describing images", "error", err, "repository", cfgRepository)
		return nil, err
	}

	if len(imgs.ImageDetails) == 0 {
		log.Error("no tagged images found", "repository", cfgRepository)
		return nil, status.Error(codes.FailedPrecondition, "No images found")
	}
	log.Debug("found images", "image count", len(imgs.ImageDetails))

	var output ecr.Image
	for _, img := range imgs.ImageDetails {
		for _, tag := range img.ImageTags {
			if *tag == cfgTag {

				output.Image = *img.RegistryId + ".dkr.ecr." + b.config.Region + ".amazonaws.com/" + cfgRepository
				output.Tag = *tag
				// TODO(kevinwang): Do we need to get architecture?
				// If we do, we can pull the image and inspect it via `cli.ImageInspectWithRaw`,
				// - prior art: /builtin/docker/builder.go -> Build
				// There is also an open issue for the ECR team to build a architecture feature into
				// the UI, which probably comes with a CLI/API change.
				// - see https://github.com/aws/containers-roadmap/issues/1591
				break
			}
		}
	}

	if output.Image == "" {
		log.Error("no matching image found", "tag", cfgTag, "repository", cfgRepository)
		return nil, status.Error(codes.FailedPrecondition, "No matching tags found")
	}

	step.Update("Using existing image: %s; TAG=%s", output.Image, output.Tag)
	step.Done()

	return &output, nil
}
