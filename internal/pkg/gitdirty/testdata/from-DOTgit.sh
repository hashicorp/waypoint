#!/bin/bash

# Converts all DOTgit folders to .git

for dotgit in `find . -type d -name DOTgit`; do
  new_dotgit=$(dirname $dotgit)/.git 
  mv $dotgit $new_dotgit
  echo $new_dotgit
done
