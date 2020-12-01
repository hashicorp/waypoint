package ami

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
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
	// AWS region to operate in
	Region string `hcl:"region"`

	// Only look for images from the given owners
	Owners []string `hcl:"owners,optional"`

	// The name of the image to search for, supports wildcards
	Name string `hcl:"name,optional"`

	// Specific filters to pass to the DescribeImages filter set
	Filters map[string]interface{} `hcl:"filters,optional"`
}

func (b *Builder) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&BuilderConfig{}), docs.FromFunc(b.BuildFunc()))
	if err != nil {
		return nil, err
	}

	doc.Description("Search for and return an existing AMI")

	doc.Output("ami.Image")

	doc.SetField(
		"region",
		"the AWS region to search in",
	)

	doc.SetField(
		"owners",
		"the set of AMI owners to restrict the search to",
	)

	doc.SetField(
		"name",
		"the name of the AMI to search for, supports wildcards",
	)

	doc.SetField(
		"filters",
		"DescribeImage specific filters to search with",
		docs.Summary(
			"the filters are always name => [value], but this api supports",
			"the ability to pass a single value as a convience. Non string",
			"values will be converted to strings",
		),
	)

	return doc, nil
}

// Config implements Configurable
func (b *Builder) Config() (interface{}, error) {
	return &b.config, nil
}

// Build
func (b *Builder) Build(
	ctx context.Context,
	ui terminal.UI,
	src *component.Source,
) (*Image, error) {
	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: b.config.Region,
	})
	if err != nil {
		return nil, err
	}
	e := ec2.New(sess)

	var (
		owners  []*string
		filters []*ec2.Filter
	)

	for _, o := range b.config.Owners {
		owners = append(owners, aws.String(o))
	}

	if b.config.Name != "" {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("name"),
			Values: []*string{aws.String(b.config.Name)},
		})
	}

	for k, v := range b.config.Filters {
		var values []*string

		switch sv := v.(type) {
		case string:
			values = append(values, aws.String(sv))
		case []interface{}:
			for _, iv := range sv {
				values = append(values, aws.String(fmt.Sprintf("%s", iv)))
			}
		default:
			values = append(values, aws.String(fmt.Sprintf("%s", sv)))
		}

		filters = append(filters, &ec2.Filter{
			Name:   aws.String(k),
			Values: values,
		})
	}

	out, err := e.DescribeImages(&ec2.DescribeImagesInput{
		Filters: filters,
		Owners:  owners,
	})

	if err != nil {
		return nil, err
	}

	if len(out.Images) == 0 {
		return nil, fmt.Errorf("no images found")
	}

	result := &Image{
		Image: *out.Images[0].ImageId,
	}

	ui.Output("Resolved AMI: %s", result.Image, terminal.WithSuccessStyle())

	return result, nil
}
