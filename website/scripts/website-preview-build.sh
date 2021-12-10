# Set the subdirectory name for the dev-portal app
PREVIEW_DIR=website-preview
# Clone the dev-portal repo
git clone "https://github.com/hashicorp/dev-portal.git" "$PREVIEW_DIR"
# TODO cd & checkout is only necessary
# TODO during proof-of-concept of dev-portal work.
# TODO will later just run off of main branch.
# Change into the cloned dev-portal directory
cd "$PREVIEW_DIR"
git checkout zs.refine-preview-modes
# Run the dev-portal content-repo build script
REPO=waypoint PREVIEW_DIR="$PREVIEW_DIR" ./scripts/content-repo-preview/build.sh
