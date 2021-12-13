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
# TODO checkout is only necessary during dev,
# TODO: will later always be main,
# TODO: but for now, checking out a specific
# TODO: branch is useful for dev.
git checkout main
# Run the dev-portal content-repo start script
REPO=waypoint PREVIEW_DIR="$PREVIEW_DIR" ./scripts/content-repo-preview/start.sh
