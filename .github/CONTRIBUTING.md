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
* Fix a [bug](https://github.com/hashicorp/waypoint/labels/type/bug).
* Implement a requested [enhancement](https://github.com/hashicorp/waypoint/labels/type/enhancement).
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

* Aim to respond promptly to any questions made by the Waypoint team on your
issue. Stale issues will be closed.

### Issue Lifecycle

1. The issue is reported.

2. The issue is verified and categorized by a Waypoint maintainer.
   Categorization is done via tags. For example, bugs are tagged as "bug".

3. Unless it is critical, the issue is left for a period of time (sometimes many
   weeks or months), giving outside contributors a chance to address the issue
   and our internal teams time to plan for inclusion in a release.

4. The issue is addressed in a pull request or commit. The issue will be
   referenced in the commit message so that the code that fixes it is clearly
   linked.

5. The issue is closed.

## Building Waypoint

If you wish to work on Waypoint itself, you'll first need [Go](https://golang.org)
installed (version 1.14+ is _required_).

Next, clone this repository and then run `make bin`. In a few moments, you'll have a working
`waypoint` executable in the Waypoint repository root and `$GOPATH/bin`:

>Note: `make bin` will build for your local machine's os/architecture.

## Making Changes to Waypoint

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

>Note: If you make any changes to the code, run `make format` to automatically format the code according to Go standards.

## Testing

Before submitting changes, run **all** tests locally by typing `make test`.
The test suite may fail if over-parallelized, so if you are seeing stochastic
failures try `GOTEST_FLAGS="-p 2 -parallel 2" make test`.
