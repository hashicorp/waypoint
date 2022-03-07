# Set the subdirectory name for the dev-portal app
PREVIEW_DIR=website-preview
# Clone the dev-portal project, if needed
if [ ! -d "$PREVIEW_DIR" ]; then
    echo "⏳ Cloning the dev-portal repo, this might take a while..."
    git clone --branch brk.feat/io-preview-cache --depth=1 https://github.com/hashicorp/dev-portal.git "$PREVIEW_DIR"
fi
# cd into the dev-portal project
cd "$PREVIEW_DIR"

# Run the dev-portal content-repo start script
REPO=waypoint DEV_IO=waypoint IS_CONTENT_PREVIEW=true HASHI_ENV=preview npm run build:deploy-preview
