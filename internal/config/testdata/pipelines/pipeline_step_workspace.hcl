# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

project = "foo"

pipeline "foo" {
  step "test" {
    image_url = "example.com/test"

    use "exec" {
      command = "bar"
    }
  }
  
  step "testworkspace" {
    image_url = "example.com/test"
    workspace = "testworkspace"
    use "exec" {
      command = "bar"
    }
  }
  
  step "othertest" {
    image_url = "example.com/test"
    depends_on = ["test"]
    use "exec" {
      command = "bar"
    }
  }
}

app "web" {
    config {
        env = {
            static = "hello"
        }
    }

    build {}

    deploy {}
}
