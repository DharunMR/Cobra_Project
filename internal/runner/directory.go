package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func RunDirectory(dir string) {

	total := 0
	passed := 0

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

		if err != nil || info.IsDir() {
			return nil
		}

		// skip test files
		if strings.HasSuffix(path, ".test.yaml") {
			return nil
		}

		// only rule yaml
		if !strings.HasSuffix(path, ".yaml") {
			return nil
		}

		testPath := strings.TrimSuffix(path, ".yaml") + ".test.yaml"

		if _, err := os.Stat(testPath); os.IsNotExist(err) {
			fmt.Println("⚠ No test file for:", path)
			return nil
		}

		fmt.Println("\n=== Running:", filepath.Base(path), "===")

		ok := RunSingle(path, testPath)

		total++

		if ok {
			passed++
		}

		return nil
	})

	fmt.Printf("\nSummary: %d/%d rule files passed\n", passed, total)
}
