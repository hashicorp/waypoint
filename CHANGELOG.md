## unreleased



## 0.3.0 (April 08, 2021)

BREAKING CHANGES:

* ui: dropped support for Internet Explorer [[GH-1075](https://github.com/hashicorp/waypoint/issues/1075)]

FEATURES:

* **GitOps: Poll for changes and automatically run `waypoint up`**. Waypoint
  can now trigger a full build, deploy, release cycle on changes detected in Git. [[GH-1109](https://github.com/hashicorp/waypoint/issues/1109)]
* **Runners: Run Waypoint operations remotely**. Runners are standalone processes that
  run operations such as builds, deploys, etc. remotely. [[GH-1167](https://github.com/hashicorp/waypoint/issues/1167)] [[GH-1171](https://github.com/hashicorp/waypoint/issues/1171)]
* **AWS Lambda**: Add support for building and deploying AWS Lambda workloads [[GH-1097](https://github.com/hashicorp/waypoint/issues/1097)]
* **Dockerless image builds**: Waypoint can now build, tag, pull, and push
  Docker images in unprivileged environments without a Docker daemon. [[GH-970](https://github.com/hashicorp/waypoint/issues/970)]
* cli: New `waypoint fmt` command will autoformat your `waypoint.hcl` files [[GH-1037](https://github.com/hashicorp/waypoint/issues/1037)]
* config: `timestamp` function allows you to avail the current date and time in your Waypoint configuration. [[GH-1255](https://github.com/hashicorp/waypoint/issues/1255)]
* ui: Add ability to create a project from the browser UI [[GH-1220](https://github.com/hashicorp/waypoint/issues/1220)]
* ui: Add ability to configure a project's git settings from the browser UI [[GH-1057](https://github.com/hashicorp/waypoint/issues/1057)]
* ui: Add ability to input a waypoint.hcl configuration from the browser UI [[GH-1253](https://github.com/hashicorp/waypoint/issues/1253)]

IMPROVEMENTS:

* cli: Require confirmation before destroying all resources [[GH-1232](https://github.com/hashicorp/waypoint/issues/1232)]
* cli: Can specify the number of deployments to prune for `up` and `release`. [[GH-1230](https://github.com/hashicorp/waypoint/issues/1230)]
* cli: Support and render new documentation subfields [[GH-1213](https://github.com/hashicorp/waypoint/issues/1213)]
* plugin/docker-pull: doesn't require Docker if no registry is configured and entrypoint injection is disabled [[GH-1198](https://github.com/hashicorp/waypoint/issues/1198)]
* plugin/k8s: Add new probe configuration options [[GH-1246](https://github.com/hashicorp/waypoint/issues/1246)]
* plugin/k8s: plugin will attempt in-cluster auth first if no kubeconfig file is specified [GH-1052] [GH-1103]
* server: Prune old deployments and jobs from server memory. This limits the number
of deployments and jobs to 10,000. The data for the old entries is still stored on disk
but it is not indexed in memory, to allow data recovery should it be needed. [[GH-1193](https://github.com/hashicorp/waypoint/issues/1193)]

BUG FIXES:

* core: default releasers initialize properly when they use HCL variables [[GH-1254](https://github.com/hashicorp/waypoint/issues/1254)]
* cli: require at least one argument [GH-1188]
* plugin/aws/alb: clamp alb name per aws limits [[GH-1225](https://github.com/hashicorp/waypoint/issues/1225)]
* ui: output failed build errors in logs [[GH-1280](https://github.com/hashicorp/waypoint/issues/1280)]

## 0.2.4 (March 18, 2021)

FEATURES:

IMPROVEMENTS:

* builtin/k8s: Include user defined labels on deploy pod [[GH-1146](https://github.com/hashicorp/waypoint/issues/1146)]
* core: Update `-from-project` in `waypoint init` to handle local projects [[GH-722](https://github.com/hashicorp/waypoint/issues/722)]
* entrypoint: dump a memory profile to `/tmp` on `SIGUSR1` [[GH-1194](https://github.com/hashicorp/waypoint/issues/1194)]

BUG FIXES:

* entrypoint: fix URL service memory leak [[GH-1200](https://github.com/hashicorp/waypoint/issues/1200)]
* builtin/k8s/release: Allow target_port to be int or string [[GH-1154](https://github.com/hashicorp/waypoint/issues/1154)]
* builtin/k8s/release: Include name field for service port release [[GH-1184](https://github.com/hashicorp/waypoint/issues/1184)]
* builtin/docker: Revert #918, ensure HostPort is randomly assigned between container deploys [[GH-1189](https://github.com/hashicorp/waypoint/issues/1189)]

## 0.2.3 (February 23, 2021)

FEATURES:

IMPROVEMENTS:

* builtin/docker: Introduce resources map for configuring cpu, memory on deploy [[GH-1116](https://github.com/hashicorp/waypoint/issues/1116)]
* internal/server: More descriptive error for unknown application deploys [[GH-973](https://github.com/hashicorp/waypoint/issues/973)]
* serverinstall/k8s: Include option to define storageClassName on install [[GH-1126](https://github.com/hashicorp/waypoint/issues/1126)]

BUG FIXES:

* builtin/docker: Fix host port mapping defined by service_port [[GH-918](https://github.com/hashicorp/waypoint/issues/918)]
* builtin/k8s: Surface pod failures on deploy [[GH-1110](https://github.com/hashicorp/waypoint/issues/1110)]
* serverinstall/nomad: Set platform as nomad at end of Install [[GH-1129](https://github.com/hashicorp/waypoint/issues/1129)]
* builtin/aws-ecs: Fix nil check on optional logging block [[GH-1120](https://github.com/hashicorp/waypoint/issues/1120)]

## 0.2.2 (February 17, 2021)

FEATURES:

IMPROVEMENTS:

* builtin/aws/ecs: Add config option for disabling the load balancer [[GH-1082](https://github.com/hashicorp/waypoint/issues/1082)]
* builtin/aws/ecs: Add awslog driver configuration [[GH-1089](https://github.com/hashicorp/waypoint/issues/1089)]
* builtin/docker: Add Binds, Labels and Networks config options for deploy [[GH-1065](https://github.com/hashicorp/waypoint/issues/1065)]
* builtin/k8s: Support multi-port application configs for deploy and release [[GH-1092](https://github.com/hashicorp/waypoint/issues/1092)]
* cli/main: Add -version flag for CLI version [[GH-1049](https://github.com/hashicorp/waypoint/issues/1049)]

BUG FIXES:

* bulitin/aws/ecs: Determine load balancer and target group for pre-existing listeners [[GH-1085](https://github.com/hashicorp/waypoint/issues/1085)]
* builtin/aws/ecs: fix listener deletion on deployment deletion [[GH-1087](https://github.com/hashicorp/waypoint/issues/1087)]
* builtin/k8s: Handle application config sync with K8s and Secrets [[GH-1073](https://github.com/hashicorp/waypoint/issues/1073)]
* cli/hostname: fix panic with no hostname arg specified [[GH-1044](https://github.com/hashicorp/waypoint/issues/1044)]
* core: Fix empty gitreftag response in config [[GH-1047](https://github.com/hashicorp/waypoint/issues/1047)]

## 0.2.1 (February 02, 2021)

FEATURES:

* **Uninstall command for all server platforms**:
Use `server uninstall` to remove the Waypoint server and artifacts from the
specified `-platform` for the active server installation. [[GH-972](https://github.com/hashicorp/waypoint/issues/972)]
* **Upgrade command for all server platforms**:
Use `server upgrade` to upgrade the Waypoint server for the
specified `-platform` for the active server installation. [[GH-976](https://github.com/hashicorp/waypoint/issues/976)]

IMPROVEMENTS:

* builtin/k8s: Allow for defined resource limits for pods [[GH-1041](https://github.com/hashicorp/waypoint/issues/1041)]
* cli: `server run` supports specifying a custom TLS certificate [[GH-951](https://github.com/hashicorp/waypoint/issues/951)]
* cli: more informative error messages on `install` [[GH-1004](https://github.com/hashicorp/waypoint/issues/1004)]
* server: store platform where server is installed to in server config [[GH-1000](https://github.com/hashicorp/waypoint/issues/1000)]
* serverinstall/docker: Start waypoint server container if stopped on install [[GH-1009](https://github.com/hashicorp/waypoint/issues/1009)]
* serverinstall/k8s: Allow using k8s context [[GH-1028](https://github.com/hashicorp/waypoint/issues/1028)]

BUG FIXES:

* builtin/aws/ami: require []string for aws-ami filters to avoid panic [[GH-1010](https://github.com/hashicorp/waypoint/issues/1010)]
* cli: ctrl-c now interrupts server connection attempts [[GH-989](https://github.com/hashicorp/waypoint/issues/989)]
* entrypoint: log disconnect messages will now only be emitted at the ERROR level if reconnection fails [[GH-930](https://github.com/hashicorp/waypoint/issues/930)]
* server: don't block startup on URL service being unavailable [[GH-950](https://github.com/hashicorp/waypoint/issues/950)]
* server: `UpsertProject` will not delete all application metadata [[GH-1027](https://github.com/hashicorp/waypoint/issues/1027)]
* server: increase timeout for hostname registration [[GH-1040](https://github.com/hashicorp/waypoint/issues/1040)]
* plugin/google-cloud-run: fix error on deploys about missing type [[GH-955](https://github.com/hashicorp/waypoint/issues/955)]

## 0.2.0 (December 10, 2020)

FEATURES:

* **Application config syncing with Kubernetes (ConfigMaps), Vault, Consul, and AWS SSM**;
Automatically sync environment variable values with remote sources and restart your
application when those values change. [[GH-810](https://github.com/hashicorp/waypoint/issues/810)]
* **Access to Artifact, Deploy Metadata**: `registry` and `deploy` configuration can use
`artifact.*` variable syntax to access metadata from the results of those stages.
The `release` configuration can use `artifact.*` and `deploy.*` to access metadata.
For example: `image = artifact.image` for Docker-based builds. [[GH-757](https://github.com/hashicorp/waypoint/issues/757)]
* **`template` Functions**: `templatefile`, `templatedir`, and `templatestring` functions
allow you to template files, directories, and strings with the variables and functions
available to your Waypoint configuration.
* **`path` Variables**: you can now use `path.project`, `path.app`, and `path.pwd` as
variables in your Waypoint file to specify paths as relative to the project (waypoint.hcl
file), app, or pwd of the CLI invocation.
* **Server snapshot/restore**: you can now use the CLI or API to take and restore online
snapshots. Restoring snapshot data requires a server restart but the restore operation
can be staged online. [[GH-870](https://github.com/hashicorp/waypoint/issues/870)]

IMPROVEMENTS:

* cli/logs: entrypoint logs can now be seen alongside app logs and are colored differently [[GH-855](https://github.com/hashicorp/waypoint/issues/855)]
* contrib/serverinstall: Automate setup of kind+k8s with metallb [[GH-845](https://github.com/hashicorp/waypoint/issues/845)]
* core: application config changes (i.e. `waypoint config set`) will now restart running applications [[GH-791](https://github.com/hashicorp/waypoint/issues/791)]
* core: add more descriptive text to include app name in `waypoint destroy` [[GH-807](https://github.com/hashicorp/waypoint/issues/807)]
* core: add better error messaging when prefix is missing from the `-raw` flag in `waypoint config set` [[GH-815](https://github.com/hashicorp/waypoint/issues/815)]
* core: align -raw flag to behave like -json flag with waypoint config set [[GH-828](https://github.com/hashicorp/waypoint/issues/828)]
* core: `waypoint.hcl` can be named `waypoint.hcl.json` and use JSON syntax [[GH-867](https://github.com/hashicorp/waypoint/issues/867)]
* install: Update flags used on server install per-platform [[GH-882](https://github.com/hashicorp/waypoint/issues/882)]
* install/k8s: support for OpenShift [[GH-715](https://github.com/hashicorp/waypoint/issues/715)]
* internal/server: Block on instance deployment becoming available [[GH-881](https://github.com/hashicorp/waypoint/issues/881)]
* plugin/aws-ecr: environment variables to be used instead of 'region' property for aws-ecr registry [[GH-841](https://github.com/hashicorp/waypoint/issues/841)]
* plugin/google-cloud-run: allow images from Google Artifact Registry [[GH-804](https://github.com/hashicorp/waypoint/issues/804)]
* plugin/google-cloud-run: added service account name field [[GH-850](https://github.com/hashicorp/waypoint/issues/850)]
* server: APIs for Waypoint database snapshot/restore [[GH-723](https://github.com/hashicorp/waypoint/issues/723)]
* website: many minor improvements were made in our plugin documentation section for this release

BUG FIXES:

* core: force killed `waypoint exec` sessions won't leave the remote process running [[GH-827](https://github.com/hashicorp/waypoint/issues/827)]
* core: waypoint exec with no TTY won't hang open until a ctrl-c [[GH-830](https://github.com/hashicorp/waypoint/issues/830)]
* core: prevent waypoint server from re-registering guest account horizon URL service [[GH-922](https://github.com/hashicorp/waypoint/issues/922)]
* cli: server config-set doesn't require a Waypoint configuration file. [[GH-819](https://github.com/hashicorp/waypoint/issues/819)]
* cli/token: fix issue where tokens could be cut off on narrow terminals [[GH-885](https://github.com/hashicorp/waypoint/issues/885)]
* plugin/aws-ecs: task_role_name is optional [[GH-824](https://github.com/hashicorp/waypoint/issues/824)]


## 0.1.5 (November 09, 2020)

FEATURES:

IMPROVEMENTS:

* plugin/google-cloud-run: set a default releaser so you don't need a `release` block [[GH-756](https://github.com/hashicorp/waypoint/issues/756)]

BUG FIXES:

* plugin/ecs: do not assign public IP on EC2 cluster [[GH-758](https://github.com/hashicorp/waypoint/issues/758)]
* plugin/google-cloud-run: less strict image validation to allow projects with slashes [[GH-760](https://github.com/hashicorp/waypoint/issues/760)]
* plugin/k8s: default releaser should create service with correct namespace [[GH-759](https://github.com/hashicorp/waypoint/issues/759)]
* entrypoint: be careful to not spawn multiple url agents [[GH-752](https://github.com/hashicorp/waypoint/issues/752)]
* cli: return error for ErrSentinel types to signal exit codes [[GH-768](https://github.com/hashicorp/waypoint/issues/768)]

## 0.1.4 (October 26, 2020)

FEATURES:

IMPROVEMENTS:

* cli/config: you can pipe in a KEY=VALUE line-delimited file via stdin to `config set` [[GH-674](https://github.com/hashicorp/waypoint/issues/674)]
* install/docker: pull server image if it doesn't exist in local Docker daemon [[GH-700](https://github.com/hashicorp/waypoint/issues/700)]
* install/nomad: added `-nomad-policy-override` flag to allow sentinel policy override on Nomad enterprise [[GH-671](https://github.com/hashicorp/waypoint/issues/671)]
* plugin/ecs: ability to specify `service_port` rather than port 3000 [[GH-661](https://github.com/hashicorp/waypoint/issues/661)]
* plugin/k8s: support for manually specifying the namespace to use [[GH-648](https://github.com/hashicorp/waypoint/issues/648)]
* plugins/nomad: support for setting docker auth credentials [[GH-646](https://github.com/hashicorp/waypoint/issues/646)]

BUG FIXES:

* cli: `server bootstrap` shows an error if a server is running in-memory [[GH-651](https://github.com/hashicorp/waypoint/issues/651)]
* cli: server connection flags take precedence over context data [[GH-703](https://github.com/hashicorp/waypoint/issues/703)]
* cli: plugins loaded in pwd work properly and don't give `PATH` errors [[GH-699](https://github.com/hashicorp/waypoint/issues/699)]
* cli: support plugins in `$HOME/.config/waypoint/plugins` even if XDG path doesn't match [[GH-707](https://github.com/hashicorp/waypoint/issues/707)]
* plugin/docker: remove intermediate containers on build [[GH-667](https://github.com/hashicorp/waypoint/issues/667)]
* plugin/docker, plugin/pack: support ssh accessible docker hosts [[GH-664](https://github.com/hashicorp/waypoint/issues/664)]

## 0.1.3 (October 19, 2020)

FEATURES:

IMPROVEMENTS:

* install/k8s: improve stability of install process by verifying stateful set, waiting for service endpoints, etc. [[GH-435](https://github.com/hashicorp/waypoint/issues/435)]
* install/k8s: detect Kind and warn about additional requirements [[GH-615](https://github.com/hashicorp/waypoint/issues/615)]
* plugin/aws-ecs: support for static environment variables [[GH-583](https://github.com/hashicorp/waypoint/issues/583)]
* plugin/aws-ecs: support for ECS secrets [[GH-583](https://github.com/hashicorp/waypoint/issues/583)]
* plugin/aws-ecs: support for configuring sidecar containers [[GH-583](https://github.com/hashicorp/waypoint/issues/583)]
* plugin/pack: can set environment variables [[GH-581](https://github.com/hashicorp/waypoint/issues/581)]
* plugin/docker: ability to target remote docker engine for deploys, automatically pull images [[GH-631](https://github.com/hashicorp/waypoint/issues/631)]
* ui: onboarding flow after redeeming an invite token is enabled and uses public release URLs [[GH-635](https://github.com/hashicorp/waypoint/issues/635)]

BUG FIXES:

* entrypoint: ensure binary is statically linked on all systems [[GH-586](https://github.com/hashicorp/waypoint/issues/586)]
* plugin/nomad: destroy works [[GH-571](https://github.com/hashicorp/waypoint/issues/571)]
* plugin/aws: load `~/.aws/config` if available and use that for auth [[GH-621](https://github.com/hashicorp/waypoint/issues/621)]
* plugin/aws-ecs: Allow `cpu` parameter for to be optional for EC2 clusters [[GH-576](https://github.com/hashicorp/waypoint/issues/576)]
* plugin/aws-ecs: don't detect inactive cluster as existing [[GH-605](https://github.com/hashicorp/waypoint/issues/605)]
* plugin/aws-ecs: fix crash if subnets are specified [[GH-636](https://github.com/hashicorp/waypoint/issues/636)]
* plugin/aws-ecs: delete ECS ALB listener on destroy [[GH-607](https://github.com/hashicorp/waypoint/issues/607)]
* plugin/google-cloud-run: Don't crash if capacity or autoscaling settings are nil [[GH-620](https://github.com/hashicorp/waypoint/issues/620)]
* install/nomad: if `-nomad-dc` flag is set, `dc1` won't be set [[GH-603](https://github.com/hashicorp/waypoint/issues/603)]
* cli: contexts will fall back to not using symlinks if symlinks aren't available [[GH-633](https://github.com/hashicorp/waypoint/issues/633)]

## 0.1.2 (October 16, 2020)

IMPROVEMENTS:

* plugin/docker: can specify alternate Dockerfile path [[GH-539](https://github.com/hashicorp/waypoint/issues/539)]
* plugin/docker-pull: support registry auth [[GH-545](https://github.com/hashicorp/waypoint/issues/545)]

BUG FIXES:

* cli: fix compatibility with Windows versions prior to 1809 [[GH-515](https://github.com/hashicorp/waypoint/issues/515)]
* cli: `waypoint config` no longer crashes [[GH-473](https://github.com/hashicorp/waypoint/issues/473)]
* cli: autocomplete in bash no longer crashes [[GH-540](https://github.com/hashicorp/waypoint/issues/540)]
* cli: fix crash for some invalid configurations [[GH-553](https://github.com/hashicorp/waypoint/issues/553)]
* entrypoint: do not block child process startup on URL service connection [[GH-544](https://github.com/hashicorp/waypoint/issues/544)]
* plugin/aws-ecs: Use an existing cluster when there is only one cluster [[GH-514](https://github.com/hashicorp/waypoint/issues/514)]

## 0.1.1 (October 15, 2020)

Initial Public Release
