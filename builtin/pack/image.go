package pack

func (i *DockerImage) Labels() map[string]string {
	return i.BuildLabels
}
