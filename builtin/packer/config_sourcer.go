package packer

import (
	"context"

	"github.com/hashicorp/go-hclog"
	packer "github.com/hashicorp/hcp-sdk-go/clients/cloud-packer-service/stable/2021-04-30/client/packer_service"
	hcpconfig "github.com/hashicorp/hcp-sdk-go/config"
	"github.com/hashicorp/hcp-sdk-go/httpclient"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	pb "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
)

type ConfigSourcer struct {
	config sourceConfig
	client packer.Client
}

type sourceConfig struct {
	ClientId       string `hcl:"client_id,optional"`
	ClientSecret   string `hcl:"client_secret,optional"`
	OrganizationId string `hcl:"organization_id,attr"`
	ProjectId      string `hcl:"project_id,attr"`
}

type reqConfig struct {
	Bucket  string `hcl:"bucket,attr"`
	Channel string `hcl:"channel,attr"`
	Region  string `hcl:"region,attr"`
	Cloud   string `hcl:"cloud,attr"`
}

// Config implements component.Configurable
func (cs *ConfigSourcer) Config() (interface{}, error) {
	return &cs.config, nil
}

// ReadFunc implements component.ConfigSourcer
func (cs *ConfigSourcer) ReadFunc() interface{} {
	return cs.read
}

// StopFunc implements component.ConfigSourcer
func (cs *ConfigSourcer) StopFunc() interface{} {
	return cs.stop
}

func (cs *ConfigSourcer) read(
	ctx context.Context,
	log hclog.Logger,
	reqs []*component.ConfigRequest,
) ([]*pb.ConfigSource_Value, error) {

	// If the user has explicitly set the client ID and secret for the config
	// sourcer, we use that. Otherwise, we use environment variables.
	opts := hcpconfig.FromEnv()
	if cs.config.ClientId != "" && cs.config.ClientSecret != "" {
		opts = hcpconfig.WithClientCredentials(cs.config.ClientId, cs.config.ClientSecret)
	}
	hcpConfig, err := hcpconfig.NewHCPConfig(opts)
	if err != nil {
		return nil, err
	}

	hcpClient, err := httpclient.New(httpclient.Config{
		HCPConfig: hcpConfig,
	})
	if err != nil {
		return nil, err
	}

	hcpPackerClient := packer.New(hcpClient, nil)
	channelParams := packer.NewPackerServiceGetChannelParams()
	channelParams.LocationOrganizationID = cs.config.OrganizationId
	channelParams.LocationProjectID = cs.config.ProjectId

	var results []*pb.ConfigSource_Value
	for _, req := range reqs {
		result := &pb.ConfigSource_Value{Name: req.Name}
		results = append(results, result)

		var packerConfig reqConfig
		// We serialize the config sourcer settings to the reqConfig struct.
		if err = mapstructure.WeakDecode(req.Config, &packerConfig); err != nil {
			result.Result = &pb.ConfigSource_Value_Error{
				Error: status.New(codes.Aborted, err.Error()).Proto(),
			}
			return nil, err
		}
		channelParams.BucketSlug = packerConfig.Bucket
		channelParams.Slug = packerConfig.Channel

		// An HCP Packer channel points to a single iteration of a bucket.
		channel, err := hcpPackerClient.PackerServiceGetChannel(channelParams, nil)
		if err != nil {
			return nil, err
		}
		log.Debug("Retrieved HCP Packer channel.")
		iteration := channel.Payload.Channel.Iteration

		// An iteration can have multiple builds, so we check for the first build
		// with the matching cloud provider and region.
		for _, build := range iteration.Builds {
			if build.CloudProvider == packerConfig.Cloud {
				log.Debug("Found build with matching cloud provider.")
				for _, image := range build.Images {
					if image.Region == packerConfig.Region {
						log.Debug("Found image with matching region.")
						result.Result = &pb.ConfigSource_Value_Value{
							// The ImageID is the Cloud Image ID or URL string
							// identifying this image for the builder that built it,
							// so this is returned to Waypoint.
							Value: image.ImageID,
						}
					}
				}
			}
		}
	}

	return results, nil
}

func (cs *ConfigSourcer) stop() error {
	return nil
}

func (cs *ConfigSourcer) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(
		docs.FromConfig(&sourceConfig{}),
		docs.RequestFromStruct(&reqConfig{}),
	)
	if err != nil {
		return nil, err
	}

	doc.Description("Read machine image information from HCP Packer.")

	doc.Example(`
// The waypoint.hcl file
project = "example-reactjs-project"

variable "image" {
  default = dynamic("packer", {
    bucket          = "nginx"
    channel         = "base"
    region          = "docker"
    cloud_provider  = "docker"
  }
  type = string
  description = "The name of the base image to use for building app Docker images."
}

app "example-reactjs" {
  build {
    use "docker" {
      dockerfile = templatefile("${path.app}"/Dockerfile, {
        base_image = var.image
      }
    }

  deploy {
    use "docker" {}
  }
}


# Multi-stage Dockerfile example
FROM node:19.2-alpine as build
WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH
COPY package.json ./
COPY package-lock.json ./
RUN npm ci --silent
RUN npm install react-scripts@3.4.1 -g --silent
COPY . ./
RUN npm run build

# ${base_image} below is the Docker repository and tag, templated to the Dockerfile
FROM ${base_image}
COPY nginx/default.conf /etc/nginx/conf.d/
COPY --from=build /app/build /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
`)

	doc.SetRequestField(
		"bucket",
		"The name of the HCP Packer bucket from which to source an image.",
	)

	doc.SetRequestField(
		"channel",
		"The name of the HCP Packer channel from which to source the latest image.",
	)

	doc.SetRequestField(
		"region",
		"The region set for the machine image's cloud provider.",
	)

	doc.SetRequestField(
		"cloud",
		"The cloud provider of the machine image to source",
	)

	doc.SetField(
		"organization_id",
		"The HCP organization ID.",
	)

	doc.SetField(
		"project_id",
		"The HCP Project ID.",
	)

	doc.SetField(
		"client_id",
		"The OAuth2 Client ID for HCP API operations.",
		docs.EnvVar("HCP_CLIENT_ID"),
	)

	doc.SetField(
		"client_secret",
		"The OAuth2 Client Secret for HCP API operations.",
		docs.EnvVar("HCP_CLIENT_SECRET"),
	)

	return doc, nil
}
