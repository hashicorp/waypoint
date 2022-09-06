project = "foo"

pipeline "foo" {
  step "test" {
    image_url = "example.com/test"

    use "exec" {
      command = "bar"
    }
  }
  
  step "testws" {
    image_url = "example.com/test"
    workspace = "testws"
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
