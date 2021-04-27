project = "foo"

app "test" {
    config {
        env = {
            static = lower(config.internal.greeting, config.internal.extra)
        }

        internal = {
            greeting = configdynamic("tfc", {})
            extra = "FOO"
        }
    }
}
