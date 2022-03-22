{ buildGoModule, fetchFromGitHub, lib }:

buildGoModule rec {
  pname = "go-protobuf-grpc-gateway";
  version = "2.8.0";

  src = fetchFromGitHub {
    owner = "grpc-ecosystem";
    repo = "grpc-gateway";
    rev = "v${version}";
    sha256 = "sha256-8eBBBYJ+tBjB2fgPMX/ZlbN3eeS75e8TAZYOKXs6hcg=";
  };

  modSha256 = "sha256-8XbFKsgmMcf363W/F1Ffh1eEh/M3NGg0zzLCZ4b5Dho=";
  vendorSha256 = "sha256-8XbFKsgmMcf363W/F1Ffh1eEh/M3NGg0zzLCZ4b5Dho=";

  subPackages = [ "protoc-gen-grpc-gateway" ];
}
