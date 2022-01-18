# Set the subdirectory name for the dev-portal app
PREVIEW_DIR=website-preview
# Clone the dev-portal project, if needed
if [ ! -d "$PREVIEW_DIR" ]; then
    echo "⏳ Cloning the dev-portal repo, this might take a while..."
    git clone --depth=1 https://github.com/hashicorp/dev-portal.git "$PREVIEW_DIR"
fi
# cd into the dev-portal project
cd "$PREVIEW_DIR"
# Run the dev-portal content-repo start script
REPO=waypoint PREVIEW_DIR="$PREVIEW_DIR" npm run start:local-preview
