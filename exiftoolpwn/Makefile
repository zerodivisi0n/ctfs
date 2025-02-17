IMAGE := exiftoolpwn
DOCKER_TAG := ghcr.io/zerodivisi0n/$(IMAGE)

docker-publish:
	docker buildx build --push \
		--platform linux/amd64,linux/arm64 \
		--tag $(DOCKER_TAG) .

docker-build:
	docker build --tag $(DOCKER_TAG) .

run: docker-build
	docker run --rm -p 3000:3000 --name $(IMAGE) $(DOCKER_TAG)
