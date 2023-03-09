<!-- This file was generated via `make gen/integrations-hcl` -->
Release target groups by attaching them to an ALB.

### Interface

- Input: **alb.TargetGroup**
- Output: **alb.Release**

### Mappers

#### Allow EC2 Deployments to be hooked up to an ALB

- Input: **ec2.Deployment**
- Output: **alb.TargetGroup**

#### Allow Lambda Deployments to be hooked up to an ALB

- Input: **lambda.Deployment**
- Output: **alb.TargetGroup**

