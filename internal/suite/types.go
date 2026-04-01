package suite

type TestSuite struct {
	Version string     `yaml:"version"`
	Tests   []TestCase `yaml:"tests"`
}

type TestCase struct {
	Name   string                 `yaml:"name"`
	Input  map[string]interface{} `yaml:"input"`
	Mock   map[string]interface{} `yaml:"mock"`
	Expect string                 `yaml:"expect"` // pass/fail
}
