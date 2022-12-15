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
