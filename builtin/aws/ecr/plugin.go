package ecr

func (i *Image) Name() string {
	return i.Image + ":" + i.Tag
}
