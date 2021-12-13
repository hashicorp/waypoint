# Set the subdirectory name for the dev-portal app
PREVIEW_DIR=website-preview
# Clone the dev-portal repo
git clone "https://github.com/hashicorp/dev-portal.git" "$PREVIEW_DIR"
# TODO cd & checkout is only necessary during dev,
# TODO: will later always be main,
# TODO: but for now, checking out a specific
# TODO: branch is useful for dev.
# Change into the cloned dev-portal directory
cd "$PREVIEW_DIR"
git checkout main
cd ..
# Run the dev-portal content-repo build script
REPO=waypoint PREVIEW_DIR="$PREVIEW_DIR" "./$PREVIEW_DIR/scripts/content-repo-preview/build.sh"
