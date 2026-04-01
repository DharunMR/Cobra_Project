package evaluator

type Evaluator interface {
	Evaluate(input map[string]interface{}, mock map[string]interface{}) (bool, error)
}
