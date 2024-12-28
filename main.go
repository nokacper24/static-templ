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

// Embed the version file
//
//go:embed .version
var versionFile embed.FS

// Constants for templ version and script paths
const (
	templVersion         = "0.2.793"
	outputScriptDirPath  = "temp"
	outputScriptFileName = "templ_static_generate_script.go"
)

// Constants for operational modes
const (
	modeBundle = "bundle"
	modeInline = "inline"
)

// Struct to hold command line flags
type flags struct {
	InputDir    string
	OutputDir   string
	Mode        string
	RunFormat   bool
	RunGenerate bool
	Debug       bool
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version", "--version":
			handleVersionCmd()
			return
		default:
			// Continue with existing flag parsing
		}
	}

	// Parse command line flags
	flags := parseFlags()

	// Prepare output directory if necessary
	if flags.Mode == modeBundle && flags.OutputDir != flags.InputDir {
		if err := clearAndCreateDir(flags.OutputDir); err != nil {
			log.Fatal("Error preparing output directory:", err)
		}
	}

	// Prepare directories and find files
	modulePath, groupedFiles := prepareDirectories(flags.InputDir)

	// Run templ fmt if specified
	if flags.RunFormat {
		runTemplFmt(groupedFiles)
	}

	// Run templ generate if specified
	if flags.RunGenerate {
		groupedFiles = runTemplGenerate(flags.InputDir)
	}

	// Find functions to call in the generated Go files
	funcs := findFunctions(groupedFiles.TemplGoFiles)

	// Create output script directory
	if err := os.MkdirAll(outputScriptDirPath, os.ModePerm); err != nil {
		log.Fatalf("Error creating temp dir: %v", err)
	}

	// Handle modes
	switch flags.Mode {
	case modeBundle:
		log.Println("Operational mode: bundle")
		handleBundleMode(funcs, modulePath, flags.InputDir, flags.OutputDir, groupedFiles.OtherFiles, flags.Debug)
	case modeInline:
		log.Println("Operational mode: inline")
		handleInlineMode(funcs, modulePath, flags.InputDir, flags.Debug)
	default:
		log.Fatalf("Unknown mode: %s", flags.Mode)
	}
}

// Handle the version command to display the version information
func handleVersionCmd() {
	versionCmd := flag.NewFlagSet("version", flag.ExitOnError)
	err := versionCmd.Parse(os.Args[2:])
	if err != nil {
		return
	}
	printVersion(getVersion(), templVersion)
}

// Parse command line flags and return a flags struct
func parseFlags() flags {
	var flags flags

	flag.StringVar(&flags.Mode, "m", "bundle", "Set the operational mode: bundle or inline.")
	flag.StringVar(&flags.InputDir, "i", "web/pages", "Specify input directory.")
	flag.StringVar(&flags.OutputDir, "o", "dist", "Specify output directory.")
	flag.BoolVar(&flags.RunFormat, "f", false, "Run templ fmt.")
	flag.BoolVar(&flags.RunGenerate, "g", false, "Run templ generate.")
	flag.BoolVar(&flags.Debug, "d", false, "Keep the generation script after completion for inspection and debugging.")
	flag.Usage = usage
	flag.Parse()

	flags.InputDir = strings.TrimRight(flags.InputDir, "/")
	flags.OutputDir = strings.TrimRight(flags.OutputDir, "/")

	return flags
}

// Prepare directories and find files
func prepareDirectories(inputDir string) (string, *finder.GroupedFiles) {
	modulePath, err := finder.FindModulePath()
	if err != nil {
		log.Fatalf("Error finding module name: %v", err)
	}

	groupedFiles, err := finder.FindFilesInDir(inputDir)
	if err != nil {
		log.Fatalf("Error finding files: %v", err)
	}

	return modulePath, groupedFiles
}

