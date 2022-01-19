project = "foo"

app "test" {
    config {
        env = {
            static = "${config.internal.greeting}"
        }

        internal = {
            greeting = "hello"
        }
    }
}
