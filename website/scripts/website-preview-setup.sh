rm -rf ./website-preview
# git clone https://github.com/hashicorp/dev-portal.git website-preview && cd website-preview && git checkout zs.cleanup-refine-migration && npm ci
git config --global credential.helper store
git clone "https://zchsh:${GITHUB_WEBSITE_PREVIEW_PAT}@github.com/hashicorp/dev-portal.git" website-preview && cd website-preview && git checkout zs.cleanup-refine-migration && npm i
