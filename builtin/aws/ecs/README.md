## AWS ECS

The AWS ECS plugin deploys an application image to an AWS [ECS cluster](https://aws.amazon.com/ecs/getting-started/).
It also launches on-demand runners to do operations remotely.

### Components

1. [Platform](/waypoint/integrations/aws-ecs/latest/components/platform)
2. [TaskLauncher](/waypoint/integrations/aws-ecs/latest/components/task)

### Related Plugins

1. [Docker](/waypoint/integrations/docker)
2. [AWS ECR](/waypoint/integrations/aws-ecr)

### Resources

1. ECS Cluster
2. IAM Execution Role 
3. IAM Task Role
4. Internal Security Group
5. External Security Group
6. Log Group
7. Service Subnets
8. ALB subnets
9. Target Group
10. ALB
11. ALB listener
12. Route53 Record
13. Task Definition
14. Service