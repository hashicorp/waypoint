## nomad-jobspec-canary (releasemanager)

Promotes a Nomad canary deployment initiated by a Nomad jobspec deployment.

If your Nomad deployment is configured to use canaries, this releaser plugin lets
you promote (or fail) the canary deployment. You may also target specific task
groups within your job for promotion, if you have multiple task groups in your canary
deployment.

-> **Note:** Using the `-prune=false` flag is recommended for this releaser. By default,
Waypoint prunes and destroys all unreleased deployments and keeps only one previous
deployment. Therefore, if `-prune=false` is not set, Waypoint may delete
your job via "pruning" a previous version. See [deployment pruning](/docs/lifecycle/release#deployment-pruning)
for more information.

### Release URL

If you want the URL of the release of your deployment to be published in Waypoint,
you must set the meta 'waypoint.hashicorp.com/release_url' in your jobspec. The
value specified in this meta field will be published as the release URL for your
application. In the future, this may source from Consul.

### Interface

### Examples

```hcl
// The waypoint.hcl file
release {
  use "nomad-jobspec-canary" {
    groups = [
      "app"
    ]
  }
}

// The app.nomad.tpl file
job "web" {
  datacenters = ["dc1"]

  group "app" {
    network {
      mode = "bridge"
      port "http" {
        to = 80
      }
    }

    // Setting a canary in the update stanza indicates a canary deployment
    update {
      max_parallel = 1
      canary       = 1
      auto_revert  = true
      auto_promote = false
      health_check = "task_states"
    }

    service {
      name = "app"
      port = 80
      connect {
        sidecar_service {}
      }
    }

    task "app" {
      driver = "docker"
      config {
        image = "${artifact.image}:${artifact.tag}"
        ports  = ["http"]
      }

      env {
        %{ for k,v in entrypoint.env ~}
        ${k} = "${v}"
        %{ endfor ~}

        // Ensure we set PORT for the URL service. This is only necessary
        // if we want the URL service to function.
        PORT = 80
      }
    }
  }

  group "app-gateway" {
    network {
      mode = "bridge"
      port "inbound" {
        static = 8080
        to     = 8080
      }
    }

    service {
      name = "gateway"
      port = "8080"

      connect {
        gateway {
          proxy {}

          ingress {
            listener {
              port = 8080
              protocol = "http"
              service {
                name  = "app"
                hosts = [ "*" ]
              }
            }
          }
        }
      }
    }
  }
  meta = {
    // Ensure we set meta for Waypoint to detect the release URL
    "waypoint.hashicorp.com/release_url" = "http://app.ingress.dc1.consul:8080"
  }
}
```

### Required Parameters

This plugin has no required parameters.

### Optional Parameters

These parameters are used in the [`use` stanza](/docs/waypoint-hcl/use) for this plugin.

#### fail_deployment

If true, marks the deployment as failed.

- Type: **bool**
- **Optional**

#### groups

List of task group names which are to be promoted.

- Type: **list of string**
- **Optional**
