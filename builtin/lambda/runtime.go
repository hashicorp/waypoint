package lambda

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type RuntimeConfig struct {
	Language     string
	DepManager   string
	AppFramework string
	ManifestName string
	SearchPaths  []string
}

var RuntimeConfigs = []RuntimeConfig{
	{
		Language:     "python",
		DepManager:   "pip",
		ManifestName: "requirements.txt",
	},
	{
		Language:     "nodejs",
		DepManager:   "npm",
		ManifestName: "package.json",
	},
	{
		Language:     "ruby",
		DepManager:   "bundler",
		ManifestName: "Gemfile",
	},
	{
		Language:     "java",
		DepManager:   "gradle",
		ManifestName: "build.gradle",
	},
	{
		Language:     "java",
		DepManager:   "gradle",
		ManifestName: "build.gradle.kts",
	},
	{
		Language:     "java",
		DepManager:   "maven",
		ManifestName: "pom.xml",
	},
	{
		Language:     "dotnet",
		DepManager:   "cli-package",
		ManifestName: ".csproj",
	},
	{
		Language:     "go",
		DepManager:   "modules",
		ManifestName: "go.mod",
	},
}

type Runtime struct {
	Name     string
	Language string
}

var Runtimes = []Runtime{
	{"python2.7", "python"},
	{"python3.6", "python"},
	{"python3.7", "python"},
	{"python3.8", "python"},
	{"nodejs4.3", "nodesjs"},
	{"nodejs6.10", "nodesjs"},
	{"nodejs8.10", "nodesjs"},
	{"nodejs10.x", "nodesjs"},
	{"nodejs12.x", "nodesjs"},
	{"ruby2.5", "ruby"},
	{"ruby2.7", "ruby"},
	{"dotnetcore2.0", "dotnet"},
	{"dotnetcore2.1", "dotnet"},
	{"go1.x", "go"},
	{"java8", "java"},
	{"java11", "java"},
}

var (
	ErrUnknownRuntime = errors.New("unknown runtime")
	ErrNoConfig       = errors.New("no compatible configuration detected")
)

func FindRuntimeConfig(runtime string, srcDir string) (*RuntimeConfig, error) {
	var lang string

	for _, rt := range Runtimes {
		if rt.Name == runtime {
			lang = rt.Language
			break
		}
	}

	if lang == "" {
		return nil, errors.Wrapf(ErrUnknownRuntime, "specified runtime: %s", runtime)
	}

	var attempts []string

	for _, rtc := range RuntimeConfigs {
		if rtc.Language != lang {
			continue
		}

		if _, err := os.Stat(filepath.Join(srcDir, rtc.ManifestName)); err == nil {
			return &rtc, nil
		}

		attempts = append(attempts, rtc.ManifestName)
	}

	return nil, errors.Wrapf(ErrNoConfig, "no manifest document detected amongst: %s", strings.Join(attempts, ", "))
}
