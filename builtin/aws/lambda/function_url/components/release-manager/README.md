<!-- This file was generated via `make gen/integrations-hcl` -->
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

