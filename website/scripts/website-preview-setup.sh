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
# TODO: this should always be main, once we're ready
git checkout zs.refine-preview-modes
# Install dependencies
# TODO maybe some way to optimize this?
# TODO likely slow with npm ci deleting node_modules
# TODO every time...
npm i --production=false
