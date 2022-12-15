parameter {
  key         = "template"
  description = <<EOT
A stanza that declares that a file or directory should be template-rendered.
EOT
  type        = "category"
  required    = true
}

parameter {
  key         = "template.path"
  description = <<EOT
The path to the file or directory to render as a template.

Templating uses the following format: https://golang.org/pkg/text/template/ Available template variables depends on the input artifact.
EOT
  type        = "string"
  required    = true
}

parameter {
  key         = "command"
  description = <<EOT
The command to execute for the deploy as a list of strings.

Each value in the list will be rendered as a template, so it may contain template directives. Additionally, the special string `<TPL>` will be replaced with the path to the rendered file-based templates. If your template path was to a file, this will be a path a file. Otherwise, it will be a path to a directory.
EOT
  type        = "list of string"
  required    = false
}

parameter {
  key         = "dir"
  description = <<EOT
The working directory to use while executing the command.

This will default to the same working directory as the Waypoint execution.
EOT
  type        = "string"
  required    = false
}

