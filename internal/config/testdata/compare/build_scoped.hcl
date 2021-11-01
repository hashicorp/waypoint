project = "foo"

app "test" {
    build {
        labels = {
            "foo" = "bar"
        }

        use "A" {}

        workspace "production" {
          use "B" {}
        }
    }
}
