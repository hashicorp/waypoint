<!-- This file was generated via `make gen/integrations-hcl` -->
Retrieve the image ID of an image whose metadata is pushed to an HCP Packer registry. The image ID is that of the HCP Packer bucket iteration assigned to the configured channel, with a matching cloud provider and region.

### Examples

```hcl
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
```

