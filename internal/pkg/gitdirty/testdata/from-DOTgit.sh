#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


# Converts all DOTgit folders to .git

# convert .git
for dotgit in `find . -type d -name DOTgit`; do
  new_dotgit=$(dirname $dotgit)/.git 
  mv $dotgit $new_dotgit
  echo $new_dotgit
done

# convert .gitignore
for dotgit in `find . -type f -name DOTgitignore`; do
  new_dotgit=$(dirname $dotgit)/.gitignore
  mv $dotgit $new_dotgit
  echo $new_dotgit
done
