project = "foo"

runner {
    enabled = "true"
    scoped_settings {
        workspace "develop" {
            data_source "git" {
                url  = "https://github.com/hashicorp/waypoint-examples.git"
                path = "docker/node-js"
                ref = "develop"
            }
        }

        workspace "staging" {
            data_source "git" {
                url  = "https://github.com/hashicorp/waypoint-examples.git"
                path = "docker/node-js"
                ref = "staging"
            }
        }

        workspace "prod" {
            data_source "git" {
                url  = "https://github.com/hashicorp/waypoint-examples.git"
                path = "docker/node-js"
                ref = "prod"
            }
        }
    }
}

app "web" {
    build {}
    deploy {}
}