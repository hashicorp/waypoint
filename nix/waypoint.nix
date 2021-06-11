{ pkgs }: {
  shell = pkgs.mkShell rec {
    name = "waypoint";

    buildInputs = [
      pkgs.docker-compose
      pkgs.go
      pkgs.go-bindata
      pkgs.grpcurl
      pkgs.nodejs-12_x
      pkgs.postgresql_12
      pkgs.protoc-gen-doc

      # Custom packages
      pkgs.protobufPin
      pkgs.go-protobuf
      pkgs.go-protobuf-json
      pkgs.go-tools
      pkgs.go-mockery
      pkgs.go-changelog
    ] ++ (with pkgs; [
      # Needed for website/
      pkgconfig autoconf automake libtool nasm autogen zlib libpng
    ]) ++ (if pkgs.stdenv.isLinux then [
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
  };
}
