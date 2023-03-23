# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

project = "foo"

pipeline "foo" {
  step "zero" {
    image_url = "example.com/test"

    use "test" {
      foo = "qubit"
    }
  }

  step "one" {
    image_url = "example.com/second"

    use "test/exec" {
      foo = "few"
      bar = "bar"
    }
  }

  step "two" {
    image_url = "example.com/different"
    depends_on = ["zero"]

    use "hunger" {
      foo = "food"
      bar = "drink"
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
