## Overview

HashiCorp Waypoint is an easy to use tool that standardizes application build, deploy, and release workflows as code. It focuses on delivering applications in a consistent and repeatable way, reducing time to deploy and allowing developers to ship applications across common developer platforms.

Waypoint uses:

- One configuration file, one common language, as code
- An end-to-end workflow to build, deploy, and release for applications

You can also use Waypoint to validate deployments across distinct environments through common execution and logging.

## Documentation & Guides

Documentation is available on the Waypoint website [here](https://waypointproject.io/).

## Development environment

This repo uses [Nix](https://nixos.org/) for managing development dependencies.

To get started, [install Nix](https://nixos.org/download.html) (macOS users, follow [these instructions](https://nixos.org/manual/nix/stable/#sect-macos-installation)).

Be sure to source `$HOME/.nix-profile/etc/profile.d/nix.sh` before continuing (this is added to ~/.bash_profile by default; if you use a different shell, be sure to configure your shell to source that file on start).

Verify nix works:

```
nix-env -i hello
hello
# Hello, world!

# uninstall hello
nix-env --rollback
```

Install `direnv` with nix:

```
nix-env -i direnv
```

Add the [appropriate direnv hook](https://direnv.net/docs/hook.html) to your shell:

```
# add this to ~/.bash_profile or similar for your shell
eval "$(direnv hook bash)"
```

_Note: Start a new shell before proceeding._

In the directory you cloned this repo to, allow `direnv` to load this project's `.envrc` file:

```
# in waypoint git dir
direnv allow
```

At this point, `direnv`'s [`nix` hook](https://github.com/direnv/direnv/wiki/Nix) will be run, and your environment will be set up.

Verify your environment has been configured correctly:

```
which go # should return a path including /nix/store
```
>>>>>>> 90e8cd85... readme: Add Development environment section
