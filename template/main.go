
package template

import (
	"fmt"
	"path"
	"path/filepath"
	"html/template"
	"os"
	"strings"
	"errors"
	"encoding/json"
	"sync"
	"github.com/da99/files.go/files"
)

const PARTIAL_PATTERN = ".partial.go.html"
type FileHandler func(string) error

func GetConfigBytes(files ...string) ([]byte, error) {
	file_path := files.First(files)
	if file_path == "" {
		return nil, errors.New("Config file not found.")
	}
	contents, read_err := os.ReadFile(file_path)
	if read_err != nil { return nil, read_err }
	return contents, nil
}

func GetConfig() (map[string]interface{}, error) {
	fin := make(map[string]interface{})

	contents, config_err := GetConfigBytes("config.json", "config/main.json")
	if config_err != nil {
		return fin, config_err
	}

	j_err := json.Unmarshal([]byte(contents), &fin)
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


// compileCmd represents the compile command
func compile(args []string) error {
	var wg sync.WaitGroup
	defer wg.Wait()

	config, c_err := GetConfig()
	if c_err != nil { return c_err }

	dirs, d_err := files.List_Shallow_Dirs(args[0])
	if d_err != nil { return d_err }

	for _, d := range dirs {
		files, err := files.List_Shallow_Files_Ext(d, "*.go.html")
		if err != nil { return err }

		tmpl, t_err := template.ParseFiles(files...)
		if t_err != nil { return err }

		for _, f := range files {
			if strings.Contains(f, PARTIAL_PATTERN) { continue; }

			fmt.Printf("-- Compiling template: %v\n", f)
			err := tmpl.ExecuteTemplate(os.Stdout, path.Base(f), config)
			if err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(1)
			}
			fmt.Println("\n")

		}
	}

	//
	//
	// for _, v := range files {
	// }
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


func CompileDir(target string, fh FileHandler) error {
	stuff, err := os.ReadDir(target)
	if err != nil { return err }
	for _, entry := range stuff {
		if entry.IsDir() && files.Is(path.Join(target, "index.go.html")) {
			fh(entry.Name())
		}
	}

	return nil
}


func LsFiles(target string) ([]string, error) {
	return filepath.Glob(filepath.Join(target, "/**/*.go.html"))
}
