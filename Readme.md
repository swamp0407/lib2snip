# lib2snip

Lib2snip is a cli tool for generating snippets for vscode.

# Install
```
$ go install github.com/swamp0407/lib2snip@latest
```

# How To Use
```
$ go run ./main.go -c ./example/config.yaml -o ./.vscode/mysnip.code-snippets
```
This commands generate mysnip.code-snippets by reading config file: ./example/config.yaml.

1 This commands read config yaml file.
```
---
file:
  python:
    - name: sample1
      path: ./sample/1.py
      prefix:
        - sample1
        - s1
      description: sample1
    - name: sample2
      path: ./sample/2.py
      description: sample2
  go:
    - name: sample3
      path: ./sample/3.go
      description: sample3

```

2 Generate output.code-snipeets.json
```
{
  "sample1": {
    "body": "def sample1():\n    print(\"sample1\")\n",
    "prefix": [
      "sample1",
      "s1"
    ],
    "description": "sample1",
    "scope": "python"
  },
  "sample2": {
    "body": "def sample2():\n    print(\"sample2\")\n",
    "prefix": "sample2",
    "description": "sample2",
    "scope": "python"
  },
  "sample3": {
    "body": "func sample3() {\n\t//\n\tfmt.Println(\"sample3\")\n}\n",
    "prefix": "sample3",
    "description": "sample3",
    "scope": "go"
  }
}
```


