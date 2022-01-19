project = "foo"

app "test" {
    config {
        env = {
            static = "${config.internal.greeting} ${config.internal.suffix}"
            extra = "extra: ${config.env.static}"
        }

        internal = {
            suffix = "ok?"
            greeting = configdynamic("foo", {})
        }
    }
}
