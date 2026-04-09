package ruletester

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	minderv1 "github.com/mindersec/minder/pkg/api/protobuf/go/minder/v1"
	rtengine "github.com/mindersec/minder/pkg/engine/v1/rtengine"
	tkv1 "github.com/mindersec/minder/pkg/testkit/v1"
	"gopkg.in/yaml.v3"
)

type TestResult struct {
	Name    string
	Passed  bool
	Message string
	Error   error
}

func RunTestSuite(ctx context.Context, rulePath, testPath string) ([]TestResult, error) {
	var results []TestResult

	rt, err := openRuleType(rulePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open rule type: %w", err)
	}

	suite, err := openTestSuite(testPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open test suite: %w", err)
	}

	if rt.Context == nil {
		rt.Context = &minderv1.Context{}
	}
	prjName := "minder-cli-test"
	rt.Context.Project = &prjName

	rtDataPath := strings.TrimSuffix(rulePath, filepath.Ext(rulePath)) + ".testdata"

	for _, tc := range suite.Tests {
		res := runSingleTest(ctx, rt, &tc, rtDataPath)
		results = append(results, res)
	}

	return results, nil
}

func runSingleTest(ctx context.Context, rt *minderv1.RuleType, tc *RuleTest, rtDataPath string) TestResult {
	res := TestResult{Name: tc.Name, Passed: false}

	var opts []tkv1.Option
	switch rt.Def.Ingest.Type {
	case "git":
		if tc.Git == nil {
			res.Error = fmt.Errorf("git test configuration is missing")
			return res
		}
		opts = append(opts, tkv1.WithGitDir(filepath.Join(rtDataPath, tc.Git.RepoBase)))
	case "rest":
		opts = append(opts, httpTestOpts(tc, rtDataPath))
	default:
		res.Error = fmt.Errorf("unsupported ingest type %s", rt.Def.Ingest.Type)
		return res
	}

	tk := tkv1.NewTestKit(opts...)
	rte, err := rtengine.NewRuleTypeEngine(ctx, rt, tk)
	if err != nil {
		res.Error = fmt.Errorf("failed to initialize rule engine: %w", err)
		return res
	}

	val := rte.GetRuleInstanceValidator()
	if err := val.ValidateRuleDefAgainstSchema(tc.Def); err != nil {
		res.Error = fmt.Errorf("rule definition schema validation failed: %w", err)
		return res
	}

	if tk.ShouldOverrideIngest() {
		rte.WithCustomIngester(tk)
	}

	if tc.Params == nil {
		tc.Params = make(map[string]any)
	}

	_, err = rte.Eval(ctx, tc.Entity.Entity, tc.Def, tc.Params, tkv1.NewVoidResultSink())

	expectedAction := tc.Expect
	if expectedAction == "allow" {
		expectedAction = "pass"
	} else if expectedAction == "deny" {
		expectedAction = "fail"
	}

	if expectedAction == "pass" {
		if err != nil {
			res.Message = fmt.Sprintf("Expected PASS, but got error: %v", err)
			return res
		}
		res.Passed = true
	} else if expectedAction == "fail" {
		if err == nil {
			res.Message = "Expected FAIL, but rule passed"
			return res
		}
		if tc.ErrorText != "" && !strings.Contains(strings.TrimSpace(err.Error()), strings.TrimSpace(tc.ErrorText)) {
			res.Message = fmt.Sprintf("Expected error text '%s', got '%v'", tc.ErrorText, err)
			return res
		}
		res.Passed = true
	} else {
		res.Message = fmt.Sprintf("Unknown expectation: %s", tc.Expect)
	}

	return res
}

func httpTestOpts(tc *RuleTest, rtDataPath string) tkv1.Option {
	if tc.HTTP == nil {
		if tc.MockIngest != nil {
			bodyBytes, _ := json.Marshal(tc.MockIngest)
			return tkv1.WithHTTP(http.StatusOK, bodyBytes, nil)
		}
		return tkv1.WithHTTP(http.StatusOK, []byte(""), nil)
	}

	if tc.HTTP.Status == 0 {
		tc.HTTP.Status = http.StatusOK
	}

	if tc.HTTP.Headers == nil {
		tc.HTTP.Headers = make(map[string]string)
	}

	if tc.HTTP.BodyFile != "" {
		data, err := os.ReadFile(filepath.Join(rtDataPath, tc.HTTP.BodyFile))
		if err == nil {
			tc.HTTP.Body = string(data)
		}
	}
	return tkv1.WithHTTP(tc.HTTP.Status, []byte(tc.HTTP.Body), tc.HTTP.Headers)
}

func openTestSuite(testPath string) (*RuleTestSuite, error) {
	file, err := os.Open(testPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	suite := &RuleTestSuite{}
	if err := yaml.NewDecoder(file).Decode(suite); err != nil {
		return nil, err
	}
	return suite, nil
}

func openRuleType(path string) (*minderv1.RuleType, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	rt := &minderv1.RuleType{}
	if err := minderv1.ParseResource(file, rt); err != nil {
		return nil, err
	}
	return rt, nil
}
