package k8s

import (
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/r3labs/diff"

	"testing"
)

func TestServicePortCoersion(t *testing.T) {
	var defaultPort uint = 9999
	var config Config

	cases := []struct {
		id     string
		config string
		expect []uint
	}{
		{"nil", "", []uint{defaultPort}},
		{"integer", `service_port = 88`, []uint{88}},
		{"array", `service_port = [88, 99]`, []uint{88, 99}},
	}

	for _, c := range cases {
		t.Run(c.id, func(t *testing.T) {
			// Should populate defaultPort value for unspecified config
			err := hclsimple.Decode("config.hcl", []byte(c.config), nil, &config)
			if err != nil {
				t.Fatal("Decode error: ", err)
			}

			ports, err := coerceConfigPortArray(config.ServicePort, defaultPort)
			if err != nil {
				t.Fatal("Coerce error: ", err)
			} else if d, _ := diff.Diff(c.expect, ports); len(d) > 0 {
				t.Errorf("Outcomes differ: %+v", d)
			}
		})
	}

}
