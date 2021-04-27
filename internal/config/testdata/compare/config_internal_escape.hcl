project = "foo"

app "test" {
    config {
        env = {
            static = file("${path.project}/config_escape_data.hcl")
            more = templatestring(file("${path.project}/config_escape_data.hcl"), {
              pass = config.internal.pass,
            })
        }

        internal = {
          pass = configdynamic("tfc", {})
        }
    }
}
