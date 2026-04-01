package ingester

func LoadGitRepo(base string) map[string]interface{} {
	// simulate repo structure
	return map[string]interface{}{
		"repo": map[string]interface{}{
			"path": base,
		},
	}
}
