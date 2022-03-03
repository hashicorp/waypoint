package lambda

import (
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/hashicorp/go-hclog"
)

// DockerArchitectureMapper maps a docker image architecture to a valid lambda architecture.
func DockerArchitectureMapper(src string, log hclog.Logger) string {
	switch src {
	case "amd64", "x86_64":
		return lambda.ArchitectureX8664
	case "arm64", "aarch64":
		return lambda.ArchitectureArm64
	default:
		log.Warn("unsupported docker architecture", "arch", src, "defaulting to:", lambda.ArchitectureX8664)
		return lambda.ArchitectureX8664
	}
}
