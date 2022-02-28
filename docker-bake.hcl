variable "TAG" {
	default = "latest"
}

group "default" {
  targets = [
    "prepare-img-tag"
  ]
}

target "prepare-img-tag" {
	dockerfile = "prepare-img-tag/Dockerfile"
	tags = ["actions-prepare-img-tag:${TAG}"]
}
