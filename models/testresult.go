package models

import "fmt"

type TestResult struct {
	File      *TestFile
	IsPassed  bool
	StdoutStr string
	StderrStr string
}

func (res *TestResult) LogLine() string {
	if res.IsPassed {
		return fmt.Sprintf("%v Passed", res.File.Name())
	} else {
		return fmt.Sprintf(
			"%v FAILED! STDERR: %v",
			res.File.Name(),
			res.StderrStr,
		)
	}
}
