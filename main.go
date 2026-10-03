package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

type executer interface {
	execute() (string, error)
}

func main() {
	proj := flag.String("p", "", "project directory")
	flag.Parse()

	if err := run(*proj, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(proj string, out io.Writer) error {
	if proj == "" {
		return fmt.Errorf("Project directory is required: %w", ErrValidation)
	}

	pipline := make([]executer, 4)

	pipline[0] = newStep(
		"go build",
		"go",
		"Go Build: SUCCESS",
		proj,
		[]string{"build", ".", "errors"},
	)

	pipline[1] = newStep(
		"go test",
		"go",
		"Go Test: SUCCESS",
		proj,
		[]string{"test", "-v"},
	)

	pipline[2] = newExceptionStep(
		"go fmt",
		"gofmt",
		"Gofmt: SUCCESS",
		proj,
		[]string{"-l", "."},
	)

	pipline[3] = newTimeoutStep(
		"git push",
		"git",
		"Git Push: SUCCESS",
		proj,
		[]string{"push", "origin", "master"},
		10*time.Second,
	)

	for _, s := range pipline {
		msg, err := s.execute()
		if err != nil {
			return err
		}

		_, err = fmt.Fprintln(out, msg)
		if err != nil {
			return err
		}
	}

	return nil
}
