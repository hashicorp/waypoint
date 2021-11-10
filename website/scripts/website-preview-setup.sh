rm -rf ./website-preview
git clone https://github.com/hashicorp/dev-portal.git website-preview && cd website-preview && git checkout zs.cleanup-refine-migration && npm i --production=false
