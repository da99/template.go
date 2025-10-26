/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cli

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

func GetConfigFile() (string, error) {
	if _, err := os.Stat("config.json"); !os.IsNotExist(err) {
		return "config.json", nil
	}
	if _, err := os.Stat("config/main.json"); !os.IsNotExist(err) {
		return "config/main.json", nil
	}
	return "", errors.New("Config file not found.")
}

func GetConfig() (map[string]interface{}, error) {
	fin := make(map[string]interface{})

	config_file, config_err := GetConfigFile()
	if config_err != nil {
		return fin, nil
	}
	contents, read_err := os.ReadFile(config_file)

	if read_err != nil {
		return fin, read_err;
	}

	j_err := json.Unmarshal([]byte(contents), &fin)
	if j_err != nil {
		return fin, nil
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

func main(args []string) error {
	switch args[0] {
	case "dirs":
		dirs, err := files.List_Shallow_Dirs(args[1])
		if err != nil { return err; }
		for _, v := range dirs { fmt.Println(v) }
	case "files":
		files, err := files.List_Shallow_Files_Ext(args[1], "*.go.html")
		if err != nil { return err; }
		for _, v := range files { fmt.Println(v) }
	default:
		return errors.New("Invalid option: ls '" + args[0] + "' '" + args[1] + "'")
	}
	return nil
}
