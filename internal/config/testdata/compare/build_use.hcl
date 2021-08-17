project = "foo"

app "test" {
    build {
        labels = { "foo" = "hello" }

        use "test" {
            foo = path.app
            bar = labels["foo"]
        }
    }
}
