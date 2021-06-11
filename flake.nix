{
  description = "HashiCorp Waypoint project";

  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        overlay = (import ./nix/overlay.nix) nixpkgs;

        pkgs = import nixpkgs {
          inherit system;
          overlays = [ overlay ];
        };

        waypoint = pkgs.callPackage ./nix/waypoint.nix {
          inherit pkgs;
        };
      in {
        devShell = waypoint.shell;
      }
    );
}
