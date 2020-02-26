package docker

// Name is the full name including the tag.
func (i *Image) Name() string {
	return i.Image + ":" + i.Tag
}
