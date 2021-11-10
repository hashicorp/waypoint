cd website-preview
rm -rf .next
DEV_IO_PROXY=waypoint IS_CONTENT_DEPLOY_PREVIEW=true npm run start
