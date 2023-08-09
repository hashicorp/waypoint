## unreleased

FEATURES:

IMPROVEMENTS:

BUG FIXES:

## 0.11.4 (August 9, 2023)

IMPROVEMENTS:

* plugin/aws/ecr-pull: Support entrypoint injection in ecr-pull builder [[GH-4847](https://github.com/hashicorp/waypoint/issues/4847)]

BUG FIXES:

* cli: Fix possible issues with deleted applications or projects failing to render in output [[GH-4867](https://github.com/hashicorp/waypoint/issues/4867)]
* runneruninstall/aws-ecs: Fix installing runners in new AWS accounts by fixing an inline policy syntax error. [[GH-4873](https://github.com/hashicorp/waypoint/issues/4873)]

## 0.11.3 (July 18, 2023)

IMPROVEMENTS:

* cli: Add config source's plugin type, scope, project, app, and workspace to
  output of `waypoint config source-get` when getting a specific config source. [[GH-4822](https://github.com/hashicorp/waypoint/issues/4822)]
* cli: Add option `all` to flag `-scope` on `waypoint config source-get` command
  to output all config sources. [[GH-4822](https://github.com/hashicorp/waypoint/issues/4822)]

BUG FIXES:

* runnerinstall/aws-ecs: Fix ODR policy for AWS ECS runners to enable adding tags
  to an ALB. [[GH-4818](https://github.com/hashicorp/waypoint/issues/4818)]
* runnerinstall/ecs: Add IAM permission required for project destruction [[GH-4840](https://github.com/hashicorp/waypoint/issues/4840)]
* runneruninstall/aws-ecs: Fix panic when uninstalling ECS runner after failing to find EFS [[GH-4829](https://github.com/hashicorp/waypoint/issues/4829)]

## 0.11.2 (June 15, 2023)

IMPROVEMENTS:

* cli,server: Introduce explicit `delete` endpoint for CLI and Server for Config
  and Config Sourcers. [[GH-4754](https://github.com/hashicorp/waypoint/issues/4754)]
* cli: Use -w flag for workspace scoping on `config set` and `config delete`,
  instead of `workspace-scope`. [[GH-4770](https://github.com/hashicorp/waypoint/issues/4770)]
* plugin/aws-ecs: Add config options for the target group protocol and protocol
  version. [[GH-4742](https://github.com/hashicorp/waypoint/issues/4742)]
* runnerinstall/nomad: Add CLI flags for setting custom CPU and memory resources. [[GH-4798](https://github.com/hashicorp/waypoint/issues/4798)]
* serverinstall/nomad: Add config flag `-nomad-host-network` for specifying the
  host network of the Waypoint server Nomad job's gRPC and HTTP (UI) ports. [[GH-4804](https://github.com/hashicorp/waypoint/issues/4804)]

BUG FIXES:

* auth: Prevent a runner token from generating a new token for a different runner [[GH-4707](https://github.com/hashicorp/waypoint/issues/4707)]
* builtin: de-dupe various hcl annotation keys [[GH-4701](https://github.com/hashicorp/waypoint/issues/4701)]
* cli: Honor runner install -platform arg [[GH-4699](https://github.com/hashicorp/waypoint/issues/4699)]
* config-sources: Return correct workspace-scoped config sources at the global
  scope, when a workspace is specified. [[GH-4774](https://github.com/hashicorp/waypoint/issues/4774)]
* config: Remove extra eval context append for parsing configs which caused a slowdown during pipeline config parsing. [[GH-4744](https://github.com/hashicorp/waypoint/issues/4744)]
* plugin/aws-ecs: Destroy the ALB only if it is managed by Waypoint. [[GH-4742](https://github.com/hashicorp/waypoint/issues/4742)]
* plugin/aws-ecs: Fix failure when destroying the target group during a release
  destroy operation. [[GH-4742](https://github.com/hashicorp/waypoint/issues/4742)]
* plugin/aws-ecs: Fix panic when settings `grpc_code` or `http_code` for a health
  check. [[GH-4742](https://github.com/hashicorp/waypoint/issues/4742)]
* plugin/aws-ecs: Set the protocol of a health check correctly. [[GH-4742](https://github.com/hashicorp/waypoint/issues/4742)]
* plugin/azure-aci: Update plugin to attempt CLI auth if environment auth fails. [[GH-4763](https://github.com/hashicorp/waypoint/issues/4763)]
* plugin/ecs: Make `alb.load_balancer_arn` optional [[GH-4729](https://github.com/hashicorp/waypoint/issues/4729)]
* plugin/ecs: Set health check timeout and interval values to compatible default
  values. [[GH-4767](https://github.com/hashicorp/waypoint/issues/4767)]
* runnerinstall/aws-ecs: Add missing permission to on-demand runner IAM policy. [[GH-4742](https://github.com/hashicorp/waypoint/issues/4742)]
* runneruninstall/aws-ecs: Fix deletion of file system for AWS ECS runner. [[GH-4792](https://github.com/hashicorp/waypoint/issues/4792)]

## 0.11.1 (May 11, 2023)

IMPROVEMENTS:

* cli: Add a `-verbose` flag to `waypoint job list` to improve relevant columns shown to user
  at a glance. [[GH-4531](https://github.com/hashicorp/waypoint/issues/4531)]
* cli: Include job `QueueTime` in output for `waypoint job list` and `waypoint job inspect`. [[GH-4531](https://github.com/hashicorp/waypoint/issues/4531)]
* cli: Introduce `waypoint runner profile edit` to edit a runners plugin config
  directly in your configured terminal editor [[GH-4594](https://github.com/hashicorp/waypoint/issues/4594)]
* cli: Update the `waypoint runner profile set` command to accept an argument
  for setting the name. This also removes the behavior where if no name was given,
  it would generate a random one. [[GH-4527](https://github.com/hashicorp/waypoint/issues/4527)]
* cli: new flags for `waypoint install` on Nomad:
  -nomad-service-address and -nomad-network-mode [[GH-4619](https://github.com/hashicorp/waypoint/issues/4619)]
* plugin/ecs: Enable custom health checks for ECS plugin. [[GH-4473](https://github.com/hashicorp/waypoint/issues/4473)]
* plugin/ecs: Update ECS releaser to verify deployment health before releasing. [[GH-4520](https://github.com/hashicorp/waypoint/issues/4520)]

BUG FIXES:

* builtin/consul: Fix request logger to properly log configured data center [[GH-4670](https://github.com/hashicorp/waypoint/issues/4670)]
* cli: Avoid panic in empty slice for runner installs platform var. [[GH-4672](https://github.com/hashicorp/waypoint/issues/4672)]
* cli: Fix load path for custom Waypoint plugins [[GH-4623](https://github.com/hashicorp/waypoint/issues/4623)]
* core: Ensure project and workspaces cannot be created with malformed names [[GH-4588](https://github.com/hashicorp/waypoint/issues/4588)]
* internal: Improve git URL string trimming when determining remote URLs [[GH-4675](https://github.com/hashicorp/waypoint/issues/4675)]
* plugin/ecs: Update ECS destroyer to wait for there to be zero listeners for the
  target group before destroying the target group. [[GH-4497](https://github.com/hashicorp/waypoint/issues/4497)]
* trigger: Ensure trigger Name is only alpha-numeric [[GH-4660](https://github.com/hashicorp/waypoint/issues/4660)]
* ui: Only show health-check “Re-run” button if project has a data source. [[GH-4553](https://github.com/hashicorp/waypoint/issues/4553)]

## 0.11.0 (February 16, 2023)

FEATURES:

* server: Add pagination protobuffs and stubs for pagination in ListProjects [[GH-4203](https://github.com/hashicorp/waypoint/issues/4203)]

IMPROVEMENTS:

* cli: `waypoint job cancel` now outputs additional insights
when trying to cancel a running job. [[GH-4294](https://github.com/hashicorp/waypoint/issues/4294)]
* cli: add socket-path flag to runner install [[GH-4246](https://github.com/hashicorp/waypoint/issues/4246)]
* core: `waypoint job list` will now retrieve paginated list of jobs to avoid
grpc data limits per request [[GH-4271](https://github.com/hashicorp/waypoint/issues/4271)]
* core: improve runner job stream error logging [[GH-3872](https://github.com/hashicorp/waypoint/issues/3872)]
* install/nomad: Allow mount options to be specified when provisioning a volume with CSI plugins [[GH-4387](https://github.com/hashicorp/waypoint/issues/4387)]
* plugin/aws: Add CORS configuration to lambda-function-url releaser [[GH-4418](https://github.com/hashicorp/waypoint/issues/4418)]
* plugin/tfc: Allow non-string tfc outputs to be used as waypoint.hcl dynamic default variables [[GH-4357](https://github.com/hashicorp/waypoint/issues/4357)]
* plugin/tfc: Allow reading all outputs from a tfc workspace with a single variable stanza [[GH-4357](https://github.com/hashicorp/waypoint/issues/4357)]
* plugins/k8s: Add `prune_whitelist` option to only prune specific resources [[GH-4345](https://github.com/hashicorp/waypoint/issues/4345)]
* plugins/k8s: Add `security_context` to the TaskLauncherConfig (on-demand runner configuration) [[GH-4346](https://github.com/hashicorp/waypoint/issues/4346)]
* server: Add UI_GetDeployment convenience endpoint [[GH-3856](https://github.com/hashicorp/waypoint/issues/3856)]
* server: Enable gRPC-Gateway on the http port of Waypoint server to add in an HTTP API for interfacing [[GH-4379](https://github.com/hashicorp/waypoint/issues/4379)]

BUG FIXES:

* aws/lambda: fix issue where deployment configuration was not injected in to
Lambda function environments, preventing waypoint-entrypoint from authenticating
with the Waypoint server [[GH-4328](https://github.com/hashicorp/waypoint/issues/4328)]
* cli/snapshot: Fix server snapshot when a config source is set. [[GH-4523](https://github.com/hashicorp/waypoint/issues/4523)]
* cli: Fix panic in `waypoint pipeline list` and `waypoint pipeline inspect` where
a pipeline run was given with no jobs. [[GH-4424](https://github.com/hashicorp/waypoint/issues/4424)]
* cli: Show full command flags when displaying help text for `waypoint runner inspect`. [[GH-4435](https://github.com/hashicorp/waypoint/issues/4435)]
* install/nomad: Fix connectivity to Waypoint server from the CLI at the end of
the Nomad server install. [[GH-4363](https://github.com/hashicorp/waypoint/issues/4363)]
* plugin/aws-ecs: Fix bringing your own alb to ecs deployments. [[GH-4457](https://github.com/hashicorp/waypoint/issues/4457)]
* plugin/k8s-apply: Update the `prune_whitelist` param to match the updated parameter
in kubectl apply, `prune_allowlist`. It also ensures this param in the plugin
is optional and not a hard requirement to use the k8s apply plugin. [[GH-4517](https://github.com/hashicorp/waypoint/issues/4517)]
* ui: fix safari bug with xterm/webgl rendering [[GH-4054](https://github.com/hashicorp/waypoint/issues/4054)]
* upgrade: Fixes a bug where pre-v0.10.4 config sources could not be updated or
deleted. [[GH-4382](https://github.com/hashicorp/waypoint/issues/4382)]


## 0.10.5 (December 15, 2022)

SECURITY:
* Waypoint now uses Go 1.19.4 to address security vulnerability (CVE-2022-41717) See the Go announcement for more details.

IMPROVEMENTS:

* cli: Respect `-remote-source` overrides for submitted job template when running
  `waypoint pipeline run`. [[GH-4319](https://github.com/hashicorp/waypoint/issues/4319)]
* config: Remove the multi-app deprecation warning.
  Please see https://discuss.hashicorp.com/t/deprecating-projects-or-how-i-learned-to-love-apps/40888/12 for more information. [[GH-4265](https://github.com/hashicorp/waypoint/issues/4265)]

BUG FIXES:

* plugin/ecs: `runner install` now creates aws policies to facilitate remotely running StopTask and WatchTask jobs [[GH-4296](https://github.com/hashicorp/waypoint/issues/4296)]

## 0.10.4 (December 08, 2022)

FEATURES:

- plugin/ecs: Accept ALB security group IDs. [[GH-4230](https://github.com/hashicorp/waypoint/issues/4230)]
- plugin/packer: A Packer config sourcer plugin to source machine image IDs from
  an HCP Packer channel. [[GH-4251](https://github.com/hashicorp/waypoint/issues/4251)]

IMPROVEMENTS:

- cli/runnerinstall: Check if runner is registered to the server before
  attempting to forget it. [[GH-3944](https://github.com/hashicorp/waypoint/issues/3944)]
- cli/runnerinstall: Delete EFS file system during ECS runner uninstall. [[GH-3944](https://github.com/hashicorp/waypoint/issues/3944)]
- cli: `project destroy` requires the `-project` or `-p` flag regardless of where it's run. [[GH-4212](https://github.com/hashicorp/waypoint/issues/4212)]
- cli: Pipeline run now shows the number of steps successfully executed for a failed run. [[GH-4268](https://github.com/hashicorp/waypoint/issues/4268)]

BUG FIXES:

- cli/context: Fix possible error when listing contexts if a non-Waypoint context file exists in the context directory. [[GH-4257](https://github.com/hashicorp/waypoint/issues/4257)]
- cli: Ensure a deploy and release URL has a scheme included if not set. [[GH-4208](https://github.com/hashicorp/waypoint/issues/4208)]
- cli: `project destroy` now successfully destroys a project created in the UI without a remote source or local waypoint.hcl file. [[GH-4212](https://github.com/hashicorp/waypoint/issues/4212)]
- plugin/nomad: Update Nomad task launcher plugin to use `entrypoint` config - fixes
  pipeline exec steps run in Nomad. [[GH-4185](https://github.com/hashicorp/waypoint/issues/4185)]
- plugin/vault: Fix usage of dynamic secrets from Vault for dynamic Waypoint app
  config. [[GH-3988](https://github.com/hashicorp/waypoint/issues/3988)]

## 0.10.3 (November 03, 2022)

FEATURES:

* plugin/consul: Consul key-value data config sourcer plugin [[GH-4045](https://github.com/hashicorp/waypoint/issues/4045)]

IMPROVEMENTS:

* cli/config-sync: Add operations flags to `config sync` command. [[GH-4143](https://github.com/hashicorp/waypoint/issues/4143)]
* cli/fmt: Add a `-check` flag that will determine if the `waypoint.hcl` is already
properly formatted, similar to `terraform fmt -check`. [[GH-4020](https://github.com/hashicorp/waypoint/issues/4020)]
* cli/pipeline_run: Show app deployment and release URLs if exist from running
pipeline. Also show input variables used. [[GH-4096](https://github.com/hashicorp/waypoint/issues/4096)]
* cli: Add prune flags to `waypoint deploy` for configuring the automatic release. [[GH-4114](https://github.com/hashicorp/waypoint/issues/4114)]
* cli: Introduce new CLI flag `-reattach` for `waypoint pipeline run` which will stream
an existing pipeline run either by the latest known run or a specific sequence id. [[GH-4042](https://github.com/hashicorp/waypoint/issues/4042)]
* cli: Only echo file name when config file is formatted with `waypoint fmt`. [[GH-4111](https://github.com/hashicorp/waypoint/issues/4111)]
* cli: Update `waypoint runner profile inspect` to show default runner profile
if no name argument supplied. [[GH-4078](https://github.com/hashicorp/waypoint/issues/4078)]
* core: Auto run a status report after a deployment or release operation rather
than only if `waypoint deploy` or `waypoint release` CLI is run. [[GH-4099](https://github.com/hashicorp/waypoint/issues/4099)]
* core: Combine git clone messages from job stream into a single message [[GH-4115](https://github.com/hashicorp/waypoint/issues/4115)]
* pipelines: Add ability to evaluate input variables in pipelines stanzas. [[GH-4132](https://github.com/hashicorp/waypoint/issues/4132)]
* ui/input-variables: Adds the ability to set an input variable as sensitive and hides its value from the list and form [[GH-4139](https://github.com/hashicorp/waypoint/issues/4139)]

BUG FIXES:

* cli/runner-profile-set: Fix panic when setting runner profile environment variables [[GH-3995](https://github.com/hashicorp/waypoint/issues/3995)]
* cli/upgrade: Update the OCI URL for the bootstrap runner profile during `server upgrade` [[GH-4175](https://github.com/hashicorp/waypoint/issues/4175)]
* cli: Fix bug where input variables were not included on pipeline run jobs. [[GH-4137](https://github.com/hashicorp/waypoint/issues/4137)]
* cli: Fix panic in `waypoint runner profile set` when no flags are specified. [[GH-4013](https://github.com/hashicorp/waypoint/issues/4013)]
* cli: Fix panic in cli for `waypoint task cancel` if attempting to cancel by run
job id with no argument. [[GH-4019](https://github.com/hashicorp/waypoint/issues/4019)]
* cli: Only show "CompleteTime" on `waypoint pipeline list` if the job has a valid
complete time. [[GH-4113](https://github.com/hashicorp/waypoint/issues/4113)]
* cli: Remove automatic uppercasing of ids, so that future runner profiles will match [[GH-4063](https://github.com/hashicorp/waypoint/issues/4063)]
* cli: Respect the -workspace flag when requesting a logstream for a deployment
by workspace [[GH-4009](https://github.com/hashicorp/waypoint/issues/4009)]
* core: Fix panic if no Use stanza found for given workspace scope on a build,
deploy, release, or registry stanza. [[GH-4112](https://github.com/hashicorp/waypoint/issues/4112)]
* core: fix panic when null value is set on a string variable [[GH-4067](https://github.com/hashicorp/waypoint/issues/4067)]
* install/nomad: Update installation with Nomad to use CSI parameters. [[GH-4157](https://github.com/hashicorp/waypoint/issues/4157)]
* pipelines: Properly mark a pipeline run as complete [[GH-4053](https://github.com/hashicorp/waypoint/issues/4053)]
* plugin/docker: fix issue with authenticating with registries when using
docker-pull [[GH-4121](https://github.com/hashicorp/waypoint/issues/4121)]

## 0.10.2 (October 03, 2022)

BREAKING CHANGES:

* plugin/helm: Add support for create_namespace and skip_crds (default is now false) [[GH-3950](https://github.com/hashicorp/waypoint/issues/3950)]

IMPROVEMENTS:

* cli/runnerinstall: Update `runner install` to set the new profile as the default
if none exists [[GH-3922](https://github.com/hashicorp/waypoint/issues/3922)]
* cli: A new `-runner-target-any` flag has been added to `runner profile set` to allow users to specify targeting any runner. [[GH-3854](https://github.com/hashicorp/waypoint/issues/3854)]
* cli: Setting configurations in a runner profile no longer resets unspecified configuration [[GH-3854](https://github.com/hashicorp/waypoint/issues/3854)]
* plugin/ecs: Implement WatchTask plugin for AWS ECS task launcher. Store ECS on-demand runner logs in the job system. [[GH-3918](https://github.com/hashicorp/waypoint/issues/3918)]
* plugin/docker: Add new `docker-ref` plugin for a build noop for referencing a Docker Image [[GH-3912](https://github.com/hashicorp/waypoint/issues/3912)]

BUG FIXES:

* cli/install: Set image pull policy configuration on Helm installation of Waypoint server and runner [[GH-3948](https://github.com/hashicorp/waypoint/issues/3948)]
* cli: Fix panic when setting project datasource to local [[GH-3972](https://github.com/hashicorp/waypoint/issues/3972)]
* core: Connected entrypoints for deleted projects now error out properly. [[GH-3949](https://github.com/hashicorp/waypoint/issues/3949)]
* core: Fix out of order job ids for `waypoint pipeline run` job stream CLI. [[GH-3946](https://github.com/hashicorp/waypoint/issues/3946)]
* plugin/nomad-jobspec: Fix deployment of periodic and system Nomad jobs. [[GH-3963](https://github.com/hashicorp/waypoint/issues/3963)]
* ui: Fix crash when encountering resources without IDs [[GH-3929](https://github.com/hashicorp/waypoint/issues/3929)]

## 0.10.1 (September 22, 2022)

IMPROVEMENTS:

* core: Update TaskCancel to also cancel WatchTask job, if it was launched. [[GH-3893](https://github.com/hashicorp/waypoint/issues/3893)]

BUG FIXES:

* plugin/nomad-jobspec: Update Nomad jobspec status check to not report partial
health for deployments with canaries [[GH-3883](https://github.com/hashicorp/waypoint/issues/3883)]
* plugin/nomad: Update Nomad task launcher plugin to respect namespace and region
configs in runner profile [[GH-3883](https://github.com/hashicorp/waypoint/issues/3883)]
* plugin: Fix panic in non-ECS plugins when destroying resources [[GH-3896](https://github.com/hashicorp/waypoint/issues/3896)]
* runner-install/k8s: Use the service account created by helm in the runner profile, rather than using the default service account. [[GH-3894](https://github.com/hashicorp/waypoint/issues/3894)]
* runner-install/kubernetes: Fix the static runner image used in the `waypoint runner install` command [[GH-3890](https://github.com/hashicorp/waypoint/issues/3890)]

## 0.10.0 (September 13, 2022)

FEATURES:

* **core: Custom Pipelines** as a Tech Preview gives users the ability to define custom pipelines to run various Waypoint actions such as a build, deploy, release, up and more for deploying applications. Waypoint will monitor and log each pipeline run, associated jobs, and show the results of executing that pipeline. [[GH-3777](https://github.com/hashicorp/waypoint/issues/3777)]
* **core: Project Destroy** Introduce a new CLI command `project destroy` to delete projects in Waypoint and destroy their associated resources [[GH-3626](https://github.com/hashicorp/waypoint/issues/3626)]
* CLI: New waypoint.hcl interactive generator, accessed with `waypoint init` when no waypoint.hcl exists in the current project [[GH-3704](https://github.com/hashicorp/waypoint/issues/3704)]

IMPROVEMENTS:

* CLI: Nomad CSI volumes names can be specified during installation and upgrades for both Waypoint runners and Waypoint server installations [[GH-3546](https://github.com/hashicorp/waypoint/issues/3546)]
* cli/runnerinstall: The runner profile created by `runner install` sets target labels
instead of a target runner ID on the runner profile, if the user supplied label flags [[GH-3755](https://github.com/hashicorp/waypoint/issues/3755)]
* cli: Add option to `waypoint logs` command to get a specific deployment's logs [[GH-3656](https://github.com/hashicorp/waypoint/issues/3656)]
* cli: Fix incorrect description for `destroy -h` command. [[GH-3580](https://github.com/hashicorp/waypoint/issues/3580)]
* cli: Implement `waypoint job get-stream` to allow users to attach to running job
streams and receieve output, or get the output from an existing job stream that
already finished. [[GH-3410](https://github.com/hashicorp/waypoint/issues/3410)]
* cli: Internal-only labels on runners are hidden from CLI output for `waypoint runner
list` and `waypoint runner inspect` [[GH-3746](https://github.com/hashicorp/waypoint/issues/3746)]
* install/runner: Allow additional arguments for `waypoint runner agent` command
to be supplied to `waypoint runner install` CLI [[GH-3746](https://github.com/hashicorp/waypoint/issues/3746)]
* internal/runnerinstall: Fix order of error checking and make error message more specific [[GH-3772](https://github.com/hashicorp/waypoint/issues/3772)]
* plugin/k8s: Add ephemeral-storage resource limits to on-demand runners through runner profiles. [[GH-3676](https://github.com/hashicorp/waypoint/issues/3676)]
* plugin/nomad: Implement WatchTask plugin for Nomad task launcher. Store Nomad on-demand runner logs in the job system. [[GH-3797](https://github.com/hashicorp/waypoint/issues/3797)]
* plugin/nomad: nomad-jobspec deployments will no longer utilize the
nomad-jobspec-canary releaser by default. [[GH-3359](https://github.com/hashicorp/waypoint/issues/3359)]
* serverinstall/ecs: Improve server upgrade process by stopping task and reducing
target group drain time [[GH-3564](https://github.com/hashicorp/waypoint/issues/3564)]
* ui: Fix header spacing on build/release detail pages [[GH-3614](https://github.com/hashicorp/waypoint/issues/3614)]
* ui: Make it clearer that remote runners are required for GitOps [[GH-3615](https://github.com/hashicorp/waypoint/issues/3615)]
* ui: Update alert copy to reflect change that we are no longer removing project definitively in 0.10.0 [[GH-3604](https://github.com/hashicorp/waypoint/issues/3604)]

BUG FIXES:

* cli/install/k8s: Fixes the k8s installer race condition with the bootstrap token [[GH-3744](https://github.com/hashicorp/waypoint/issues/3744)]
* cli/runnerinstall: Fix output error message for missing arguments on Nomad runner installation [[GH-3761](https://github.com/hashicorp/waypoint/issues/3761)]
* cli/runnerinstall: The runner profile created by `runner install` no longer sets
the profile as the default, and appends the runner ID to the name of the profile
for uniqueness [[GH-3755](https://github.com/hashicorp/waypoint/issues/3755)]
* cli/serverinstall/nomad: Add service discovery provider configuration
to server install for Nomad. [[GH-3500](https://github.com/hashicorp/waypoint/issues/3500)]
* cli: Do not set runner profile defaultness to false if default flag not specified. [[GH-3702](https://github.com/hashicorp/waypoint/issues/3702)]
* cli: Fix issue where the CLI would exit with no error or action taken if a local
  waypoint.hcl file was invalid [[GH-3657](https://github.com/hashicorp/waypoint/issues/3657)]
* cli: Fix order of CLI outputs during install [[GH-3729](https://github.com/hashicorp/waypoint/issues/3729)]
* cli: Only show runner profile default deletion hint if more than 0 default
profiles are detected on upgrades. [[GH-3688](https://github.com/hashicorp/waypoint/issues/3688)]
* cli: Set plugin configuration on runner profile created during `runner install` [[GH-3699](https://github.com/hashicorp/waypoint/issues/3699)]
* cli: Set the subnet and security group ID configs for the ECS task launcher plugin
during an ECS runner install [[GH-3701](https://github.com/hashicorp/waypoint/issues/3701)]
* cli: `deployment destroy` will now attempt to destroy any known resources for failed
deployments. Previously, `destroy` would skip unsuccessful deployments. [[GH-3602](https://github.com/hashicorp/waypoint/issues/3602)]
* cli: `runner profile delete` deletes by name rather than the user-invisible ID. [[GH-3803](https://github.com/hashicorp/waypoint/issues/3803)]
* cli: `runner profile set` does not create entries with duplicate names. [[GH-3803](https://github.com/hashicorp/waypoint/issues/3803)]
* core/runner: Default to odr oci image url for runner profiles created via runner install [[GH-3800](https://github.com/hashicorp/waypoint/issues/3800)]
* core/runner: Server no longer panics when a runner stopped after it is forgotten. [[GH-3756](https://github.com/hashicorp/waypoint/issues/3756)]
* core: prune releases which released deployments that are being pruned during
a release [[GH-3730](https://github.com/hashicorp/waypoint/issues/3730)]
* internal/pkg/gitdirty: Fix an issue detecting dirty local source code when using git@ remotes. [[GH-3636](https://github.com/hashicorp/waypoint/issues/3636)]
* plugin/nomad-jobspec-canary: Retrieve the status of the correct Nomad job for the
job resource [[GH-3766](https://github.com/hashicorp/waypoint/issues/3766)]
* ui: Add label to “follow logs” button [[GH-3732](https://github.com/hashicorp/waypoint/issues/3732)]
* ui: Make UI more resilient to invalid state JSON [[GH-3786](https://github.com/hashicorp/waypoint/issues/3786)]
* upgrade: Nomad server upgrade upgrade now detects runners with the name
`waypoint-runner` or `waypoint-static-runner` [[GH-3804](https://github.com/hashicorp/waypoint/issues/3804)]
* upgrade: Update ECS server upgrade to respect -ecs-server-image flag when the
existing server image is hashicorp/waypoint:latest [[GH-3820](https://github.com/hashicorp/waypoint/issues/3820)]

## 0.9.1 (July 28, 2022)

IMPROVEMENTS:

* core: Adds string `replace` function for HCL configs [[GH-3522](https://github.com/hashicorp/waypoint/issues/3522)]
* plugin/nomad: Convert Nomad job for ODR to batch job [[GH-3468](https://github.com/hashicorp/waypoint/issues/3468)]

BUG FIXES:

* internal: Fix deprecation warning being fatal [[GH-3605](https://github.com/hashicorp/waypoint/issues/3605)]
* plugin/tfc: Fix HCL field, `refresh_interval` [[GH-3524](https://github.com/hashicorp/waypoint/issues/3524)]

## 0.9.0 (July 05, 2022)

FEATURES:

* **Waypoint Task Tracking**: Waypoint now tracks the lifecycle of on-demand
runner tasks through a new internal core concept `Task`. As ODR jobs run, Task
will keep track of what part the jobs are at for better debugging and on-demand
runner insight. [[GH-3203](https://github.com/hashicorp/waypoint/issues/3203)]
* cli: **New `runner install` and `runner uninstall` commands** to install/uninstall Waypoint runners to a specified platform [[GH-3335](https://github.com/hashicorp/waypoint/issues/3335)]
* cli: **New `runner profile delete` command** to delete a Waypoint runner profile [[GH-3474](https://github.com/hashicorp/waypoint/issues/3474)]
* cli: **Refactor k8s server install to use Helm** [[GH-3335](https://github.com/hashicorp/waypoint/issues/3335)]
* core: **Add ability to have cli and runners use OAuth2 to get an auth token** [[GH-3298](https://github.com/hashicorp/waypoint/issues/3298)]
* plugin/aws-ecr-pull: **Introduces an `aws-ecr-pull` builder plugin**
that enables using AWS ECR images that are built outside of Waypoint. [[GH-3396](https://github.com/hashicorp/waypoint/issues/3396)]
* plugin/lambda-function-url: **Adds a new plugin and `releaser` component.** This leverages Lambda URLs. [[GH-3187](https://github.com/hashicorp/waypoint/issues/3187)]

IMPROVEMENTS:

* cli: Show list of existing default runner profiles on post-upgrade to warn user
that only one runner profile should be default. [[GH-3497](https://github.com/hashicorp/waypoint/issues/3497)]
* cli: Add interactive input for server upgrade, server uninstall, and destroy commands [[GH-3238](https://github.com/hashicorp/waypoint/issues/3238)]
* cli: Remove unused flag `runner-profile` in `waypoint project apply` [[GH-3318](https://github.com/hashicorp/waypoint/issues/3318)]
* core/api: include commit message in datasource/git response [[GH-3457](https://github.com/hashicorp/waypoint/issues/3457)]
* core: on-demand runner logs are now captured from the underlying platform
and stored in the job system. [[GH-3306](https://github.com/hashicorp/waypoint/issues/3306)]
* core: show helpful errors when using invalid runner profile plugin config [[GH-3465](https://github.com/hashicorp/waypoint/issues/3465)]
* plugin/nomad: Support Nomad service discovery in Nomad platform plugin [[GH-3461](https://github.com/hashicorp/waypoint/issues/3461)]
* runner: runners will now accept and execute multiple jobs concurrently
if multiple jobs are available. On-demand runners continue to execute exactly
one job since they are purpose launched for single job execution. [[GH-3300](https://github.com/hashicorp/waypoint/issues/3300)]
* server: Introduce basic server-side metric collections around operations [[GH-3440](https://github.com/hashicorp/waypoint/issues/3440)]
* serverinstall/k8s: By default, do not set a mem or cpu limit or request for
the default runner profile installed. [[GH-3475](https://github.com/hashicorp/waypoint/issues/3475)]
* ui: Updated UI of breadcrumbs and UX to include current page [[GH-3166](https://github.com/hashicorp/waypoint/issues/3166)]
* upgrade: Warn user if default k8s runner profile has incorrect plugin configs [[GH-3503](https://github.com/hashicorp/waypoint/issues/3503)]

DEPRECATIONS:

* config: More than one app stanza within a waypoint.hcl file is deprecated, and will be removed in 0.10. Please see https://discuss.hashicorp.com/t/deprecating-projects-or-how-i-learned-to-love-apps/40888 for more information.

  Since the initial version of Waypoint, the product has supported the ability to configure multiple apps within a single waypoint.hcl file. This functionality is deprecated and will be removed in the next release. The vast majority of users are not using this functionality and it served mostly as a source of confusion for users. For users who are using a monorepo pattern, we plan to add better workflows for you.

  With a waypoint.hcl now focused on the configuration of a single application, the concept of a project within the Waypoint data model will be removed, moving applications to the top level. This is already how users talk about using Waypoint and we are confident that it will improve the overall understanding of Waypoint as well.

  If you have questions about this change in functionality, we invite you to discuss with us at https://discuss.hashicorp.com/t/deprecating-projects-or-how-i-learned-to-love-apps/40888 or https://github.com/hashicorp/waypoint/issues.
  
  [[GH-3466](https://github.com/hashicorp/waypoint/pull/3466)]

BUG FIXES:

* cli: fix git dirty check that was broken for some versions of the git cli [[GH-3432](https://github.com/hashicorp/waypoint/issues/3432)]
* cli: fix panic when running status report on app with zero prior deployments [[GH-3425](https://github.com/hashicorp/waypoint/issues/3425)]
* core: Fix a rare panic when generating an invite token. [[GH-3505](https://github.com/hashicorp/waypoint/issues/3505)]
* plugin/docker: Ensure that the docker task launcher does not require a resources
block to be set when attempting to load a task config to launch a task. [[GH-3486](https://github.com/hashicorp/waypoint/issues/3486)]
* plugin/docker: fix issue with remote operations for `docker-pull` builder [[GH-3398](https://github.com/hashicorp/waypoint/issues/3398)]
* plugin/k8s: Properly parse kubernetes task launcher config on plugin invoke. [[GH-3475](https://github.com/hashicorp/waypoint/issues/3475)]
* upgrade: Update existing runner profile during server upgrade & change naming convention of initial runner profile [[GH-3490](https://github.com/hashicorp/waypoint/issues/3490)]


## 0.8.2 (May 19, 2022)

IMPROVEMENTS:

* cli: Output message if no runners are found for 'runner list' [[GH-3266](https://github.com/hashicorp/waypoint/issues/3266)]
* core: git cloning now supports recursively cloning submodules [[GH-3351](https://github.com/hashicorp/waypoint/issues/3351)]
* plugin/aws-lambda: Add platform config validation. [[GH-3193](https://github.com/hashicorp/waypoint/issues/3193)]
* plugin/aws-lambda: add support for lambda storage size [[GH-3213](https://github.com/hashicorp/waypoint/issues/3213)]
* plugin/k8s: Add CPU and memory resource limits to on-demand runners through runner profiles and at install time. These resource limits follow the same format that kubernetes expects within `spec.containers[].resources`. [[GH-3307](https://github.com/hashicorp/waypoint/issues/3307)]
* plugin/lambda: Add `static_environment` to deploy plugin [[GH-3282](https://github.com/hashicorp/waypoint/issues/3282)]
* plugin/nomad-jobspec: Add configuration option to parse jobspec as HCL1 instead of HCL2 [[GH-3287](https://github.com/hashicorp/waypoint/issues/3287)]
* plugin/nomad: Support Consul & Vault tokens for job submission [[GH-3222](https://github.com/hashicorp/waypoint/issues/3222)]

BUG FIXES:

* builtin/k8s: Ensure pod.container.static_environment is applied [[GH-3197](https://github.com/hashicorp/waypoint/issues/3197)]
* cli: Fix missing bootstrap hint with server run command [[GH-3196](https://github.com/hashicorp/waypoint/issues/3196)]
* cli: Prevent panic when releasing unsuccessful deployments [[GH-3207](https://github.com/hashicorp/waypoint/issues/3207)]
* cli: Show better error message when there are no Waypoint contexts when attempting
to open the UI [[GH-3262](https://github.com/hashicorp/waypoint/issues/3262)]
* install/nomad: Add support for CSI params & secrets to Nomad install [[GH-3279](https://github.com/hashicorp/waypoint/issues/3279)]
* install/nomad: Fix DB directory for Nomad install [[GH-3261](https://github.com/hashicorp/waypoint/issues/3261)]
* internal/cli: Fix `waypoint exec` workspace selection [[GH-3226](https://github.com/hashicorp/waypoint/issues/3226)]
* plugin/docker: support remote operations for docker-pull plugin [[GH-3253](https://github.com/hashicorp/waypoint/issues/3253)]
* plugin/k8s: Ensure `container=docker` environment variable is set for Kaniko
to properly detect running inside a container, which prevented on-demand
runners from working on Kubernetes 1.23. [[GH-3322](https://github.com/hashicorp/waypoint/issues/3322)]
* server: fix issue cleaning up tasks in Kubernetes that completed successfully [[GH-3299](https://github.com/hashicorp/waypoint/issues/3299)]

## 0.8.1 (April 08, 2022)

BUG FIXES:

* server: Fix runner database by setting proper runner bucket for initialization. Sever upgrades to 0.8.0 would previously fail before this fix. [[GH-3200](https://github.com/hashicorp/waypoint/issues/3200)]

## 0.8.0 (April 07, 2022)

FEATURES:

* **Targetable Runners with Labels**: Waypoint now supports runner profiles that target specific on-demand runners by labels.
Projects and/or Apps can be configured to use a specific runner profile, identified by name.
The runner profile will target all operations to a specific on-demand runner by ID, or labels on the runner. [[GH-3145](https://github.com/hashicorp/waypoint/issues/3145)]
* **cli:** Introduce a new CLI command for job management and inspection
`waypoint job`. [[GH-3067](https://github.com/hashicorp/waypoint/issues/3067)]
* **core, cli:** Support setting variables to `sensitive` and obfuscate those values in outputs [[GH-3138](https://github.com/hashicorp/waypoint/issues/3138)]
* **plugin/nomad:** Nomad jobspec canary promotion releaser [[GH-2938](https://github.com/hashicorp/waypoint/issues/2938)]

IMPROVEMENTS:

* cli: Add `waypoint runner inspect` command [[GH-3004](https://github.com/hashicorp/waypoint/issues/3004)]
* cli: Add a way for `waypoint context create` to set the Waypoint server platform [[GH-3055](https://github.com/hashicorp/waypoint/issues/3055)]
* cli: Display count of instance connections in deployment status reports [[GH-3043](https://github.com/hashicorp/waypoint/issues/3043)]
* cli: Introduce a `-dangerously-force` flag to attempt to force cancel a job [[GH-3102](https://github.com/hashicorp/waypoint/issues/3102)]
* cli: Print operation sequence ids and pushed artifact id [[GH-3081](https://github.com/hashicorp/waypoint/issues/3081)]
* cli: `runner list` shows runner labels [[GH-3133](https://github.com/hashicorp/waypoint/issues/3133)]
* cli: `runner profile set` deprecates the `-env-vars` flag in favor of the `-env-var` flag instead. [[GH-3136](https://github.com/hashicorp/waypoint/issues/3136)]
* cli: `runner profile` command set now supports target runner labels [[GH-3145](https://github.com/hashicorp/waypoint/issues/3145)]
* core: runners automatically reconnect on startup if the server is unavailable
or becomes unavailable during the startup process [[GH-3087](https://github.com/hashicorp/waypoint/issues/3087)]
* core: runners can have labels, which are used for targeting and metadata [[GH-2954](https://github.com/hashicorp/waypoint/issues/2954)]
* core: runners can tolerate a server outage during job execution and will
wait for the server to come back online [[GH-3119](https://github.com/hashicorp/waypoint/issues/3119)]
* plugin/aws-ecr: Output `architecture` from Docker image input [[GH-3046](https://github.com/hashicorp/waypoint/issues/3046)]
* plugin/ecs: Add `cpu_architecture` aws-ecs parameter
to support deploying Docker images built by the Apple M1 chip on ECS [[GH-3068](https://github.com/hashicorp/waypoint/issues/3068)]
* plugin/aws-lambda: Default lambda architecture to Docker/ECR image architecture [[GH-3046](https://github.com/hashicorp/waypoint/issues/3046)]
* plugin/docker: Output `architecture` from Docker builder [[GH-3046](https://github.com/hashicorp/waypoint/issues/3046)]
* plugin/k8s: don't error if previous deployment is not found during cleanup [[GH-3070](https://github.com/hashicorp/waypoint/issues/3070)]
* plugin/nomad: Resource manager for Nomad jobspec [[GH-2938](https://github.com/hashicorp/waypoint/issues/2938)]

BUG FIXES:

* ceb: Fix connecting to servers with TLS verification [[GH-3167](https://github.com/hashicorp/waypoint/issues/3167)]
* cli: Fix panic in `waypoint plugin` CLI [[GH-3095](https://github.com/hashicorp/waypoint/issues/3095)]
* cli: Fix panic when attempting to reinstall autocomplete [[GH-2986](https://github.com/hashicorp/waypoint/issues/2986)]
* cli: Fix the -set-default flag on `waypoint context create` [[GH-3044](https://github.com/hashicorp/waypoint/issues/3044)]
* core: Ensure remote runners have dynamic config sources overrides for
evaluating defaults for job variables. [[GH-3171](https://github.com/hashicorp/waypoint/issues/3171)]
* core: Fix panic when running `waypoint build` remotely outside of project directory. [[GH-3165](https://github.com/hashicorp/waypoint/issues/3165)]
* core: Fix panic where on-demand runner config was nil before starting task [[GH-3054](https://github.com/hashicorp/waypoint/issues/3054)]
* plugin/alb: Handle DuplicateListener errors from aws-alb releaser [[GH-3035](https://github.com/hashicorp/waypoint/issues/3035)]
* plugin/aws-alb: Use Route53 zone id when destroying a resource record [[GH-3076](https://github.com/hashicorp/waypoint/issues/3076)]
* plugin/docker: Add Docker auth support for builder and platform, and add config options for
docker-pull auth and registry auth [[GH-2895](https://github.com/hashicorp/waypoint/issues/2895)]
* plugin/k8s: clean up pending pods from cancelled jobs [[GH-3143](https://github.com/hashicorp/waypoint/issues/3143)]
* plugin/k8s: fix issue when destroying multiple deployments [[GH-3111](https://github.com/hashicorp/waypoint/issues/3111)]
* plugin/nomad: Fix Nomad job namespace when using ODRs [[GH-2896](https://github.com/hashicorp/waypoint/issues/2896)]
* server: Fix project poll job singleton id to only include application on project
poll jobs if exist. Otherwise, only include the workspace and project. [[GH-3158](https://github.com/hashicorp/waypoint/issues/3158)]
* ui: Fix missing links to resources [[GH-3172](https://github.com/hashicorp/waypoint/issues/3172)]
* ui: fix missing link on release detail page [[GH-3142](https://github.com/hashicorp/waypoint/issues/3142)]

## 0.7.2 (February 24, 2022)

FEATURES:

* **Targetable Runners:** Allow apps and projects to target specific runner profiles.
Allow runner profiles to target specific remote runners. [[GH-2862](https://github.com/hashicorp/waypoint/issues/2862)]
* **Introduce executing Trigger URL configurations via HTTP**. Users can
start a trigger via HTTP and stream the job event stream output directly over
http. [[GH-2970](https://github.com/hashicorp/waypoint/issues/2970)]

IMPROVEMENTS:

* plugin/docker: Add parameter to disable the build cache [[GH-2953](https://github.com/hashicorp/waypoint/issues/2953)]

BUG FIXES:

* cli: Fix panic on nil value for project [[GH-2968](https://github.com/hashicorp/waypoint/issues/2968)]
* cli: Replace panic with message when attempting to `config get -app` without a `-project` flag while outside a project directory [[GH-3039](https://github.com/hashicorp/waypoint/issues/3039)]
* cli: requires -app flag if `config set -scope=app` is set [[GH-3040](https://github.com/hashicorp/waypoint/issues/3040)]
* server: Cache Horizon hostname URL lookup when listing deployments in the
`UI_ListDeployments bundle`. Now we look up the deployment URL once, and craft
the deployment URLs based on the original hostname lookup. [[GH-2950](https://github.com/hashicorp/waypoint/issues/2950)]
* ui: fixed issue with focus jumping back to the skip link on automatic refresh [[GH-3019](https://github.com/hashicorp/waypoint/issues/3019)]

## 0.7.1 (January 25, 2022)

IMPROVEMENTS:

* ui: Automatically select the appropriate workspace [[GH-2835](https://github.com/hashicorp/waypoint/issues/2835)]

BUG FIXES:

* build: Add arm64 ceb to released build [[GH-2945](https://github.com/hashicorp/waypoint/issues/2945)]
* plugin/nomad: Fix Nomad job namespace when using ODRs [[GH-2896](https://github.com/hashicorp/waypoint/issues/2896)]
* ui: Ensure logs update correctly when switching between deployments [[GH-2901](https://github.com/hashicorp/waypoint/issues/2901)]
* ui: Limit number of deployments requested [[GH-2930](https://github.com/hashicorp/waypoint/issues/2930)]
* ui: Update empty logs message [[GH-2925](https://github.com/hashicorp/waypoint/issues/2925)]

## 0.7.0 (January 13, 2022)

FEATURES:

* **Source variable values from remote systems:** The default value for input variables
can now use the `configdynamic` function to source data from Vault, Consul,
Terraform Cloud, and more. This is a pluggable system. [[GH-2889](https://github.com/hashicorp/waypoint/issues/2889)]
* **Workspace and Label-scoped Configuration:** Build, registry, deploy, and
release configurations can now be changed depending on the current workspace
or label sets. This can be used to alter configurations between environments
(staging, production, etc.) or metadata (region, etc.). [[GH-2699](https://github.com/hashicorp/waypoint/issues/2699)]
* **core: Introduce RunTriggers to Waypoint Server and CLI.** Triggers can be configured
ahead of time to execute lifecycle operations on demand through the Waypoint API.
Currently, only the gRPC API is supported, but in the future an HTTP endpoint
will be added to be used within CI. [[GH-2840](https://github.com/hashicorp/waypoint/issues/2840)]
* ui/auth: users can now authenticate through the UI using an OIDC provider [[GH-2688](https://github.com/hashicorp/waypoint/issues/2688)]
* ui: Add Exec Terminal to the web UI [[GH-2849](https://github.com/hashicorp/waypoint/issues/2849)]
* ui: Added a tab with an overview of the resources provisioned by operations [[GH-2777](https://github.com/hashicorp/waypoint/issues/2777)]
* ui: Added timeline component to artifact details pages [[GH-2793](https://github.com/hashicorp/waypoint/issues/2793)]
* ui: Update UI of builds and releases tab [[GH-2852](https://github.com/hashicorp/waypoint/issues/2852)]
* ui: add workspace switcher to app pages [[GH-2674](https://github.com/hashicorp/waypoint/issues/2674)]
* ui: reformatted app overview page and header of artifact details page [[GH-2606](https://github.com/hashicorp/waypoint/issues/2606)]

IMPROVEMENTS:

* cli: Add alias for -app and -workspace flags [[GH-2700](https://github.com/hashicorp/waypoint/issues/2700)]
* cli: Add new `workspace create` command [[GH-2797](https://github.com/hashicorp/waypoint/issues/2797)]
* cli: Deprecate -remote flag for operations, replace with -local flag [[GH-2771](https://github.com/hashicorp/waypoint/issues/2771)]
* cli: Enhance warning for project flag mismatches when project parsed from config [[GH-2815](https://github.com/hashicorp/waypoint/issues/2815)]
* cli: Report where each operation runs (locally vs remotely) [[GH-2795](https://github.com/hashicorp/waypoint/issues/2795)]
* cli: Warn if about to perform a remote operation with a dirty local git state [[GH-2799](https://github.com/hashicorp/waypoint/issues/2799)]
* install/nomad: Ensure static runner has started during install by validating its
running status for a few seconds once it is in a "running" state. [[GH-2698](https://github.com/hashicorp/waypoint/issues/2698)]
* plugin/docker: inject `arm64` Waypoint entrypoints for arm images [[GH-2692](https://github.com/hashicorp/waypoint/issues/2692)]
* plugin/ecs: Implement the destruction of AWS resources created when deploying a workspace [[GH-2684](https://github.com/hashicorp/waypoint/issues/2684)]
* plugin/pack: detect non-Intel Docker server and show a warning [[GH-2692](https://github.com/hashicorp/waypoint/issues/2692)]
* serverinstall/ecs: Add permissions to the ECS runner IAM policy to allow the removal of security groups and de-registration of task [[GH-2684](https://github.com/hashicorp/waypoint/issues/2684)]
* serverinstall: Set Nomad's ODR profile name to "nomad" [[GH-2713](https://github.com/hashicorp/waypoint/issues/2713)]
* ui: Improved UX of screen readers' transition between pages [[GH-2837](https://github.com/hashicorp/waypoint/issues/2837)]
* ui: Updated list items UI on deployments tab [[GH-2879](https://github.com/hashicorp/waypoint/issues/2879)]
* ui: Updated the deployments tab UI/UX [[GH-2773](https://github.com/hashicorp/waypoint/issues/2773)]
* ui: upgraded icons to flight icons library [[GH-2681](https://github.com/hashicorp/waypoint/issues/2681)]

BREAKING CHANGES:

* core: `configdynamic` has been renamed to `dynamic`. The existing function
name continues to work but is deprecated and may be removed in a future version. [[GH-2892](https://github.com/hashicorp/waypoint/issues/2892)]
* plugin/docker: `img`-based Dockerless builds are no longer supported.
Dockerless builds are still fully supported via Kaniko and on-demand
runners that shipped in Waypoint 0.6. Static runners without access to
a Docker daemon can no longer build images. [[GH-2534](https://github.com/hashicorp/waypoint/issues/2534)]

BUG FIXES:

* cli: Added check for empty deployUrl in output for release switch case [[GH-2755](https://github.com/hashicorp/waypoint/issues/2755)]
* cli: Fix issue where users could not disable project polling from the CLI [[GH-2673](https://github.com/hashicorp/waypoint/issues/2673)]
* core: fix issue where runners would fail but not shut down [[GH-2571](https://github.com/hashicorp/waypoint/issues/2571)]
* ui: Fix edge case issue where users would not be redirected to the authentication screen if no api token was set [[GH-2696](https://github.com/hashicorp/waypoint/issues/2696)]
* ui: Logs & Exec Terminals resize smoothly [[GH-2890](https://github.com/hashicorp/waypoint/issues/2890)]

## 0.6.3 (December 10, 2021)

SECURITY UPDATE:

* Update Go build version to 1.16.12 per [Go security release](https://groups.google.com/g/golang-announce/c/hcmEScgc00k?pli=1)

## 0.6.2 (November 4, 2021)

FEATURES:

* plugin/nomad: **Introduce Nomad On-Demand Runner support for Waypoint server.** Launch
tasks to build containers in short lived runners. [[GH-2593](https://github.com/hashicorp/waypoint/issues/2593)]

IMPROVEMENTS:

* cli/login: Allow login from-kubernetes command to work with non-default namespace installations [[GH-2575](https://github.com/hashicorp/waypoint/issues/2575)]
* serverinstall/nomad: Update install helper to always setup a Consul service
with a backend and ui service tag. [[GH-2597](https://github.com/hashicorp/waypoint/issues/2597)]

BUG FIXES:

* cli: Use values from filter flags when listing deployments and releases [[GH-2672](https://github.com/hashicorp/waypoint/issues/2672)]
* cli: `waypoint status` and `waypoint status -app` no longer display destroyed deployments [[GH-2564](https://github.com/hashicorp/waypoint/issues/2564)]
* core: Fix a panic where a custom Waypoint plugin would panic if the plugin
did not properly implement a Registry component with AccessInfoFunc() [[GH-2532](https://github.com/hashicorp/waypoint/issues/2532)]
* serverinstall/k8s: Clean up rbac resources on uninstall, and do not error when existing rbac resources are detected during server upgrade. [[GH-2654](https://github.com/hashicorp/waypoint/issues/2654)]
* ui: avoid loading *all* status reports [[GH-2562](https://github.com/hashicorp/waypoint/issues/2562)]
* ui: improve docker reference parsing [[GH-2518](https://github.com/hashicorp/waypoint/issues/2518)]

## 0.6.1 (October 21, 2021)

BUG FIXES:

* cli: Fix `project apply` to set runner profiles by name [[GH-2489](https://github.com/hashicorp/waypoint/issues/2489)]
* cli: Fix displaying config variables set with pre-0.6.0 Waypoint [[GH-2535](https://github.com/hashicorp/waypoint/issues/2535)]
* cli: Fix panic in logs and exec commands [[GH-2526](https://github.com/hashicorp/waypoint/issues/2526)]
* cli: Fix issue where sending ctrl-c to the CLI could block all subsiquent operations on the app/project/workspace for two minutes [[GH-2513](https://github.com/hashicorp/waypoint/issues/2513)]
* serverinstall/ecs: Fix potential panic in some ECS On-Demand Runner releases [[GH-2533](https://github.com/hashicorp/waypoint/issues/2533)]
* serverinstall/ecs: Update ODR role permissions to enable status reports [[GH-2543](https://github.com/hashicorp/waypoint/issues/2543)]

## 0.6.0 (October 14, 2021)

FEATURES:

* cli: Add new commands `workspace inspect` and `workspace list` to view and list
workspaces. [[GH-2385](https://github.com/hashicorp/waypoint/issues/2385)]
* cli: Allow `install` cmd to support pass-through flags to `server run` [[GH-2328](https://github.com/hashicorp/waypoint/issues/2328)]
* config: Specify configuration (env vars and files) for runners while executing
operations related to a specific to that project or application. [[GH-2237](https://github.com/hashicorp/waypoint/issues/2237)]
* config: Specify configuration that is scoped to deployments in certain workspaces
or label sets. [[GH-2237](https://github.com/hashicorp/waypoint/issues/2237)]
* config: `labels` variable for accessing the label set of an operation [[GH-2065](https://github.com/hashicorp/waypoint/issues/2065)]
* config: New functions `selectormatch` and `selectorlookup` for working with
label selectors [[GH-2065](https://github.com/hashicorp/waypoint/issues/2065)]
* core/server: Allow exporting of grpc server traces and stats by introducing OpenCensus and DataDog telemetry for Waypoint Server by request [[GH-2402](https://github.com/hashicorp/waypoint/issues/2402)]
* core: Runner configuration can now write to files [[GH-2201](https://github.com/hashicorp/waypoint/issues/2201)]
* core: Runner configuration can use dynamic configuration sources [[GH-2201](https://github.com/hashicorp/waypoint/issues/2201)]
* platform/nomad: Add persistent data volumes to nomad deploy [[GH-2282](https://github.com/hashicorp/waypoint/issues/2282)]
* plugin/docker: Add ability to build images with kaniko inside an ondemand runner [[GH-2056](https://github.com/hashicorp/waypoint/issues/2056)]
* plugin/helm: A new plugin "helm" can deploy using Helm charts. [[GH-2336](https://github.com/hashicorp/waypoint/issues/2336)]
* plugin/k8s: Report events on failed pods when a deployment fails [[GH-2399](https://github.com/hashicorp/waypoint/issues/2399)]
* plugin/k8s: Allows users to add sidecar containers to apps using the k8s plugin config. [[GH-2428](https://github.com/hashicorp/waypoint/issues/2428)]
* plugin/pack: Add ability to build images with kaniko inside an ondemand runner [[GH-2056](https://github.com/hashicorp/waypoint/issues/2056)]
* runner: Add ability to build images without needing a containarization API [[GH-2056](https://github.com/hashicorp/waypoint/issues/2056)]
* runner: Adds ondemand runners, single job runner processes launched via the task API [[GH-2056](https://github.com/hashicorp/waypoint/issues/2056)]
* ui: Allow config variables to be managed in the browser UI [[GH-1915](https://github.com/hashicorp/waypoint/issues/1915)]
* ui: Deployment resources [[GH-2317](https://github.com/hashicorp/waypoint/issues/2317)]
* ui: Release resources [[GH-2386](https://github.com/hashicorp/waypoint/issues/2386)]
* ui: Overview section added + Docker container information displayed [[GH-2352](https://github.com/hashicorp/waypoint/issues/2352)]

IMPROVEMENTS:

* cli/serverinstall/k8s: Add new cluster role and binding to allow nodeport services to work [[GH-2412](https://github.com/hashicorp/waypoint/issues/2412)]
* cli/serverinstall/k8s: Fix a problem where deployments would be marked as "Degraded", but were actually fine. [[GH-2412](https://github.com/hashicorp/waypoint/issues/2412)]
* cli: Add new context subcommand "set" to set the workspace value for the current
context. [[GH-2353](https://github.com/hashicorp/waypoint/issues/2353)]
* cli: Remove unused arg and use sequence ID for CLI message in `release` [[GH-2426](https://github.com/hashicorp/waypoint/issues/2426)]
* cli: Return help on malformed command [[GH-2444](https://github.com/hashicorp/waypoint/issues/2444)]
* cli: Update base commands to default to all apps within the project if project has more than one application [[GH-2413](https://github.com/hashicorp/waypoint/issues/2413)]
* cli: Use default log level of debug instead of trace on server install [[GH-2325](https://github.com/hashicorp/waypoint/issues/2325)]
* cli: `server run` can now create a non-TLS HTTP listener. This listener
redirects to HTTPS unless X-Forwarded-Proto is https. [[GH-2347](https://github.com/hashicorp/waypoint/issues/2347)]
* cli: `login` subcommand defaults server port to 9701 if it isn't set [[GH-2320](https://github.com/hashicorp/waypoint/issues/2320)]
* config: `gitrefpretty` no longer requires `git` to be installed [[GH-2371](https://github.com/hashicorp/waypoint/issues/2371)]
* config: Input variables (`variable`) can now use an `env` key to specify
alternate environment variable names to source variable values from [[GH-2362](https://github.com/hashicorp/waypoint/issues/2362)]
* core: Automatically remotely init projects with a Git data source [[GH-2145](https://github.com/hashicorp/waypoint/issues/2145)]
* core: HTTP requests from Kubernetes probes are logged at a trace level rather than info [[GH-2348](https://github.com/hashicorp/waypoint/issues/2348)]
* core: Easier to understand error messages when using incompatible plugins [[GH-2143](https://github.com/hashicorp/waypoint/issues/2143)]
* core: Server with custom TLS certificates will automatically reload and rotate
the TLS certificates when they change on disk [[GH-2346](https://github.com/hashicorp/waypoint/issues/2346)]
* plugin/docker: Add support for multi-stage Dockerfile builds [[GH-1992](https://github.com/hashicorp/waypoint/issues/1992)]
* plugin/k8s: Add new ability to release by creating an ingress resource to route
traffic to a service backend from an ingress controller. [[GH-2261](https://github.com/hashicorp/waypoint/issues/2261)]
* plugin/k8s: Introduce a new config option `autoscale`, which creates a horizontal
pod autoscaler for each deployment. [[GH-2309](https://github.com/hashicorp/waypoint/issues/2309)]
* plugin/k8s: Introduce a new config option `cpu` and `memory` for defining
resource limits and requests for a pod in a deployment. [[GH-2309](https://github.com/hashicorp/waypoint/issues/2309)]
* plugin/k8s: Use sequence number in k8s deployment name for improved traceability to waypoint deployments. [[GH-2296](https://github.com/hashicorp/waypoint/issues/2296)]
* ui: Display project remote initialization state [[GH-2145](https://github.com/hashicorp/waypoint/issues/2145)]
* ui: Gitops users not using Git polling can run "Up" from the browser [[GH-2331](https://github.com/hashicorp/waypoint/issues/2331)]
* ui: Improve design of status row on Build/Deployment/Release detail pages [[GH-2036](https://github.com/hashicorp/waypoint/issues/2036)]
* ui: Improve tab styles for dark mode [[GH-2053](https://github.com/hashicorp/waypoint/issues/2053)]
* ui: Toggle checkboxes are nicely styled in dark mode [[GH-2410](https://github.com/hashicorp/waypoint/issues/2410)]
* ui: Improve the input field for server-side HCL file on the settings page [[GH-2168](https://github.com/hashicorp/waypoint/issues/2168)]
* ui: The rendering of the Application and Operation Logs has been greatly improved [[GH-2356](https://github.com/hashicorp/waypoint/issues/2356)]

BUG FIXES:

* cli: Fix a panic in `waypoint status` when no successful release is available [[GH-2436](https://github.com/hashicorp/waypoint/issues/2436)]
* cli: Fix logic on when a rocket indicator shows in `release list` [[GH-2426](https://github.com/hashicorp/waypoint/issues/2426)]
* config: Fix dynamic config vars targeting files. [[GH-2416](https://github.com/hashicorp/waypoint/issues/2416)]
* entrypoint: Fix issue injecting waypoint-entrypoint multiple times [[GH-2447](https://github.com/hashicorp/waypoint/issues/2447)]
* plugin/docker: Resolve image identifiers properly [[GH-2067](https://github.com/hashicorp/waypoint/issues/2067)]
* plugin/docker: Support SSH hosts for entrypoint injection [[GH-2277](https://github.com/hashicorp/waypoint/issues/2277)]
* plugin/k8: Setup Kubernetes services for different workspaces properly [[GH-2399](https://github.com/hashicorp/waypoint/issues/2399)]
* server: Adds API validation to ensure server doesn't panic when given an empty
request body [[GH-2273](https://github.com/hashicorp/waypoint/issues/2273)]
* server: Validate GetDeployment request has a valid request body to avoid a server
panic. [[GH-2269](https://github.com/hashicorp/waypoint/issues/2269)]
* ui: Fixed config variable duplication when renaming [[GH-2421](https://github.com/hashicorp/waypoint/issues/2421)]
* ui: Notification messages display nicely when containing long words such as URLs [[GH-2411](https://github.com/hashicorp/waypoint/issues/2411)]

## 0.5.2 (September 09, 2021)

FEATURES:

* cli: Add a new command for inspecting project information, `waypoint project inspect`. [[GH-2055](https://github.com/hashicorp/waypoint/issues/2055)]

IMPROVEMENTS:

* cli/status: Include a way to refresh project application statuses for deployments and releases with a '-refresh' flag before showing the status view [[GH-2081](https://github.com/hashicorp/waypoint/issues/2081)]
* cli: Add functionality to list releases with `waypoint release list` command [[GH-2082](https://github.com/hashicorp/waypoint/issues/2082)]
* core: App status polling will always queue status reports refresh jobs for latest deployment and release if present [[GH-2039](https://github.com/hashicorp/waypoint/issues/2039)]
* plugin/aws-ecs: Allow configuration of ALB subnets independently of service subnets [[GH-2205](https://github.com/hashicorp/waypoint/issues/2205)]
* plugin/aws-ecs: Allow public ip assignment for tasks to be disabled [[GH-2205](https://github.com/hashicorp/waypoint/issues/2205)]
* plugin/aws-ecs: Deployments delete their resources on failure. [[GH-2098](https://github.com/hashicorp/waypoint/issues/2098)]
* plugin/aws-ecs: Error messages contain additional context [[GH-2098](https://github.com/hashicorp/waypoint/issues/2098)]
* plugin/aws-ecs: Improve security of ECS tasks by restricting ingress to the ALB [[GH-2098](https://github.com/hashicorp/waypoint/issues/2098)]
* plugin/aws-ecs: More complete list of resources displayed in `waypoint deploy` logs [[GH-2098](https://github.com/hashicorp/waypoint/issues/2098)]
* plugin/aws-ecs: Support for status reports, enabling `waypoint status` for ECS deployments [[GH-2098](https://github.com/hashicorp/waypoint/issues/2098)]
* plugin/aws: Add ability to pass IAM policy ARNs for attaching to task role [[GH-1935](https://github.com/hashicorp/waypoint/issues/1935)]

BUG FIXES:

* internal/config: Fix parsing of complex HCL types in `-var-file` [[GH-2217](https://github.com/hashicorp/waypoint/issues/2217)]
* plugin/aws-ecs: Fix panic when specifying a sidecar without a health check [[GH-2098](https://github.com/hashicorp/waypoint/issues/2098)]
* plugin/nomad: Only use non-empty job.StatusDescription for HealthMessage [[GH-2093](https://github.com/hashicorp/waypoint/issues/2093)]
* server/singleprocess: Stop returning error when polling an app with no deployment or release [[GH-2204](https://github.com/hashicorp/waypoint/issues/2204)]
* ui: Fix leaky project repository settings being reused when creating a new project [[GH-2250](https://github.com/hashicorp/waypoint/issues/2250)]

## 0.5.1 (August 19, 2021)

IMPROVEMENTS:

* cli/status: Display '(unknown)' when the time a status report was generated is
not known [[GH-2047](https://github.com/hashicorp/waypoint/issues/2047)]
* cli/uninstall: Remove hard requirement on platform flag, attempt to read server
platform from server context. Platform flag overrides anything set in a server
platform context [[GH-2052](https://github.com/hashicorp/waypoint/issues/2052)]

BUG FIXES:

* plugin/aws/alb: Always set the generated time for a status report [[GH-2048](https://github.com/hashicorp/waypoint/issues/2048)]
* plugin/aws/ecs: Fix destroy non-latest deployments in ECS [[GH-2054](https://github.com/hashicorp/waypoint/issues/2054)]
* ui: Prevent deletion of git/input variable settings when saving the other [[GH-2057](https://github.com/hashicorp/waypoint/issues/2057)]


## 0.5.0 (August 12, 2021)

FEATURES:

* **Status Reports:** Waypoint now has multiple improvements to support status
checks for deployed resources. See `Improvements` for more.
* **Input variables:** Waypoint now allows users to parameterize the waypoint.hcl
file through an input variable system. [[GH-1548](https://github.com/hashicorp/waypoint/issues/1548)]
* **OIDC Authentication and User Accounts:** Waypoint now has a user account system
and can be configured to sign up and log in users using any OIDC provider
such as Google, GitLab, etc. [[GH-1831](https://github.com/hashicorp/waypoint/issues/1831)]
* cli: can login with a token using the new `waypoint login` command [[GH-1848](https://github.com/hashicorp/waypoint/issues/1848)]
* cli: new "waypoint user" CLI for user management [[GH-1864](https://github.com/hashicorp/waypoint/issues/1864)]
* core: platform plugins may now advertise deployment-specific URLs [[GH-1387](https://github.com/hashicorp/waypoint/issues/1387)]
* ui: Show deployment URL if available [[GH-1739](https://github.com/hashicorp/waypoint/issues/1739)]
* ui: added button on individual artifact (deployments + releases) page for on demand health checks [[GH-1911](https://github.com/hashicorp/waypoint/issues/1911)]

IMPROVEMENTS:

* **Status reports:** server: Continuously generate a status report for an application after the initial
deployment or release for projects backed by a git data source [[GH-1801](https://github.com/hashicorp/waypoint/issues/1801)]
* cli: Adds a `-git-path` flag to `waypoint project apply` [[GH-2013](https://github.com/hashicorp/waypoint/issues/2013)]
* core: git poller setting to optionally ignore changes outside of the configured project path [[GH-1821](https://github.com/hashicorp/waypoint/issues/1821)]
* entrypoint: Can disable `waypoint exec` only by setting the
`WAYPOINT_CEB_DISABLE_EXEC` environment variable to a truthy value. [[GH-1973](https://github.com/hashicorp/waypoint/issues/1973)]
* plugin/aws/alb: Update ALB Releaser to use new SDK Resource Manager [[GH-1648](https://github.com/hashicorp/waypoint/issues/1648)]
* plugin/aws/ecs: Add ability to specify security group IDs [[GH-1919](https://github.com/hashicorp/waypoint/issues/1919)]
* plugin/docker: Enables image build for specified platform [[GH-1949](https://github.com/hashicorp/waypoint/issues/1949)]
* plugin/google: Add non-blocking alert if unable to delete revision [[GH-2005](https://github.com/hashicorp/waypoint/issues/2005)]
* plugin/google: Implement DestroyWorkspace to cleanup all created resources [[GH-2005](https://github.com/hashicorp/waypoint/issues/2005)]
* plugin/google: Update Google Cloud platform to use SDK Resource Manager [[GH-2005](https://github.com/hashicorp/waypoint/issues/2005)]
* plugin/google: Update error message to be more helpful [[GH-1958](https://github.com/hashicorp/waypoint/issues/1958)]
* plugin/k8s: Include deployment and release resources in `waypoint status` output. [[GH-2024](https://github.com/hashicorp/waypoint/issues/2024)]
* plugin/k8s: Update K8s Releaser to use SDK Resource Manager [[GH-1938](https://github.com/hashicorp/waypoint/issues/1938)]
* plugin/k8s: Updates release status report to check k8s service status [[GH-2024](https://github.com/hashicorp/waypoint/issues/2024)]
* plugin/nomad: Add consul service optional flags [[GH-2033](https://github.com/hashicorp/waypoint/issues/2033)]
* plugin/nomad: Update Nomad platform to use SDK Resource Manager [[GH-1941](https://github.com/hashicorp/waypoint/issues/1941)]
* ui: Add ability to add input variables in the project settings UI [[GH-1658](https://github.com/hashicorp/waypoint/issues/1658)]
* ui: Add dynamic page titles [[GH-1916](https://github.com/hashicorp/waypoint/issues/1916)]
* ui: Add git commit SHAs to operations in the browser UI [[GH-2025](https://github.com/hashicorp/waypoint/issues/2025)]
* ui: Update authentication page with new supported `waypoint user token` command. [[GH-2006](https://github.com/hashicorp/waypoint/issues/2006)]

BUG FIXES:

* plugin/aws-alb: Fix issue destroying when Target Group still in use [[GH-1648](https://github.com/hashicorp/waypoint/issues/1648)]
* plugin/docker: Fix docker buildkit build failures [[GH-1937](https://github.com/hashicorp/waypoint/issues/1937)]
* plugin/nomad: Fix case where nomad error would be ignored during a status check [[GH-1723](https://github.com/hashicorp/waypoint/issues/1723)]
* plugin/k8s: Fix `ports` configurability [[GH-1650](https://github.com/hashicorp/waypoint/issues/1650)]
* serverinstall/ecs: Handle errors when resources are already destroyed [[GH-1984](https://github.com/hashicorp/waypoint/issues/1984)]
* ui: Display and read/write base64 strings correctly in SSH and Hcl inputs [[GH-1967](https://github.com/hashicorp/waypoint/issues/1967)]


## 0.4.2 (July 22, 2021)

FEATURES:

* plugin: Allow debugging of plugins with tools like delve [[GH-1716](https://github.com/hashicorp/waypoint/issues/1716)]

IMPROVEMENTS:

* serverinstall/k8s: Add information to cli output for upgrade path [[GH-1886](https://github.com/hashicorp/waypoint/issues/1886)]
* ui: Incorporate pushed artifacts into build display [[GH-1840](https://github.com/hashicorp/waypoint/issues/1840)]

BUG FIXES:

* plugin/aws/ecs: Validate memory and cpu values [[GH-1872](https://github.com/hashicorp/waypoint/issues/1872)]
* plugin/nomad: Fix broken -nomad-runner-memory and -nomad-server-memory flags [[GH-1895](https://github.com/hashicorp/waypoint/issues/1895)]
* serverinstall/ecs: Validate memory and cpu values [[GH-1872](https://github.com/hashicorp/waypoint/issues/1872)]

## 0.4.1 (July 1, 2021)

FEATURES:

* config: Add `${app.name}` variable [[GH-1709](https://github.com/hashicorp/waypoint/issues/1709)]

IMPROVEMENTS:

* cli: Fix incorrect description for `hostname list` command [[GH-1628](https://github.com/hashicorp/waypoint/issues/1628)]
* core: Correct parsing of boolean environment variables [[GH-1699](https://github.com/hashicorp/waypoint/issues/1699)]
* plugin/aws-alb: Update ALB Releaser to use new SDK Resource Manager [[GH-1648](https://github.com/hashicorp/waypoint/issues/1648)]
* ui: Add reporting on status of a release [[GH-1657](https://github.com/hashicorp/waypoint/issues/1657)]

BUG FIXES:

* builtin/k8s: Fix `ports` configurability [[GH-1650](https://github.com/hashicorp/waypoint/issues/1650)]
* cli: Fix issue parsing string slice inputs [[GH-1669](https://github.com/hashicorp/waypoint/issues/1669)]
* cli: Ignore error on Unimplemented for health checks [[GH-1596](https://github.com/hashicorp/waypoint/issues/1596)]
* cli: Fix crash that could occur when running commands outside the context of a project with an hcl config file. [[GH-1710](https://github.com/hashicorp/waypoint/issues/1710)]
* cli: Prevent use of operation flags on `runner agent` command [[GH-1708](https://github.com/hashicorp/waypoint/issues/1708)]
* cli: Set runner poll interval default for runner defined in waypoint.hcl [[GH-1690](https://github.com/hashicorp/waypoint/issues/1690)]
* cli: List deployments shows status for each deployment [[GH-1594](https://github.com/hashicorp/waypoint/issues/1594)]
* core: Fix crash that could occur when using `templatefile` with certain HCL files [[GH-1679](https://github.com/hashicorp/waypoint/issues/1679)]
* plugin/aws-alb: Fix issue destroying when Target Group still in use [[GH-1648](https://github.com/hashicorp/waypoint/issues/1648)]
* plugin/docker: Fix issue falling back to `img` for builds when docker daemon not present [[GH-1685](https://github.com/hashicorp/waypoint/issues/1685)]
* plugin/nomad: Fix case where Nomad error would be ignored during a status check [[GH-1723](https://github.com/hashicorp/waypoint/issues/1723)]
* server/client: Configure keepalive properties to RPC connections to persist connections even after inactivity [[GH-1735](https://github.com/hashicorp/waypoint/issues/1735)]
* server/runner: Correctly exit liveness listener when connection is closed [[GH-1732](https://github.com/hashicorp/waypoint/issues/1732)]
* serverinstall/k8s: Accept k8s-namespace when uninstalling server [[GH-1730](https://github.com/hashicorp/waypoint/issues/1730)]

## 0.4.0 (June 03, 2021)

FEATURES:

* **Mutable Deployments**: Waypoint now supports the concept of "mutable" deployments
where a deployment updates an existing resource rather than creating something
new. New plugins in this release include deploying from a Nomad job file,
Kubernetes apply from a directory or file, and more. [[GH-1298](https://github.com/hashicorp/waypoint/issues/1298)]
* **Status Reports**: Waypoint now supports a new feature for reporting
on the health of deployments or releases. Waypoint surfaces a deployment
and or releases status by relying on an existing platform for health checks.
A status is responsible for reporting the health of a deployed service by
representing its states as Ready, Alive, Partial, Down, or Unknown.
Platform health reporting lets teams take action quickly depending on the health
of their applications. Currently, the Kubernetes, Nomad, AWS ALB, and Docker built-in
plugins support the new Status reporting, with more support on the way! [[GH-1488](https://github.com/hashicorp/waypoint/issues/1488)]
* config: Add ability to define internal variables and compose variables together via templating [[GH-1382](https://github.com/hashicorp/waypoint/issues/1382)]
* config: Add ability to write configuration values as files rather than environment variables. [[GH-1395](https://github.com/hashicorp/waypoint/issues/1395)]
* config: `jsonnetfile` and `jsonnetdir` functions for processing Jsonnet files
and converting them to JSON. [[GH-1360](https://github.com/hashicorp/waypoint/issues/1360)]
* plugin/aws-alb: Report on status of releases [[GH-1567](https://github.com/hashicorp/waypoint/issues/1567)]
* plugin/docker: Add reporting on status of a deployment [[GH-1487](https://github.com/hashicorp/waypoint/issues/1487)]
* plugin/k8s: A new plugin "kubernetes-apply" that is able to deploy any typical
directory of YAML or JSON files to Kubernetes [[GH-1357](https://github.com/hashicorp/waypoint/issues/1357)]
* plugin/k8s: Add reporting on status of a deployment and release [[GH-1547](https://github.com/hashicorp/waypoint/issues/1547)]
* plugin/nomad: A new plugin "nomad-jobspec" for deploying a Nomad job specification directly. [[GH-1299](https://github.com/hashicorp/waypoint/issues/1299)]
* plugin/nomad: Add reporting on status of a deployment [[GH-1554](https://github.com/hashicorp/waypoint/issues/1554)]
* server/ecs: Use `--platform=ecs` to install waypoint server to
AWS ECS using Fargate [[GH-1564](https://github.com/hashicorp/waypoint/issues/1564)]
* ui: Add reporting on status of a deployment [[GH-1559](https://github.com/hashicorp/waypoint/issues/1559)]
* ui: Mutable deployments support on web dashboard. [[GH-1549](https://github.com/hashicorp/waypoint/issues/1549)]

IMPROVEMENTS:

* internal/config: add `${workspace}` variable [[GH-1419](https://github.com/hashicorp/waypoint/issues/1419)]
* plugin/pack: Support for non-default process types [[GH-1475](https://github.com/hashicorp/waypoint/issues/1475)]
* plugins/docker: Add support build context [[GH-1490](https://github.com/hashicorp/waypoint/issues/1490)]

BUG FIXES:

* plugin/k8s: destroy deployment on error [[GH-1528](https://github.com/hashicorp/waypoint/issues/1528)]
* plugin/pack: Upgrade pack package to fix downloading remote buildpacks [[GH-1452](https://github.com/hashicorp/waypoint/issues/1452)]
* server: Fix a bug that sometimes sends duplicate cancellation messages [[GH-1538](https://github.com/hashicorp/waypoint/issues/1538)]
* server: fix order of log lines when showing logs from multiple instances [[GH-1441](https://github.com/hashicorp/waypoint/issues/1441)]
* ui/checkbox-inputs-safari: Custom Inputs in the browser Ui now render properly on all supported browsers [[GH-1312](https://github.com/hashicorp/waypoint/issues/1312)]
* ui: unread log count resets after scrolling [[GH-1373](https://github.com/hashicorp/waypoint/issues/1373)]

BREAKING CHANGES:

* plugin/netlify: Removed the netlify plugin [[GH-1525](https://github.com/hashicorp/waypoint/issues/1525)]

## 0.3.1 (April 20, 2021)

IMPROVEMENTS:

* cli: Make `purge` default and remove flag for Nomad uninstall [[GH-1326](https://github.com/hashicorp/waypoint/issues/1326)]
* cli: Show usage example on `waypoint context use` command [[GH-1325](https://github.com/hashicorp/waypoint/issues/1325)]
* cli: version command now shows the server version [[GH-1364](https://github.com/hashicorp/waypoint/issues/1364)]
* entrypoint: can change log level using the `WAYPOINT_LOG_LEVEL` env var, which can also be set with `waypoint config` [[GH-1330](https://github.com/hashicorp/waypoint/issues/1330)]
* entrypoint: default log level changed to DEBUG [[GH-1330](https://github.com/hashicorp/waypoint/issues/1330)]
* plugin/nomad: Add CPU and Memory resource options for server and runner installs, and app deploys [[GH-1318](https://github.com/hashicorp/waypoint/issues/1318)]
* plugin/nomad: Allow for auth soft fail on serverinstall for server image [[GH-1106](https://github.com/hashicorp/waypoint/issues/1106)]
* ui: Improve the design of the Project Settings forms [[GH-1335](https://github.com/hashicorp/waypoint/issues/1335)]

BUG FIXES:

* cli: connections with TLS without insecure flag properly connect [[GH-1307](https://github.com/hashicorp/waypoint/issues/1307)]
* cli: server bootstrap will not give auth token errors [[GH-1320](https://github.com/hashicorp/waypoint/issues/1320)]
* plugin/aws/ecs: Route 53 "A" Type record properly created when not found for domain name [[GH-1256](https://github.com/hashicorp/waypoint/issues/1256)]
* plugin/nomad: use namespace config option for deploy [[GH-1300](https://github.com/hashicorp/waypoint/issues/1300)]

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
