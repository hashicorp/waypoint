let
  sources = import ./sources.nix;
in
import sources.nixpkgs {
  overlays = [
    (self: super: {
      # niv = (import sources.niv { pkgs = self; }).niv;
      go = super.go_1_14.overrideAttrs (
        old: rec {
          version = "1.14.5";
          src = super.fetchurl {
            url = "https://dl.google.com/go/go${version}.src.tar.gz";
            sha256 = "0p1i80j3dk597ph5h6mvvv8p7rbzwmxdfb6558amcpkkj060hk6a";
          };
        }
      );
    })
  ];
}
