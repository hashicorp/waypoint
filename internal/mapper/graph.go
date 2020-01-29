package mapper

// Graph represents a set of M values. A Graph can be used to find or test
// a path through mappers to go from type A to B to ... F where each transition
// requires another mapper.
type Graph struct{}

// Path returns a path through the graph to get the desired values.
//
// An even count of values should be given with each pair being an M
// and a desired string name associated with that M. This represents the
// values you know ahead of time that you want to have.
//
// You must also provide the list of values that you have access to now,
// this is the source set of values. List these values after the (M, name)
// pairs. Once the first non-M even-index argument is found, that value plus
// all following are expected to be input values.
//
// Example:
//
//   Path(m1, "a", m2, "b", m3, "c", input1, input2)
//
//
// A non-nil error will be returned if no paths are found or if multiple
// paths are found of the same length. If multiple paths are found, a shortest
// path is always assumed to be the valid value if that shortest path is
// a subset of the longer path.
func (g *Graph) Path(values ...interface{}) (*Path, error) {
	return nil, nil
}

type Path struct{}

func (p *Path) Target(m *M, n string) (interface{}, error) {
	return nil, nil
}

/*
Path(
	BuilderM.Type(), "pack",
	RegistryM.Type(), "docker",
	PlatformM.Type(), "gcr",
)
*/
