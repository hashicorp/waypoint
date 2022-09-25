---
layout: docs
page_title: 'Status'
description: |-
  How to implement the Status component for a Waypoint plugin
---

# Status

https://pkg.go.dev/github.com/hashicorp/waypoint-plugin-sdk/component#Status

The `Status` interface is responsible for reporting on the current health
for resources which have been created by the waypoint `deploy`
and `release` phase. The intended use is to leverage the platforms existing
health check features to determine the overall health.

Status can only be implemented in `Platform` and `ReleaseManager` components and is implemented through the following interface.

```go
type Status interface {
  // StatusReportFunc should return a proto.StatusReport that details the
  // result of the most recent health check for a deployment.
  StatusFunc() interface{}
}
```

`StatusFunc` requires that you return a StatusReport proto message describing
the reported health of the deployment or release. If your plugin pulls a report
from an external platform, you'll need to mark the report as `External`.

Each report has a top level `Health` which should be the health of the deployment
or release your plugin created. Each report can also describe the resources
involved and their reported health.

https://pkg.go.dev/github.com/hashicorp/waypoint-plugin-sdk/proto/gen#StatusReport

```go
type StatusReport struct {

  // a collection of resources for a deployed application
  Resources []*StatusReport_Resource `protobuf:"bytes,1,rep,name=resources,proto3" json:"resources,omitempty"`
  // the current overall health state for a deployment
  Health StatusReport_Health `protobuf:"varint,2,opt,name=health,proto3,enum=hashicorp.waypoint.sdk.StatusReport_Health" json:"health,omitempty"`
  // a simple human readable message detailing the Health state
  HealthMessage string `protobuf:"bytes,3,opt,name=health_message,json=healthMessage,proto3" json:"health_message,omitempty"`
  // the time when this report was generated
  GeneratedTime *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=generated_time,json=generatedTime,proto3" json:"generated_time,omitempty"`
  // where the health check was performed. External means not executed by Waypoint,
  // but by the platform deployed to.
  External bool `protobuf:"varint,5,opt,name=external,proto3" json:"external,omitempty"`
  // contains filtered or unexported fields
}
```

When picking what kind of `Health` your deployment or release has, a Waypoint
status report has a few options to fit its status into:

```go
const (
  StatusReport_UNKNOWN StatusReport_Health = 0 // We cannot determine the status of the resource
  StatusReport_ALIVE   StatusReport_Health = 1 // The resource exists but might not be ready to handle requests
  StatusReport_READY   StatusReport_Health = 2 // The resource exists and is ready to handle requests
  StatusReport_DOWN    StatusReport_Health = 3 // The resource is down and not responding
  StatusReport_MISSING StatusReport_Health = 5 // We're expecting it to exist, but it does not.
  StatusReport_PARTIAL StatusReport_Health = 4 // Some resources in deployment are OK, others are not OK
)
```

When adding a `StatusFunc` to your plugin, you might define one like the following
example in the `waypoint-plugin-examples` repo with the `filepath` plugin. In
this case, the `filepath` plugin creates a file when Waypoint runs a deployment.
The example status check below is not an `External` check, because it did
the health check itself rather than pulling from an external report
(such as pod status health in Kubernetes).

```go
// You might have other imports, this is for a mostly complete example
import (
  "context"
  "os"
  "path/filepath"

  "github.com/hashicorp/waypoint-plugin-sdk/component"
  sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
  "github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

// The StatusFunc here tells the Waypoint plugin SDK which function in this
// plugin to invoke when Waypoint runs through its Status check
func (d *Deploy) StatusFunc() interface{} {
  return d.status
}

func (d *Deploy) status(
  ctx context.Context,
  ji *component.JobInfo,
  deploy *Deployment,
  ui terminal.UI,
) (*sdk.StatusReport, error) {
  sg := ui.StepGroup()
  s := sg.Add("Checking the status of the file...")

  report := &sdk.StatusReport{}
  report.External = false

  if _, err := os.Stat(deploy.Path); err == nil {
    s.Update("File is ready")
    report.Health = sdk.StatusReport_READY
  } else {
    st := ui.Status()
    defer st.Close()
    st.Step(terminal.StatusError, "File is missing")
    s.Status(terminal.StatusError)

    report.Health = sdk.StatusReport_MISSING
  }
  s.Done()

  return report, nil
}
```
