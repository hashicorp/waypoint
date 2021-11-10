# Delete cached .next, if present
# TODO: would be nice to leverage caching,
# TODO: deleted for now as it makes debugging harder
rm -rf .next
# Clone the dev-portal repo in ./website-preview
git config --global credential.helper store
git clone "https://zchsh:${GITHUB_WEBSITE_PREVIEW_PAT}@github.com/hashicorp/dev-portal.git" website-preview && cd website-preview && git checkout zs.cleanup-refine-migration && npm i --production=false
# Copy all local content into the cloned dev-portal directory
# TODO: may be able to avoid, fix in docs-page
# TODO: the "include partials" remark plugin path needs
# TODO: to be an arg of renderPageMdx, not fixed. ref:
# TODO: https://github.com/hashicorp/react-components/blob/a546b9dc9a874df81a324b9e54e4fb1034b79a5d/packages/docs-page/render-page-mdx.js#L14
cp -R ./content/ ./website-preview/content
# Change into the cloned dev-portal directory
cd ./website-preview
# Build the site
DEV_IO_PROXY=waypoint IS_CONTENT_DEPLOY_PREVIEW=true npm run build
# Copy .next build output folder into project root,
# so that Vercel's NextJS preset picks up on the build output
cp -R .next/ ../.next
# Merge the local images (all in ./public/img) into
# the shared dev-portal public folder
cp -R ./public/img/** ./website-preview/public/img/
# Delete the local public folder and replace it
# with the now-combined dev-portal public folder
rm -rf ../public
cp -R public/ ../public
