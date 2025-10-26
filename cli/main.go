/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"fmt"
	"github.com/da99/cli.go/args"
	"github.com/da99/files.go/files"
	"os"
)

func PrintError(x error) {
	if x != nil {
		fmt.Printf("Error: %v\n", x)
		os.Exit(1)
	}
}

func main() {

	if args.IsMatch("-h", 0) {
		fmt.Println("-h -- Print this message.")
		fmt.Println("ls [dirs|files] [dir]")
		fmt.Println("compile [dirs|files] [dir]")
		os.Exit(0)
	}

	if args.IsMatch("ls dirs", 1) {
		dirs, err := files.List_Shallow_Dirs(args.CAPTURE[0])
		PrintError(err)
		for _, v := range dirs { fmt.Println(v) }
		os.Exit(0)
	}

	if args.IsMatch("ls files", 1) {
		files, err := files.List_Shallow_Files_Ext(args.CAPTURE[0], "*.go.html")
		PrintError(err)
		for _, v := range files { fmt.Println(v) }
		os.Exit(0)
	}

	args.Fail()
}
