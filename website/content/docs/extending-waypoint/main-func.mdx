---
layout: docs
page_title: 'Initializing the SDK'
description: |-
  How Waypoint plugins work
---

# Initializing the SDK

A plugin is a Go binary containing one or more components. A single plugin may be responsible for Deploying the
application (Platform), archiving application artifacts (Registry), releasing it (ReleaseManager) or more. In fact, a
single plugin can be responsible for the entire Waypoint lifecycle.

Technically, plugins are [gRPC Go-Plugins](https://github.com/hashicorp/go-plugin); however, as a plugin developer, the
Waypoint SDK abstracts this complexity for you. To implement a plugin, you can use the `sdk.Main` function. You pass
Main the components you would like to register for the plugin. You do not need to specify the type of component when you
register a component; the Waypoint SDK inspects the interfaces you have added and automatically hooks it up to the
correct part of the life cycle

```go
func main() {
  sdk.Main(sdk.WithComponents(
    &registry.Registry{},
    &deploy.Deploy{},
    &release.Releaser{},
  ))
}
```
