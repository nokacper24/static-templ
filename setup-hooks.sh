#!/bin/bash

# Set the hooks directory to the `.githooks` folder in the repository
git config core.hooksPath .githooks

# Make the pre-commit hook executable
chmod +x .githooks/pre-commit

echo "Git hooks have been set up."
