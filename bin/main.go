
package main

import (
	"fmt"
	"os"
	"github.com/da99/cli.go/args"
	"github.com/da99/cli.go/files"
)

func main() {

	if args.IsMatch("-h", 0) {
		fmt.Println("  -h|help|--help   -- This message.")
		fmt.Println("  compile (dir)    -- Compile all files or just the `dir`.")
		return
	}

	if args.IsMatch("compile", 0) {
		dirs, err := files.List_Shallow_Dirs("public/section")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		for _, d := range dirs {
			fmt.Println(d)
		}
		return
	}

	args.Fail()
}
