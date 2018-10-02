package main

import (
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/ryanjoneil/tsppd-dd/ddo"
	"github.com/ryanjoneil/tsppd-dd/tsppd/solvers/sequential"
	"github.com/ryanjoneil/tsppd-dd/tsppd/solvers/successor"
)

func main() {
	flags := parseFlags()
	flags.validate()

	if flags.cpuprof() != "" {
		f, err := os.Create(flags.cpuprof())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	problem := readProblem(flags.input())

	var root ddo.State
	if flags.form() == "sequential" {
		root = sequential.CreateRootState(
			problem,
			flags.infer(),
			flags.relax(),
			flags.ordering(),
			flags.width(),
			flags.verbosity(),
		)
	} else if flags.form() == "successor" {
		root = successor.CreateRootState(
			problem,
			flags.infer(),
			flags.relax(),
			flags.ordering(),
			flags.width(),
			flags.verbosity(),
		)
	}

	output := createOutput(flags, problem)
	solver := ddo.CreateSolver(root, output.write)
	solver.Batch = flags.batch()
	solver.Workers = flags.workers()
	solver.MaxMillis = flags.maxmillis()
	solver.MaxNodes = flags.maxnodes()

	solver.Minimize()

	if flags.memprof() != "" {
		f, err := os.Create(flags.memprof())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		pprof.WriteHeapProfile(f)
		defer pprof.StopCPUProfile()
	}
}
