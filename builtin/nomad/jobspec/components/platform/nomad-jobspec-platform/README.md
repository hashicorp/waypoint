<!-- This file was generated via `make gen/integrations-hcl` -->
Deploy to a Nomad cluster from a pre-existing Nomad job specification file.

This plugin lets you use any pre-existing Nomad job specification file to
deploy to Nomad. This deployment is able to support all the features of Waypoint.
You may use Waypoint's [templating features](/waypoint/docs/waypoint-hcl/functions/template)
to template the Nomad jobspec with information such as the artifact from
a previous build step, entrypoint environment variables, etc.

### Artifact Access

You may use Waypoint's [templating features](/waypoint/docs/waypoint-hcl/functions/template)
to access information such as the artifact from the build or push stages.
An example below shows this by using `templatefile` mixed with
variables such as `artifact.image` to dynamically configure the
Docker image within the Nomad job specification.

-> **Note:** If using [Nomad interpolation](/nomad/docs/runtime/interpolation) in your jobspec file,
and the `templatefile` function in your waypoint.hcl file, any interpolated values must be escaped with a second 
`$`. For example: `$${meta.metadata}` instead of `${meta.metadata}`.

### Entrypoint Functionality

Waypoint [entrypoint functionality](/waypoint/docs/entrypoint#functionality) such
as logs, exec, app configuration, and more require two properties to be true:

1. The running image must already have the Waypoint entrypoint installed
  and configured as the entrypoint. This should happen in the build stage.

2. Proper environment variables must be set so the entrypoint knows how
  to communicate to the Waypoint server. **This step happens in this
  deployment stage.**

**Step 2 does not happen automatically.** You must manually set the entrypoint
environment variables using the [templating feature](/waypoint/docs/waypoint-hcl/functions/template).
One of the examples below shows the entrypoint environment variables being
injected.

-> **Note:** The Waypoint entrypoint and the [Nomad entrypoint functionality](/nomad/docs/drivers/docker#entrypoint) 
cannot be used simultaneously. In order to use the features of the Waypoint entrypoint, the Nomad entrypoint must not be used in your jobspec.

### URL Service

If you want your workload to be accessible by the
[Waypoint URL service](/waypoint/docs/url), you must set the PORT environment variable
within your job and be using the Waypoint entrypoint (documented in the
previous section).

The PORT environment variable should be the port that your web service
is listening on that the URL service will connect to. See one of the examples
below for more details.

### Interface

- Input: **docker.Image**
- Output: **jobspec.Deployment**

### Examples

```hcl
// The waypoint.hcl file
deploy {
  use "nomad-jobspec" {
    // Templated to perhaps bring in the artifact from a previous
    // build/registry, entrypoint env vars, etc.
    jobspec = templatefile("${path.app}/app.nomad.tpl")
  }
}

// The app.nomad.tpl file
job "web" {
  datacenters = ["dc1"]

  group "app" {
    task "app" {
      driver = "docker"

      config {
        image = "${artifact.image}:${artifact.tag}"
      }

      env {
        %{ for k,v in entrypoint.env ~}
        ${k} = "${v}"
        %{ endfor ~}

        // Ensure we set PORT for the URL service. This is only necessary
        // if we want the URL service to function.
        PORT = 3000
      }
    }
  }
}
```

