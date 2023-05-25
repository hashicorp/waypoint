project = "foo"

pipeline "foo" {
  step "test" {
    image_url = "example.com/test"

    use "exec" {
      command = "bar"
    }
  }
}


pipeline "bar" {
  step "test2" {
    image_url = "example.com/test"

    use "exec" {
      command = "bar"
    }
  }
}

pipeline "foofoo" {
  step "test2" {
    image_url = "example.com/test"

    use "exec" {
      command = "bar"
    }
  }
}

pipeline "barbar" {
  step "test2" {
    image_url = "example.com/test"

    use "exec" {
      command = "bar"
    }
  }
}

pipeline "foobar" {
  step "test2" {
    image_url = "example.com/test"

    use "exec" {
      command = "bar"
    }
  }
}

pipeline "naming-is-hard" {
  step "test2" {
    image_url = "example.com/test"

    use "exec" {
      command = "bar"
    }
  }
}

pipeline "hey-we-made-it" {
  step "test2" {
    image_url = "example.com/test"

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
