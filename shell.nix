{ pkgsPath ? <nixpkgs> }:

let
  # First we setup our overlays. These are overrides of the official nix packages.
  # We do this to pin the versions we want to use of the software that is in
  # the official nixpkgs repo.
  pkgs = import pkgsPath {
    overlays = [(self: super: {

      go = super.go.overrideAttrs ( old: rec {
        version = "1.14.5";
        src = super.fetchurl {
          url = "https://dl.google.com/go/go${version}.src.tar.gz";
          sha256 = "0p1i80j3dk597ph5h6mvvv8p7rbzwmxdfb6558amcpkkj060hk6a";
        };
      });

      go-protobuf = super.go-protobuf.overrideAttrs ( old: rec {
        version = "1.3.5";
        src = super.fetchFromGitHub {
          owner = "golang";
          repo = "protobuf";
          rev = "v${version}";
          sha256 = "1gkd1942vk9n8kfzdwy1iil6wgvlwjq7a3y5jc49ck4lz9rhmgkq";
        };

        modSha256 = "0jjjj9z1dhilhpc8pq4154czrb79z9cm044jvn75kxcjv6v5l2m5";
      });

    })];
  };
in with pkgs; let
  go-protobuf-json = buildGoModule rec {
    pname = "go-protobuf-json";
    version = "069933b8c8344593ed8905d46d59c6647c886f47";

    src = fetchFromGitHub {
      owner = "mitchellh";
      repo = "protoc-gen-go-json";
      rev = "069933b8c8344593ed8905d46d59c6647c886f47";
      sha256 = "1q5s2pfdxxzvdqghmbw3y2w5nl7wa4x15ngahfarjhahwqsbfsx4";
    };

    modSha256 = "01wrk2qhrh74nkv6davfifdz7jq6fcl3snn4w2g7vr8p0incdlcf";
  };

  go-tools = buildGoModule rec {
    pname = "go-tools";
    version = "57a9e4404bf7b38f22bbca9af3ddf0dee8e76a04";

    src = fetchFromGitHub {
      owner = "golang";
      repo = "tools";
      rev = "57a9e4404bf7b38f22bbca9af3ddf0dee8e76a04";
      sha256 = "1zih0v855vkr5j1rvahbbfd1w7rjf5rrgm20ra0b34nw7656x88h";
    };

    modSha256 = "1pijbkp7a9n2naicg21ydii6xc0g4jm5bw42lljwaks7211ag8k9";

    subPackages = [ "cmd/stringer" ];
  };

  go-mockery = buildGoModule rec {
    pname = "go-mockery";
    version = "1.1.2";

    src = fetchFromGitHub {
      owner = "vektra";
      repo = "mockery";
      rev = "v${version}";
      sha256 = "16yqhr92n5s0svk31yy3k42764fas694mnqqcny633yi0wqb876a";
    };

    modSha256 = "0wyzfmhk7plazadbi26rzq3w9cmvqz2dd5jsl6kamw53ps5yh536";

    subPackages = [ "cmd/mockery" ];
  };
in pkgs.mkShell rec {
  name = "waypoint";

  # The packages in the `buildInputs` list will be added to the PATH in our shell
  buildInputs = [
    pkgs.go
    pkgs.go-bindata
    pkgs.go-protobuf
    pkgs.protobuf3_11
    pkgs.postgresql_12
    go-protobuf-json
    go-tools
    go-mockery
  ];

  # Extra env vars
  PGHOST = "localhost";
  PGPORT = "5432";
  PGDATABASE = "noop";
  PGUSER = "postgres";
  PGPASSWORD = "postgres";
  DATABASE_URL = "postgresql://${PGUSER}:${PGPASSWORD}@${PGHOST}:${PGPORT}/${PGDATABASE}?sslmode=disablie";
}
