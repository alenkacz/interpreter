package repl

import (
	"bufio"
	"fmt"
	"github.com/alenkacz/interpreter-book/pkg/token"
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
		fmt.Print(line)
		t := tokenizer.New(line)
		for tok := t.NextToken(); tok.Type != token.EOF; tok = t.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
