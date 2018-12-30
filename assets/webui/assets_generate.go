// +build ignore

package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

const filenameVfsGen = "assets_vfsgen.go"

func main() {
	fs := http.Dir("dist")

	log.Printf("webui assets_generate.go")
	err := vfsgen.Generate(fs, vfsgen.Options{
		Filename:     filenameVfsGen,
		PackageName:  "webui",
		VariableName: "Assets",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
