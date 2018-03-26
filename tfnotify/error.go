package main

import (
	"fmt"
	"os"
)

type ErrorFormatter interface {
	Format(s fmt.State, verb rune)
}

type ExitCoder interface {
	error
	ExitCode() int
}

type ExitError struct {
	exitCode int
	err      error
}

func NewExitError(exitCode int, err error) *ExitError {
	return &ExitError{
		exitCode: exitCode,
		err:      err,
	}
}

func (ee *ExitError) Error() string {
	if ee.err == nil {
		return ""
	}
	return fmt.Sprintf("%v", ee.err)
}

func (ee *ExitError) ExitCode() int {
	return ee.exitCode
}

func RunExitCoder(err error) int {
	success := 0
	if err == nil {
		return success
	}

	if exitErr, ok := err.(ExitCoder); ok {
		if err.Error() != "" {
			if _, ok := exitErr.(ErrorFormatter); ok {
				fmt.Fprintf(os.Stderr, "%+v\n", err)
			} else {
				fmt.Fprintln(os.Stderr, err)
			}
		}
		return exitErr.ExitCode()
	}

	if _, ok := err.(error); ok {
		fmt.Fprintf(os.Stderr, "[ERROR] %v\n", err)
		return 1
	}

	return success
}
