package metrics

func StatusStringFromError(err error) string {
	if err != nil {
		return "error"
	}

	return "success"
}
