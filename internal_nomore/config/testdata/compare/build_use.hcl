project = "foo"

app "test" {
    build {
        use "test" {
            foo = path.app
        }
    }
}
