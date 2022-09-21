package dockerref

import (
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	wpdocker "github.com/hashicorp/waypoint/builtin/docker"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

// Builder uses `docker build` to build a Docker iamge.
type Builder struct {
	config BuilderConfig
}

// BuildFunc implements component.Builder
func (b *Builder) BuildFunc() interface{} {
	return b.Build
}

// BuildODRFunc implements component.BuilderODR
func (b *Builder) BuildODRFunc() interface{} {
	return b.Build
}

// BuilderConfig is the configuration structure for the registry.
type BuilderConfig struct {
	// Image to pull
	Image string `hcl:"image,attr"`
	Tag   string `hcl:"tag,attr"`
}

func (b *Builder) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&BuilderConfig{}), docs.FromFunc(b.BuildFunc()))
	if err != nil {
		return nil, err
	}

	doc.Description(`
Use an existing, pre-built Docker image without modifying it.
`)

	doc.Example(`
build {
  use "docker-ref" {
    image = "gcr.io/my-project/my-image"
    tag   = "abcd1234"
  }
}
`)

	doc.Input("component.Source")
	doc.Output("docker.Image")

	doc.SetField(
		"image",
		"The image to pull.",
		docs.Summary(
			"This should NOT include the tag (the value following the ':' in a Docker image).",
			"Use `tag` to define the image tag.",
		),
	)

	doc.SetField(
		"tag",
		"The tag of the image to pull.",
	)

	return doc, nil
}

// Config implements Configurable
func (b *Builder) Config() (interface{}, error) {
	return &b.config, nil
}

// We use the struct form of arguments so that we can access named
// values (such as "HasRegistry").
type buildArgs struct {
	argmapper.Struct

	UI terminal.UI
}

func (b *Builder) Build(args buildArgs) (*wpdocker.Image, error) {
	// Pull all the args out to top-level values. This is mostly done
	// cause the struct was added later, but also because these are very common.
	ui := args.UI
	sg := ui.StepGroup()
	defer sg.Wait()
	step := sg.Add("")
	defer func() {
		if step != nil {
			step.Abort()
		}
	}()

	result := &wpdocker.Image{
		Image:    b.config.Image,
		Tag:      b.config.Tag,
		Location: &wpdocker.Image_Docker{Docker: &empty.Empty{}},
	}

	step.Update("Using Docker image in remote registry: %s", result.Name())
	step.Done()

	result.Location = &wpdocker.Image_Registry{Registry: &wpdocker.Image_RegistryLocation{}}
	return result, nil
}
