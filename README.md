# Static Templ

Build static HTML websites with file based routing, using [templ](https://github.com/a-h/templ)! All components are pre-rendered at build time, and the resulting HTML files can be served using any static file server.

## Installation

```bash
go install github.com/nokacper24/static-templ@latest
```

## Usage

```bash
Usage of static-templ:
static-templ [flags] [subcommands]

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
  static-templ -i web/demos -o output

  # Specify input directory, run templ generate and output to default directory
  static-templ -i web/demos -g=true

  # Display the version information
  static-templ version
```

## Assumptions

Templ components that will be turned into html files must be **exported**, and take **no arguments**. If these conditions are not met, the component will be ignored. Your components must be in the *input* directory, their path will be mirrored in the *output* directory.

All files other than `.go` and `.templ` files will be copied to the output directory, preserving the directory structure. This allows you to include any assets and reference them using relative paths.

## Example project

## Background

Templ does support rendering components into files, as shown in their [documentation](https://templ.guide/static-rendering/generating-static-html-files-with-templ/). I wanted to avoid manually writing code for each page.

`static-templ` creates a script that renders the desired components, writes them into files, executes the script, and cleans up afterward.

## Contribution

Contributions are welcome! If you have suggestions for improvements or new features, please submit an issue or a pull request.

Before submitting a pull request, please follow these steps to ensure a smooth and consistent development process:

### Setting Up Git Hooks

We use Git hooks to automate versioning and ensure code quality. After cloning the repository, you must set up the Git hooks by running the following script. This step ensures that the hooks are properly installed and executed on every commit.

1. Clone the repository:

    ```bash
    git clone https://github.com/nokacper24/static-templ.git
    cd static-templ
    ```

2. Run the setup script to install the Git hooks:

    **For Unix-based systems (Linux, macOS):**

    ```bash
      ./setup-hooks.sh
    ```

    **For Windows systems:**

    ```cmd
    setup-hooks.bat
    ```

By running the appropriate setup script, you ensure that the pre-commit hook is properly installed. This hook will automatically update the version number in the `.version` file and stage it for commit.
