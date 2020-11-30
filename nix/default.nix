let
  sources = import ./sources.nix;
in
import sources.nixpkgs {
  overlays = [
    (self: super: {
      # niv = (import sources.niv { pkgs = self; }).niv;
      go = super.go_1_15.overrideAttrs (
        old: rec {
          version = "1.15.5";
          src = super.fetchurl {
            url = "https://dl.google.com/go/go${version}.src.tar.gz";
            sha256 = "1wc43h3pmi92r6ypmh58vq13vm44rl1di09assz3xdwlry86n1y1";
          };
        }
      );
    })
  ];
}
