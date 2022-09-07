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

  step "pipe_nested" {
    pipeline "nested" {
      step "test_nested" {
        image_url = "example.com/test"

        use "exec" {
          command = "nested"
        }
      }
    }
  }
  
  step "pipe_nested_workspace" {
    workspace = "testworkspace"
    pipeline "nested_workspace" {
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
      step "test_nested_override" {
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
