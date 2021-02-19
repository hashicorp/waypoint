let
  sources = import ./sources.nix;
in
import sources.nixpkgs {
  overlays = [
    (self: super: {
      # niv = (import sources.niv { pkgs = self; }).niv;

      go = super.go_1_16;
      buildGoModule = super.buildGo116Module;

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
