package template

import (
	"fmt"
	"path"
	"path/filepath"
	"html/template"
	"os"
	"os/exec"
	"strings"
	// "errors"
	"encoding/json"
	// "sync"
	"github.com/da99/files.go/files"
	"github.com/da99/cli.go/exit"
)

const PARTIAL_PATTERN = ".partial.go.html"
type FileHandler func(string) error

func must_exist(str_path string) bool {
	if !files.Is(str_path) {
		fmt.Fprintf(os.Stderr, "Does not exist: %v\n", str_path)
		os.Exit(1)
	}
	return true
}


func shell_lines(cmd_str string) []string {
	raw := exec.Command("sh", "-c", cmd_str)

	output, o_err := raw.Output()
	exit.PrintError(o_err)

	return strings.Split(strings.TrimSpace(string(output)), "\n")
}

func List_Template_Files(str_dir string) []string {
	must_exist(str_dir)
	return shell_lines("find " + str_dir + " -type f -name '*.go.html' -and -not -name '*.partial.go.html' | sort")
}

func List_Dirs(str_dir string) []string {
	must_exist(str_dir)
	return shell_lines("find " + str_dir + " -type f -name '*.go.html' -and -not -name '*.partial.go.html' | xargs dirname | sort | uniq")
}

func GetConfigBytes(raw_files ...string) ([]byte, error) {
	file_path := files.First(raw_files...)
	if file_path == "" {
		return nil, nil
	}
	contents, read_err := os.ReadFile(file_path)
	if read_err != nil { return nil, read_err }
	return contents, nil
}

func FileHTML(raw_path string) string {
	return strings.Replace(raw_path, ".go.html", ".html", 1)
}

func GetConfig() (map[string]interface{}, error) {
	fin := make(map[string]interface{})

	contents, config_err := GetConfigBytes("config.json", "config/main.json")
	if config_err != nil {
		return fin, config_err
	}

	if contents == nil { return fin, nil }

	j_err := json.Unmarshal(contents, &fin)
	if j_err != nil {
		return fin, j_err
	}

	return fin, nil
}

func CompileFile(fp string) error {
	fmt.Println("Compiling: " + fp)
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		fmt.Println("errorr")
		return err
	}
	return tmpl.Execute(os.Stdout, "http://www.lewrockwell.com/>a?")
}


func IsPartial(fp string) bool {
	return strings.Contains(fp, PARTIAL_PATTERN)
}

func CompileDir(str_dir string) error {
	// var wg sync.WaitGroup
	// defer wg.Wait()

	config, c_err := GetConfig()
	if c_err != nil { return c_err }

	all_dirs := List_Dirs(str_dir)

	for _, d := range all_dirs {
		all_files, err := files.List_Shallow_Files_Ext(d, "*.go.html")
		if err != nil { return err }

		tmpl, t_err := template.ParseFiles(all_files...)
		if t_err != nil { return err }

		for _, f := range all_files {
			if IsPartial(f) { continue; }

			fmt.Printf("-- Compiling template: %v\n", f)
			new_file_path := FileHTML(f)
			filer, ferr := os.Create(new_file_path)
			if ferr != nil {
				filer.Close()
				return ferr
			}

			err := tmpl.ExecuteTemplate(filer, path.Base(f), config)
			filer.Close()

			if err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(1)
			}

			fmt.Println("Wrote: " + new_file_path)
		}
	}

	return nil
}

func print_it(str string) error {
	_, err := fmt.Printf("File: %v\n", str)
	return err
}

func each_line(matches []string, fh FileHandler) error {
	for _, v := range matches {
		err := fh(v)
		if err != nil {
			return err
		}
	}
	return nil
}



func LsFiles(target string) ([]string, error) {
	return filepath.Glob(filepath.Join(target, "/**/*.go.html"))
}