// Run templ fmt command
func runTemplFmt(groupedFiles *finder.GroupedFiles) {
	done := make(chan struct{})
	go func() {
		err := generator.RunTemplFmt(groupedFiles.TemplFiles, done)
		if err != nil {
			log.Fatalf("Failed to run 'templ fmt' command: %v", err)
		}
	}()
	<-done
	log.Println("Completed running 'templ fmt'")
}

// Run templ generate command
func runTemplGenerate(inputDir string) *finder.GroupedFiles {
	done := make(chan struct{})
	go func() {
		err := generator.RunTemplGenerate(done)
		if err != nil {
			log.Fatalf("Failed to run 'templ generate' command: %v", err)
		}
	}()
	<-done
	log.Println("Completed running 'templ generate'")

	groupedFiles, err := finder.FindFilesInDir(inputDir)
	if err != nil {
		log.Fatal("Error finding _templ.go files after templ generate completion:", err)
	}
	return groupedFiles
}

// Find functions in the templ Go files
func findFunctions(templGoFiles []string) []finder.FunctionToCall {
	funcs, err := finder.FindFunctionsInFiles(templGoFiles)
	if err != nil {
		log.Fatalf("Error finding functions: %v", err)
	} else if len(funcs) < 1 {
		log.Fatalf(`No components found`)
	}
	return funcs
}

// Run the generated script
func runGeneratedScript(debug bool) {
	cmd := exec.Command("go", "run", getOutputScriptPath())
	if err := cmd.Start(); err != nil {
		log.Fatalf("Error starting script: %v", err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatalf("Error running script: %v", err)
	}

	if !debug {
		if err := os.RemoveAll(outputScriptDirPath); err != nil {
			log.Fatalf("Error removing script folder: %v", err)
		}
	}
}

func usage() {
	output := fmt.Sprintf(`Usage of %[1]v:
%[1]v [flags] [subcommands]

Flags:
  -m  Set the operational mode: bundle or inline. (default "bundle").
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

// Get the version from the embedded version file
func getVersion() string {
	content, err := versionFile.ReadFile(".version")
	if err != nil {
		log.Fatalf("Error reading version file: %v", err)
	}
	return strings.TrimSpace(string(content))
}

// Print the version information
func printVersion(version, templVersion string) {
	templModulePath := "github.com/a-h/templ"
	fmt.Printf("Version: %s (built with %s@v%s)\n", version, templModulePath, templVersion)
}

// Clear and create the specified directory
func clearAndCreateDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return os.MkdirAll(dir, os.ModePerm)
}

func handleBundleMode(funcs []finder.FunctionToCall, modulePath, inputDir, outputDir string, otherFiles []string, debug bool) {
	if err := copyFilesIntoOutputDir(otherFiles, inputDir, outputDir); err != nil {
		log.Fatalf("Error copying files: %v", err)
	}
	if err := generator.GenerateForBundleMode(getOutputScriptPath(), finder.FindImports(funcs, modulePath), funcs, inputDir, outputDir); err != nil {
		log.Fatalf("Error generating script when mode=pages: %v", err)
	}
	runGeneratedScript(debug)
}

// Handle components mode
func handleInlineMode(funcs []finder.FunctionToCall, modulePath, inputDir string, debug bool) {
	if err := generator.GenerateForInlineMode(getOutputScriptPath(), finder.FindImports(funcs, modulePath), funcs, inputDir); err != nil {
		log.Fatalf("Error generating script when mode=components: %v", err)
	}
	runGeneratedScript(debug)
}

// Copy a file from one path to another
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

// Copy all files from the input directory to the output directory
func copyFilesIntoOutputDir(files []string, inputDir string, outputDir string) error {
	for _, f := range files {
		if err := copyFile(f, strings.Replace(f, inputDir, outputDir, 1)); err != nil {
			return err
		}
	}
	return nil
}

// Get the output script path
func getOutputScriptPath() string {
	return filepath.Join(outputScriptDirPath, outputScriptFileName)
}
