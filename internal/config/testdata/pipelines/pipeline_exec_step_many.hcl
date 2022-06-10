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
      args = ["foo", "fighters"]
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
