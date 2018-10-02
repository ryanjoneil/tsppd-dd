package main

import (
	"flag"
	"fmt"
	"os"
)

type flags struct {
	_batch     *int
	_cpuprof   *string
	_form      *string
	_infer     *string
	_input     *string
	_maxmillis *uint64
	_maxnodes  *uint64
	_memprof   *string
	_ordering  *string
	_output    *string
	_relax     *string
	_seed      *int64
	_verbosity *uint
	_width     *uint
	_workers   *int
}

func parseFlags() *flags {
	flags := &flags{
		_batch:     flag.Int("batch", 1, "batch size for parallelization"),
		_cpuprof:   flag.String("cpuprof", "", "cpu profile output"),
		_form:      flag.String("form", "", "formulation {sequential, successor}"),
		_infer:     flag.String("infer", "none", "inference dual {ap, none}"),
		_input:     flag.String("input", "-", "input json file"),
		_maxmillis: flag.Uint64("maxmillis", 0, "max milliseconds for search"),
		_maxnodes:  flag.Uint64("maxnodes", 0, "max nodes and fails for search"),
		_memprof:   flag.String("memprof", "", "mem profile output"),
		_ordering:  flag.String("ordering", "", "successor={greedy, input, regret}"),
		_output:    flag.String("output", "", "{csv, csv-header}"),
		_relax:     flag.String("relax", "none", "relaxation dual sequential={dd, none}"),
		_verbosity: flag.Uint("verbosity", 0, "solver verbosity (0 = quiet, 1 = solutions, 2 = layer construction)"),
		_width:     flag.Uint("width", 0, "diagram width"),
		_workers:   flag.Int("workers", 1, "number of workers"),
	}
	flag.Parse()
	return flags
}

func (f *flags) validate() {
	if *f._batch < 1 {
		fmt.Fprintln(os.Stderr, fmt.Errorf("batch size must be >= 1"))
		os.Exit(1)
	}

	if f.form() != "sequential" && f.form() != "successor" {
		fmt.Fprintln(os.Stderr, fmt.Errorf("valid formulation required"))
		os.Exit(1)
	}

	if f.infer() != "none" && f.infer() != "ap" {
		fmt.Fprintln(os.Stderr, fmt.Errorf("invalid inference dual form"))
		os.Exit(1)
	}

	orderings := map[string]map[string]bool{
		"successor": map[string]bool{
			"input":  true,
			"greedy": true,
			"regret": true,
		},
	}

	if orderings[f.form()] != nil && !orderings[f.form()][f.ordering()] {
		fmt.Fprintln(os.Stderr, fmt.Errorf(f.form()+" form requires valid decision ordering"))
		os.Exit(1)
	}

	if f.relax() != "none" && (f.form() != "sequential" && f.relax() != "dd") {
		fmt.Fprintln(os.Stderr, fmt.Errorf("invalid relaxation dual form"))
		os.Exit(1)
	}

	if *f._workers < 1 {
		fmt.Fprintln(os.Stderr, fmt.Errorf("workers must be >= 1"))
		os.Exit(1)
	}
}

func (f *flags) batch() int {
	return *f._batch
}

func (f *flags) cpuprof() string {
	return *f._cpuprof
}

func (f *flags) form() string {
	return *f._form
}

func (f *flags) infer() string {
	return *f._infer
}

func (f *flags) input() string {
	return *f._input
}

func (f *flags) maxmillis() uint64 {
	return *f._maxmillis
}

func (f *flags) maxnodes() uint64 {
	return *f._maxnodes
}

func (f *flags) memprof() string {
	return *f._memprof
}

func (f *flags) ordering() string {
	return *f._ordering
}

func (f *flags) output() string {
	return *f._output
}

func (f *flags) relax() string {
	return *f._relax
}

func (f *flags) verbosity() uint {
	return *f._verbosity
}

func (f *flags) width() uint {
	return *f._width
}

func (f *flags) workers() int {
	return *f._workers
}
