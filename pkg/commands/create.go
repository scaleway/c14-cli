package commands

type create struct {
	Base
}

func Create() Command {
	return &create{}
}

func (c *create) GetName() string {
	return "create"
}

func (c *create) Parse() error {
	return nil
}

func (c *create) Run() error {
	return nil
}
