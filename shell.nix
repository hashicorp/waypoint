let
  # First we setup our overlays. These are overrides of the official nix packages.
  # We do this to pin the versions we want to use of the software that is in
  # the official nixpkgs repo.
  pkgs = import ./nix;
in with pkgs; let
  go-protobuf = buildGoModule rec {
    pname = "go-protobuf";
    version = "v1.4.2";

    src = fetchFromGitHub {
      owner = "golang";
      repo = "protobuf";
      rev = "v1.4.2";
      sha256 = "0m5z81im4nsyfgarjhppayk4hqnrwswr3nix9mj8pff8x9jvcjqw";
    };

    modSha256 = "0lnk1zpl6y9vnq6h3l42ssghq6iqvmixd86g2drpa4z8xxk116wf";
    vendorSha256 = "04w9vhkrwb2zfqk73xmhignjyvjqmz1j93slkqp7v8jj2dhyla54";

    subPackages = [ "protoc-gen-go" ];
  };

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
    vendorSha256 = "1hx31gr3l2f0nc8316c9ipmk1xx435g732msr5b344rcfcfrlaxh";
  };

  go-tools = buildGoModule rec {
    pname = "go-tools";
    version = "35839b7038afa36a6c000733552daa1f5ce1e838";

    src = fetchFromGitHub {
      owner = "golang";
      repo = "tools";
      rev = "35839b7038afa36a6c000733552daa1f5ce1e838";
      sha256 = "1gnqf62s7arqk807gadp4rd2diz1g0v2khwv9wsb50y8k9k4dfqs";
    };

    modSha256 = "1pijbkp7a9n2naicg21ydii6xc0g4jm5bw42lljwaks7211ag8k9";
    vendorSha256 = "0i2fhaj2fd8ii4av1qx87wjkngip9vih8v3i9yr3h28hkq68zkm5";

    subPackages = [ "cmd/stringer" ];

    # This has to be enabled because the stringer tests recompile itself
    # so it needs a valid reference to `go`
    allowGoReference = true;
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

    buildFlagsArray = ''
      -ldflags=
      -s -w -X github.com/vektra/mockery/mockery.SemVer=${version}
    '';

    modSha256 = "0wyzfmhk7plazadbi26rzq3w9cmvqz2dd5jsl6kamw53ps5yh536";
    vendorSha256 = "0fai4hs3q822dg36a2zrxb191f71xdpafapn6ymi1w9dx68navcb";

    subPackages = [ "cmd/mockery" ];
  };

  go-changelog = buildGoModule rec {
    pname = "go-changelog";
    version = "56335215ce3a8676ba7153be7c444daadcb132c7";

    src = fetchFromGitHub {
      owner = "hashicorp";
      repo = "go-changelog";
      rev = "56335215ce3a8676ba7153be7c444daadcb132c7";
      sha256 = "0z6ysz4x1rim09g9knbc5x5mrasfk6mzsi0h7jn8q4i035y1gg2j";
    };

    vendorSha256 = "1pahh64ayr885kv9rd5i4vh4a6hi1w583wch9n1ncvnckznzsdbg";

    subPackages = [ "cmd/changelog-build" ];
  };
in pkgs.mkShell rec {
  name = "waypoint";

  # The packages in the `buildInputs` list will be added to the PATH in our shell
  buildInputs = [
    pkgs.docker-compose
    pkgs.go
    pkgs.go-bindata
    pkgs.grpcurl
    pkgs.niv
    pkgs.nodejs-12_x
    pkgs.protobuf3_11
    pkgs.postgresql_12
    go-protobuf
    go-protobuf-json
    go-tools
    go-mockery
    go-changelog
  ] ++ (with pkgs; [
    # Needed for website/
    pkgconfig autoconf automake libtool nasm autogen zlib libpng
  ]) ++ (if stdenv.isLinux then [
    # On Linux we use minikube as the primary k8s testing platform
    pkgs.minikube
  ] else []);

  # workaround for npm/gulp dep compilation
  # https://github.com/imagemin/optipng-bin/issues/108
  shellHook = ''
    LD=$CC
  '';

  # Extra env vars
  PGHOST = "localhost";
  PGPORT = "5432";
  PGDATABASE = "noop";
  PGUSER = "postgres";
  PGPASSWORD = "postgres";
  DATABASE_URL = "postgresql://${PGUSER}:${PGPASSWORD}@${PGHOST}:${PGPORT}/${PGDATABASE}?sslmode=disable";
}
