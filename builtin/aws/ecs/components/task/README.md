<!-- This file was generated via `make gen/integrations-hcl` -->
Launch an ECS task for on-demand tasks from the Waypoint server.

This will use the standard AWS environment variables and IAM Role information to
source authentication information for AWS, using the configured task role.
If no task role name is specified, Waypoint will create one with the required
permissions.

### Interface

