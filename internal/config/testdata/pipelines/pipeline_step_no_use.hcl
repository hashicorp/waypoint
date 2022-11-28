project = "foo"

pipeline "foo" {
  step "bad" {
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
