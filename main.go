package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/huderlem/poryscript/emitter"
	"github.com/huderlem/poryscript/lexer"
	"github.com/huderlem/poryscript/parser"
	"github.com/huderlem/poryscript/types"
)

const version = "2.8.0"

type mapOption map[string]string

func (opt mapOption) String() string {
	return ""
}

func (opt mapOption) Set(value string) error {
	result := strings.SplitN(value, "=", 2)
	if len(result) != 2 {
		return fmt.Errorf("expected key-value option to be separate by '=', but got '%s' instead", value)
	}
	opt[result[0]] = result[1]
	return nil
}

type options struct {
	gen                types.Gen
	inputFilepath      string
	outputFilepath     string
	fontWidthsFilepath string
	optimize           bool
	compileSwitches    map[string]string
}

func parseOptions() options {
	helpPtr := flag.Bool("h", false, "show poryscript help information")
	versionPtr := flag.Bool("v", false, "show version of poryscript")
	genPtr := flag.Int("gen", 3, "decompilation project generation (pokecrystal = 2, pokeruby/emerald/firered = 3). Defaults to 3")
	inputPtr := flag.String("i", "", "input poryscript file (leave empty to read from standard input)")
	outputPtr := flag.String("o", "", "output script file (leave empty to write to standard output)")
	fontsPtr := flag.String("fw", "font_widths.json", "font widths config JSON file")
	optimizePtr := flag.Bool("optimize", true, "optimize compiled script size (To disable, use '-optimize=false')")
	compileSwitches := make(mapOption)
	flag.Var(compileSwitches, "s", "set a compile-time switch. Multiple -s options can be set. Example: -s VERSION=RUBY -s LANGUAGE=GERMAN")
	flag.Parse()

	if *helpPtr == true {
		flag.Usage()
		os.Exit(0)
	}

	if *versionPtr == true {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}

	gen := types.GEN3
	if *genPtr == 2 {
		gen = types.GEN2
	}

	return options{
		gen:                gen,
		inputFilepath:      *inputPtr,
		outputFilepath:     *outputPtr,
		fontWidthsFilepath: *fontsPtr,
		optimize:           *optimizePtr,
		compileSwitches:    compileSwitches,
	}
}

func getInput(filepath string) (string, error) {
	var bytes []byte
	var err error
	if filepath == "" {
		bytes, err = ioutil.ReadAll(os.Stdin)
	} else {
		bytes, err = ioutil.ReadFile(filepath)
	}

	return string(bytes), err
}

func writeOutput(output string, filepath string) error {
	if filepath == "" {
		fmt.Print(output)
	} else {
		f, err := os.Create(filepath)
		if err != nil {
			return err
		}

		_, err = io.WriteString(f, output)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	log.SetFlags(0)
	options := parseOptions()
	input, err := getInput(options.inputFilepath)
	if err != nil {
		log.Fatalf("PORYSCRIPT ERROR: %s\n", err.Error())
	}

	parser := parser.New(lexer.New(input, options.gen), options.gen, options.fontWidthsFilepath, options.compileSwitches)
	program, err := parser.ParseProgram()
	if err != nil {
		log.Fatalf("PORYSCRIPT ERROR: %s\n", err.Error())
	}

	emitter := emitter.New(program, options.gen, options.optimize)
	result, err := emitter.Emit()
	if err != nil {
		log.Fatalf("PORYSCRIPT ERROR: %s\n", err.Error())
	}
	err = writeOutput(result, options.outputFilepath)
	if err != nil {
		log.Fatalf("PORYSCRIPT ERROR: %s\n", err.Error())
	}
}
