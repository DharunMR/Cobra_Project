package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/yourusername/minder-ruletest-cli/pkg/ruletester"
)

var (
	ruleFile string
	testFile string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "minder",
		Short: "Minder - Security Rule and Profile Manager",
	}

	testCmd := &cobra.Command{
		Use:   "test",
		Short: "Run tests against a specific Minder rule definition",
		Run:   runTestCmd,
	}

	testCmd.Flags().StringVarP(&ruleFile, "rule", "f", "", "Path to the rule definition YAML")
	testCmd.Flags().StringVarP(&testFile, "test", "t", "", "Path to the test suite YAML")
	_ = testCmd.MarkFlagRequired("rule")
	_ = testCmd.MarkFlagRequired("test")

	rootCmd.AddCommand(testCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runTestCmd(cmd *cobra.Command, args []string) {
	zerolog.SetGlobalLevel(zerolog.FatalLevel)

	fmt.Printf("\n==================== Minder Rule Test Session Starts ====================\n")
	fmt.Printf("Rule Target : %s\n", ruleFile)
	fmt.Printf("Test Suite  : %s\n", testFile)
	fmt.Printf("-------------------------------------------------------------------------\n\n")

	startTime := time.Now()

	results, err := ruletester.RunTestSuite(context.Background(), ruleFile, testFile)
	if err != nil {
		color.Red("Fatal Error: %v\n", err)
		os.Exit(1)
	}

	var total, passed, failed int

	for _, res := range results {
		total++
		if res.Error != nil {
			failed++
			color.Red("  ✗ %s -> ERROR: %v\n", res.Name, res.Error)
		} else if !res.Passed {
			failed++
			color.Red("  ✗ %s -> FAILED: %s\n", res.Name, res.Message)
		} else {
			passed++
			color.Green("  ✓ %s\n", res.Name)
		}
	}

	duration := time.Since(startTime)
	fmt.Printf("\n========================= Short Test Summary Info =========================\n")

	if failed > 0 {
		color.Red("FAILED %d/%d tests in %v\n", failed, total, duration)
		os.Exit(1)
	} else {
		color.Green("PASSED %d tests in %v\n", passed, duration)
		os.Exit(0)
	}
}
