package ptypes

import (
	"strings"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// TestProject returns a valid project for tests.
func TestProject(t testing.T, src *pb.Project) *pb.Project {
	t.Helper()

	if src == nil {
		src = &pb.Project{}
	}

	require.NoError(t, mergo.Merge(src, &pb.Project{
		Name: "test",
	}))

	return src
}

// Type wrapper around the proto type so that we can add some methods.
type Project struct{ *pb.Project }

// App returns the index of the app with the given name or -1 if its not found.
func (p *Project) App(n string) int {
	n = strings.ToLower(n)
	for i, app := range p.Applications {
		if strings.ToLower(app.Name) == n {
			return i
		}
	}

	return -1
}

// ValidateProject validates the project structure.
func ValidateProject(p *pb.Project) error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Name, validation.By(isEmpty)),
	)
}
