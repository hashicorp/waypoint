<!-- This file was generated via `make gen/integrations-hcl` -->
Launch a Nomad job for on-demand tasks from the Waypoint server.

This will use the standard Nomad environment used for with the server install
to launch on demand Nomad jobs for Waypoint server tasks.

### Interface

### Examples

```hcl
task {
	use "nomad" {}
}
```

