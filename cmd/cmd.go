package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/swamp0407/lib2snip/config"
	"github.com/swamp0407/lib2snip/entities"

	"strings"
)

func genSnippets(scope string, files []config.File, configPath string) []entities.Snippet {
	var snippets []entities.Snippet
	for _, file := range files {
		if !strings.HasPrefix(file.Path, "/") {
			file.Path = filepath.Join(filepath.Dir(configPath), file.Path)
			buf, err := os.ReadFile(file.Path)
			if err != nil {
				fmt.Println(err)
				continue
			}

			snippet := entities.Snippet{
				Name:        file.Name,
				Prefix:      file.Prefix,
				Scope:       scope,
				Description: file.Description,
				Body:        string(buf),
			}
			snippets = append(snippets, snippet)
		}
	}
	return snippets
}

func genSnippetsWithScope(c config.Config) map[string][]entities.Snippet {
	snippetsWithScope := map[string][]entities.Snippet{}
	for scope, files := range c.Scope {
		snippets := genSnippets(scope, files, c.ConfigFile)
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

	c, err := config.ParseConfigFile(configFile)
	if err != nil {
		panic(err)
	}

	snippetsWithScope := genSnippetsWithScope(c)

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

	outputSnippets(snippetsWithScope, filename)
}

func init() {
	flag.String("c", "config.yaml", "config file")
	// flag.Bool("f", false, "overwrite snippets file") # not implemented
	// flag.Bool("debug", false, "debug mode") # not implemented
	flag.String("o", "output.code-snippets", "output filename")
}
