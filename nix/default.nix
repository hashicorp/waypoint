let
  sources = import ./sources.nix;
in
import sources.nixpkgs {
  overlays = [
    (self: super: {
      # niv = (import sources.niv { pkgs = self; }).niv;
      go = super.go_1_15.overrideAttrs (
        old: rec {
          version = "1.15.3";
          src = super.fetchurl {
            url = "https://dl.google.com/go/go${version}.src.tar.gz";
            sha256 = "1228nv4vyzbqv768dl0bimsic47x9yyqld61qbgqqk75f0jn0sl9";
          };
        }
      );
    })
  ];
}
