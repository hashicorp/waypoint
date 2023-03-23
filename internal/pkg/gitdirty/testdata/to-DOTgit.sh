#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


# Converts all .git folders to DOTgit

# convert .git
for dotgit in `find . -type d -name .git`; do
  new_dotgit=$(dirname $dotgit)/DOTgit 
  mv $dotgit $new_dotgit
  echo $new_dotgit
done

# convert .gitignore
for dotgit in `find . -type f -name .gitignore`; do
  new_dotgit=$(dirname $dotgit)/DOTgitignore 
  mv $dotgit $new_dotgit
  echo $new_dotgit
done

