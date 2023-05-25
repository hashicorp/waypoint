project = "foo"

pipeline "foo" {
  step "test" {
    image_url = "example.com/test"

    use "test" {
      foo = "bar"
    }
  }
}


pipeline "bar" {
  step "test2" {
    image_url = "example.com/test"

    use "test" {
      foo = "bar"
    }
  }
}

pipeline "foofoo" {
  step "test2" {
    image_url = "example.com/test"

    use "test" {
      foo = "bar"
    }
  }
}

pipeline "barbar" {
  step "test2" {
    image_url = "example.com/test"

    use "test" {
      foo = "bar"
    }
  }
}

pipeline "foobar" {
  step "test2" {
    image_url = "example.com/test"

    use "test" {
      foo = "bar"
    }
  }
}

pipeline "naming-is-hard" {
  step "test2" {
    image_url = "example.com/test"

    use "test" {
      foo = "bar"
    }
  }
}

pipeline "hey-we-made-it" {
  step "test2" {
    image_url = "example.com/test"

    use "test" {
      foo = "bar"
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
