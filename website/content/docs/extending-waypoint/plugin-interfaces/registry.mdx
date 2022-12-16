---
layout: docs
page_title: 'Registry'
description: |-
  How to implement the Registry component for a Waypoint plugin
---

# Registry

https://pkg.go.dev/github.com/hashicorp/waypoint-plugin-sdk/component#Registry

The registry component handles the storage of the built assets in an artifact registry such as Docker registry,
GitHub releases, or Artifactory.

![Authenticator](/img/extending-waypoint/build.png)

To build a plugin which allows the storage of assets you need to implement the `PushFunc` method from the
`Registry` interface in your component.

```go
type Registry interface {
  PushFunc() interface{}
}
```

The signature for the function returned by PushFunc accepts the usual Default Mappers in addition, the data
model which was returned as the first output parameter from BuildFunc can also be specified. The output parameters
from this function are a data model which contains the details of the published artifact and an error message.

```go
type Registry struct {
  // Other component fields
}

func (r *Registry) PushFunc() interface{} {
  return r.push
}

func (r *Registry) push(
  ctx context.Context,
  img *builder.Binary,
  ui terminal.UI,
) (*Artifact, error) {
```

The data model Artifact shown in the previous example will be made available to be injected into a function in a
later stage of the Waypoint lifecycle.
