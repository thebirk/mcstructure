//go:build generate

package main

import (
	"encoding/json"
	"net/http"
	"os"
	"text/template"

	strcase "github.com/stoewer/go-strcase"
)

type block struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	DisplayName  string `json:"displayName"`
	VariableName string
}

func downloadBlocks() []block {
	resp, err := http.Get("https://raw.githubusercontent.com/PrismarineJS/minecraft-data/refs/heads/master/data/pc/1.21.3/blocks.json")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var blocks []block
	if err := json.NewDecoder(resp.Body).Decode(&blocks); err != nil {
		panic(err)
	}

	for i := range blocks {
		blocks[i].VariableName = strcase.UpperCamelCase(blocks[i].Name)
	}

	return blocks
}

const tmpl = `
package blocks

type Block struct {
	Id int
	Name string
	DisplayName string
	NamespacedName string
}

{{- range . }}
var {{ .VariableName }} = Block{ Id: {{.Id}}, Name: "{{.Name}}", DisplayName: "{{.DisplayName}}", NamespacedName: "minecraft:{{.Name}}"}{{ end }}
`

//go:generate go run $GOFILE
//go:generate go fmt blocks.go
func main() {
	blocks := downloadBlocks()

	f, err := os.OpenFile("blocks.go", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := template.Must(template.New("").Parse(tmpl)).Execute(f, blocks); err != nil {
		panic(err)
	}
}
