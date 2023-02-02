integration {
  name        = "Nomad Jobspec Canary"
  description = "The Nomad Jobspec Canary plugin promotes a Nomad canary deployment initiated by a Nomad jobspec deployment."
  identifier  = "waypoint/nomad-jobspec-canary"
  components  = ["release-manager"]
  flags       = ["builtin"]
  license {
    type = "MPL-2.0"
    url  = "https://github.com/hashicorp/waypoint/blob/main/LICENSE"
  }
}
