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

### Required Parameters

These parameters are used in `dynamic` for sourcing [configuration values](/docs/app-config/dynamic) or [input variable values](/docs/waypoint-hcl/variables/dynamic).

#### path

The path for the parameter to read from the parameter store.

- Type: **string**

### Optional Parameters

This plugin has no optional parameters.

### Source Parameters

The parameters below are used with `waypoint config source-set` to configure
the behavior this plugin. These are _not_ used in `dynamic` calls. The
parameters used for `dynamic` are in the previous section.

#### Required Source Parameters

This plugin has no required source parameters.

#### Optional Source Parameters

##### access_key

This is the AWS access key. It must be provided, but it can also be sourced from the `AWS_ACCESS_KEY_ID` environment variable, or via a shared credentials file if `profile` is specified.

- Type: **string**
- **Optional**

##### assume_role_arn

Amazon Resource Name (ARN) of the IAM Role to assume.

- Type: **string**
- **Optional**

##### assume_role_duration_seconds

Number of seconds to restrict the assume role session duration.

- Type: **int**
- **Optional**

##### assume_role_external_id

External identifier to use when assuming the role.

- Type: **string**
- **Optional**

##### assume_role_policy

IAM Policy JSON describing further restricting permissions for the IAM Role being assumed.

- Type: **string**
- **Optional**

##### assume_role_session_name

Session name to use when assuming the role.

- Type: **string**
- **Optional**

##### iam_endpoint

Custom endpoint address for the IAM service.

- Type: **string**
- **Optional**

##### insecure

Explicitly allow the provider to perform "insecure" SSL requests.

- Type: **bool**
- **Optional**
- Default: false

##### max_retries

This is the maximum number of times an API call is retried, in the case where requests are being throttled or experiencing transient failures. The delay between the subsequent API calls increases exponentially.

- Type: **int**
- **Optional**
- Default: 25

##### profile

This is the AWS profile name as set in the shared credentials file.

- Type: **string**
- **Optional**

##### region

This is the AWS region. It must be provided, but it can also be sourced from the `AWS_DEFAULT_REGION` environment variables, or via a shared credentials file if profile is specified.

- Type: **string**
- **Optional**

##### secret_key

This is the AWS secret key. It must be provided, but it can also be sourced from the `AWS_SECRET_ACCESS_KEY` environment variable, or via a shared credentials file if `profile` is specified.

- Type: **string**
- **Optional**

##### shared_credentials_file

This is the path to the shared credentials file. If this is not set and a profile is specified, `~/.aws/credentials` will be used.

- Type: **string**
- **Optional**

##### skip_credentials_validation

Skip the credentials validation via the STS API. Useful for AWS API implementations that do not have STS available or implemented.

- Type: **bool**
- **Optional**

##### skip_metadata_api_check

Skip the AWS Metadata API check. Useful for AWS API implementations that do not have a metadata API endpoint. Setting to true prevents Terraform from authenticating via the Metadata API. You may need to use other authentication methods like static credentials, configuration variables, or environment variables.

- Type: **bool**
- **Optional**

##### skip_requesting_account_id

Skip requesting the account ID. Useful for AWS API implementations that do not have the IAM, STS API, or metadata API.

- Type: **bool**
- **Optional**

##### sts_endpoint

Custom endpoint for the STS service.

- Type: **string**
- **Optional**

##### token

- Type: **string**
- **Optional**
