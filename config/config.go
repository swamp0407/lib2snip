package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/swamp0407/lib2snip/entities"
	"gopkg.in/yaml.v2"
)

func (c *Config) parseYamlConfig() error {
	buf, err := os.ReadFile(c.ConfigFile)
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(buf, c); err != nil {
		return errors.New("error parsing Yaml config file")
	}
	return nil
}

func (c *Config) parseJsonConfig() error {
	panic("not implemented")
}

type File struct {
	Name        string      `yaml:"name"`
	Path        string      `yaml:"path"`
	Prefix      interface{} `yaml:"prefix"` // string or []string
	Description string      `yaml:"description"`
}
type Scope = map[string][]File

type Config struct {
	File       Scope `yaml:"file"`
	ConfigFile string
}

type VsOutput = map[string][]entities.Snippet

func NewConfig(configFile string) Config {
	return Config{ConfigFile: configFile}
}

func (c *Config) ParseConfigFile() error {
	var err error
	if c.ConfigFile == "" {
		return fmt.Errorf("config file is required")
	}

	if strings.HasSuffix(c.ConfigFile, ".yaml") || strings.HasSuffix(c.ConfigFile, ".yml") {
		err = c.parseYamlConfig()
	} else if strings.HasSuffix(c.ConfigFile, ".json") {
		err = c.parseJsonConfig()
	} else {
		return fmt.Errorf("config file must be yaml or json")
	}
	return err
}

func validateAndEvaluatePrefix(prefix_i interface{}) (string, error) {
	var prefix string
	switch prefix_i := prefix_i.(type) {
	case string:
		prefix = prefix_i
	case []interface{}:
		prefixList := make([]string, 0)
		for _, v := range prefix_i {
			if _, ok := v.(string); !ok {
				return "", errors.New("prefix is not valid")
			}

			prefixList = append(prefixList, v.(string))
		}
		prefix = strings.Join(prefixList, ",")
	default:
		return "", nil
	}
	return prefix, nil
}

func readFromPath(path string) (*[]byte, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}

func (c *Config) customReadFromPath(path string) (*[]byte, error) {
	if !strings.HasPrefix(path, "/") {
		path = filepath.Join(filepath.Dir(c.ConfigFile), path)
	}
	return readFromPath(path)

}

func (c *Config) convertFile2Snippet(file File) (*entities.Snippet, error) {
	var prefix string
	var name = file.Name
	var description = file.Description
	var path = file.Path
	if name == "" {
		return nil, errors.New("name is required")
	}
	if path == "" {
		return nil, errors.New("path is required")
	}
	fmt.Println("path: ", path)
	buf, err := c.customReadFromPath(path)
	if err != nil {
		return nil, err
	}
	prefix, err = validateAndEvaluatePrefix(file.Prefix)
	if err != nil {
		return nil, errors.New("prefix is not valid")
	}
	if prefix == "" {
		prefix = name
	}

	snippet := entities.Snippet{
		Name:        name,
		Body:        string(*buf),
		Prefix:      prefix,
		Description: description,
	}
	return &snippet, nil
}

func (c *Config) genSnippets(files []File) []entities.Snippet {
	var snippets []entities.Snippet
	for _, file := range files {
		snippet, err := c.convertFile2Snippet(file)
		if err != nil {
			fmt.Println(err)
			continue
		}
		snippets = append(snippets, *snippet)
	}
	return snippets
}

func (c *Config) GenVsOutput() VsOutput {
	output := make(VsOutput)
	for scope, files := range c.File {
		snippets := c.genSnippets(files)
		output[scope] = snippets
	}
	return output
}
