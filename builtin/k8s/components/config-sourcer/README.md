<!-- This file was generated via `make gen/integrations-hcl` -->
Read configuration values from Kubernetes ConfigMap or Secret resources. Note that to read a config value from a Secret, you must set `secret = true`. Otherwise Waypoint will load a dynamic value from a ConfigMap.

### Examples

```hcl
config {
  env = {
    PORT = dynamic("kubernetes", {
	  name = "my-config-map"
	  key = "port"
	})

    DATABASE_PASSWORD = dynamic("kubernetes", {
	  name = "database-creds"
	  key = "password"
	  secret = true
	})
  }
}
```

