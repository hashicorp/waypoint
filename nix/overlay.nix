final: prev: {
  # This is the pinned protoc version we have for this project.
  protobufPin = prev.protobuf3_15;

  devShell = final.callPackage ./waypoint.nix { };

  # Need to manually do this since 1.16 is still the default
  go = final.go_1_17;

  go-protobuf = prev.callPackage ./go-protobuf.nix { };

  go-protobuf-json = prev.callPackage ./go-protobuf-json.nix { };

  go-tools = prev.callPackage ./go-tools.nix { };

  go-mockery = prev.callPackage ./go-mockery.nix { };

  go-changelog = prev.callPackage ./go-changelog.nix { };
}
