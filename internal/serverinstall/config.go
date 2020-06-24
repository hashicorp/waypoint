package serverinstall

import (
	"path"
	"strings"
	"text/template"

	"github.com/hashicorp/waypoint/internal/serverinstall/datagen"
)

// Config is the configuration for the Kubernetes installation template.
type Config struct {
	Namespace          string
	ServiceName        string
	ServerImage        string
	ServiceAnnotations map[string]string
	ImagePullSecret    string
}

// Render renders the installation files with the given configuration.
func Render(config *Config) (string, error) {
	const prefix = "k8s-install"

	files, err := datagen.AssetDir(prefix)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	for _, name := range files {
		bs, err := datagen.Asset(path.Join(prefix, name))
		if err != nil {
			return "", err
		}

		tmpl, err := template.New("root").Parse(string(bs))
		if err != nil {
			return "", err
		}

		if err := tmpl.Execute(&result, config); err != nil {
			return "", err
		}
	}

	return result.String(), nil
}
