## aws-ec2 (platform)

Deploy the application into an AutoScaling Group on EC2.

### Interface

- Input: **ami.Image**
- Output: **ec2.Deployment**

### Required Parameters

These parameters are used in the [`use` stanza](/docs/waypoint-hcl/use) for this plugin.

#### count

How many EC2 instances to configure the ASG with.

The fields here (desired, min, max) map directly to the typical ASG configuration.

- Type: **ec2.countConfig**

#### instance_type

The EC2 instance type to deploy.

- Type: **string**

#### region

The AWS region to deploy into.

- Type: **string**

#### service_port

The TCP port on the instances that the app will be running on.

- Type: **int**

### Optional Parameters

These parameters are used in the [`use` stanza](/docs/waypoint-hcl/use) for this plugin.

#### extra_ports

Additional TCP ports to allow into the EC2 instances.

These additional ports are usually used to allow secondary services, such as ssh.

- Type: **list of int**
- **Optional**

#### key

The name of an SSH Key to associate with the instances, as preconfigured in EC2.

- Type: **string**
- **Optional**

#### security_groups

Additional security groups to attached to the EC2 instances.

This plugin creates security groups that match the above ports by default. this field allows additional security groups to be specified for the instances.

- Type: **list of string**
- **Optional**

#### subnet

The subnet to place the instances into.

- Type: **string**
- **Optional**
- Default: a public subnet in the dafault VPC

### Output Attributes

Output attributes can be used in your `waypoint.hcl` as [variables](/docs/waypoint-hcl/variables) via [`artifact`](/docs/waypoint-hcl/variables/artifact) or [`deploy`](/docs/waypoint-hcl/variables/deploy).

#### public_dns

- Type: **string**

#### public_ip

- Type: **string**

#### region

- Type: **string**

#### service_name

- Type: **string**

#### target_group_arn

- Type: **string**
