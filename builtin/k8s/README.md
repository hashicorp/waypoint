The Kubernetes plugin can deploy a Docker image of an application to Kubernetes,
expose the Deployment with a Kubernetes Service, and source configuration from
a Kubernetes Secret or ConfigMap. It also launches on-demand runners to do
operations remotely.

### Components

1. [Platform](/waypoint/integrations/hashicorp/kubernetes/latest/components/platform/kubernetes-platform)
2. [ReleaseManager](/waypoint/integrations/hashicorp/kubernetes/latest/components/release-manager/kubernetes-release-manager)
3. [ConfigSourcer](/waypoint/integrations/hashicorp/kubernetes/latest/components/config-sourcer/kubernetes-config-sourcer)
4. [TaskLauncher](/waypoint/integrations/hashicorp/kubernetes/latest/components/task/kubernetes-task)

### Related Plugins

1. [Docker](/waypoint/integrations/hashicorp/docker)

### Resources

#### Platform

1. Kubernetes Deployment
2. Kubernetes Autoscaler

#### Release Manager

1. Kubernetes Service
2. Kubernetes Ingress
