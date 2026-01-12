package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"rainlang/analyzer"
	"rainlang/debug"
	"rainlang/ir"
	"rainlang/lexer"
	"rainlang/parser"
	"rainlang/runner"
)

func main() {
	filePtr := flag.String("file", "", "file that you want to compile")
	outPtr := flag.String("out", "", "output C code file")
	langPtr := flag.String("lang", "c", "target language [c, go, flood(compiled rainlang bytecode)]")
	floodPtr := flag.String("flood", "", "flood bytecode file")

	flag.Parse()

	if *floodPtr != "" {
		bytecode, err := os.ReadFile(*floodPtr)
		if err != nil {
			fmt.Printf("An error occured while reading the file %s\n", *floodPtr)
			panic(err)
		}
		runFlood(bytecode)
		return
	}

	if *filePtr != "" {
		readAndExecute(*filePtr, outPtr, *langPtr)
		return
	}

	repl()
}

func runFlood(bytecode []byte) {
	ast, err := runner.DecodeFloodBytecode(bytecode)
	if err != nil {
		debug.PrintAST(ast)
		fmt.Println(err)
		return
	}
	debug.PrintAST(ast)

	fmt.Println("===============")

	err = runner.Run(ast)
	if err != nil {
		panic(err)
	}
}

func readAndExecute(filePath string, outFile *string, lang string) {
	code, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("An error occured while reading the file %s\n", filePath)
		panic(err)
	}

	if *outFile == "" {
		fmt.Println("No output file specified")
		return
	}

	switch lang {
	case "c", "C":
		lang = "c"
	case "go", "Go", "GO", "golang", "GOLANG", "goLang", "Golang":
		lang = "go"
	case "flood", "Flood", "FLOOD":
		lang = "flood"
	default:
		fmt.Printf("Unknown language %s\n", lang)
		return
	}

	out := run(string(code), lang, false)
	err = os.WriteFile(*outFile, []byte(out), 0644)
	if err != nil {
		panic(err)
	}
}

func repl() {
	fmt.Printf(`
        .
       .@.
      .@@@.
     :@@@@@:          RAINLANG V0.1.0
    .@@@@@@@.
    :@@@@@@@:
     '@@@@@'
       '"'

`)
	for {
		fmt.Print(">>> ")

		var code string

		in := bufio.NewReader(os.Stdin)
		code, err := in.ReadString('\n')
		if err != nil {
			panic(err)
		}

		run(code, "flood", true)
	}
}

func run(code string, lang string, replMode bool) []byte {
	tokens, err := lexer.Tokenizer(code)
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}

	debug.PrintTokens(tokens)
	fmt.Println("===============")

	ast, err := parser.Parse(tokens)
	if err != nil {
		debug.PrintAST(ast)
		fmt.Println(err)
		return []byte{}
	}

	debug.PrintAST(ast)
	fmt.Println("===============")

	err = analyzer.Analyze(ast)
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}

	fmt.Println("===============")

	var result []byte
	switch lang {
	case "c":
		outCode, err := ir.GenerateCCode(ast)
		if err != nil {
			fmt.Println(err)
			return []byte{}
		}
		result = []byte(outCode)
	case "go":
		outCode, err := ir.GenerateGoCode(ast)
		if err != nil {
			fmt.Println(err)
			return []byte{}
		}
		result = []byte(outCode)
	case "flood":
		result, err = ir.GenerateFloodBytecode(ast)
		if err != nil {
			fmt.Println(err)
			return []byte{}
		}
	}

	fmt.Println(result)

	if replMode {
		fmt.Println("===============")
		ast, err := runner.DecodeFloodBytecode(result)
		if err != nil {
			debug.PrintAST(ast)
			fmt.Println(err)
			return []byte{}
		}
		debug.PrintAST(ast)

		fmt.Println("===============")

		err = runner.Run(ast)
		if err != nil {
			fmt.Println(err)
			return []byte{}
		}
	}

	return result
}
