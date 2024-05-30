package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/nokacper24/templ-static-generator/internal/finder"
	"github.com/nokacper24/templ-static-generator/internal/generator"
)

func main() {
	dirPath := "dist/"
	inputPath := "web/pages/"
	outputScriptPath := "templ_static_generate_script.go"

	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Fatal("err creating dirs:", err)
	}

	modName, err := finder.FindModulePath()
	if err != nil {
		log.Fatal("Error finding module name:", err)
	}

	dirs, err := finder.FindDirsWithGoFiles(inputPath)
	if err != nil {
		log.Fatal("Error finding go dirs:", err)
	}

	funcs, err := finder.FindFuncsToCall(dirs)
	if err != nil {
		log.Fatal("Error finding funcs:", err)
	}

	importsMap := map[string]bool{}
	for _, f := range funcs {
		importsMap[f.Location("")] = true
	}
	var importsSlice []string
	for imp := range importsMap {
		importsSlice = append(importsSlice, fmt.Sprintf("%s/%s", modName, imp))
	}

	err = os.RemoveAll(dirPath)
	if err != nil {
		log.Fatal("err removing files", err)
	}

	err = generator.Generate(outputScriptPath, importsSlice, funcs)
	if err != nil {
		log.Fatal("err generating script", err)
	}

	cmd := exec.Command("go", "run", outputScriptPath)
	err = cmd.Start()
	if err != nil {
		log.Fatal("err starting script", err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Fatal("err running script", err)
	}

	err = os.Remove(outputScriptPath)
	if err != nil {
		log.Fatal("err removing enerated script file", err)
	}
}
