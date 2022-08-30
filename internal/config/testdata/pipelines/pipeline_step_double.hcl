project = "foo"

pipeline "foo" {
  step "bad" {
    pipeline "bad" {
    }
    use "bad" {
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
