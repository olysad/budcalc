/*
Copyright © 2026 olysad <lacak3736@gmail.com>
*/
package utils

import (
	"bytes"
	"fmt"
	"github.com/muesli/termenv"
	"os"
	"text/template"
)

func Wchoose() {
	var i string
	o := termenv.NewOutput(os.Stdout)
	f := o.TemplateFuncs()

	str := `{{- Bold ( Foreground "#ff0000" "Warning!危险操作！" ) -}}{{ Foreground "#00ff00" "您确认要继续吗？[Y/n]" -}}`

	tpl, _ := template.New("w").Funcs(f).Parse(str)

	var tag bool
	for !tag {
		var p bytes.Buffer
		_ = tpl.Execute(&p, nil)
		fmt.Printf("\x1b[1A\r\x1b[K%v", &p)

		fmt.Scanln(&i)
		switch i {
		case "n":
			os.Exit(0)
		case "Y":
			tag = true
		default:
			tag = false
		}
	}
}
