package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/nokacper24/static-templ/internal/finder"
	"github.com/nokacper24/static-templ/internal/generator"
)

//go:embed .version
var versionFile embed.FS

const templVersion = "0.2.731"

const (
	outputScriptDirPath  string = "temp"
	outputScriptFileName string = "templ_static_generate_script.go"
)

func main() {
	versionCmd := flag.NewFlagSet("version", flag.ExitOnError)

	if len(os.Args) < 2 {
		usage()
		return
	}

	switch os.Args[1] {
	case "version", "--version":
		versionCmd.Parse(os.Args[2:])
		printVersion(getVersion(), templVersion)
		return
	default:
		// Continue with existing flag parsing
	}

	var inputDir, outputDir string
	var runFormat, runGenerate, debug bool

	flag.StringVar(&inputDir, "i", "web/pages", "Specify input directory.")
	flag.StringVar(&outputDir, "o", "dist", "Specify output directory.")
	flag.BoolVar(&runFormat, "f", false, "Run templ fmt.")
	flag.BoolVar(&runGenerate, "g", false, "Run templ generate.")
	flag.BoolVar(&debug, "d", false, "Keep the generation script after completion for inspection and debugging.")
	flag.Usage = usage
	flag.Parse()

	inputDir = strings.TrimRight(inputDir, "/")
	outputDir = strings.TrimRight(outputDir, "/")

	if outputDir != inputDir {
		if err := clearAndCreateDir(outputDir); err != nil {
			log.Fatal("Error preparing output directory:", err)
		}
	}

	modulePath, err := finder.FindModulePath()
	if err != nil {
		log.Fatal("Error finding module name:", err)
	}

	groupedFiles, err := finder.FindFilesInDir(inputDir)
	if err != nil {
		log.Fatal("Error finding files:", err)
	}

	if runFormat {
		done := make(chan struct{})
		go func() {
			err := generator.RunTemplFmt(groupedFiles.TemplFiles, done)
			if err != nil {
				log.Fatalf("failed to run 'templ fmt' command: %v", err)
			}
		}()
		<-done
		log.Println("completed running 'templ fmt'")
	}

	if runGenerate {
		done := make(chan struct{})
		go func() {
			err := generator.RunTemplGenerate(done)
			if err != nil {
				log.Fatalf("failed to run 'templ generate' command: %v", err)
			}
		}()
		<-done
		log.Println("completed running 'templ generate'")
		// run groupedFiles again after templ generate command execution to get TemplGoFiles updated
		groupedFiles, err = finder.FindFilesInDir(inputDir)
		if err != nil {
			log.Fatal("Error finding _templ.go files after templ generate completion:", err)
		}
	}

	funcs, err := finder.FindFunctionsInFiles(groupedFiles.TemplGoFiles)
	if err != nil {
		log.Fatal("Error finding funcs:", err)
	} else if len(funcs) < 1 {
		log.Fatalf(`No components found in "%s"`, inputDir)
	}

	if err = os.MkdirAll(outputScriptDirPath, os.ModePerm); err != nil {
		log.Fatal("err creating temp dir:", err)
	}

	if err = copyFilesIntoOutputDir(groupedFiles.OtherFiles, inputDir, outputDir); err != nil {
		log.Fatal("err copying files:", err)
	}

	if err = generator.Generate(
		getOutputScriptPath(),
		finder.FindImports(funcs, modulePath),
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

	if !debug {
		if err = os.RemoveAll(outputScriptDirPath); err != nil {
			log.Fatal("err removing script folder", err)
		}
	}
}

// Display usage information
func usage() {
	output := fmt.Sprintf(`Usage of %[1]v:
%[1]v [flags] [subcommands]

Flags:
  -i  Specify input directory (default "web/pages").
  -o  Specify output directory (default "dist").
  -f  Run templ fmt.
  -g  Run templ generate.
  -d  Keep the generation script after completion for inspection and debugging.

Subcommands:
  version  Display the version information.

Examples:
  # Specify input and output directories
  %[1]v -i web/demos -o output

  # Specify input directory, run templ generate and output to default directory
  %[1]v -i web/demos -g=true

  # Display the version information
  %[1]v version
`, os.Args[0])

	fmt.Println(output)
}

func getVersion() string {
	content, err := versionFile.ReadFile(".version")
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(content))
}

func printVersion(version, templVersion string) {
	templModulePath := "github.com/a-h/templ"
	fmt.Printf("Version: %s (built with %s@v%s)\n", version, templModulePath, templVersion)
}

func clearAndCreateDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return os.MkdirAll(dir, os.ModePerm)
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

func copyFilesIntoOutputDir(files []string, inputDir string, outputDir string) error {
	for _, f := range files {
		if err := copyFile(f, strings.Replace(f, inputDir, outputDir, 1)); err != nil {
			return err
		}
	}
	return nil
}

func getOutputScriptPath() string {
	return filepath.Join(outputScriptDirPath, outputScriptFileName)
}
