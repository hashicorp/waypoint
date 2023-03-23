# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# This file was generated via `make gen/integrations-hcl`
parameter {
  key           = "builder"
  description   = "The buildpack builder image to use"
  type          = "string"
  required      = false
  default_value = "heroku/buildpacks:20"
}

parameter {
  key         = "buildpacks"
  description = "The exact buildpacks to use\nIf set, the builder will run these buildpacks in the specified order. They can be listed using several [URI formats](https://buildpacks.io/docs/app-developer-guide/specific-buildpacks)."
  type        = "list of string"
  required    = false
}

parameter {
  key         = "disable_entrypoint"
  description = "if set, the entrypoint binary won't be injected into the image\nThe entrypoint binary is what provides extended functionality such as logs and exec. If it is not injected at build time the expectation is that the image already contains it"
  type        = "bool"
  required    = false
}

parameter {
  key         = "ignore"
  description = "file patterns to match files which will not be included in the build\nEach pattern follows the semantics of .gitignore. This is a summarized version:\n\n1. A blank line matches no files, so it can serve as a separator\n\t for readability.\n\n2. A line starting with # serves as a comment. Put a backslash (\"\\\")\n\t in front of the first hash for patterns that begin with a hash.\n\n3. Trailing spaces are ignored unless they are quoted with backslash (\"\\\").\n\n4. An optional prefix \"!\" which negates the pattern; any matching file\n\t excluded by a previous pattern will become included again. It is not\n\t possible to re-include a file if a parent directory of that file is\n\t excluded. Git doesn’t list excluded directories for performance reasons,\n\t so any patterns on contained files have no effect, no matter where they\n\t are defined. Put a backslash (\"\\\") in front of the first \"!\" for\n\t patterns that begin with a literal \"!\", for example, \"\\!important!.txt\".\n\n5. If the pattern ends with a slash, it is removed for the purpose of the\n\t following description, but it would only find a match with a directory.\n\t In other words, foo/ will match a directory foo and paths underneath it,\n\t but will not match a regular file or a symbolic link foo (this is\n\t consistent with the way how pathspec works in general in Git).\n\n6. If the pattern does not contain a slash /, Git treats it as a shell glob\n\t pattern and checks for a match against the pathname relative to the\n\t location of the .gitignore file (relative to the top level of the work\n\t tree if not from a .gitignore file).\n\n7. Otherwise, Git treats the pattern as a shell glob suitable for\n\t consumption by fnmatch(3) with the FNM_PATHNAME flag: wildcards in the\n\t pattern will not match a / in the pathname. For example,\n\t \"Documentation/*.html\" matches \"Documentation/git.html\" but not\n\t \"Documentation/ppc/ppc.html\" or \"tools/perf/Documentation/perf.html\".\n\n8. A leading slash matches the beginning of the pathname. For example,\n\t \"/*.c\" matches \"cat-file.c\" but not \"mozilla-sha1/sha1.c\".\n\n9. Two consecutive asterisks (\"**\") in patterns matched against full\n\t pathname may have special meaning:\n\n\t\ti.   A leading \"**\" followed by a slash means match in all directories.\n\t\t\t\t For example, \"** /foo\" matches file or directory \"foo\" anywhere,\n\t\t\t\t the same as pattern \"foo\". \"** /foo/bar\" matches file or directory\n\t\t\t\t \"bar\" anywhere that is directly under directory \"foo\".\n\n\t\tii.  A trailing \"/**\" matches everything inside. For example, \"abc/**\"\n\t\t\t\t matches all files inside directory \"abc\", relative to the location\n\t\t\t\t of the .gitignore file, with infinite depth.\n\n\t\tiii. A slash followed by two consecutive asterisks then a slash matches\n\t\t\t\t zero or more directories. For example, \"a/** /b\" matches \"a/b\",\n\t\t\t\t \"a/x/b\", \"a/x/y/b\" and so on.\n\n\t\tiv.  Other consecutive asterisks are considered invalid."
  type        = "list of string"
  required    = false
}

parameter {
  key         = "process_type"
  description = "The process type to use from your Procfile. if not set, defaults to `web`\nThe process type is used to control over all container modes, such as configuring it to start a web app vs a background worker"
  type        = "string"
  required    = false
}

parameter {
  key         = "static_environment"
  description = "environment variables to expose to the buildpack\nthese environment variables should not be run of the mill configuration variables, use waypoint config for that. These variables are used to control over all container modes, such as configuring it to start a web app vs a background worker"
  type        = "map of string to string"
  required    = false
}

