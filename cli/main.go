package main

import (
	"fmt"
	"os"
	"path"
	"github.com/da99/cli.go/args"
	"github.com/da99/template.go/template"
	"github.com/da99/cli.go/exit"
	// go_tmpl "html/template"
	// "github.com/da99/files.go/files"
)


func main() {

	if args.IsMatch("-h", 0) {
		fmt.Println("-h -- Print this message.")
		fmt.Println("ls [dirs|files] [dir]")
		fmt.Println("compile [dir|file]")
		os.Exit(0)
	}

	if args.IsMatch("ls dirs", 1) {
		lines := template.List_Dirs(args.CAPTURE[0])
		for _, v := range lines { fmt.Println(v) }
		os.Exit(0)
	}

	if args.IsMatch("ls files", 1) {
		files := template.List_Template_Files(args.CAPTURE[0])
		for _, v := range files { fmt.Println(v) }
		os.Exit(0)
	}

	if args.IsMatch("compile", 1) {
		target := args.CAPTURE[0]
		stat, err := os.Stat(target)
		exit.PrintError(err)

		if !stat.IsDir() {
			target = path.Dir(target)
		}

		cerr := template.Compile_Dir(target)
		exit.PrintError(cerr)

		os.Exit(0)
	}

	args.Fail()
}
