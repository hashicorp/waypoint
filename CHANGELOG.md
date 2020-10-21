## unreleased

FEATURES:

IMPROVEMENTS:

* plugin/ecs: ability to specify `service_port` rather than port 3000 [GH-661]
* plugin/k8s: support for manually specifying the namespace to use [GH-648]
* plugins/nomad: support for setting docker auth credentials [GH-646]

BUG FIXES:

* cli: `server bootstrap` shows an error if a server is running in-memory [GH-651]
* plugin/docker: remove intermediate containers on build [GH-667]
* plugin/docker, plugin/pack: support ssh accessible docker hosts [GH-664]

## 0.1.3 (October 19, 2020)

FEATURES:

IMPROVEMENTS:

* install/k8s: improve stability of install process by verifying stateful set, waiting for service endpoints, etc. [GH-435]
* install/k8s: detect Kind and warn about additional requirements [GH-615]
* plugin/aws-ecs: support for static environment variables [GH-583]
* plugin/aws-ecs: support for ECS secrets [GH-583]
* plugin/aws-ecs: support for configuring sidecar containers [GH-583]
* plugin/pack: can set environment variables [GH-581]
* plugin/docker: ability to target remote docker engine for deploys, automatically pull images [GH-631]
* ui: onboarding flow after redeeming an invite token is enabled and uses public release URLs [GH-635]

BUG FIXES:

* entrypoint: ensure binary is statically linked on all systems [GH-586]
* plugin/nomad: destroy works [GH-571]
* plugin/aws: load `~/.aws/config` if available and use that for auth [GH-621]
* plugin/aws-ecs: Allow `cpu` parameter for to be optional for EC2 clusters [GH-576]
* plugin/aws-ecs: don't detect inactive cluster as existing [GH-605]
* plugin/aws-ecs: fix crash if subnets are specified [GH-636]
* plugin/aws-ecs: delete ECS ALB listener on destroy [GH-607]
* plugin/google-cloud-run: Don't crash if capacity or autoscaling settings are nil [GH-620]
* install/nomad: if `-nomad-dc` flag is set, `dc1` won't be set [GH-603]
* cli: contexts will fall back to not using symlinks if symlinks aren't available [GH-633]

## 0.1.2 (October 16, 2020)

IMPROVEMENTS:

* plugin/docker: can specify alternate Dockerfile path [GH-539]
* plugin/docker-pull: support registry auth [GH-545]

BUG FIXES:

* cli: fix compatibility with Windows versions prior to 1809 [GH-515]
* cli: `waypoint config` no longer crashes [GH-473]
* cli: autocomplete in bash no longer crashes [GH-540]
* cli: fix crash for some invalid configurations [GH-553]
* entrypoint: do not block child process startup on URL service connection [GH-544]
* plugin/aws-ecs: Use an existing cluster when there is only one cluster [GH-514]

## 0.1.1 (October 15, 2020)

Initial Public Release
