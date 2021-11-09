rm -rf ./website-preview
git config --global credential.helper store
git clone "https://zchsh:$GITHUB_PAT_FOR_CLONE@github.com/hashicorp/dev-portal.git" website-preview
cd website-preview
git checkout zs.cleanup-refine-migration
npm ci
