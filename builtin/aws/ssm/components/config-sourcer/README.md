# /aws/ssm/components/config-sourcera## aws-ssm (configsourcer)

Read configuration values from AWS SSM Parameter Store.

### Examples

```hcl
config {
  env = {
    PORT = dynamic("aws-ssm", {
      path = "port"
    })
  }
}
```
