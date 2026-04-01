package suite

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadSuite(path string) (*TestSuite, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var suite TestSuite
	err = yaml.Unmarshal(data, &suite)
	return &suite, err
}
