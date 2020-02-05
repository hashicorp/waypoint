package cli

type UpCommand struct {
}

func (c *UpCommand) Run([]string) int {
	return 0
}

func (c *UpCommand) Synopsis() string {
	return ""
}

func (c *UpCommand) Help() string {
	return ""
}
