project = "foo"

app "test" {
    config {
        env = {
            DATABASE_URL = configdynamic("vault", {
                path = "foo/"
            })
        }
    }
}
