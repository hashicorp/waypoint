app "api" {
  config {
    env = { "foo" = "bar" }

    workspace "dev" {
      env = { "bar" = "baz" }
    }
  }
}
