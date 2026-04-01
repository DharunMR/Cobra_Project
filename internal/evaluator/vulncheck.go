package evaluator

import "fmt"

type VulnCheck struct{}

func NewVulnCheck() *VulnCheck {
	return &VulnCheck{}
}

func (v *VulnCheck) Evaluate(input map[string]interface{}, mock map[string]interface{}) (bool, error) {

	if mock == nil {
		return false, nil
	}

	db := mock["vulnerability_db"].(map[string]interface{})

	pr := input["pull_request"].(map[string]interface{})
	diff := pr["diff"].(map[string]interface{})
	files := diff["files"].([]interface{})

	for _, f := range files {
		file := f.(map[string]interface{})
		deps := file["added_dependencies"].([]interface{})

		for _, d := range deps {
			dep := d.(map[string]interface{})
			key := fmt.Sprintf("%s@%s", dep["name"], dep["version"])

			if vuln, ok := db[key]; ok {
				vulnData := vuln.(map[string]interface{})
				if vulnData["vulnerable"].(bool) {
					return true, fmt.Errorf("vulnerability found: %s", key)
				}
			}
		}
	}

	return false, nil
}
