package jobs

import (
	"os"
	"path/filepath"
	"strings"
)

type Job struct {
	UniqueName string
}

func getJobsBaseDir() string {
	return filepath.Join("./", "/jobs")
}

func LoadJobsFromDisk() ([]*Job, error) {
	path := getJobsBaseDir()
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, err
	}

	jobs := []*Job{}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		filename := e.Name()
		splitName := strings.Split(filename, ".")
		extension := splitName[len(splitName)-1]
		if extension != "job" {
			continue
		}
	}

	return jobs, nil
}

// Job file example:

// #- conf:
// - name=This is some job
// - restart=1
//
// #- fetch:
// some code to fetch the data
//
// #- prerun:
// executes before run
//
// #- run:
// some code to run the job (should be blocking)
//
