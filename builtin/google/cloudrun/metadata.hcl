integration {
  name        = "Google Cloud Run"
  description = "The Google Cloud Run plugin deploys a container to Google Cloud Run."
  identifier  = "waypoint/google-cloud-run"
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
  component {
    type = "platform"
    name = "Google Cloud Run Platform"
    slug = "google-cloud-run-platform"
  }
  component {
    type = "release-manager"
    name = "Google Cloud Run Release Manager"
    slug = "google-cloud-run-release-manager"
  }
}
