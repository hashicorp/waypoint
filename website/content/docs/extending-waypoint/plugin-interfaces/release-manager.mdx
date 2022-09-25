---
layout: docs
page_title: 'ReleaseManager'
description: |-
  How to implement the ReleaseManager component for a Waypoint plugin
---

# ReleaseManager

https://pkg.go.dev/github.com/hashicorp/waypoint-plugin-sdk/component#ReleaseManager

The ReleaseManager component is responsible for taking a deployment and making it active, this could be as
simple as exposing it using a public load balancer or it may be a gradual and phased canary deployment.

![Release Manager](/img/extending-waypoint/release-manager.png)

To create a ReleaseManager component you implement the ReleaseManager interface in your component.

```go
type ReleaseManager interface {
  // ReleaseFunc should return the method handle for the "release" operation.
  ReleaseFunc() interface{}
}
```

`ReleaseManager` has a single method which you must implement which returns a function called by Waypoint.
The signature for the function returned by your implementation of ReleaseFunc can accept all the standard parameters
and in addition you can request the output value from the Deployment component. Return parameters for the function
are a Waypoint value and an error.

```go
type Releaser struct {
  // Other component fields
}

func (r *Releaser) ReleaseFunc() interface{} {
	return r.Release
}

func (r *Releaser) Release(
  ctx context.Context,
  log hclog.Logger,
  ui terminal.UI,
  target *Deployment,
) (*Release, error)
```
