package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
)

const command_name = "tapl"

type test struct {
	script   string
	options  []string
	ok       string
	ok_regex string
	rc       int
}

func TestTapl(t *testing.T) {
	tests := []test{
		//BASIC ELEMENTS
		{script: "", rc: 1, ok: ""},
		{script: "true", ok: "true\n"},
		{script: "false", ok: "false\n"},
		{script: "0", ok: "0\n"},

		// IF STATEMENT
		{script: "if true then 0 else false", ok: "0\n"},
		{script: "if false then 0 else false", ok: "false\n"},
		{script: "if true then true else 0", ok: "true\n"},
		{script: "if false then true else 0", ok: "0\n"},

		// NUMBER FUNCS
		{script: "succ 0", ok: "1\n"},
		{script: "pred 0", ok: "-1\n"},
		{script: "succ succ succ 0", ok: "3\n"},
		{script: "succ pred succ 0", ok: "1\n"},

		// PREDICATE FUNCS
		{script: "iszero 0", ok: "true\n"},
		{script: "iszero succ 0", ok: "false\n"},
		{script: "iszero pred 0", ok: "false\n"},
		{script: "iszero succ pred succ pred 0", ok: "true\n"},
		{script: "if iszero 0 then succ 0 else pred 0", ok: "1\n"},
		{script: "if iszero succ 0 then succ 0 else pred 0", ok: "-1\n"},

		// Consts function
		{script: "0", options: []string{"-c"}, ok: "0\nConsts=>[0]\n"},
		{script: "true", options: []string{"-c"}, ok: "true\nConsts=>[true]\n"},
		{script: "false", options: []string{"-c"}, ok: "false\nConsts=>[false]\n"},
		{script: "pred 0", options: []string{"-c"}, ok: "-1\nConsts=>[0]\n"},
		{script: "succ 0", options: []string{"-c"}, ok: "1\nConsts=>[0]\n"},
		{script: "iszero 0", options: []string{"-c"}, ok: "true\nConsts=>[0]\n"},
		{script: "if iszero 0 then true else false", options: []string{"-c"}, ok: "true\nConsts=>[0 true false]\n"},
		{script: "if iszero 0 then true else if iszero 0 then true else false", options: []string{"-c"}, ok: "true\nConsts=>[0 true false]\n"},

		// Size function
		{script: "0", options: []string{"-s"}, ok: "0\nSize=>1\n"},
		{script: "true", options: []string{"-s"}, ok: "true\nSize=>1\n"},
		{script: "false", options: []string{"-s"}, ok: "false\nSize=>1\n"},
		{script: "pred 0", options: []string{"-s"}, ok: "-1\nSize=>2\n"},
		{script: "succ 0", options: []string{"-s"}, ok: "1\nSize=>2\n"},
		{script: "iszero 0", options: []string{"-s"}, ok: "true\nSize=>2\n"},
		{script: "if iszero 0 then succ 0 else pred 0", options: []string{"-s"}, ok: "1\nSize=>6\n"},

		// Depth function
		{script: "0", options: []string{"-d"}, ok: "0\nDepth=>1\n"},
		{script: "true", options: []string{"-d"}, ok: "true\nDepth=>1\n"},
		{script: "false", options: []string{"-d"}, ok: "false\nDepth=>1\n"},
		{script: "pred 0", options: []string{"-d"}, ok: "-1\nDepth=>2\n"},
		{script: "succ 0", options: []string{"-d"}, ok: "1\nDepth=>2\n"},
		{script: "iszero 0", options: []string{"-d"}, ok: "true\nDepth=>2\n"},
		{script: "if iszero 0 then succ 0 else pred 0", options: []string{"-d"}, ok: "1\nDepth=>2\n"},

		//OPTIONS
		/*
			{script: "1", options: []string{"-a"}, ok_regex: `ast.NumExpr{Literal:"1"}`},
			{script: "1", options: []string{"-i"}, ok_regex: `define i32 @main\(\) {`},
			{script: "for i=0;i<1;i++{print i}", options: []string{"-d"}, ok_regex: `debug option`},
			{script: "1", options: []string{"-n"}, ok: ""},
		*/
	}

	//realStdin := os.Stdin
	realStdout := os.Stdout
	realStderr := os.Stderr
	case_number := 0

	for _, test := range tests {
		case_number++
		wg := &sync.WaitGroup{}

		//fmt.Printf("TEST[%v] %v\n", case_number, test.script)

		// OUT PIPE
		readFromOut, writeToOut, err := os.Pipe()
		if err != nil {
			//os.Stdin = realStdin
			os.Stderr = realStderr
			t.Fatal("Pipe error:", err)
		}
		os.Stdout = writeToOut
		//logger.Print("pipe out created")

		// Read Stdout goroutine
		readerOut := bufio.NewScanner(readFromOut)
		chanOut := make(chan string)
		wg.Add(1)
		go func() {
			for readerOut.Scan() {
				chanOut <- readerOut.Text()
			}
			close(chanOut)
			wg.Done()
			return
		}()

		// Run Script goroutine
		wg.Add(1)
		go func() {

			os.Args = []string{command_name}
			os.Args = append(os.Args, test.options...)
			if test.script != "" {
				os.Args = append(os.Args, test.script)
			}
			rc := _main()
			//fmt.Fprintf(realStdout, "case:%d os.Args=%v *print_ast=%v\n", case_number, os.Args, *print_ast)
			if rc != test.rc && !strings.Contains(test.ok, "error") {
				t.Errorf("return code want:%v get:%v case:%v\n", test.rc, rc, test.script)
			}

			//close(chanDone) //NG
			writeToOut.Close()
			wg.Done()
		}()

		// Get Result
		var resultOut string
	LOOP:
		for {
			select {
			case dataOut, ok := <-chanOut:
				if !ok {
					break LOOP
				}
				dataOut = strings.TrimSpace(dataOut)
				resultOut = fmt.Sprintf("%s%s%s", resultOut, dataOut, "\n")
			}
		}

		// Result Check
		//fmt.Fprintf(realStdout, "result:[%v]\ttest.ok:[%v]\n", resultOut, test.ok)
		if test.ok != "" && resultOut != strings.Replace(test.ok, "\r", "", -1) { //replace for Windows
			t.Errorf("Case:[%v] received: %v - expected: %v - runSource: %v", case_number, resultOut, test.ok, test.script)
		}
		if test.ok_regex != "" {
			r := regexp.MustCompile(test.ok_regex)
			if !r.MatchString(resultOut) {
				t.Errorf("Case:[%v] received: %v - expected(regexp): %v - runSource: %v", case_number, resultOut, test.ok_regex, test.script)
			}
		}

		wg.Wait()
		//readFromIn.Close()
		//writeToIn.Close()
		readFromOut.Close()
		writeToOut.Close()
		//os.Stdin = realStdin
		os.Stderr = realStderr
		os.Stdout = realStdout
	}

}
