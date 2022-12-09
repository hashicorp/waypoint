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

### Required Parameters

These parameters are used in the [`use` stanza](/docs/waypoint-hcl/use) for this plugin.

#### region

The AWS region for the ECS cluster.

- Type: **string**

### Optional Parameters

These parameters are used in the [`use` stanza](/docs/waypoint-hcl/use) for this plugin.

#### architecture

The instruction set architecture that the function supports. Valid values are: "x86_64", "arm64".

- Type: **string**
- **Optional**
- Default: x86_64

#### iam_role

An IAM Role specified by ARN that will be used by the Lambda at execution time.

- Type: **string**
- **Optional**
- Default: created automatically

#### memory

The amount of memory, in megabytes, to assign the function.

- Type: **int**
- **Optional**
- Default: 256

#### static_environment

Environment variables to expose to the lambda function.

Environment variables that are meant to configure the application in a static way. This might be to control an image that has multiple modes of operation, selected via environment variable. Most configuration should use the waypoint config commands.

- Type: **map of string to string**
- **Optional**

#### storagemb

The storage size (in MB) of the Lambda function's `/tmp` directory. Must be a value between 512 and 10240.

- Type: **int**
- **Optional**
- Default: 512

#### timeout

The number of seconds a function has to return a result.

- Type: **int**
- **Optional**
- Default: 60

### Output Attributes

Output attributes can be used in your `waypoint.hcl` as [variables](/docs/waypoint-hcl/variables) via [`artifact`](/docs/waypoint-hcl/variables/artifact) or [`deploy`](/docs/waypoint-hcl/variables/deploy).

#### func_arn

- Type: **string**

#### id

- Type: **string**

#### region

- Type: **string**

#### target_group_arn

- Type: **string**

#### ver_arn

- Type: **string**

#### version

- Type: **string**
