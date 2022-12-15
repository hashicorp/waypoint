# Required
parameter {
  key         = "count"
  description = <<EOT
How many EC2 instances to configure the ASG with.

The fields here (desired, min, max) map directly to the typical ASG configuration.
EOT
  type        = "ec2.countConfig"
  required    = true
}
parameter {
  key         = "instance_type"
  description = <<EOT
The EC2 instance type to deploy.
EOT
  type        = "string"
  required    = true
}
parameter {
  key         = "region"
  description = <<EOT
The AWS region to deploy into.
EOT
  type        = "string"
  required    = true
}
parameter {
  key         = "service_port"
  description = <<EOT
The TCP port on the instances that the app will be running on.
EOT
  type        = "int"
  required    = true
}

# Optional
parameter {
  key         = "extra_ports"
  description = <<EOT
Additional TCP ports to allow into the EC2 instances.

These additional ports are usually used to allow secondary services, such as ssh.
EOT
  type        = "list of int"
  required    = false
}
parameter {
  key         = "key"
  description = <<EOT
The name of an SSH Key to associate with the instances, as preconfigured in EC2.
EOT
  type        = "string"
  required    = false
}
parameter {
  key         = "security_groups"
  description = <<EOT
Additional security groups to attached to the EC2 instances.

This plugin creates security groups that match the above ports by default. this field allows additional security groups to be specified for the instances.
EOT
  type        = "list of string"
  required    = false
}
parameter {
  key           = "subnet"
  description   = <<EOT
The subnet to place the instances into.
EOT
  type          = "string"
  required      = false
  default_value = "a public subnet in the dafault VPC"
}
