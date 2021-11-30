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
	@rm -rf ./website-preview

# Use this if you have run `build-image` to use the locally built image
# rather than our CI-generated image to test dependency changes.

# COMMAND EXPLANATIONS
# (for my own Docker noob benefit)
#
# --volume "$(shell pwd):/website" \ 
# mount current directory (shell pwd) to the /website directory in the container
# (for prebuilt image, we'll want to mount the content dir into the container)
#
# --volume "/website/website-preview" \
# no second argument... which is why this works... but why?
# current directory is website, which might have duplicate node_modules
# command above would "overwrite" the node_modules above
# (this is called an "anonymous volume", seems to force Docker to use
# the node_modules on the container, rather than the one on the host cwd)
website:
	@echo "==> Starting website in Docker..."
	@docker run \
		--interactive \
		--rm \
		--tty \
		--workdir "/website" \
		--volume "$(shell pwd)/content:/website/content" \
		--volume "$(shell pwd)/data:/website/data" \
		--volume "$(shell pwd)/public:/website/public" \
		--publish "3000:3000" \
		hashicorp-waypoint-website-local \
		bash
# TODO right now I just run ./scripts/website-preview-start.sh
# TODO via the interactive terminal that "make website" starts up,
# TODO and that seems to work just fine...
# TODO but running ./scripts/website-preview-start.sh directly doesn't
# TODO work... says the ./website-preview dir isn't found, even though
# TODO when using the aforementioned interactive shell it finds it just fine.
# TODO I need to figure this out so thhat ./scripts/website-preview-start.sh can
# TODO run directly rather than requiring interaction.

.DEFAULT_GOAL := website
.PHONY: build build-image website