package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

func parseYaml(configFile string) (map[string]interface{}, error) {
	buf, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	var config map[string]interface{}
	if err = yaml.Unmarshal(buf, &config); err != nil {
		return nil, errors.New("error parsing config file")
	}
	return config, nil
}

func parseJson(configFile string) (map[string]interface{}, error) {
	panic("not implemented")
}

type File struct {
	Name        string
	Path        string
	Prefix      string
	Description string
}

// var _ = File{name: "a", path: "b", prefix: "c", description: "d"}

type Config struct {
	Scope      map[string][]File
	ConfigFile string
}

var _ = ParseConfigFile

func ParseConfigFile(configFile string) (Config, error) {
	if strings.HasSuffix(configFile, ".yaml") || strings.HasSuffix(configFile, ".yml") {
		config, err := parseYaml(configFile)
		if err != nil {
			return Config{}, err
		}
		if config["file"] == nil {
			return Config{}, fmt.Errorf("file key not found in config file")
		}
		file := config["file"].(map[interface{}]interface{})
		if file["scope"] == nil {
			return Config{}, fmt.Errorf("scope key not found in config file")
		}
		scopes := file["scope"].(map[interface{}]interface{})
		var configScope = map[string][]File{}
		for k, v := range scopes {
			scope := k.(string)
			vv := v.([]interface{})
			var files []File
			for _, vvv := range vv {
				vvvv := vvv.(map[interface{}]interface{})
				var prefix string
				var name string
				var path string
				var description string
				if vvvv["name"] == nil {
					return Config{}, fmt.Errorf("name key not found in config file")
				}
				name = vvvv["name"].(string)

				switch vvvv["prefix"].(type) {
				case string:
					prefix = vvvv["prefix"].(string)
				case []interface{}:
					prefixL := vvvv["prefix"].([]interface{})
					var prefixL2 []string
					for _, v := range prefixL {
						prefixL2 = append(prefixL2, v.(string))
					}
					prefix = strings.Join(prefixL2, ",")
				default:
					prefix = name
				}

				if vvvv["path"] == nil {
					return Config{}, fmt.Errorf("path key not found in config file")
				}
				path = vvvv["path"].(string)

				if vvvv["description"] == nil {
					description = ""
				}
				description = vvvv["description"].(string)

				f := File{
					Name:        name,
					Path:        path,
					Prefix:      prefix,
					Description: description,
				}
				// fmt.Println(scope, f)
				files = append(files, f)
				// fmt.Println(files)
			}
			configScope[scope] = files

		}
		// fmt.Println(configScope)

		return Config{Scope: configScope, ConfigFile: configFile}, nil
	} else if strings.HasSuffix(configFile, ".json") {
		config, err := parseJson(configFile)
		if err != nil {
			return Config{}, err
		}
		_ = config
		return Config{}, nil
	}
	return Config{}, fmt.Errorf("config file must be yaml")
}
