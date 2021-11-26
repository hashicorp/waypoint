# # Default: run this if working on the website locally to run in watch mode.
# website:
# 	@echo "==> Downloading latest Docker image..."
# 	@docker pull hashicorp/vault-website
# 	@echo "==> Starting website in Docker..."
# 	@docker run \
# 		--interactive \
# 		--rm \
# 		--tty \
# 		--workdir "/website" \
# 		--volume "$(shell pwd):/website" \
# 		--volume "/website/node_modules" \
# 		--publish "3000:3000" \
# 		hashicorp/vault-website \
# 		npm start

# # This command will generate a static version of the website to the "out" folder.
# build:
# 	@echo "==> Downloading latest Docker image..."
# 	@docker pull hashicorp/vault-website
# 	@echo "==> Starting build in Docker..."
# 	@docker run \
# 		--interactive \
# 		--rm \
# 		--tty \
# 		--workdir "/website" \
# 		--volume "$(shell pwd):/website" \
# 		--volume "/website/node_modules" \
# 		hashicorp/vault-website \
# 		npm run static

# Generate a new local Docker image using the
# adjacent Dockerfile.
# local Docker image with the dependency changes included.
setup:
	@echo "==> Cloning hashicorp/dev-portal..."
	@rm -rf ./website-preview
	@git clone https://github.com/hashicorp/dev-portal.git website-preview
	@echo "==> Building Docker image..."
	@docker build -t hashicorp-waypoint-website-local .

# Use this if you have run `build-image` to use the locally built image
# rather than our CI-generated image to test dependency changes.
website:
	@echo "==> Starting website in Docker..."
	@docker run \
		--interactive \
		--rm \
		--tty \
		--workdir "/website" \
		--volume "$(shell pwd):/website" \
		--volume "/website/website-preview" \
		--publish "3000:3000" \
		hashicorp-waypoint-website-local \
		ls -a && cd website-preview && npm i --production=false && cd .. && npm run start

.DEFAULT_GOAL := website
.PHONY: build build-image website

# docker run --interactive --rm --tty --workdir "/website" --volume "$(pwd):/website" --publish "3000:3000" hashicorp-waypoint-website-local ls -a && cd website-preview && npm i --production=false