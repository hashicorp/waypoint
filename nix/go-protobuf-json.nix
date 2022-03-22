{ buildGoModule, fetchFromGitHub, lib }:

buildGoModule rec {
  pname = "go-protobuf-json";
  version = "45822525aa9c2f948c1da1d0724fe042dd7a20e4";

  src = fetchFromGitHub {
    owner = "mitchellh";
    repo = "protoc-gen-go-json";
    rev = version;
    sha256 = "sha256-ifcof7GPV33ABFVAh9og8QML4MJQWbqC/E5wUL0QmJE=";
  };

  modSha256 = "sha256-bbSYkD+KqpGUNzyANR7OMq4wjnBNDZHOlQ6bGZSoo3A=";
  vendorSha256 = "sha256-bbSYkD+KqpGUNzyANR7OMq4wjnBNDZHOlQ6bGZSoo3A=";
}
