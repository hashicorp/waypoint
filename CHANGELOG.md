## unreleased changes

FEATURES:

IMPROVEMENTS:

* install/k8s: improve stability of install process by verifying stateful set, waiting for service endpoints, etc. [GH-435]
* install/k8s: detect Kind and warn about additional requirements [GH-615]
* plugin/pack: can set environment variables [GH-581]
* plugin/docker: ability to target remote docker engine for deploys, automatically pull images [GH-631]

BUG FIXES:

* entrypoint: ensure binary is statically linked on all systems [GH-586]
* plugin/nomad: destroy works [GH-571]
* plugin/aws-ecs: Allow `cpu` parameter for to be optional for EC2 clusters [GH-576]
* plugin/aws-ecs: don't detect inactive cluster as existing [GH-605]
* plugin/google-cloud-run: Don't crash if capacity or autoscaling settings are nil [GH-620]
* install/nomad: if `-nomad-dc` flag is set, `dc1` won't be set [GH-603]

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
