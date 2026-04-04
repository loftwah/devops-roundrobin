output "release_name" {
  description = "Helm release name for the application platform."
  value       = helm_release.platform.name
}

