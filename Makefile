export DOCKER_CLI_EXPERIMENTAL=enabled

.PHONY: multiarch
multiarch:
	faas-cli build --shrinkwrap --filter add-post
	docker buildx create --use --name=multiarch --node=multiarch
	docker buildx build --platform linux/amd64,linux/arm/v7,linux/arm64 --output "type=image,push=true" --tag alexellis2/add-post:0.1.0 build/add-post/ --build-arg GO111MODULE=on
