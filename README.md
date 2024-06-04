# Static Templ
Build static HTML websites with file based routing, using [templ](https://github.com/a-h/templ)! All components are pre-rendered at build time, and the resulting HTML files can be served using any static file server.

## Installation

```bash
go install github.com/nokacper24/static-templ
```

## Usage

```bash
templ generate
static-templ -i web/pages
```
This will generate static html files in `dist` directory. You can specify a different output directory using `-o` flag.



## Assumptions
Templ components that will be turned into html files must be **exported**, and take **no arguments**. If these conditions are not met, the component will be ignored. Your components must be in the *input* directory, their path will be mirrored in the *output* directory.

All files other than `.go` and `.templ` files will be copied to the output directory, preserving the directory structure. This allows you to include any assets and reference them using relative paths.

## Example project


## Background
Templ does support rendering components into files, as shown in their [documentation](https://templ.guide/static-rendering/generating-static-html-files-with-templ/). I wanted to avoid writing code to do so individually for each page.

It works by automatically creating a script that will render wanted components and write them intoto files, executes the script and cleans up after itself.