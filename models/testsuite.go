package models

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type TestSuite struct {
	Path    string
	Project *Project
}

func (ts *TestSuite) Name() string {
	return filepath.Base(ts.Path)
}

func (ts *TestSuite) Run() ([]TestResult, error) {
	tests, err := ts.GetTests()
	if err != nil {
		return nil, err
	}

	var futures []chan TestResult

	for _, test := range tests {
		c := make(chan TestResult)
		go func(t TestFile) {
			result, err := t.Run()
			if err != nil {
				log.Println(err)
			}
			c <- result
		}(test)
		futures = append(futures, c)
	}

	var results []TestResult

	for _, c := range futures {
		results = append(results, <-c)
	}

	return results, nil
}

func (ts *TestSuite) GetTests() ([]TestFile, error) {
	files, err := ioutil.ReadDir(ts.Path)
	if err != nil {
		return nil, err
	}

	var result []TestFile

	for _, file := range files {
		if file.IsDir() || strings.HasPrefix(file.Name(), "_") {
			continue
		}
		if path, err := ts.absPath(file); err == nil {
			result = append(result, TestFile{Path: path, Suite: ts})
		} else {
			return nil, err
		}
	}
	return result, nil
}

func (ts *TestSuite) absPath(fi os.FileInfo) (string, error) {
	return filepath.Abs(filepath.Join(ts.Path, fi.Name()))
}
