package repl

import (
	"bufio"
	"fmt"
	"github.com/alenkacz/interpreter-book/pkg/eval"
	"github.com/alenkacz/interpreter-book/pkg/parser"
	"github.com/alenkacz/interpreter-book/pkg/tokenizer"
	"io"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(">> ")
		scanned := scanner.Scan()
		if !scanned {
			if scanner.Err() != nil {
				fmt.Fprintf(out, "Error when reading input: %v", scanner.Err())
				return
			}
		}

		line := scanner.Text()
		t := tokenizer.New(line)
		p := parser.New(t)
		ast := p.ParseProgram()

		if len(p.Errors) != 0 {
			printParserErrors(out, p.Errors)
			continue
		}

		fmt.Fprintf(out, "%s\n", eval.Eval(ast).Print())
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Woops! An error while parsing the program!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
