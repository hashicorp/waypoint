project = "foo"

app "web" {
    build "pack" {
        registry "docker" {
            name = "gcr.io/mitchellh-test/myapp:latest"
        }
    }

    deploy "google-cloud-run" {}
}
