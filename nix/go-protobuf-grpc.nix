{ buildGoModule, fetchFromGitHub, lib }:

buildGoModule rec {
  pname = "go-protobuf-grpc";
  version = "1.2.0";

  src = fetchFromGitHub {
    owner = "grpc";
    repo = "grpc-go";
    rev = "cmd/protoc-gen-go-grpc/v${version}";
    sha256 = "sha256-pIHMykwwyu052rqwSV5Js0+JCKCNLrFVN6XQ3xjlEOI=";
  };

  modSha256 = "sha256-yxOfgTA5IIczehpWMM1kreMqJYKgRT5HEGbJ3SeQ/Lg=";
  vendorSha256 = "sha256-yxOfgTA5IIczehpWMM1kreMqJYKgRT5HEGbJ3SeQ/Lg=";

  modRoot = "cmd/protoc-gen-go-grpc";
}
