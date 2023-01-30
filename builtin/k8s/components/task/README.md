<!-- This file was generated via `make gen/integrations-hcl` -->
Launch a Kubernetes pod for on-demand tasks from the Waypoint server.

This will use the standard Kubernetes environment variables to source
authentication information for Kubernetes. If this is running within Kubernetes
itself (typical for a Kubernetes-based installation), it will use the pod's
service account unless other auth is explicitly given. This allows the task
launcher to work by default.

### Examples

```hcl
task {
	use "kubernetes" {}
}
```

