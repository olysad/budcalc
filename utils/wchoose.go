/*
Copyright © 2026 olysad <lacak3736@gmail.com>
*/
package utils

import (
	"fmt"
	"github.com/muesli/termenv"
	"os"
)

func Wchoose() {
	var i string
	o := termenv.NewOutput(os.Stdout)
	t1 := o.String("Warning!").Foreground(o.Color("#ff0000")).Bold()
	t2 := o.String("危险操作！您确认要继续吗？").Foreground(o.Color("#FFFF00"))
	t3 := o.String("[Y/n]").Foreground(o.Color("#00ff00"))

	var tag bool
	for !tag {
		fmt.Printf("%v%v%v", t1, t2, t3)
		fmt.Scanln(&i)
		if i == "n" {
			os.Exit(0)
		} else if i == "Y" {
			tag = true
		} else {
			tag = false
		}
	}
}
