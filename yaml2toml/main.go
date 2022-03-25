package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

var (
	out = flag.StringP("output", "o", "", "output file")
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: yaml2toml [options] [file]")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	log.SetPrefix("yaml2toml: ")
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() > 1 {
		usage()
	}

	of := os.Stdout
	if *out != "" {
		f, err := os.Create(*out)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		of = f
	}

	w := bufio.NewWriter(of)
	defer w.Flush()

	if flag.NArg() == 0 {
		if err := convert(w, os.Stdin); err != nil {
			log.Fatal(err)
		}
	} else {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		if err := convert(w, f); err != nil {
			log.Fatal(err)
		}
	}
}

func convert(w io.Writer, r io.Reader) error {
	var v any
	if err := yaml.NewDecoder(r).Decode(&v); err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(v); err != nil {
		return err
	}

	// Trim final newliens.
	b := buf.Bytes()
	for i := len(b) - 1; i > 0; i-- {
		if b[i-1] != '\n' && b[i] == '\n' {
			w.Write(b[:i+1])
			break
		}
	}

	return nil
}
