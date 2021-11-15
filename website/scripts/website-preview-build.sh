# Delete cached .next, if present
# TODO: would be nice to leverage caching,
# TODO: deleted for now as it makes debugging harder
echo "Files before .next cache delete:"
ls -a
echo "Deleting .next cache..."
rm -rf .next
# Clone the dev-portal repo in ./website-preview
git config --global credential.helper store
git clone "https://zchsh:${GITHUB_WEBSITE_PREVIEW_PAT}@github.com/hashicorp/dev-portal.git" website-preview
# Copy all local content into the cloned dev-portal directory
# TODO: may be able to avoid, fix in docs-page
# TODO: the "include partials" remark plugin path needs
# TODO: to be an arg of renderPageMdx, not fixed. ref:
# TODO: https://github.com/hashicorp/react-components/blob/a546b9dc9a874df81a324b9e54e4fb1034b79a5d/packages/docs-page/render-page-mdx.js#L14
cp -R ./content/ ./website-preview/content
# Merge the local images (all in ./public/img) into
# the shared dev-portal public folder
cp -R ./public/img/** ./website-preview/public/img/
# Delete the local public folder and replace it
# with the now-combined dev-portal public folder
rm -rf ./public
cp -R ./website-preview/public/ ./public
# Change into the cloned dev-portal directory
cd ./website-preview
# Install dependencies
git checkout zs.cleanup-refine-migration
npm i --production=false
# Delete other products' pages,
# these will just increase build times
rm -rf ./src/pages/_proxied-dot-io/boundary
# Build the site
# TODO: is there a way to remove these two manully set env vars?
# TODO: maybe using Vercel's System Environment variables?
npm run build
# Copy .next build output folder into project root,
# so that Vercel's NextJS preset picks up on the build output
cp -R .next/ ../.next
