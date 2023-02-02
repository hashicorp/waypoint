integration {
  name        = "Google Cloud Run"
  description = "The Google Cloud Run plugin deploys a container to Google Cloud Run."
  identifier  = "waypoint/google-cloud-run"
  components  = ["platform", "release-manager"]
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
}
