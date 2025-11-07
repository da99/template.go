
package main

import (
	"fmt"
	"github.com/da99/cli.go/args"
	// "github.com/da99/cli.go/files"
	"github.com/da99/template.go/template"
)

func main() {

	if args.IsMatch("-h", 0) {
		fmt.Println("  -h|help|--help   -- This message.")
		fmt.Println("  compile [dir]    -- Compile all files in `dir`.")
		return
	}

	if args.IsMatch("compile", 1) {
		dir := args.CAPTURE[0]
		template.Compile_All(dir)
		return
	}

	args.Fail()
}
