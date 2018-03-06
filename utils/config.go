package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type LoadConfig struct {
	config map[string]interface{}
	params map[string]interface{}
}

func NewLoadConfig(argsWithoutProg []string) *LoadConfig {
	p := new(LoadConfig)
	p.parseCommandLineArgs(argsWithoutProg)

	if !p.IssetParam("config") {
		p.load("migrate.json")
		return p
	}

	p.load(p.GetParam("config"))
	return p
}

func (c *LoadConfig) parseCommandLineArgs(argsWithoutProg []string) {
	var params = make(map[string]interface{})
	for _, v := range argsWithoutProg {
		values := strings.Split(v, "=")
		if len(values) > 1 {
			reg := regexp.MustCompile("[-]+")
			replaceStr := reg.ReplaceAllString(values[0], "")
			params[replaceStr] = values[1]
		}
	}
	c.params = params
}

func (c *LoadConfig) load(url string) {
	file, _ := ioutil.ReadFile(url)
	if err := json.Unmarshal(file, &c.config); err != nil {
		panic(fmt.Sprintf("%v", "Config not found"))
		fmt.Println("Config not found")
		os.Exit(1)
	}
}

func (c *LoadConfig) GetAll() map[string]interface{} {
	return c.config
}

func (c *LoadConfig) GetStr(key string) string {
	return MapGetString(c.config, key)
}

func (c *LoadConfig) GetBool(key string) bool {
	return MapGetBool(c.config, key)
}

func (c *LoadConfig) GetParam(key string) string {
	return MapGetString(c.params, key)
}

func (c *LoadConfig) IssetParam(key string) bool {
	// if _, ok := params["name"]; !ok {}
	return MapContain(c.params, key)
}
