# (OPTIONAL) Overrides the copywrite config schema version
# Default: 1
schema_version = 1

project {
  # (OPTIONAL) SPDX-compatible license identifier
  # Leave blank if you don't wish to license the project
  # Default: "MPL-2.0"
  license = "BUSL-1.1"

  # (OPTIONAL) Represents the year that the project initially began
  # Default: <the year the repo was first created>
  copyright_year = 2023

  # (OPTIONAL) A list of globs that should not have copyright or license headers .
  # Supports doublestar glob patterns for more flexibility in defining which
  # files or folders should be ignored
  # Default: []
  header_ignore = [
    "**_test.go",
    "**/testdata/**",
    "test-e2e/**",
    "ui/tests/**",
    "ui/mirage/**",
    ".circleci/**",
    ".github/**",
    ".release/**",
    ".vscode/**",
    "ci/**",
    "builtin/**/parameters.hcl",
    "builtin/**/outputs.hcl",
    "**/node_modules/**",
    "website/scripts/**",
    "thirdparty/**",
    "website/.next/**",
    "website/.vscode/**",
    "website/scripts/**",
    "website/website-preview**",
    // packages copied from other sources
    "internal/pkg/spinner/**",
    "ui/lib/api-common-protos/**",
    "ui/lib/grpc-web/**",
    "internal/pkg/jsonpb/**",
    "internal/pkg/defaults/**",
    "internal/pkg/copy/**", 
    "website/public/ie-custom-properties.js",
  ]
}

# More information about configuration options is available in [the documentation](https://github.com/hashicorp/copywrite#config-structure).
