## AWS ALB

The AWS ALB plugin releases applications deployed to AWS by attaching [target
groups](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-target-groups.html)
to an [ALB](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/introduction.html).

### Components

1. [ReleaseManager](./components/release-manager/README.md)

### Related Plugins

1. [AWS EC2](../ec2/README.md)
2. [AWS Lambda](../lambda/README.md)

### Resources

1. Security group
2. Load Balancer
3. Listener
4. Record Set