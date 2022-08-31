project = "foo"

pipeline "foo" {
  step "test" {
    image_url = "example.com/test"

    use "exec" {
      command = "bar"
    }
  }

  step "pipe" {
    pipeline "nested" {
      step "test_nested" {
        image_url = "example.com/test"

        use "exec" {
          command = "nested"
        }
      }
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
