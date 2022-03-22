{ buildGoModule, fetchFromGitHub, lib }:

buildGoModule rec {
  pname = "go-protobuf";
  version = "1.27.1";

  src = fetchFromGitHub {
    owner = "protocolbuffers";
    repo = "protobuf-go";
    rev = "v${version}";
    sha256 = "sha256-wkUvMsoJP38KMD5b3Fz65R1cnpeTtDcVqgE7tNlZXys=";
  };

  modSha256 = "sha256-yb8l4ooZwqfvenlxDRg95rqiL+hmsn0weS/dPv/oD2Y=";
  vendorSha256 = "sha256-yb8l4ooZwqfvenlxDRg95rqiL+hmsn0weS/dPv/oD2Y=";

  subPackages = [ "cmd/protoc-gen-go" ];
}
