project = "foo"

config {
  env = {
    parent = "1"
  }
}

app "test" {
    config {
        env = {
            child = "2"
        }
    }
}
