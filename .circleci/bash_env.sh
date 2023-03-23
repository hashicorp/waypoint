#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


export GIT_COMMIT=$(git rev-parse --short HEAD)
export GIT_DIRTY=$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)
export GIT_DESCRIBE=$(git describe --tags --always --match "v*")
export GIT_IMPORT=github.com/hashicorp/waypoint/internal/version
export GOLDFLAGS="-X ${GIT_IMPORT}.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X ${GIT_IMPORT}.GitDescribe=${GIT_DESCRIBE}"
export CGO_ENABLED=0
