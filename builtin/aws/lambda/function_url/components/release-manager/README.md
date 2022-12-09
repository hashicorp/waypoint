## lambda-function-url (releasemanager)

Create an AWS Lambda function URL.

### Interface

- Input: **lambda.Deployment**
- Output: **lambda.Release**

### Examples

```hcl
release {
	use "lambda-function-url" {
		auth_type = "NONE"
	}
}
```

### Required Parameters

This plugin has no required parameters.

### Optional Parameters

These parameters are used in the [`use` stanza](/docs/waypoint-hcl/use) for this plugin.

#### auth_type

The Lambda function URL auth type.

The AuthType parameter determines how Lambda authenticates or authorizes requests to your function URL. Must be either `AWS_IAM` or `NONE`.

- Type: **string**
- **Optional**
- Default: NONE

#### principal

The principal to use when auth_type is `AWS_IAM`.

The Principal parameter specifies the principal that is allowed to invoke the function.

- Type: **string**
- **Optional**
- Default: \*
