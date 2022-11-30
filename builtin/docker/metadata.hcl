integration {
  name = "Docker"
  description = "Run and deploy Docker containers"
  identifier = "waypoint/docker"
  components = [ "builder", "platform", "registry", "task" ]
  docs {
    process_docs = true
    readme_location = "./README.md"
  }
}
