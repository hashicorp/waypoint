{ buildGoModule, fetchFromGitHub }:

buildGoModule rec {
  pname = "go-protobuf-swagger";
  version = "48e5be386c15aed147eeb9787f21e54faaacd531";

  src = fetchFromGitHub {
    owner = "evanphx";
    repo = "grpc-gateway";
    rev = "48e5be386c15aed147eeb9787f21e54faaacd531";
    sha256 = "sha256-4jZAsn5DjBI+qVS750aTOPUb7fOtKfuC4sOcwdcr5nY=";
  };

  modSha256 = "sha256-iCeuVuRh00of65Oe1XDFZUQM0PdTCoUOBli+oEaXtg8=";
  vendorSha256 = "sha256-iCeuVuRh00of65Oe1XDFZUQM0PdTCoUOBli+oEaXtg8=";

  subPackages = [ "protoc-gen-swagger" ];
}
