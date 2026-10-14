package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"os"
)

func main() {
	mo := flag.NewFlagSet("预算管家", flag.ContinueOnError)
	moNew := mo.BoolP("new", "n", false, "新建一个账本")
	mo.Lookup("new").NoOptDefVal = "true"
	err := mo.Parse(os.Args[1:])
	if err == flag.ErrHelp {
		fmt.Println("预算管家 - 您的终端预算会计v0.0.1")
		os.Exit(1)
	}
	fmt.Println(*moNew)
}
