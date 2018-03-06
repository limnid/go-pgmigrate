package utils

import (
	"log"
	"io/ioutil"
	"encoding/json"
)

func Serialize(values interface{}, model interface{}) {
	var data []byte

	stringValue, isString := values.(string)
	mapValue, isMap := values.(map[string]interface{})
	sliceValue, isSlice := values.([]interface{})

	if isSlice {
		jsonString, _ := json.Marshal(sliceValue)
		data = jsonString
	}

	if isMap {
		jsonString, _ := json.Marshal(mapValue)
		data = jsonString
	}

	if isMap {
		jsonString, _ := json.Marshal(mapValue)
		data = jsonString
	}

	if isString {
		data = []byte(stringValue)
	}

	if err := json.Unmarshal(data, &model); err != nil {
		println(err)
	}
}

func StringInMap(a string) bool {
	visitedURL := map[string]bool {
		"http://www.google.com": true,
		"https://paypal.com": true,
	}
	if visitedURL[a] {
		return true
	}
	return false
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func Check(err error) {
	if err != nil {
		panic(err)
		log.Fatal(err)
	}
}

func CreateFile(path string, content interface{}) {
	if content == nil {
		content = ""
	}
	d1 := []byte(content.(string))
	err_write := ioutil.WriteFile(path, d1, 0644)
	Check(err_write)
}