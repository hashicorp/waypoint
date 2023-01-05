<!-- This file was generated via `make gen/integrations-hcl` -->
Deploy to Kubernetes from a Helm chart. The Helm chart can be a local path
or a chart in a repository.

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
These must be passed in using Helm values (i.e. the chart must make
environment variables configurable).

This is documented in more detail with a full example in the
[Kubernetes Helm Deployment documentation](/waypoint/docs/platforms/kubernetes/helm-deploy).

#### URL Service

If you want your workload to be accessible by the
[Waypoint URL service](/waypoint/docs/url), you must set the PORT environment variable
within the pod with your web service and also be using the Waypoint
entrypoint (documented in the previous section).

The PORT environment variable should be the port that your web service
is listening on that the URL service will connect to. See one of the examples
below for more details.

### Interface

### Examples

```hcl
// Configuring an image to match the build. This assumes the chart
// has a value named "deployment.image".
deploy {
  use "helm" {
    chart = "${path.app}/chart"

    set {
      name  = "deployment.image"
      value = artifact.name
    }
  }
}
```

