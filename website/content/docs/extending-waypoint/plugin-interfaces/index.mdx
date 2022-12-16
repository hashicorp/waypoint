---
layout: docs
page_title: 'Plugin Interfaces'
description: |-
  The various component interfaces for a Waypoint plugin
---

# Plugin Interfaces and Components

A Waypoint plugin is a binary which implements one or more Waypoint components, each of which are related
to a different part of the lifecycle. There are 6 different components which can be implemented; these are shown below along
with the Waypoint command that triggers them.

![Plugin Components](/img/extending-waypoint/components.png)

# Implementing Components

To extend a particular part of the Waypoint application lifecycle, you create a component which is a Go
struct which implements the correct component interface.

For example, if you want to create a plugin that responds to build commands, create a
component that implements the Builder interface.

```go
type Builder interface {
  BuildFunc() interface{}
}
```

The Builder interface has a single method `BuildFunc` which has the return type of an interface. All of the plugin interfaces
in the Waypoint SDK are not called directly; instead, they require that you return a function. The following example shows
how the Builder interface could be implemented on a component.

```go
type Builder struct {
  // Other component fields
}

func (b *Builder) BuildFunc() interface{} {
  return b.Build
}

func (b *Builder) Build(
  ctx context.Context,
  log hclog.Logger,
  ui terminal.UI,
) (*Binary, error) {
  return nil, nil
}
```

There is no specific signature for the function you return from `BuildFunc`. The Waypoint SDK automatically injects the
specified parameters. In the previous example, the signature defines three input parameters and two return parameters.
As a plugin author, you determine which parameters you want the Waypoint SDK to inject for you. These are made up of
the Default Parameters and the custom Output Value that are returned from other components.

The output parameters are more strict and differ from interface to interface.
In this example for `BuildFunc`, you are required to return Output Values as a Go struct serializable to a Protocol
Buffer binary object and an `error`. The Output Value you return from one lifecycle function can be used
by the next in the chain. In the instance that `error` is not `nil`, the plugin execution will stop and the error will be
returned to the user.

More details on `Output Values`, `Default Parameters` and the specific details for each interface component can be found
in the respective documentation.
