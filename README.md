# Static Templ

Build static HTML websites with file based routing, using [templ](https://github.com/a-h/templ)! All components are pre-rendered at build time, and the resulting HTML files can be served using any static file server.

## Installation

```bash
go install github.com/nokacper24/static-templ@latest
```

## Usage

```bash
Usage of static-templ:
static-templ [options]

Options:
  -i  Specify input directory (default "web/pages").
  -o  Specify output directory (default "dist").
  -f  Run templ fmt.
  -g  Run templ generate.

Examples:
  # Specify input and output directories
  static-templ -i web/demos -o output

  # Specify input directory, run templ generate and output to default directory
  static-templ -i web/demos -g=true
```

## Assumptions

Templ components that will be turned into html files must be **exported**, and take **no arguments**. If these conditions are not met, the component will be ignored. Your components must be in the *input* directory, their path will be mirrored in the *output* directory.

All files other than `.go` and `.templ` files will be copied to the output directory, preserving the directory structure. This allows you to include any assets and reference them using relative paths.

## Example project

## Background

Templ does support rendering components into files, as shown in their [documentation](https://templ.guide/static-rendering/generating-static-html-files-with-templ/). I wanted to avoid writing code to do so manually for each page.

static-templ creates a script that will render wanted components and write them intoto files, executes it and cleans up after itself.
