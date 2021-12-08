# Delete cached .next, if present
# TODO: would be nice to leverage caching,
# TODO: deleted for now as it makes debugging harder
echo "Files before .next cache delete:"
ls -a
echo "Deleting .next cache..."
rm -rf .next
echo "Done"
# Clone the dev-portal repo in ./website-preview
git clone "https://github.com/hashicorp/dev-portal.git" website-preview
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
# TODO in rendering them...
# TODO however, this requires product repo awareness.
# TODO perhaps we could achieve a similar effect
# TODO with logic in the dev-portal docs routes
# TODO which would return an empty array
# TODO from getStaticPaths if the product
# TODO is not the target product (based on env check)
# TODO during a build
rm -rf ./src/pages/_proxied-dot-io/boundary
# Build the site
# Note that DEV_IO and IS_CONTENT_PREVIEW are set
# in Vercel configuration for the project
npm run build
# Copy .next build output folder into project root,
# so that Vercel's NextJS preset picks up on the build output
echo "Listing files after build..."
ls -a
echo "Copying .next output to project root..."
cp -R .next/ ../.next
echo "Done."
