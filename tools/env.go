package tools

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var envs = make(map[string]string)

// GetAllEnv get the whole envs
func GetAllEnv() map[string]string {
	return envs
}

func GetEnv(key string) string {
	return envs[key]
}

// LoadConfigFromFile load env config file
func LoadConfigFromFile(path string) error {
	path = filepath.Clean(path)
	filePath := filepath.Join("/", path)

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		// skip empty line
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		// skip comment
		if line[0] == '#' {
			continue
		}

		params := strings.Split(line, "=")
		fmt.Printf("Params: %+v\n", params)

		// skip error format
		if len(params) < 2 {
			continue
		}

		var aggVal string
		for i := 1; i < len(params); i++ {
			if i != 1 {
				aggVal += "="
			}
			aggVal += params[i]
		}

		key, val := strings.TrimSpace(params[0]), strings.TrimSpace(aggVal)
		envs[key] = strings.Trim(val, "\"")
	}

	fmt.Printf("Load config: %s\n", path)

	return nil
}
