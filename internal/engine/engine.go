package engine

import (
	"errors"

	"minder-test/internal/evaluator"
)

type Engine struct {
	EvalType string
	RegoCode string
}

func NewEngine(evalType, regoCode string) *Engine {
	return &Engine{
		EvalType: evalType,
		RegoCode: regoCode,
	}
}

func (e *Engine) Eval(input map[string]interface{}, mock map[string]interface{}) (bool, error) {

	switch e.EvalType {

	case "rego":
		r := evaluator.NewRegoEvaluator(e.RegoCode)
		return r.Evaluate(input, mock)

	case "vulncheck":
		v := evaluator.NewVulnCheck()
		return v.Evaluate(input, mock)

	default:
		return false, errors.New("unsupported eval type")
	}
}
