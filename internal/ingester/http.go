package ingester

func MockHTTP(status int, body string) map[string]interface{} {
	return map[string]interface{}{
		"http": map[string]interface{}{
			"status": status,
			"body":   body,
		},
	}
}
