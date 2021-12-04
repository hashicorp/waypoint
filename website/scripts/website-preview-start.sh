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
