# Set the subdirectory name for the dev-portal app
PREVIEW_DIR=website-preview
# Clone the dev-portal project, if needed
if [ ! -d "$PREVIEW_DIR" ]; then
    echo "⏳ Cloning the dev-portal repo, this might take a while..."
    git clone https://github.com/hashicorp/dev-portal.git "$PREVIEW_DIR"
fi
# cd into the dev-portal project
cd "$PREVIEW_DIR"
# Ensure we're on the branch intended for preview mode
# TODO: this should always be main,
# TODO: ie no checkout command, once we're ready
git checkout zs.refine-preview-modes
# Run the dev-portal content-repo start script
REPO=waypoint PREVIEW_DIR="$PREVIEW_DIR" ./scripts/content-repo-preview/start.sh
