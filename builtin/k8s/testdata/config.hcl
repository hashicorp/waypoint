service_account = "service-account"
namespace       = "default"
image_secret    = "regcred"

security_context {
  run_as_non_root = true
  run_as_user     = 1000
  fs_group        = 1000
}

scratch_path = [
  "/home/waypoint",
  "/tmp",
]

mount_secrets {
  mount_path  = "/opt/kubernetes/develop"
  secret_name = "waypoint-kubeconfig-develop"
  sub_path    = "kubeconfig"
  default_mode = 420
}

mount_secrets {
  mount_path  = "/opt/kubernetes/staging"
  secret_name = "waypoint-kubeconfig-staging"
  sub_path    = "kubeconfig"
  default_mode = 420
}

mount_secrets {
  mount_path  = "/opt/kubernetes/prod"
  secret_name = "waypoint-kubeconfig-prod"
  sub_path    = "kubeconfig"
  default_mode = 420
}

mount_secrets {
  mount_path  = "/home/waypoint/.docker/config.json"
  secret_name = "regcred"
  sub_path    = ".dockerconfigjson"
  default_mode = 420
}

env_from_secret = {
  "SSH_KEY_CONTENT" = {
    name = "git"
    key  = "ssh-key"
  }

  "SSH_KNOWN_HOSTS_CONTENT" = {
    name = "git"
    key  = "known-hosts"
  }
}

cpu {
  request = "20m"
  limit   = "200m"
}

memory {
  request = "128Mi"
  limit   = "512Mi"
}
