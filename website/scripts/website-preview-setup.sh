# TODO: would be nice to only reset, clone,
# TODO: install deps, etc if website-preview
# TODO: does not already exist
# Clean up existing clone
rm -rf ./website-preview
# Clone project, and cd into it
git clone https://github.com/hashicorp/dev-portal.git website-preview
cd website-preview
# Ensure we're on main (can probably remove this)
git checkout main
# Install dependencies
npm i --production=false
