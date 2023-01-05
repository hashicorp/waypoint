<!-- This file was generated via `make gen/integrations-hcl` -->
Copies files to a specific directory.

### Interface

### Examples

```hcl
build {
  use "files" {}
  registry {
	use "files" {
	  path = "/path/to/file"
	}
  }
}
```

