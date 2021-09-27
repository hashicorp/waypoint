slug=$(display)
ref=$(display)
.PHONY: docs/add-version
# Adds a new version to the version-manifest file, which drives display of past documentation versions on the site.
docs/add-version:
	@echo "==> Adding version $(display) to the manifest..."
	@if [ $(shell git diff --cached | wc -l) -ne 0 ]; then echo "⚠️  staged files found, please commit or unstage these files to add a version" && exit 1; fi
	@docker run \
		--rm \
		-it \
		--volume "$(shell pwd):/website" \
		--workdir "/website" \
		docker.mirror.hashicorp.services/node:14-alpine \
		npx @hashicorp/platform-cli docs:add-version $(display) $(slug) $(ref)
	git add version-manifest.json
	git commit -m "add version $(display) to version-manifest.json"


.PHONY: docs/remove-version
# Removes a version from the version-manifest file, which drives display of past documentation versions on the site.
docs/remove-version:
	@echo "==> Removing version $(display) from the manifest..."
	@if [ $(shell git diff --cached | wc -l) -ne 0 ]; then echo "⚠️  staged files found, please commit or unstage these files to add a version" && exit 1; fi
	@docker run \
		--rm \
		-it \
		--volume "$(shell pwd):/website" \
		--workdir "/website" \
		docker.mirror.hashicorp.services/node:14-alpine \
		npx @hashicorp/platform-cli docs:remove-version $(display)
	git add version-manifest.json
	git commit -m "remove version $(display) from version-manifest.json"