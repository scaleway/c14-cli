package commands

import "fmt"

type test struct {
	Base
	testFlags
}

type testFlags struct {
}

// test returns a new command "test"
func Test() Command {
	ret := &test{}
	ret.Init(Config{
		UsageLine:   "test",
		Description: "",
		Help:        "",
		Examples:    "",
	})
	return ret
}

func (t *test) GetName() string {
	return "test"
}

func (t *test) Run(args []string) (err error) {
	if err = t.InitAPI(); err != nil {
		return err
	}
	keys, err := t.OnlineAPI.GetSSHKeys()
	fmt.Println(keys, err)
	return
}
