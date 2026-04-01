package evaluator

import (
	"context"
	"fmt"

	"github.com/open-policy-agent/opa/rego"
)

type RegoEvaluator struct {
	policy string
}

func NewRegoEvaluator(policy string) *RegoEvaluator {
	return &RegoEvaluator{
		policy: policy,
	}
}

func (r *RegoEvaluator) Evaluate(input map[string]interface{}, mock map[string]interface{}) (bool, error) {

	ctx := context.Background()

	// 🔥 Build query
	query, err := rego.New(
		rego.Query("data.minder.allow"),
		rego.Module("policy.rego", r.policy),
	).PrepareForEval(ctx)

	if err != nil {
		return false, fmt.Errorf("rego compile error: %w", err)
	}

	results, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return false, fmt.Errorf("rego eval error: %w", err)
	}

	// 🧠 Interpret result
	if len(results) == 0 || len(results[0].Expressions) == 0 {
		return true, fmt.Errorf("rule denied (no allow result)")
	}

	allowed, ok := results[0].Expressions[0].Value.(bool)
	if !ok {
		return true, fmt.Errorf("invalid allow result")
	}

	// 🔥 Minder semantics:
	// allow = true → PASS
	// allow = false → FAIL

	if allowed {
		return false, nil // PASS
	}

	return true, fmt.Errorf("rule denied")
}
