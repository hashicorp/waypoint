## null (configsourcer)

Simple configuration values for experimentation or testing.

### Examples

```hcl
config {
  env = {
    "STATIC" = configdynamic("null", {
      static_value = "hello"
    })

    "FROM_CONFIG" = configdynamic("null", {
      config_key = "foo"
    })
  }
}
```

### Source Parameters

The parameters below are used with `waypoint config source-set` to configure
the behavior this plugin. These are _not_ used in `configdynamic` calls. The
parameters used for `configdynamic` are in the previous section.

#### Required Source Parameters

This plugin has no required source parameters.

#### Optional Source Parameters

##### values

A mapping of key to value of values that can be pulled with `config_key`.

These values can be sourced using the `config_key` attribute as as `configdynamic` argument. See the `config_key` documentation for more information on why this is useful.

- Type: **map of string to string**
- **Optional**
