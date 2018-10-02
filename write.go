package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ryanjoneil/tsppd-dd/ddo"
	"github.com/ryanjoneil/tsppd-dd/tsppd"
)

type output struct {
	flags   *flags
	problem *tsppd.Problem
	writer  *csv.Writer
}

func createOutput(f *flags, problem *tsppd.Problem) *output {
	var writer *csv.Writer

	if f.verbosity() == 1 {
		fmt.Print("instance        size   form        infer  relax  ")
		fmt.Println("ordering  width     batch  workers  clock    cpu      primal    optimal  nodes     fails")
		for i := 0; i < 150; i++ {
			fmt.Print("=")
		}
		fmt.Println()

	} else if strings.Contains(f.output(), "csv") {
		writer = csv.NewWriter(os.Stdout)
		if strings.Contains(f.output(), "header") {
			writer.Write([]string{
				"instance",
				"size",
				"form",
				"infer",
				"relax",
				"ordering",
				"width",
				"batch",
				"workers",
				"maxmillis",
				"maxnodes",
				"clock",
				"cpu",
				"primal",
				"optimal",
				"nodes",
				"fails",
				"path",
			})
		}

	}

	return &output{
		flags:   f,
		problem: problem,
		writer:  writer,
	}
}

func (o *output) write(b *ddo.Bounds, stats ddo.Statistics) {
	solution := b.Primal.(tsppd.State).Solution()

	if o.flags.verbosity() == 1 {
		fmt.Printf(
			"%-16s%-7d%-12s%-7s%-7s%-10s%-10d%-7d%-9d%-9.3f%-9.3f%-10d%-9t%-10d%-10d\n",
			solution.Problem.Name,
			len(solution.Problem.Nodes),
			o.flags.form(),
			o.flags.infer(),
			o.flags.relax(),
			o.flags.ordering(),
			o.flags.width(),
			o.flags.batch(),
			o.flags.workers(),
			stats.ClockSeconds,
			stats.CPUSeconds,
			b.PrimalBound(),
			stats.Optimal,
			stats.Nodes,
			stats.Fails,
		)

	} else if o.writer != nil {
		o.writer.Write([]string{
			o.problem.Name,
			strconv.Itoa(len(solution.Problem.Nodes)),
			o.flags.form(),
			o.flags.infer(),
			o.flags.relax(),
			o.flags.ordering(),
			strconv.Itoa(int(o.flags.width())),
			strconv.Itoa(int(o.flags.batch())),
			strconv.Itoa(int(o.flags.workers())),
			strconv.FormatUint(o.flags.maxmillis(), 10),
			strconv.FormatUint(o.flags.maxnodes(), 10),
			fmt.Sprintf("%.10f", stats.ClockSeconds),
			fmt.Sprintf("%.10f", stats.CPUSeconds),
			strconv.FormatInt(b.PrimalBound(), 10),
			strconv.FormatBool(stats.Optimal),
			strconv.FormatUint(stats.Nodes, 10),
			strconv.FormatUint(stats.Fails, 10),
			strings.Join(solution.Path, " "),
		})
		o.writer.Flush()
	}
}
