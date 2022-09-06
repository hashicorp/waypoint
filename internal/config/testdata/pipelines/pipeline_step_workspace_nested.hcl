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

  step "pipe_normal" {
    pipeline "nested" {
      step "test_nested" {
        image_url = "example.com/test"

        use "exec" {
          command = "nested"
        }
      }
    }
  }
  
  step "pipe_changed" {
    workspace = "testworkspace"
    pipeline "nested" {
      step "test_nested" {
        image_url = "example.com/test"

        use "exec" {
          command = "nested"
        }
      }
      step "test_nested_dontoverride" {
        workspace = "dontoverride"
        image_url = "example.com/test"

        use "exec" {
          command = "nested"
        }
      }
      step "test_nested_overridejk" {
        image_url = "example.com/test"

        use "exec" {
          command = "nested"
        }
      }
    }
  }
  
  step "normal" {
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
