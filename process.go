package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// == Direct input processing ========================================

func processInputFiles(stringTemplate string, inputFiles []string, outputs []string, g *Gomplate) error {
	inputs, err := readInputs(stringTemplate, inputFiles)
	if err != nil {
		return err
	}

	if len(outputs) == 0 {
		outputs = []string{"-"}
	}

	n := 0
	for _, input := range inputs {
		output := ""
		if !input.partial {
			output = outputs[n]
			n++
		}
		if err := renderTemplate(g, input, output); err != nil {
			return err
		}
	}
	return nil
}

// == Recursive input dir processing ======================================

func processInputDir(input string, output string, g *Gomplate) error {
	input = filepath.Clean(input)
	output = filepath.Clean(output)

	// assert tha input path exists
	si, err := os.Stat(input)
	if err != nil {
		return err
	}

	// read directory
	entries, err := ioutil.ReadDir(input)
	if err != nil {
		return err
	}

	// ensure output directory
	if err = os.MkdirAll(output, si.Mode()); err != nil {
		return err
	}

	// process or dive in again
	for _, entry := range entries {
		nextInPath := filepath.Join(input, entry.Name())
		nextOutPath := filepath.Join(output, entry.Name())

		if entry.IsDir() {
			err := processInputDir(nextInPath, nextOutPath, g)
			if err != nil {
				return err
			}
		} else {
			input, err := readInput(nextInPath)
			if err != nil {
				return err
			}
			if input.partial {
				nextOutPath = "" // Don't create files for partials
			}
			if err := renderTemplate(g, input, nextOutPath); err != nil {
				return err
			}
		}
	}
	return nil
}

// == File handling ================================================

func isPartialFilename(filename string) bool {
	return strings.HasPrefix(filepath.Base(filename), "_")
}

type Input struct {
	filename string
	text     string
	partial  bool
}

func textInput(text string) Input {
	return Input{text: text}
}

func fileInput(filename, text string) (input Input) {
	return Input{filename: filename, text: text, partial: isPartialFilename(filename)}
}

func readInputs(inputStr string, files []string) ([]Input, error) {
	if inputStr != "" {
		return []Input{textInput(inputStr)}, nil
	}
	if len(files) == 0 {
		files = []string{"-"}
	}
	ins := make([]Input, len(files))

	for n, filename := range files {
		input, err := readInput(filename)
		if err != nil {
			return nil, err
		}
		ins[n] = input
	}
	return ins, nil
}

func readInput(filename string) (Input, error) {
	var err error
	var inFile *os.File
	if filename == "-" {
		inFile = os.Stdin
	} else {
		inFile, err = os.Open(filename)
		if err != nil {
			return Input{}, fmt.Errorf("failed to open %s\n%v", filename, err)
		}
		// nolint: errcheck
		defer inFile.Close()
	}
	bytes, err := ioutil.ReadAll(inFile)
	if err != nil {
		err = fmt.Errorf("read failed for %s\n%v", filename, err)
		return Input{}, err
	}
	return fileInput(filename, string(bytes)), nil
}

type nullWriteCloser struct{}

func (nullWriteCloser) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (nullWriteCloser) Close() error {
	return nil
}

func openOutFile(filename string) (out io.WriteCloser, err error) {
	if filename == "" {
		return nullWriteCloser{}, nil
	}
	if filename == "-" {
		return os.Stdout, nil
	}
	return os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
}
