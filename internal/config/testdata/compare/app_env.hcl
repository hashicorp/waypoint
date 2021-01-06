project = "foo"

app "bar" {
    path = "./bar"

    labels = {
        "env": env["APP_ENV"]
    }
}
