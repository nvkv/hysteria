package models

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Project struct {
	Path        string
	TestTimeout int
}

func (p *Project) Name() string {
	return filepath.Base(p.Path)
}

func (p *Project) absPath(fi os.FileInfo) (string, error) {
	return filepath.Abs(filepath.Join(p.Path, fi.Name()))
}

func (p *Project) GetTestSuites() ([]TestSuite, error) {
	files, err := ioutil.ReadDir(p.Path)
	if err != nil {
		return nil, err
	}

	var results []TestSuite
	for _, file := range files {
		if file.IsDir() == false || strings.HasPrefix(file.Name(), ".") {
			continue
		}
		if path, err := p.absPath(file); err == nil {
			results = append(results, TestSuite{Path: path, Project: p})
		} else {
			return nil, err
		}
	}
	return results, nil
}
