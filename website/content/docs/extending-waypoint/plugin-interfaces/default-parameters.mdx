---
layout: docs
page_title: 'Default Parameters'
description: |-
  Default parameters for Waypoint plugin component methods
---

# Default Parameters

In addition to user defined return values you can also request waypoint injects default parameters.
At present there are 11 possible parameters which you can inject, these are:

- context.Context
- \*component.Source
- \*component.JobInfo
- \*component.DeploymentConfig
- hclog.Logger
- terminal.UI
- \*component.LabelSet

## context.Context

https://pkg.go.dev/context#Context

Context is a default parameter and is used by Waypoint to notify you when the work that component is
performing should be cancelled. As with any Go code calling the `Err()` function on the Context will return
an error if the Context has already been cancelled, you can also use the channel returned by the `Done()` function
to notify you when Waypoint cancels the current operation.

```go
if ctx.Err() != nil {
  return nil, status.Error(codes.Internal, "Context cancelled error message")
}
```

## \*component.Source

https://pkg.go.dev/github.com/hashicorp/waypoint-plugin-sdk/component#Source

Source is a data struct which provides the current application name `App`, and `Path` which is the directory
containing the current Waypoint file.

```go
type Source struct {
  App  string
  Path string
}
```

## \*component.JobInfo

https://pkg.go.dev/github.com/hashicorp/waypoint-plugin-sdk/component#JobInfo

JobInfo is a data struct which provides details related to the current Job.

```go
type JobInfo struct {
  // Id is the ID of the job that is executing this plugin operation.
  // If this is empty then it means that the execution is happening
  // outside of a job.
  Id string

  // Local is true if the operation is running locally on a machine
  // alongside the invocation. This can be used to determine if you can
  // do things such as open browser windows, read user files, etc.
  Local bool

  // Workspace is the workspace that this job is executing in. This should
  // be used by plugins to properly isolate resources from each other.
  Workspace string
}
```

## \*component.DeploymentConfig

https://pkg.go.dev/github.com/hashicorp/waypoint-plugin-sdk/component#DeploymentConfig

DeploymentConfig contains configuration which the entrypoint binary requires in order to run correctly.
For example if you deploy an application to Google Cloud Run which is using the builtin Docker builder,
the entrypoint is automatically bundled into the container. In order for the entrypoint to function correctly
it needs to be configured correctly.

```go
type DeploymentConfig struct {
  Id                    string
  ServerAddr            string
  ServerTls             bool
  ServerTlsSkipVerify   bool
  EntrypointInviteToken string
}
```

To simplify the configuration of the entrypoint environment variables DeploymentConfig also has a function which
returns a map of the correct environment variables and their values.

```go
func (c *DeploymentConfig) Env() map[string]string
```

## hclog.Logger

https://pkg.go.dev/github.com/hashicorp/go-hclog#Logger

Logger provides you with the ability to write to the Waypoint standard logger using different log levels such as
Info, Debug, Trace, etc. By default, log output is disabled when running a waypoint command. The environment variable
`WAYPOINT_LOG_LEVEL=[debug, trace, info, warning]` can be set to enable log output in the CLI.

```go
  log.Info(
    "Start build",
    "src", src,
    "job", job,
    "projectDataDir", projectDir.DataDir(),
    "projectCacheDir", projectDir.CacheDir(),
    "appDataDir", appDir.DataDir(),
    "appCacheDir", appDir.CacheDir(),
    "componentDataDir", componentDir.DataDir(),
    "componentCacheDir", componentDir.CacheDir(),
    "labels", labels,
  )
```

## terminal.UI

https://pkg.go.dev/github.com/hashicorp/waypoint-plugin-sdk/terminal#UI

UI allows you to build rich output for your plugin, more details on using terminal.UI can be found in the
Interacting with the UI section.

## \*component.LabelSet

https://pkg.go.dev/github.com/hashicorp/waypoint-plugin-sdk/component#LabelSet

LabelSet allows you to read the labels defined at an app level on your Waypoint configuration.

```go
app "wpmini" {
  labels = {
    "service" = "wpmini",
    "env"     = "dev"
  }
#...
```

The LabelSet struct exposes a single field Labels which is of type `map[string]string`, this collection contains
all the labels defined in the configuration.

```go
Labels map[string]string
```
