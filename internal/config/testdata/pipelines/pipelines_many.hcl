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

app "web" {
    config {
        env = {
            static = "hello"
        }
    }

    build {}

    deploy {}
}
