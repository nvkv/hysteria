package models

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strconv"
)

type TestFile struct {
	Path    string
	Suite   *TestSuite
	Timeout int
}

func (tf *TestFile) Name() string {
	return filepath.Base(tf.Path)
}

func (tf *TestFile) Run() (TestResult, error) {

	result := TestResult{
		File:      tf,
		IsPassed:  false,
		StdoutStr: "",
		StderrStr: "",
	}

	if tf.Suite != nil && tf.Suite.Project != nil {
		tf.Timeout = tf.Suite.Project.TestTimeout
	} else {
		tf.Timeout = 30 // Fallback default value for script timout
	}

	/**
	 * First of all, make test file executable, just in case it is not
	 */
	chmodErr := exec.Command("chmod", "+x", tf.Path).Run()
	if chmodErr != nil {
		return result, chmodErr
	}

	testCmd := exec.Command("timeout", strconv.Itoa(tf.Timeout), tf.Path)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	testCmd.Stdout = &stdout
	testCmd.Stderr = &stderr

	testErr := testCmd.Run()

	result.StdoutStr = stdout.String()
	result.StderrStr = stderr.String()

	result.IsPassed = testErr == nil

	return result, nil
}
