project = "foo"

app "foo" {
    build {
        workspace "foo" {
            use "docker" {}
        }

        label "bar" {}
    }

    deploy {}
}
