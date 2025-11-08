package template

import (
	"fmt"

	"github.com/da99/cli.go/run"
	"github.com/da99/cli.go/exit"
	"github.com/da99/cli.go/config"
	"github.com/da99/cli.go/files"
	"os"
	"path/filepath"
	go_template "html/template"
	"strings"
	"slices"
)

const DOT_PARTIAL = ".partial.go.html"
const DOT_LAYOUT  = ".layout.go.html"

func Remove_Dot_Go(raw_path string) string {
	return strings.Replace(raw_path, ".go.html", ".html", 1)
}

func Is_Partial(str_path string) bool {
	return strings.LastIndex(str_path, DOT_PARTIAL) > 1
}

func Is_Layout(str_path string) bool {
	return strings.Contains(str_path, "/layout/") ||
	strings.Contains(str_path, "/layouts/") ||
	strings.LastIndex(str_path, DOT_LAYOUT) > 1
}

func List_Templates_In_Dir(dir string) []string {
	s_args := fmt.Sprintf("-mindepth 1 -maxdepth 1 -type f -name *.go.html -and -not -name *%s -and -not -name *%s", DOT_PARTIAL, DOT_LAYOUT)
	args := strings.Fields(s_args)
	return run.Cmd_Args("find", append([]string{dir}, args...)...)
}

func List_Template_Dirs(dir string, raw_exclude ...string) []string {
	raw_dirs := run.One_Line_Script("find " + dir + "  -maxdepth 2 -mindepth 2 -type f -name '*.go.html' | xargs dirname | sort -u")

	excludes := files.Clean_Paths(raw_exclude...)

	var dirs []string

	for _, d := range raw_dirs {
		if slices.Contains(excludes, filepath.Clean(d)) {
			continue
		}
		dirs = append(dirs, d)
	}

	return dirs
}

func List_All_Dot_Go(dir string) []string {
	return run.One_Line_Script(`find ` +  dir + ` -maxdepth 1 -mindepth 1 -type f -name '*.go.html'`)
}

func List_Related_Files(dir string, layouts_partials ...string) []string {
	var related []string
	for _, sub_dir := range layouts_partials {
		related = append(related, List_All_Dot_Go(sub_dir)...)
	}

	dir_files := run.Cmd_Args("find", dir, "-type", "f", "-name", "*.partial.go.html", "-or", "-name", "*.layout.go.html")
	return append(related, dir_files...)
}

func Compile_Template(config_json map[string]interface{}, related_files []string, f string) error {
	fmt.Printf("Template: %v\n", f)
	fmt.Printf("New File: %v\n", Remove_Dot_Go(f))

	t, err := go_template.ParseFiles(append([]string{f}, related_files...)...)
	if err != nil { exit.PrintError(fmt.Errorf("Parsing files: %v, %v", f, related_files), err) }

	new_file, c_err  := os.Create(Remove_Dot_Go(f))
	if c_err != nil { exit.PrintError(c_err) }
	defer new_file.Close()

	return t.Execute(new_file, config_json)
}

func Compile_All(config_file string, dir string, related_dirs ...string) {
	config_json, err := config.Get_Config(config_file)
	exit.Print_Msg(err, "Config could not be retrieved.")

	for i, d := range List_Template_Dirs(dir) {
		fmt.Printf("%v - %v\n", i, d)
		related_files := List_Related_Files(d, related_dirs...)
		for _, f := range List_Templates_In_Dir(d) {
			ct_err := Compile_Template(config_json, related_files, f)
			exit.PrintError(ct_err)
		}
	}
}


