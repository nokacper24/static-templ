package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/nokacper24/static-templ/internal/finder"
	"github.com/nokacper24/static-templ/internal/generator"
)

const (
	outputScriptDirPath  string = "temp"
	outputScriptFileName string = "templ_static_generate_script.go"
)

func main() {
	var inputDir string
	var outputDir string
	flag.StringVar(&inputDir, "i", "web/pages", `Specify input directory.`)
	flag.StringVar(&outputDir, "o", "dist", `Specify output directory.`)
	flag.Parse()
	inputDir = strings.TrimRight(inputDir, "/")
	outputDir = strings.TrimRight(outputDir, "/")

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatal("err creating dirs:", err)
	}

	modName, err := finder.FindModulePath()
	if err != nil {
		log.Fatal("Error finding module name:", err)
	}

	allFiles, err := finder.FindAllFiles(inputDir)
	if err != nil {
		log.Fatal("Error finding go dirs:", err)
	}

	groupedFiles := allFiles.ToGroupedFiles()

	funcs, err := finder.FindFunctionsInFiles(groupedFiles.TemplGoFiles)
	if err != nil {
		log.Fatal("Error finding funcs:", err)
	} else if len(funcs) < 1 {
		log.Fatalf(`No components found in "%s"`, inputDir)
	}

	importsMap := map[string]bool{}
	for _, f := range funcs {
		importsMap[f.DirPath()] = true
	}
	var importsSlice []string
	for imp := range importsMap {
		importsSlice = append(importsSlice, fmt.Sprintf("%s/%s", modName, imp))
	}

	if err = os.RemoveAll(outputDir); err != nil {
		log.Fatal("err removing files", err)
	}

	if err = os.MkdirAll(outputScriptDirPath, os.ModePerm); err != nil {
		log.Fatal("err creating temp dir:", err)
	}

	for _, f := range groupedFiles.OtherFiles {
		if err := copyFile(f, strings.Replace(f, inputDir, outputDir, 1)); err != nil {
			log.Fatal("could not copy file: ", err)
		}
	}

	if err = generator.Generate(
		getOutputScriptPath(),
		importsSlice,
		funcs,
		inputDir,
		outputDir,
	); err != nil {
		log.Fatal("err generating script", err)
	}

	cmd := exec.Command("go", "run", getOutputScriptPath())
	if err = cmd.Start(); err != nil {
		log.Fatal("err starting script", err)
	}
	if err = cmd.Wait(); err != nil {
		log.Fatal("err running script", err)
	}

	if err = os.Remove(getOutputScriptPath()); err != nil {
		log.Fatal("err removing enerated script file", err)
	}
}

func getOutputScriptPath() string {
	return fmt.Sprintf("%s/%s", outputScriptDirPath, outputScriptFileName)
}

func copyFile(fromPath string, toPath string) error {
	if err := os.MkdirAll(path.Dir(toPath), os.ModePerm); err != nil {
		return err
	}

	src, err := os.ReadFile(fromPath)
	if err != nil {
		return err
	}

	if err = os.WriteFile(toPath, src, 0644); err != nil {
		return err
	}
	return nil
}
