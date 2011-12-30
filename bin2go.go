package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var (
	in  = flag.String("in", "", "use this file instead of the stdin for input")
	out = flag.String("out", "", "use this file instead of the stdout for output")
	pkg = flag.String("pkg", "", "prepend package clause specifying this package")
)

// Exit error codes
const (
	NO_ERROR = iota
	WRONG_ARGS
	INPUT_FAIL
	OUTPUT_FAIL
)

func printUsage() {
	fmt.Printf("usage: %s [-in=<path>] [-out=<path>] [-pkg=<name>] <varname>\n", os.Args[0])
	flag.PrintDefaults()
}

func printUsageAndExit() {
	flag.Usage()
	os.Exit(WRONG_ARGS)
}

func readInput() []byte {
	var data []byte
	var err error

	if *in != "" {
		data, err = ioutil.ReadFile(*in)
	} else {
		data, err = ioutil.ReadAll(os.Stdin)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read input: %s\n", err)
		os.Exit(INPUT_FAIL)
	}
	return data
}

func checkOutputFailure(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write output: %s\n", err)
		os.Exit(OUTPUT_FAIL)
	}
}

func writeData(data []byte, out io.Writer) {
	varname := flag.Arg(0)

	// write header
	_, err := fmt.Fprintf(out, "var %s = []byte{\n\t", varname)
	checkOutputFailure(err)

	lastbytei := len(data) - 1
	n := 8
	for i, b := range data {
		// write single byte
		_, err = fmt.Fprintf(out, "0x%.2x,", b)
		checkOutputFailure(err)

		n += 6

		// if this is not the last byte
		if i != lastbytei {
			// be readable, break line after 78 characters
			if n >= 78 {
				_, err = fmt.Fprint(out, "\n\t")
				checkOutputFailure(err)

				n = 8
			} else {
				// if we're not breaking the line, insert space
				// after ','
				_, err = fmt.Fprint(out, " ")
				checkOutputFailure(err)
			}

		}
	}
	_, err = fmt.Fprint(out, "\n}\n")
	checkOutputFailure(err)
}

func writeOutput(data []byte) {
	var output *bufio.Writer

	// prepare "output"
	if *out != "" {
		file, err := os.Create(*out)
		checkOutputFailure(err)
		defer file.Close()

		output = bufio.NewWriter(file)
	} else {
		output = bufio.NewWriter(os.Stdout)
	}

	// write package clause if any
	if *pkg != "" {
		_, err := fmt.Fprintf(output, "package %s\n\n", *pkg)
		checkOutputFailure(err)
	}

	// write data
	writeData(data, output)

	// flush
	err := output.Flush()
	checkOutputFailure(err)
}

func main() {
	flag.Usage = printUsage
	flag.Parse()

	if flag.NArg() != 1 {
		printUsageAndExit()
	}

	data := readInput()
	writeOutput(data)
}
