---
layout: docs
page_title: 'Destroy'
description: |-
  How to implement the Destroy component for a Waypoint plugin
---

# Destroy

https://pkg.go.dev/github.com/hashicorp/waypoint-plugin-sdk/component#Destroyer

The `Destroyer` interface is responsible for removing any resources which have been created by the waypoint `deploy`
and `release` phase.

![Platform](/img/extending-waypoint/destroy.png)

Destroy can only be implemented in `Platform` and `ReleaseManager` components and is implemented through the following interface.

```go
type Destroyer interface {
  // DestroyFunc should return the method handle for the destroy operation.
  DestroyFunc() interface{}
}
```

`DestroyFunc` requires you return a function, like other interface functions Waypoint will automatically inject any of the
standard parameters. In addition you can specify the output value returned from the `DeployFunc`, or from `ReleaseFunc`,
the details of which you can use to clean up any deployments. As shown in the example below, the function signature for
a `DestroyFunc` function only has a single output parameter which is an error used to signal if the destroy operation succeeded.

```go
type Deploy struct {
  // Other component fields
}

func (d *Deploy) DestroyFunc() interface{} {
  return d.Destroy
}

func (d *Deploy) Destroy(
  ctx context.Context,
  ui terminal.UI,
  deployment *Deployment,
) error {
  st := ui.Status()
  defer st.Close()

  err := os.RemoveAll(d.config.Directory)
  if err != nil {
    st.Step(terminal.ErrorStyle, fmt.Sprintf("Unable to remove deployments %s", d.config.Directory))
    return err
  }

  st.Step(terminal.StatusOK, fmt.Sprintf("Removed deployments %s", d.config.Directory))
  return nil
}
```
