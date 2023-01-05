<!-- This file was generated via `make gen/integrations-hcl` -->
Deploy a container to Google Cloud Run.

### Interface

### Examples

```hcl
project = "wpmini"

app "wpmini" {
  labels = {
    "service" = "wpmini",
    "env"     = "dev"
  }

  build {
    use "pack" {}

    registry {
      use "docker" {
        image = "gcr.io/waypoint-project-id/wpmini"
        tag   = "latest"
      }
    }
  }

  deploy {
    use "google-cloud-run" {
      project  = "waypoint-project-id"
      location = "europe-north1"

      port = 5000

      static_environment = {
        "NAME" : "World"
      }

      capacity {
        memory                     = 128
        cpu_count                  = 2
        max_requests_per_container = 10
        request_timeout            = 300
      }

	  service_account_name = "cloudrun@waypoint-project-id.iam.gserviceaccount.com"

      auto_scaling {
        max = 10
      }

      cloudsql_instances = ["waypoint-project-id:europe-north1:sql-instance"]

      vpc_access {
        connector = "custom-vpc-connector"
        egress = "all"
      }
    }
  }

  release {
    use "google-cloud-run" {}
  }
}
```

