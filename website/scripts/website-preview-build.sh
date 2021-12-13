# Set the subdirectory name for the dev-portal app
PREVIEW_DIR=website-preview
# Clone the dev-portal repo
git clone "https://github.com/hashicorp/dev-portal.git" "$PREVIEW_DIR"
# Run the dev-portal content-repo build script
REPO=waypoint PREVIEW_DIR="$PREVIEW_DIR" "./$PREVIEW_DIR/scripts/content-repo-preview/build.sh"
