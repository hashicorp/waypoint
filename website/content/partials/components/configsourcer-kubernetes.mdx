## kubernetes (configsourcer)

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

### Required Parameters

These parameters are used in `dynamic` for sourcing [configuration values](/docs/app-config/dynamic) or [input variable values](/docs/waypoint-hcl/variables/dynamic).

#### key

The key in the ConfigMap or Secret to read the value from.

ConfigMaps and Secrets store data in key/value format. This specifies the key to read from the resource. If you want multiple values you must specify multiple dynamic values.

- Type: **string**

#### name

The name of the ConfigMap of Secret.

- Type: **string**

### Optional Parameters

These parameters are used in `dynamic` for sourcing [configuration values](/docs/app-config/dynamic) or [input variable values](/docs/waypoint-hcl/variables/dynamic).

#### namespace

The namespace to load the ConfigMap or Secret from.

By default this will use the namespace of the running pod. If this config source is used outside of a pod, this will use the namespace from the kubeconfig.

- Type: **string**
- **Optional**

#### secret

This must be set to true to read from a Secret. If it is false we read from a ConfigMap.

- Type: **bool**
- **Optional**
