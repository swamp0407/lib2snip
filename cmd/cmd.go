package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/swamp0407/lib2snip/config"

	"strings"
)

func outputJson(v config.VsOutput, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
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

func hasValidSuffixV(filename string) bool {
	return filepath.Ext(filename) == ".code-snippets"
}

func processOFlag() string {
	oflag := flag.Lookup("o").Value.String()
	filedir := filepath.Dir(oflag)
	filebase := filepath.Base(oflag)

	if filebase == ".code-snippets" {
		filebase = "snippets.code-snippets"
	}
	if !hasValidSuffixV(oflag) {
		filebase = filebase + ".code-snippets"
	}

	var splitFiledir []string = strings.Split(filedir, "/")
	if splitFiledir[len(splitFiledir)-1] != ".vscode" {
		filedir = filepath.Join(filedir, ".vscode")
	}
	if _, err := os.Stat(filedir); os.IsNotExist(err) {
		fmt.Printf("dir %s not exist, Can create it?", filedir)
		if ask4confirm() {
			os.MkdirAll(filedir, 0755)
		} else {
			os.Exit(1)
		}
	}
	filename := filepath.Join(filedir, filebase)
	if _, err := os.Stat(filename); err == nil {
		fmt.Printf("file %s exist, Can overwrite it?", filename)
		if !ask4confirm() {
			os.Exit(1)
		}
	}
	return filename
}

func Run() {
	flag.Parse()
	configFile := flag.Lookup("c").Value.String()
	c := config.NewConfig(configFile)

	if err := c.ParseConfigFile(); err != nil {
		fmt.Println(err)
		usage()
		os.Exit(1)
	}

	v := c.GenVsOutput()
	filename := processOFlag()
	outputJson(v, filename)

}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])

	flag.PrintDefaults()
}

func init() {
	flag.String("c", "config.yaml", "config file")
	// flag.Bool("f", false, "overwrite snippets file") # not implemented
	flag.Bool("debug", false, "debug mode")
	flag.String("o", "output.code-snippets", "output filename")
}
