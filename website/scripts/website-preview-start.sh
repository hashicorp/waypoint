# TODO: clean this up, add comments
cd ./website-preview
git clean -f
rm -rf .next
cp -R ../content/ ./content
cp -R ../public/img/** ./public/img/
DEV_IO_PROXY=waypoint IS_CONTENT_DEPLOY_PREVIEW=true ./node_modules/.bin/next
