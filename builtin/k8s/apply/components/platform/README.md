<!-- This file was generated via `make gen/integrations-hcl` -->
Deploy Kubernetes resources directly from a single file or a directory of YAML
or JSON files.

This plugin lets you use any pre-existing set of Kubernetes resource files
to deploy to Kubernetes. This plugin supports all the features of Waypoint.
You may use Waypoint's [templating features](/waypoint/docs/waypoint-hcl/functions/template)
to template the resources with information such as the artifact from
a previous build step, entrypoint environment variables, etc.

### Requirements

This plugin requires "kubectl" to be installed since this plugin works by
subprocessing to "kubectl apply". Other Waypoint Kubernetes plugins sometimes
use the API directly but this plugin requires "kubectl".

"kubectl" must also be configured to access your Kubernetes cluster. You
may specify an alternate kubeconfig file using the "kubeconfig" configuration
parameter. If this isn't specified, the default kubectl lookup paths will be
used.

### Artifact Access

You may use Waypoint's [templating features](/waypoint/docs/waypoint-hcl/functions/template)
to access information such as the artifact from the build or push stages.
An example below shows this by using `templatedir` mixed with
variables such as `artifact.image` to dynamically configure the
Docker image within a Kubernetes Deployment.

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

### URL Service

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
// The waypoint.hcl file
deploy {
  use "kubernetes-apply" {
    // Templated to perhaps bring in the artifact from a previous
    // build/registry, entrypoint env vars, etc.
    path        = templatedir("${path.app}/k8s")
    prune_label = "app=myapp"
	prune_whitelist = [
		"apps/v1/Deployment",
		"apps/v1/ReplicaSet"
  	]
  }
}

// ./k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
  labels:
    app: myapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
      - name: myapp
        image: ${artifact.image}:${artifact.tag}
        env:
          %{ for k,v in entrypoint.env ~}
          - name: ${k}
            value: "${v}"
          %{ endfor ~}

          # Ensure we set PORT for the URL service. This is only necessary
          # if we want the URL service to function.
          - name: PORT
            value: "3000"
```

