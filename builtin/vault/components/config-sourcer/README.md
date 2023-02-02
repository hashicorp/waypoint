<!-- This file was generated via `make gen/integrations-hcl` -->
Read configuration values from Vault.

### Examples

```hcl
# Setting an input variable dynamically with Vault
variable "my_api_key" {
  default = dynamic("vault", {
    path = "secret/data/keys"
    key  = "/data/my_api_key"
  })
  type        = string
  sensitive   = true
  description = "my api key from vault"
}

# Setting a dynamic variable for an environment variable
config {
  env = {
    "DATABASE_USERNAME" = dynamic("vault", {
      path = "database/creds/my-role"
      key = "username"
    })

    "DATABASE_PASSWORD" = dynamic("vault", {
      path = "database/creds/my-role"
      key = "password"
    })

    # KV Version 2
    "PASSWORD_FOO" = dynamic("vault", {
      path = "secret/data/my-secret"
      key = "/data/password"  # key must be prefixed with "/data" (see below)
    })

    # KV Version 1
    "PASSWORD_BAR" = dynamic("vault", {
      path = "kv1/my-secret"
      key = "password"
    })
  }
}
```

