project = "foo"

app "test" {
    build {
        labels = {
            "foo" = "bar"
        }

        use "docker" {}

        registry {
          use "A" {}

          workspace "production" {
            use "B" {}
          }
        }
    }
}
