# This file was generated via `make gen/integrations-hcl`
parameter {
  key         = "count"
  description = "how many EC2 instances to configure the ASG with\nthe fields here (desired, min, max) map directly to the typical ASG configuration"
  type        = "ec2.countConfig"
  required    = true
}

parameter {
  key         = "extra_ports"
  description = "additional TCP ports to allow into the EC2 instances\nthese additional ports are usually used to allow secondary services, such as ssh"
  type        = "list of int"
  required    = false
}

parameter {
  key         = "instance_type"
  description = "the EC2 instance type to deploy"
  type        = "string"
  required    = true
}

parameter {
  key         = "key"
  description = "the name of an SSH Key to associate with the instances, as preconfigured in EC2"
  type        = "string"
  required    = false
}

parameter {
  key         = "region"
  description = "the AWS region to deploy into"
  type        = "string"
  required    = true
}

parameter {
  key         = "security_groups"
  description = "additional security groups to attached to the EC2 instances\nthis plugin creates security groups that match the above ports by default. this field allows additional security groups to be specified for the instances"
  type        = "list of string"
  required    = false
}

parameter {
  key         = "service_port"
  description = "the TCP port on the instances that the app will be running on"
  type        = "int"
  required    = true
}

parameter {
  key           = "subnet"
  description   = "the subnet to place the instances into"
  type          = "string"
  required      = false
  default_value = "a public subnet in the dafault VPC"
}

