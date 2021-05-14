let
  sources = import ./sources.nix;
in
import sources.nixpkgs {
  overlays = [
    (self: super: {
      # niv = (import sources.niv { pkgs = self; }).niv;

      go = super.go_1_16;
      buildGoModule = super.buildGo116Module;

      # This is the pinned protoc version we have for this project.
      protobufPin = super.callPackage (
        sources.nixpkgs + "/pkgs/development/libraries/protobuf/generic-v3.nix"
      ) {
        version = "3.15.8";
        sha256 = "1q3k8axhq6g8fqczmd6kbgzpdplrrgygppym4x1l99lzhplx9rqv";
      };

      /* in case we need to override again:
      go = super.go_1_15.overrideAttrs (
        old: rec {
          version = "1.15.7";
          src = super.fetchurl {
            url = "https://dl.google.com/go/go${version}.src.tar.gz";
            sha256 = "8631b3aafd8ecb9244ec2ffb8a2a8b4983cf4ad15572b9801f7c5b167c1a2abc";
          };
        }
      );
      */
    })
  ];
}
