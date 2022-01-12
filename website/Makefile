.PHONY: website
.DEFAULT_GOAL := website

website:
	@echo "==> Downloading latest Docker image..."
	@docker pull hashicorp/dev-portal
	@echo "==> Starting website..."
	@docker run \
		-it \
		--publish "3000:3000" \
		--rm \
		--tty \
		--volume "$(shell pwd)/content:/app/content" \
		--volume "$(shell pwd)/public:/app/public" \
		--volume "$(shell pwd)/data:/app/data" \
		--volume "$(shell pwd)/redirects.js:/app/redirects.js" \
		--volume "next-dir:/app/website-preview/.next" \
		--volume "$(shell pwd)/.env:/app/.env" \
		-e "REPO=waypoint" \
		hashicorp/dev-portal
