package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/tesujiro/taplGo/ch03/debug"
	"github.com/tesujiro/taplGo/ch03/parser"
	"github.com/tesujiro/taplGo/ch03/vm"
)

var (
	print_ast, no_exec, dbg bool
	constsF, sizeF, depthF  bool
)

func main() {
	os.Exit(_main())
}

func _main() int {

	//parser.TraceLexer()

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover:", err)
			for depth := 0; ; depth++ {
				_, file, line, ok := runtime.Caller(depth)
				if !ok {
					break
				}
				log.Printf("=>%d: %v:%d", depth, file, line)
			}
		}
	}()

	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	f.BoolVar(&print_ast, "a", false, "print AST")
	f.BoolVar(&no_exec, "n", false, "no execution")
	//f.BoolVar(&dbg, "d", false, "debug option")
	f.BoolVar(&constsF, "c", false, "Consts function")
	f.BoolVar(&sizeF, "s", false, "Size function")
	f.BoolVar(&depthF, "d", false, "Deph function")

	f.Parse(os.Args[1:])
	args := f.Args()

	if len(args) < 1 {
		fmt.Println("No expression error!")
		//fmt.Printf("ex: %v '(1+1)*3+10' ; echo $?\n", os.Args[0])
		return 1
	}

	if dbg {
		debug.On()
		fmt.Println("debug option")
	} else {
		debug.Off()
	}

	for _, source := range args {
		//fmt.Printf("source: %v\n", source)
		//result := runScript(strings.NewReader(source))
		result := runScript(source)
		if result != 0 {
			return result
		}
	}
	return 0
}

func runScript(source string) int {

	//env := vm.NewEnv()
	ast, parseError := parser.ParseSrc(source)
	if parseError != nil {
		fmt.Printf("%v\n", parseError)
		return 1
	}
	if print_ast {
		parser.Dump(ast)
	}

	val, err := vm.Run(ast)
	if err != nil {
		fmt.Printf("Runtime error: %v \n", err)
		return 1
	}
	fmt.Printf("%v\n", val)

	if constsF {
		val, err := vm.Consts(ast)
		if err != nil {
			fmt.Printf("Function error: %v \n", err)
		}
		fmt.Printf("%v\n", val)
	}

	if sizeF {
		val, err := vm.Size(ast)
		if err != nil {
			fmt.Printf("Function error: %v \n", err)
		}
		fmt.Printf("%v\n", val)
	}

	if depthF {
		val, err := vm.Depth(ast)
		if err != nil {
			fmt.Printf("Function error: %v \n", err)
		}
		fmt.Printf("%v\n", val)
	}

	if no_exec {
		return 0
	}

	return 0
}
