cd website-preview
rm -rf .next
DEV_IO_PROXY=waypoint ENABLE_VERSIONED_DOCS=false npm run start
