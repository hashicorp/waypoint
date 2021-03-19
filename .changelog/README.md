# How To Use

Waypoint uses `go-changelog` to generate its changelog on release:

* https://github.com/hashicorp/go-changelog

To install, run the following command:

```
go get github.com/hashicorp/go-changelog/cmd/changelog-build
```

## CHANGELOG entry examples

CHANGELOG entries are expected to be txt files created inside this folder
`.changelog`. The file name is expected to be the same issue number that will
be linked when the CHANGELOG is generated. So for example, if your issue is
\#1234, your file name would be `.changelog/1234.txt`.

While for git commit messages, we expect the leading subject to be more specific
as to the section it updates, for example a change with k8s might be:

```
builtin/k8s: Add support for feature Y

This commit adds support for feature Y....
```

The changelog entry should be more user facing friendly, so it would instead read:

~~~
```release-note:improvement
plugin/k8s: Add support for feature Y
```
~~~

Below are some examples of how to generate a CHANGELOG entry with your pull
request.

### Improvement

~~~
```release-note:improvement
server: Add new option for configs
```
~~~

### Feature

~~~
```release-note:feature
plugin/nomad: New feature integration
```
~~~

### Bug

~~~
```release-note:bug
plugin/docker: Fix broken code
```
~~~

### Multiple Entries

~~~
```release-note:bug
plugin/docker: Fix broken code
```

```release-note:bug
plugin/nomad: Fix broken code
```

```release-note:bug
plugin/k8s: Fix broken code
```
~~~

### Long Description with Markdown

~~~
```release-note:feature
cli: Lorem ipsum dolor `sit amet`, _consectetur_ adipiscing elit, **sed** do
eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim
veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo
consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse
cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non
proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
```
~~~

## How to generate CHANGELOG entries for release

Below is an example for running `go-changelog` to generate a collection of
entries. It will generate output that can be inserted into CHANGELOG.md.

For more information as to what each flag does, make sure to run `changelog-build -help`.

```
changelog-build -last-release v0.5.0 -entries-dir .changelog/ -changelog-template changelog.tmpl -note-template note.tmpl -this-release 86b6b38faa7c69f26f1d4c71e271cd4285daadf9
```

