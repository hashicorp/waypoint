## unreleased

FEATURES:

* cli: New `waypoint fmt` command will autoformat your `waypoint.hcl` files [#1037]
* server: Add ability for server to poll a project via git and automatically deploy [#1109]
* runner: Start a standalone runner process along side the server to enable CLI-less actions [#1167] [#1171]
* builtin/aws/lambda: Add support for AWS Lambda deploys [#1097]

IMPROVEMENTS:

* plugin/docker: support for building, pulling, and pushing Docker images without a Docker daemon available. [#970]
* plugin/k8s: plugin will attempt in-cluster auth first if no kubeconfig file is specified [#1052] [#1103]

BUG FIXES:

BREAKING CHANGES:

* ui: dropped support for Internet Explorer [#1075]

## 0.2.3 (February 23, 2021)

FEATURES:

IMPROVEMENTS:

* builtin/docker: Introduce resources map for configuring cpu, memory on deploy [#1116]
* internal/server: More descriptive error for unknown application deploys [#973]
* serverinstall/k8s: Include option to define storageClassName on install [#1126]

BUG FIXES:

* builtin/docker: Fix host port mapping defined by service_port [#918]
* builtin/k8s: Surface pod failures on deploy [#1110]
* serverinstall/nomad: Set platform as nomad at end of Install [#1129]
* builtin/aws-ecs: Fix nil check on optional logging block [#1120]

## 0.2.2 (February 17, 2021)

FEATURES:

IMPROVEMENTS:

* builtin/aws/ecs: Add config option for disabling the load balancer [#1082]
* builtin/aws/ecs: Add awslog driver configuration [#1089]
* builtin/docker: Add Binds, Labels and Networks config options for deploy [#1065]
* builtin/k8s: Support multi-port application configs for deploy and release [#1092]
* cli/main: Add -version flag for CLI version [#1049]

BUG FIXES:

* bulitin/aws/ecs: Determine load balancer and target group for pre-existing listeners [#1085]
* builtin/aws/ecs: fix listener deletion on deployment deletion [#1087]
* builtin/k8s: Handle application config sync with K8s and Secrets [#1073]
* cli/hostname: fix panic with no hostname arg specified [#1044]
* core: Fix empty gitreftag response in config [#1047]

## 0.2.1 (February 02, 2021)

FEATURES:

* **Uninstall command for all server platforms**:
Use `server uninstall` to remove the Waypoint server and artifacts from the
specified `-platform` for the active server installation. [#972]
* **Upgrade command for all server platforms**:
Use `server upgrade` to upgrade the Waypoint server for the
specified `-platform` for the active server installation. [#976]

IMPROVEMENTS:

* builtin/k8s: Allow for defined resource limits for pods [#1041]
* cli: `server run` supports specifying a custom TLS certificate [#951]
* cli: more informative error messages on `install` [#1004]
* server: store platform where server is installed to in server config [#1000]
* serverinstall/docker: Start waypoint server container if stopped on install [#1009]
* serverinstall/k8s: Allow using k8s context [#1028]

BUG FIXES:

* builtin/aws/ami: require []string for aws-ami filters to avoid panic [#1010]
* cli: ctrl-c now interrupts server connection attempts [#989]
* entrypoint: log disconnect messages will now only be emitted at the ERROR level if reconnection fails [#930]
* server: don't block startup on URL service being unavailable [#950]
* server: `UpsertProject` will not delete all application metadata [#1027]
* server: increase timeout for hostname registration [#1040]
* plugin/google-cloud-run: fix error on deploys about missing type [#955]

## 0.2.0 (December 10, 2020)

FEATURES:

* **Application config syncing with Kubernetes (ConfigMaps), Vault, Consul, and AWS SSM**;
Automatically sync environment variable values with remote sources and restart your
application when those values change. [#810]
* **Access to Artifact, Deploy Metadata**: `registry` and `deploy` configuration can use 
`artifact.*` variable syntax to access metadata from the results of those stages. 
The `release` configuration can use `artifact.*` and `deploy.*` to access metadata.
For example: `image = artifact.image` for Docker-based builds. [#757]
* **`template` Functions**: `templatefile`, `templatedir`, and `templatestring` functions
allow you to template files, directories, and strings with the variables and functions
available to your Waypoint configuration. 
* **`path` Variables**: you can now use `path.project`, `path.app`, and `path.pwd` as
variables in your Waypoint file to specify paths as relative to the project (waypoint.hcl
file), app, or pwd of the CLI invocation.
* **Server snapshot/restore**: you can now use the CLI or API to take and restore online
snapshots. Restoring snapshot data requires a server restart but the restore operation 
can be staged online. [#870]

IMPROVEMENTS:

* cli/logs: entrypoint logs can now be seen alongside app logs and are colored differently [#855]
* contrib/serverinstall: Automate setup of kind+k8s with metallb [#845]
* core: application config changes (i.e. `waypoint config set`) will now restart running applications [#791]
* core: add more descriptive text to include app name in `waypoint destroy` [#807]
* core: add better error messaging when prefix is missing from the `-raw` flag in `waypoint config set` [#815]
* core: align -raw flag to behave like -json flag with waypoint config set [#828]
* core: `waypoint.hcl` can be named `waypoint.hcl.json` and use JSON syntax [#867]
* install: Update flags used on server install per-platform [#882]
* install/k8s: support for OpenShift [#715]
* internal/server: Block on instance deployment becoming available [#881]
* plugin/aws-ecr: environment variables to be used instead of 'region' property for aws-ecr registry [#841]
* plugin/google-cloud-run: allow images from Google Artifact Registry [#804]
* plugin/google-cloud-run: added service account name field [#850]
* server: APIs for Waypoint database snapshot/restore [#723]
* website: many minor improvements were made in our plugin documentation section for this release

BUG FIXES:

* core: force killed `waypoint exec` sessions won't leave the remote process running [#827]
* core: waypoint exec with no TTY won't hang open until a ctrl-c [#830]
* core: prevent waypoint server from re-registering guest account horizon URL service [#922]
* cli: server config-set doesn't require a Waypoint configuration file. [#819]
* cli/token: fix issue where tokens could be cut off on narrow terminals [#885]
* plugin/aws-ecs: task_role_name is optional [#824]


## 0.1.5 (November 09, 2020)

FEATURES:

IMPROVEMENTS:

* plugin/google-cloud-run: set a default releaser so you don't need a `release` block [#756]

BUG FIXES:

* plugin/ecs: do not assign public IP on EC2 cluster [#758]
* plugin/google-cloud-run: less strict image validation to allow projects with slashes [#760]
* plugin/k8s: default releaser should create service with correct namespace [#759]
* entrypoint: be careful to not spawn multiple url agents [#752]
* cli: return error for ErrSentinel types to signal exit codes [#768]

## 0.1.4 (October 26, 2020)

FEATURES:

IMPROVEMENTS:

* cli/config: you can pipe in a KEY=VALUE line-delimited file via stdin to `config set` [#674]
* install/docker: pull server image if it doesn't exist in local Docker daemon [#700]
* install/nomad: added `-nomad-policy-override` flag to allow sentinel policy override on Nomad enterprise [#671]
* plugin/ecs: ability to specify `service_port` rather than port 3000 [#661]
* plugin/k8s: support for manually specifying the namespace to use [#648]
* plugins/nomad: support for setting docker auth credentials [#646]

BUG FIXES:

* cli: `server bootstrap` shows an error if a server is running in-memory [#651]
* cli: server connection flags take precedence over context data [#703]
* cli: plugins loaded in pwd work properly and don't give `PATH` errors [#699]
* cli: support plugins in `$HOME/.config/waypoint/plugins` even if XDG path doesn't match [#707]
* plugin/docker: remove intermediate containers on build [#667]
* plugin/docker, plugin/pack: support ssh accessible docker hosts [#664]

## 0.1.3 (October 19, 2020)

FEATURES:

IMPROVEMENTS:

* install/k8s: improve stability of install process by verifying stateful set, waiting for service endpoints, etc. [#435]
* install/k8s: detect Kind and warn about additional requirements [#615]
* plugin/aws-ecs: support for static environment variables [#583]
* plugin/aws-ecs: support for ECS secrets [#583]
* plugin/aws-ecs: support for configuring sidecar containers [#583]
* plugin/pack: can set environment variables [#581]
* plugin/docker: ability to target remote docker engine for deploys, automatically pull images [#631]
* ui: onboarding flow after redeeming an invite token is enabled and uses public release URLs [#635]

BUG FIXES:

* entrypoint: ensure binary is statically linked on all systems [#586]
* plugin/nomad: destroy works [#571]
* plugin/aws: load `~/.aws/config` if available and use that for auth [#621]
* plugin/aws-ecs: Allow `cpu` parameter for to be optional for EC2 clusters [#576]
* plugin/aws-ecs: don't detect inactive cluster as existing [#605]
* plugin/aws-ecs: fix crash if subnets are specified [#636]
* plugin/aws-ecs: delete ECS ALB listener on destroy [#607]
* plugin/google-cloud-run: Don't crash if capacity or autoscaling settings are nil [#620]
* install/nomad: if `-nomad-dc` flag is set, `dc1` won't be set [#603]
* cli: contexts will fall back to not using symlinks if symlinks aren't available [#633]

## 0.1.2 (October 16, 2020)

IMPROVEMENTS:

* plugin/docker: can specify alternate Dockerfile path [#539]
* plugin/docker-pull: support registry auth [#545]

BUG FIXES:

* cli: fix compatibility with Windows versions prior to 1809 [#515]
* cli: `waypoint config` no longer crashes [#473]
* cli: autocomplete in bash no longer crashes [#540]
* cli: fix crash for some invalid configurations [#553]
* entrypoint: do not block child process startup on URL service connection [#544]
* plugin/aws-ecs: Use an existing cluster when there is only one cluster [#514]

## 0.1.1 (October 15, 2020)

Initial Public Release
