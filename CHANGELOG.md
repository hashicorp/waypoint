## 0.2.0 (December 10, 2020)

FEATURES:

* **Application config syncing with Kubernetes (ConfigMaps), Vault, Consul, and AWS SSM**;
Automatically sync environment variable values with remote sources and restart your
application when those values change. [GH-810]
* **Access to Artifact, Deploy Metadata**: `registry` and `deploy` configuration can use 
`artifact.*` variable syntax to access metadata from the results of those stages. 
The `release` configuration can use `artifact.*` and `deploy.*` to access metadata.
For example: `image = artifact.image` for Docker-based builds. [GH-757]
* **`template` Functions**: `templatefile`, `templatedir`, and `templatestring` functions
allow you to template files, directories, and strings with the variables and functions
available to your Waypoint configuration. 
* **`path` Variables**: you can now use `path.project`, `path.app`, and `path.pwd` as
variables in your Waypoint file to specify paths as relative to the project (waypoint.hcl
file), app, or pwd of the CLI invocation.
* **Server snapshot/restore**: you can now use the CLI or API to take and restore online
snapshots. Restoring snapshot data requires a server restart but the restore operation 
can be staged online. [GH-870]

IMPROVEMENTS:

* cli/logs: entrypoint logs can now be seen alongside app logs and are colored differently [GH-855]
* contrib/serverinstall: Automate setup of kind+k8s with metallb [GH-845]
* core: application config changes (i.e. `waypoint config set`) will now restart running applications [GH-791]
* core: add more descriptive text to include app name in `waypoint destroy` [GH-807]
* core: add better error messaging when prefix is missing from the `-raw` flag in `waypoint config set` [GH-815]
* core: align -raw flag to behave like -json flag with waypoint config set [GH-828]
* core: `waypoint.hcl` can be named `waypoint.hcl.json` and use JSON syntax [GH-867]
* install: Update flags used on server install per-platform [GH-882]
* install/k8s: support for OpenShift [GH-715]
* internal/server: Block on instance deployment becoming available [GH-881]
* plugin/aws-ecr: environment variables to be used instead of 'region' property for aws-ecr registry [GH-841]
* plugin/google-cloud-run: allow images from Google Artifact Registry [GH-804]
* plugin/google-cloud-run: added service account name field [GH-850]
* server: APIs for Waypoint database snapshot/restore [GH-723]
* website: many minor improvements were made in our plugin documentation section for this release

BUG FIXES:

* core: force killed `waypoint exec` sessions won't leave the remote process running [GH-827]
* core: waypoint exec with no TTY won't hang open until a ctrl-c [GH-830]
* cli: server config-set doesn't require a Waypoint configuration file. [GH-819]
* cli/token: fix issue where tokens could be cut off on narrow terminals [GH-885]
* plugin/aws-ecs: task_role_name is optional [GH-824]


## 0.1.5 (November 09, 2020)

FEATURES:

IMPROVEMENTS:

* plugin/google-cloud-run: set a default releaser so you don't need a `release` block [GH-756]

BUG FIXES:

* plugin/ecs: do not assign public IP on EC2 cluster [GH-758]
* plugin/google-cloud-run: less strict image validation to allow projects with slashes [GH-760]
* plugin/k8s: default releaser should create service with correct namespace [GH-759]
* entrypoint: be careful to not spawn multiple url agents [GH-752]
* cli: return error for ErrSentinel types to signal exit codes [GH-768]

## 0.1.4 (October 26, 2020)

FEATURES:

IMPROVEMENTS:

* cli/config: you can pipe in a KEY=VALUE line-delimited file via stdin to `config set` [GH-674]
* install/docker: pull server image if it doesn't exist in local Docker daemon [GH-700]
* install/nomad: added `-nomad-policy-override` flag to allow sentinel policy override on Nomad enterprise [GH-671]
* plugin/ecs: ability to specify `service_port` rather than port 3000 [GH-661]
* plugin/k8s: support for manually specifying the namespace to use [GH-648]
* plugins/nomad: support for setting docker auth credentials [GH-646]

BUG FIXES:

* cli: `server bootstrap` shows an error if a server is running in-memory [GH-651]
* cli: server connection flags take precedence over context data [GH-703]
* cli: plugins loaded in pwd work properly and don't give `PATH` errors [GH-699]
* cli: support plugins in `$HOME/.config/waypoint/plugins` even if XDG path doesn't match [GH-707]
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
