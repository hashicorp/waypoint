# (OPTIONAL) Overrides the copywrite config schema version
# Default: 1
schema_version = 1

project {
  # (OPTIONAL) SPDX-compatible license identifier
  # Leave blank if you don't wish to license the project
  # Default: "MPL-2.0"
  license = "MPL-2.0"

  # (OPTIONAL) Represents the year that the project initially began
  # Default: <the year the repo was first created>
  copyright_year = 2020

  # (OPTIONAL) A list of globs that should not have copyright or license headers .
  # Supports doublestar glob patterns for more flexibility in defining which
  # files or folders should be ignored
  # Default: []
  header_ignore = [
    "**_test.go",
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
  ]
}

# More information about configuration options is available in [the documentation](https://github.com/hashicorp/copywrite#config-structure).
