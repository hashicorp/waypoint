output "observed_generation" {
  value = google_cloud_run_service.default.status[0].observed_generation
}

output "desired_generation" {
  value = google_cloud_run_service.default.metadata[0].generation
}

output "url" {
  value = google_cloud_run_service.default.status[0].url
}

