The Packer plugin retrieves the image ID of an image whose metadata is pushed
to an [HCP Packer](https://cloud.hashicorp.com/products/packer) registry. The
image ID is that of the HCP Packer bucket iteration assigned to the configured
channel, with a matching cloud provider and region.

### Components

1. [ConfigSourcer](/waypoint/integrations/hashicorp/packer/latest/components/config-sourcer/packer-config-sourcer)
