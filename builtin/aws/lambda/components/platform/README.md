## aws-lambda (platform)

Deploy functions as OCI Images to AWS Lambda.

### Interface

- Input: **ecr.Image**
- Output: **lambda.Deployment**

### Examples

```hcl
deploy {
	use "aws-lambda" {
		region = "us-east-1"
		memory = 512
	}
}
```
