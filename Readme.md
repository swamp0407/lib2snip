# lib2snip

Lib2snip is a cli tool for generating snippets for vscode.


# Install
```
$ go install github.com/swamp0407/lib2snip
```

# How To Use
```
$ go run ./main.go -c ./example/config.yaml -o ./.vscode/mysnip.code-snippets
```
This commands generate mysnip.code-snippets by reading config file: ./example/config.yaml.

