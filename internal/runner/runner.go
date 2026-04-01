package runner

import (
	"os"

	"gopkg.in/yaml.v3"

	"minder-test/internal/engine"
	"minder-test/internal/output"
	"minder-test/internal/suite"
)

func extractRule(rulePath string) (evalType string, regoCode string) {

	data, _ := os.ReadFile(rulePath)

	var raw map[string]interface{}
	yaml.Unmarshal(data, &raw)

	def := raw["def"].(map[string]interface{})
	eval := def["eval"].(map[string]interface{})

	evalType = eval["type"].(string)

	if evalType == "rego" {
		regoBlock := eval["rego"].(map[string]interface{})
		regoCode = regoBlock["def"].(string)
	}

	return
}

func RunSingle(rulePath, testPath string) bool {

	evalType, regoCode := extractRule(rulePath)

	suiteData, err := suite.LoadSuite(testPath)
	if err != nil {
		panic(err)
	}

	eng := engine.NewEngine(evalType, regoCode)

	passed := 0

	for _, t := range suiteData.Tests {

		result, err := eng.Eval(t.Input, t.Mock)

		success := false

		if t.Expect == "pass" && !result {
			success = true
		}

		if t.Expect == "fail" && result {
			success = true
		}

		output.PrintResult(t.Name, success, err)

		if success {
			passed++
		}
	}

	output.PrintSummary(passed, len(suiteData.Tests))

	return passed == len(suiteData.Tests)
}
