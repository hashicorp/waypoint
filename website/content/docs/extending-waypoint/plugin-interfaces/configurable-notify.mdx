---
layout: docs
page_title: 'ConfigurableNotify'
description: |-
  How to implement the ConfigurableNotify component for a Waypoint plugin
---

# ConfigurableNotify

https://pkg.go.dev/github.com/hashicorp/waypoint-plugin-sdk/component#ConfigurableNotify

`ConfigurableNotify` is an optional interface you can implement to receive a call back after the configuration
has been decoded by the Waypoint SDK. It has a single input parameter which is the configuration reference you
return from the `Config` method. Returning an error from `ConfigSet` would stop execution of the Waypoint operation.

```go
type ConfigurableNotify interface {
  Configurable

  // ConfigSet is called with the value of the configuration after
  // decoding is complete successfully.
  ConfigSet(interface{}) error
}
```

`ConfigSet` can be used to validate configuration before it is used, the following example shows an implementation of
ConfigurableNotify which does just that.

```go
func (b *Builder) ConfigurableNotify(config interface{}) error {
  c, ok := config.(*BuildConfig)
  if !ok {
    return fmt.Errorf("Expected type BuildConfig")
  }

  // validate the config
  _, err := os.Stat(c.Source)
  if err != nil {
    return fmt.Errorf("Source folder does not exist")
  }

  // config validated ok
  return nil
}
```
