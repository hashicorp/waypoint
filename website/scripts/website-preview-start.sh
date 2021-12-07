# TODO: would be nice to only reset, clone,
# TODO: install deps, etc if website-preview
# TODO: does not already exist
# Clean up existing clone
# TODO: actually, don't do this,
# TODO: just git clean and then checkout the
# TODO: target branch (if the repo already exists)
# Clone project, and cd into it
# IF ./website-preview does NOT exist...
FILE=website-preview
if [ ! -d "$FILE" ]; then
    echo "Cloning the dev-portal repo, this might take a while..."
    git clone https://github.com/hashicorp/dev-portal.git website-preview
fi
# cd into the project
cd website-preview
# Ensure we're on the branch intended for preview mode
# TODO: this should always be main,
# TODO: ie no checkout command, once we're ready
git checkout zs.refine-preview-modes
# TODO: move the below into a script within dev-portal,
# TODO: which we can now run since we'll
# TODO: have cloned it and checked out latest main
# Install dependencies
# TODO maybe some way to optimize this?
# TODO likely slow with npm ci deleting node_modules
# TODO every time...
npm i --production=false
# TODO: clean this up, add comments
cd ./website-preview
git restore --staged .
git restore .
git clean -f -d
git checkout zs.refine-preview-modes
git pull
rm -rf .next
# cp -R ../content/ ./content
# TODO set up watcher to sync all files
# TODO under website/public into website/website-preview/public
cp -R ../public/img/** ./public/img/
DEV_IO_PROXY=waypoint IS_CONTENT_DEPLOY_PREVIEW=true ./node_modules/.bin/next
