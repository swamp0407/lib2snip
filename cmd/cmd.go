package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"lib2snip/entities"
	"os"
	"path/filepath"

	"errors"
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
	name        string
	path        string
	prefix      string
	description string
}

var _ = File{name: "a", path: "b", prefix: "c", description: "d"}

type Config struct {
	Scope      map[string][]File
	configFile string
}

func parseFile(configFile string) (Config, error) {
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
					name:        name,
					path:        path,
					prefix:      prefix,
					description: description,
				}
				// fmt.Println(scope, f)
				files = append(files, f)
				// fmt.Println(files)
			}
			configScope[scope] = files

		}
		// fmt.Println(configScope)

		return Config{Scope: configScope, configFile: configFile}, nil
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

func genSnippets(scope string, files []File, configPath string) []entities.Snippet {
	var snippets []entities.Snippet
	for _, file := range files {
		if !strings.HasPrefix(file.path, "/") {
			file.path = filepath.Join(filepath.Dir(configPath), file.path)
			buf, err := os.ReadFile(file.path)
			if err != nil {
				fmt.Println(err)
				continue
			}

			snippet := entities.Snippet{
				Name:        file.name,
				Prefix:      file.prefix,
				Scope:       scope,
				Description: file.description,
				Body:        string(buf),
			}
			snippets = append(snippets, snippet)
		}
	}
	return snippets
}

func genSnippetsWithScope(config Config) map[string][]entities.Snippet {
	snippetsWithScope := map[string][]entities.Snippet{}
	for scope, files := range config.Scope {
		snippets := genSnippets(scope, files, config.configFile)
		snippetsWithScope[scope] = snippets
	}
	return snippetsWithScope
}
func outputJson(snippetsWithScope map[string][]entities.Snippet, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	var Vssnippets = map[string]VsSnip{}
	for _, snippets := range snippetsWithScope {
		for _, snippet := range snippets {
			vssnippet := VsSnip{
				Prefix:      snippet.Prefix,
				Body:        snippet.Body,
				Description: snippet.Description,
				Scope:       snippet.Scope,
			}
			Vssnippets[snippet.Name] = vssnippet
		}
	}

	return enc.Encode(Vssnippets)
}

//	func outputYaml(snippetsWithScope map[string][]entities.Snippet, filename string) error {
//		panic("not implemented")
//	}
type VsSnip struct {
	Prefix      string `json:"prefix"`
	Body        string `json:"body"`
	Description string `json:"description"`
	Scope       string `json:"scope"`
}

func outputSnippets(snippetsWithScope map[string][]entities.Snippet, filename string) error {
	if strings.HasSuffix(filename, ".code-snippets") {
		return outputJson(snippetsWithScope, filename)
	}
	//  else if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
	// 	return outputYaml(snippetsWithScope, filename)
	// }
	return fmt.Errorf("output file must be .code-snippets")
}

func ask4confirm() bool {
	var s string

	fmt.Printf("(y/N): ")
	_, err := fmt.Scan(&s)
	if err != nil {
		panic(err)
	}

	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	if s == "y" || s == "yes" {
		return true
	}
	return false
}

func Run() {
	flag.Parse()
	configFile := flag.Lookup("c").Value.String()

	config, err := parseFile(configFile)
	if err != nil {
		panic(err)
	}

	snippetsWithScope := genSnippetsWithScope(config)

	oflag := flag.Lookup("o").Value.String()

	filedir := filepath.Dir(oflag)
	filebase := filepath.Base(oflag)
	if !strings.HasSuffix(filebase, ".code-snippets") {
		fmt.Println("output file must be .code-snippets")
		os.Exit(1)
	}
	if filebase == ".code-snippets" {
		filebase = "snippets.code-snippets"
	}
	var filename string
	var splitFiledir []string = strings.Split(filedir, "/")
	if splitFiledir[len(splitFiledir)-1] != ".vscode" {
		filedir = filepath.Join(filedir, ".vscode")
	}
	if _, err := os.Stat(filedir); os.IsNotExist(err) {
		fmt.Printf("dir %s not exist, Can create it?", filepath.Dir(filename))
		if ask4confirm() {
			os.MkdirAll(filepath.Dir(filename), 0755)
		} else {
			os.Exit(1)
		}
	}
	filename = filepath.Join(filedir, filebase)

	fmt.Println(snippetsWithScope["python"][1].Body)
	outputSnippets(snippetsWithScope, filename)
}

func init() {
	flag.String("c", "config.yaml", "config file")
	// flag.Bool("f", false, "overwrite snippets file") # not implemented
	// flag.Bool("debug", false, "debug mode") # not implemented
	flag.String("o", "output.code-snippets", "output filename")
}
