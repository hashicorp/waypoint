# Delete cached .next, if present
# TODO: would be nice to leverage caching,
# TODO: deleted for now as it makes debugging harder
echo "Files before .next cache delete:"
ls -a
echo "Deleting .next cache..."
rm -rf .next
echo "Done"
# Clone the dev-portal repo in ./website-preview
git config --global credential.helper store
git clone "https://zchsh:${GITHUB_WEBSITE_PREVIEW_PAT}@github.com/hashicorp/dev-portal.git" website-preview
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
git checkout zs.refine-preview-modes
npm i --production=false
# Delete other products' docs pages,
# these will just increase build times
# TODO instead of doing thing, should perhaps
# TODO delete /pages/_proxied-dot-io/boundary
# TODO entirely? Such routes will be redirected
# TODO anyways, so not point wasting time
# TODO in rendering them.
rm -rf ./src/pages/_proxied-dot-io/boundary/docs
rm -rf ./src/pages/_proxied-dot-io/boundary/api-docs
# Build the site
# TODO: is there a way to remove these two manully set env vars?
# TODO: maybe using Vercel's System Environment variables?
npm run build
# Copy .next build output folder into project root,
# so that Vercel's NextJS preset picks up on the build output
echo "Listing files after build..."
ls -a
# echo "Copying .next output to project root..."
# cp -R .next/ ../.next
echo "Done."
