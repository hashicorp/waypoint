# Contributing to Waypoint

>**Note:** We take Waypoint's security and our users' trust very seriously.
>If you believe you have found a security issue in Waypoint, please responsibly
>disclose by contacting us at security@hashicorp.com.

**First:** if you're unsure or afraid of _anything_, just ask or submit the
issue or pull request anyways. You won't be yelled at for giving your best
effort. The worst that can happen is that you'll be politely asked to change
something. We appreciate any sort of contributions, and don't want a wall of
rules to get in the way of that.

That said, if you want to ensure that a pull request is likely to be merged,
talk to us! A great way to do this is in issues themselves. When you want to
work on an issue, comment on it first and tell us the approach you want to take.

## Getting Started

### Some Ways to Contribute

* Report potential bugs.
* Suggest product enhancements.
* Increase our test coverage.
* Fix a [bug](https://github.com/hashicorp/waypoint/labels/bug).
* Implement a requested [enhancement](https://github.com/hashicorp/waypoint/labels/enhancement).
* Improve our guides and documentation.

### Reporting an Issue:

>Note: Issues on GitHub for Waypoint are intended to be related to bugs or feature requests.
>Questions should be directed to other community resources such as the [forum](https://discuss.hashicorp.com/)

* Make sure you test against the latest released version. It is possible we
already fixed the bug you're experiencing. However, if you are on an older
version of Waypoint and feel the issue is critical, do let us know.

* Check existing issues (both open and closed) to make sure it has not been
reported previously.

* Provide a reproducible test case. If a contributor can't reproduce an issue,
then it dramatically lowers the chances it'll get fixed. If we can't reproduce
an issue long enough, we are usually forced to close the issue.

* As part of the test case, please include any Waypoint configurations
(`waypoint.hcl`), build configs such as Dockerfiles, etc. Log output with
log level set with verbose flags (at least `-vv`) is helpful too.

* If the issue is related to the browser UI, please also include the name 
and version of the browser and any extensions that may be interacting 
with the UI

* Aim to respond promptly to any questions made by the Waypoint team on your
issue. Stale issues will be closed.

### Issue Lifecycle

1. The issue is reported.
1. The issue is verified and categorized by a Waypoint maintainer.
   Categorization is done via tags. For example, bugs are tagged as "bug".
1. Unless it is critical, the issue is left for a period of time (sometimes many
   weeks or months), giving outside contributors a chance to address the issue
   and our internal teams time to plan for inclusion in a release.
1. Once someone commits to addressing the issue, it is put in the target milestone
   (e.g. `0.5.x`) and assigned to the person who has committed to the work.
   1. If you'd like to work on an open issue, please double-check first that no
   one else is currently working on it.
   1. It's also best to give a quick overview of how you plan to solve the issue
   in a comment. That way people have a chance to inform you of any potential
   hurdles or unknown complications you may run into.
1. The issue is addressed in a pull request or commit. The issue will be
   referenced in the commit message so that the code that fixes it is clearly
   linked.
1. The issue is closed.

## Building Waypoint

If you wish to work on Waypoint itself, you'll first need [Go](https://golang.org)
installed (version 1.14+ is _required_).

[go-bindata](https://github.com/go-bindata/go-bindata) is a binary dependency
that must be on your PATH to build Waypoint. This 
[repository version](https://github.com/kevinburke/go-bindata/) may be installed with:
`brew install go-bindata`

Next, clone this repository and then run the following commands:
* `make bin` will build the binary for your local machine's os/architecture
* (optional) `make install` will copy that executable binary to `$GOPATH/bin/waypoint`
* `make docker/server` will build the docker image of the server with the tag `waypoint:dev`

Once those steps are complete, you can install the waypoint server you just built. To do
this on docker you would run:
```
waypoint install -platform=docker -accept-tos -docker-server-image=waypoint:dev
```

>Note: If you didn't run `make install` then you should use `path/to/waypoint` 
in place of `waypoint`.

## Making Changes to Waypoint

>Note: See [Issue Lifecycle](#issue-lifecycle) for more info on recognizing when issues are already
a work-in-progress

Run `make tools` to install the list of tools in ./tools/tools.go.
>Note: If notice you have a large set of diffs due to upgrading the version of 
>a tool, it is best to separate out the upgrade into its own PR.

The first step to making changes is to fork Waypoint. Afterwards, the easiest way
to work on the fork is to set it as a remote of the Waypoint project:

1. Navigate to `$GOPATH/src/github.com/hashicorp/waypoint`
2. Rename the existing remote's name: `git remote rename origin upstream`.
3. Add your fork as a remote by running
   `git remote add origin <github url of fork>`. For example:
   `git remote add origin https://github.com/myusername/waypoint`.
4. Checkout a feature branch: `git checkout -t -b new-feature`
5. Make changes
6. Push changes to the fork when ready to submit PR:
   `git push -u origin new-feature`

By following these steps you can push to your fork to create a PR, but the code on disk still
lives in the spot where the go cli tools are expecting to find it.

If the scope of the code change requires it, follow the [Changelog Guide](/.github/CHANGELOG_GUIDE.md) to add an entry.

>Note: If you make any changes to the code, run `make format` to automatically format the code according to Go standards.

## Opening a PR

1. Title the PR with a helpful prefix following the pattern `parent/child`, 
e.g. `builtin/k8s` or `internal/core`
1. Include helpful information in the description
   * For a bug fix, explain or show an example of the behavior before and 
  after the change
   * If applicable, include information on how to test manually
1. Request review from either `waypoint-core` or `waypoint-frontend` based on 
your changes

>Note: the auto-labeler will assign other labels after you open the PR, based 
>on what files have changes.

## Testing

Before submitting changes, run **all** tests locally by typing `make test`.
The test suite may fail if over-parallelized, so if you are seeing stochastic
failures try `GOTEST_FLAGS="-p 2 -parallel 2" make test`.
