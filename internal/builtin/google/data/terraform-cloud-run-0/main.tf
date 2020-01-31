resource "google_cloud_run_service" "default" {
  name     = var.name
  location = var.location
  project  = var.project

  template {
    metadata {
      annotations = {
        devflow = 1
      }
    }

    spec {
      containers {
        image = var.image
      }
    }
  }
}

data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth" {
  location    = google_cloud_run_service.default.location
  project     = google_cloud_run_service.default.project
  service     = google_cloud_run_service.default.name

  policy_data = data.google_iam_policy.noauth.policy_data
}
