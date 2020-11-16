project = "foo"

app "web" {
    config {
        env = {
            static = "hello"
        }
    }

    build {}

    deploy {}
}
