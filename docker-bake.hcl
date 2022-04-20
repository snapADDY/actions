variable "TAG" {
	default = "latest"
}

group "default" {
  targets = [
    "bandit-github",
    "prepare-img-tag"
  ]
}

target "bandit-github" {
	context = "bandit-github"
	dockerfile = "Dockerfile"
	tags = ["actions-bandit-github:${TAG}"]
}

target "prepare-img-tag" {
	dockerfile = "prepare-img-tag/Dockerfile"
	tags = ["actions-prepare-img-tag:${TAG}"]
}
