#!/bin/bash

# Converts all .git folders to DOTgit

for dotgit in `find . -type d -name .git`; do
  new_dotgit=$(dirname $dotgit)/DOTgit 
  mv $dotgit $new_dotgit
  echo $new_dotgit
done
