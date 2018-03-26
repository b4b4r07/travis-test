package terraform

import (
	"fmt"
	"regexp"
	"strings"
)

// Parser is an interface for parsing terraform execution result
type Parser interface {
	Parse(body string) ParseResult
}

// ParseResult represents the result of parsed terraform execution
type ParseResult struct {
	Result   string
	ExitCode int
	Error    error
}

// DefaultParser is a parser for terraform commands
type DefaultParser struct {
}

// FmtParser is a parser for terraform fmt
type FmtParser struct {
	Pass *regexp.Regexp
	Fail *regexp.Regexp
}

// PlanParser is a parser for terraform plan
type PlanParser struct {
	Pass *regexp.Regexp
	Fail *regexp.Regexp
}

// ApplyParser is a parser for terraform apply
type ApplyParser struct {
	Pass *regexp.Regexp
	Fail *regexp.Regexp
}

// NewDefaultParser is DefaultParser initializer
func NewDefaultParser() *DefaultParser {
	return &DefaultParser{}
}

// NewFmtParser is FmtParser initialized with its Regexp
func NewFmtParser() *FmtParser {
	return &FmtParser{
		Fail: regexp.MustCompile(`(?m)^(diff a/)`),
	}
}

// NewPlanParser is PlanParser initialized with its Regexp
func NewPlanParser() *PlanParser {
	return &PlanParser{
		Pass: regexp.MustCompile(`(?m)^(Plan: \d|No changes.)`),
		Fail: regexp.MustCompile(`(?m)^(Error: )`),
	}
}

// NewApplyParser is ApplyParser initialized with its Regexp
func NewApplyParser() *ApplyParser {
	return &ApplyParser{
		Pass: regexp.MustCompile(`(?m)^(Apply complete!)`),
		Fail: regexp.MustCompile(`(?m)^(Error: Error applying plan:)`),
	}
}

// Parse returns ParseResult related with terraform commands
func (p *DefaultParser) Parse(body string) ParseResult {
	return ParseResult{
		Result:   body,
		ExitCode: ExitPass,
		Error:    nil,
	}
}

// Parse returns ParseResult related with terraform fmt
func (p *FmtParser) Parse(body string) ParseResult {
	result := ParseResult{}
	if p.Fail.MatchString(body) {
		result.Result = "There is diff in your .tf file (need to be formatted)"
		result.ExitCode = ExitFail
	}
	return result
}

// Parse returns ParseResult related with terraform plan
func (p *PlanParser) Parse(body string) ParseResult {
	var exitCode int
	switch {
	case p.Pass.MatchString(body):
		exitCode = ExitPass
	case p.Fail.MatchString(body):
		exitCode = ExitFail
	default:
		return ParseResult{
			Result:   "",
			ExitCode: ExitFail,
			Error:    fmt.Errorf("no matched"),
		}
	}
	var result string
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		if p.Pass.MatchString(line) || p.Fail.MatchString(line) {
			result = line
		}
	}
	return ParseResult{
		Result:   result,
		ExitCode: exitCode,
		Error:    nil,
	}
}

// Parse returns ParseResult related with terraform apply
func (p *ApplyParser) Parse(body string) ParseResult {
	var exitCode int
	switch {
	case p.Pass.MatchString(body):
		exitCode = ExitPass
	case p.Fail.MatchString(body):
		exitCode = ExitFail
	default:
		return ParseResult{
			Result:   "",
			ExitCode: ExitFail,
			Error:    fmt.Errorf("no matched"),
		}
	}
	var result string
	lines := strings.Split(body, "\n")
	var i int
	for idx, line := range lines {
		if p.Pass.MatchString(line) || p.Fail.MatchString(line) {
			i = idx
			break
		}
	}
	result = strings.Join(trimLastNewline(lines[i:]), "\n")
	return ParseResult{
		Result:   result,
		ExitCode: exitCode,
		Error:    nil,
	}
}

func trimLastNewline(s []string) []string {
	if len(s) == 0 {
		return s
	}
	last := len(s) - 1
	if s[last] == "" {
		return s[:last]
	}
	return s
}
