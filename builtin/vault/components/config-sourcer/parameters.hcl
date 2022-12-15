parameter {
  key         = "key"
  description = <<EOT
The key name that exists at the specified Vault `path` parameter.
The value can be a direct key such as `password` or it can be a [JSON pointer](https://tools.ietf.org/html/rfc6901) string to retrieve a nested value. When using the Vault KV [Version 2](https://www.vaultproject.io/docs/secrets/kv/kv-v2) secret backend, the key must be prefixed with an additional string of `/data`. For example, `/data/password`. When using the Vault KV [Version 1](https://www.vaultproject.io/docs/secrets/kv/kv-v1) secret backend, the key can be a direct key name such as `password`. This is because the Vault KV API returns different data structures in its response depending on the Vault KV version the key is stored in. Therefore, the `/data` prefix is required for keys stored in the Vault KV `Version 2` secret backend in order to retrieve its nested value using JSON pointer string.

EOT
  type        = "string"
  required    = true

}

parameter {
  key         = "path"
  description = <<EOT
The Vault path to read the secret.
Within a single application, multiple dynamic values that use the same path will only read the value once. This allows multiple keys from a single secret to be extracted into multiple values. The example above shows this functionality by reading the username and password into separate values. When using the Vault KV secret backend, the path is usually `<mount>/data/<key>`. For example, if you wrote data with `vault kv put secret/myapp` then the key for Waypoint must be `secret/data/myapp`. This can be confusing but is caused by the fact that the Vault API is what Waypoint uses and the Vault CLI does this automatically for KV.

EOT
  type        = "string"
  required    = true

}

