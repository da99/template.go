/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"fmt"
	"github.com/da99/cli.go/args"
	"github.com/da99/files.go/files"
	"os"
	"os/exec"
	"strings"
	"path"
	go_tmpl "html/template"
	"github.com/da99/template.go/template"
)

func PrintError(x error) {
	if x != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", x)
		os.Exit(1)
	}
}

func shell_lines(cmd_str string) []string {
	raw := exec.Command("bash", "-c", cmd_str)

	output, o_err := raw.Output()
	PrintError(o_err)

	return strings.Split(strings.TrimSpace(string(output)), "\n")
}

func must_exist(str_path string) bool {
	if !files.Is(str_path) {
		fmt.Fprintf(os.Stderr, "Does not exist: %v\n", str_path)
		os.Exit(1)
	}
	return true
}

func t_files(str_dir string) []string {
	must_exist(str_dir)
	return shell_lines("find " + str_dir + " -type f -name '*.go.html' -and -not -name '*.partial.go.html' | sort")
}

func t_dirs(str_dir string) []string {
	must_exist(str_dir)
	return shell_lines("find " + str_dir + " -type f -name '*.go.html' -and -not -name '*.partial.go.html' | xargs dirname | sort | uniq")
}

func main() {

	if args.IsMatch("-h", 0) {
		fmt.Println("-h -- Print this message.")
		fmt.Println("ls [dirs|files] [dir]")
		fmt.Println("compile [dir|file]")
		os.Exit(0)
	}

	if args.IsMatch("ls dirs", 1) {
		lines := t_dirs(args.CAPTURE[0])
		for _, v := range lines { fmt.Println(v) }
		os.Exit(0)
	}

	if args.IsMatch("ls files", 1) {
		files := t_files(args.CAPTURE[0])
		for _, v := range files { fmt.Println(v) }
		os.Exit(0)
	}

	if args.IsMatch("compile", 1) {
		target := args.CAPTURE[0]
		stat, err := os.Stat(target)
		PrintError(err)

		if !stat.IsDir() {
			target = path.Dir(target)
		}

		all_dirs := t_dirs(target)
		config, c_err := template.GetConfig()
		PrintError(c_err)

		for _, d := range all_dirs {
			all_files, l_err := files.List_Shallow_Files_Ext(d, "*.go.html")
			PrintError(l_err)

			tmpl, t_err := go_tmpl.ParseFiles(all_files...)
			PrintError(t_err)

			for _, f := range all_files {
				if template.IsPartial(f) { continue; }

				fmt.Printf("-- Compiling template: %v\n", f)
				err := tmpl.ExecuteTemplate(os.Stdout, path.Base(f), config)
				PrintError(err)
				fmt.Println("")
			}
		}
		os.Exit(0)
	}

	args.Fail()
}
