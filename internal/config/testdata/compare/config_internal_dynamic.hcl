project = "foo"

app "test" {
    config {
        env = {
            static = "${config.internal.greeting} ${config.internal.suffix}"
        }

        internal = {
            suffix = "ok?"
            greeting = configdynamic("foo", {})
        }
    }
}
