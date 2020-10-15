name: "Check nix shell and populate cache"
on:
  pull_request:
  push:
jobs:
  tests:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2.3.3
    - uses: cachix/install-nix-action@v11
    # TODO: Uncomment if you want a cachix cache
    # - uses: cachix/cachix-action@v6
    #   with:
    #     name: waypoint
    #     signingKey: '${{ secrets.CACHIX_SIGNING_KEY }}'
    #     # Only needed for private caches
    #     #authToken: '${{ secrets.CACHIX_AUTH_TOKEN }}'
    - run: nix-shell --run "echo OK"
