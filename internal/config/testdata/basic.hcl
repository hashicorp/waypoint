project = "foo"

app "web" {
    build {
        use "pack" {}

        registry {
            use "docker" {
                name = "gcr.io/mitchellh-test/myapp:latest"
            }
        }
    }

    deploy {
        use "google-cloud-run" {}
    }
}
