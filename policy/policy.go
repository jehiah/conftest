package policy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	nonTestRegoFiles = func(name string) bool {
		return filepath.Ext(name) == ".rego" && !strings.HasSuffix(name, "_test.rego")
	}
	allRegoFiles = func(name string) bool {
		return filepath.Ext(name) == ".rego"
	}
)

// ReadFiles returns all of the policy files (not including tests)
// at the given path(s) including its subdirectories.
func ReadFiles(paths ...string) ([]string, error) {
	return getPolicyFiles(paths, nonTestRegoFiles)
}

// ReadFilesWithTests returns all of the policies and test files
// at the given path(s) including its subdirectories.
// Test files are Rego files that have a suffix of _test.rego
func ReadFilesWithTests(paths ...string) ([]string, error) {
	return getPolicyFiles(paths, allRegoFiles)
}

func getPolicyFiles(paths []string, filter func(string) bool) ([]string, error) {
	var files []string
	for _, path := range paths {
		err := filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("walk path: %w", err)
			}

			if info.IsDir() {
				return nil
			}

			if filter(info.Name()) {
				if info.Size() == 0 {
					return fmt.Errorf("empty policy found in %s", currentPath)
				}

				files = append(files, currentPath)
			}

			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("search rego files: %w", err)
		}
	}

	if len(files) < 1 {
		return nil, fmt.Errorf("no policies found in %v", paths)
	}

	return files, nil
}
