## kubernetes (platform)

Deploy the application into a Kubernetes cluster using Deployment objects.

### Interface

### Examples

```hcl
use "kubernetes" {
	image_secret = "registry_secret"
	replicas = 3
	probe_path = "/_healthz"
}
```
