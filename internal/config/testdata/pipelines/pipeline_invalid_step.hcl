project = "foo"

pipeline "foo" {
  step "test" {
    image_url = "example.com/test"

    use "invalid" {
      test = "bar"
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
