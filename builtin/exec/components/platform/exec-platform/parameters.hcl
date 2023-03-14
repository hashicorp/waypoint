# This file was generated via `make gen/integrations-hcl`
parameter {
  key         = "command"
  description = "The command to execute for the deploy as a list of strings.\nEach value in the list will be rendered as a template, so it may contain template directives. Additionally, the special string `<TPL>` will be replaced with the path to the rendered file-based templates. If your template path was to a file, this will be a path a file. Otherwise, it will be a path to a directory."
  type        = "list of string"
  required    = false
}

parameter {
  key         = "dir"
  description = "The working directory to use while executing the command.\nThis will default to the same working directory as the Waypoint execution."
  type        = "string"
  required    = false
}

parameter {
  key         = "template"
  description = "A stanza that declares that a file or directory should be template-rendered."
  type        = "category"
  required    = true
}

parameter {
  key         = "template.path"
  description = "The path to the file or directory to render as a template.\nTemplating uses the following format: https://golang.org/pkg/text/template/ Available template variables depends on the input artifact."
  type        = "string"
  required    = true
}

