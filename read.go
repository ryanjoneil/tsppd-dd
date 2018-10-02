package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ryanjoneil/tsppd-dd/tsppd"
)

func readProblem(input string) *tsppd.Problem {
	var b []byte
	var err error
	if input == "-" {
		b, err = ioutil.ReadAll(os.Stdin)
	} else {
		b, err = ioutil.ReadFile(input)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	problem, err := tsppd.Decode(b)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return &problem
}
