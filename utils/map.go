package utils

func MapReduceTokens(tokens *[]string, list_limit int) *map[int][]string {
	token_list := make(map[int][]string)
	list_counter := 1
	list_id := 1
	for _, token := range *tokens {
		if list_counter > list_limit {
			list_id++
			list_counter = 1
		}
		token_list[list_id] = append(token_list[list_id], token)
		list_counter++
	}
	return &token_list
}

func MapGetInterface(values map[string]interface{}, key string) interface{} {
	if value, err := values[key]; err == true {
		return value
	}
	return nil
}

func MapGetString(values map[string]interface{}, key string) string {
	if value, err := values[key]; err == true {
		return value.(string)
	}
	return ""
}

func MapGetInt(values map[string]interface{}, key string) int {
	if value, err := values[key]; err == true {
		return value.(int)
	}
	return 0
}

func MapGetFloat(values map[string]interface{}, key string) float64 {
	if value, err := values[key]; err == true {
		return value.(float64)
	}
	return 0.0
}

func MapGetBool(values map[string]interface{}, key string) bool {
	if value, err := values[key]; err == true {
		return value.(bool)
	}
	return false
}

func MapContain(values map[string]interface{}, key string) bool {
	if _, err := values[key]; err == true {
		return true
	} else {
		return false
	}
}
