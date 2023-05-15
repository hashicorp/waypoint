// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	// The HCP Client ID to authenticate to HCP
	ClientId string `hcl:"client_id,optional"`

	// The HCP Client Secret to authenticate to HCP
	ClientSecret string `hcl:"client_secret,optional"`

	// The HCP Organization ID to authenticate to in HCP
	OrganizationId string `hcl:"organization_id,attr"`

	// The HCP Project ID within the organization to authenticate to in HCP
	ProjectId string `hcl:"project_id,attr"`
}

type reqConfig struct {
	// The name of the HCP Packer registry bucket from which to source an image
	Bucket string `hcl:"bucket,attr"`

	// The name of the HCP Packer registry bucket channel from which to source an image
	Channel string `hcl:"channel,attr"`

	// The region of the machine image to be pulled
	Region string `hcl:"region,attr"`

	// The cloud provider of the machine image to be pulled
	Cloud string `hcl:"cloud,attr"`
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
			continue
		}
		channelParams.BucketSlug = packerConfig.Bucket
		channelParams.Slug = packerConfig.Channel

		// An HCP Packer channel points to a single iteration of a bucket.
		channel, err := hcpPackerClient.PackerServiceGetChannel(channelParams, nil)
		if err != nil {
			result.Result = &pb.ConfigSource_Value_Error{
				Error: status.New(codes.Internal, err.Error()).Proto(),
			}
			continue
		}
		log.Debug("retrieved HCP Packer channel", "channel", channel.Payload.Channel.Slug)
		iteration := channel.Payload.Channel.Iteration

		// An iteration can have multiple builds, so we check for the first build
		// with the matching cloud provider and region.
		for _, build := range iteration.Builds {
			if build.CloudProvider == packerConfig.Cloud {
				log.Debug("found build with matching cloud provider",
					"cloud provider", build.CloudProvider,
					"build ID", build.ID)
				for _, image := range build.Images {
					if image.Region == packerConfig.Region {
						log.Debug("found image with matching region",
							"region", image.Region,
							"image ID", image.ID)
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

	doc.Description("Retrieve the image ID of an image whose metadata is pushed " +
		"to an HCP Packer registry. The image ID is that of the HCP Packer bucket " +
		"iteration assigned to the configured channel, with a matching cloud provider " +
		"and region.")

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
