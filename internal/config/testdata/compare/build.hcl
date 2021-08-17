project = "foo"

app "test" {
    build {
        labels = {
            "foo" = "bar"
        }

        hook {
            when = labels["foo"]
            command = ["echo", "foo"]
        }
    }
}
