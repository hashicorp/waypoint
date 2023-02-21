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

  modSha256 = "sha256-AW2Gn/mlZyLMwF+NpK59eiOmQrYWW/9HPjbunYc9Ij4=";
  vendorSha256 = "sha256-AW2Gn/mlZyLMwF+NpK59eiOmQrYWW/9HPjbunYc9Ij4=";

  subPackages = [ "protoc-gen-grpc-gateway" ];
}
